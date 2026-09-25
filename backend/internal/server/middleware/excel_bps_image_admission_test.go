package middleware

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type bpsImageTestSettings struct {
	enabled     bool
	err         error
	maxRequests int
}

func (s bpsImageTestSettings) GetExcelBPSImageRelaySettings(context.Context) (service.ExcelBPSImageRelaySettings, error) {
	return service.ExcelBPSImageRelaySettings{Enabled: s.enabled, MaxRequests: s.maxRequests}, s.err
}

type bpsImageCountingBody struct {
	reads  *atomic.Int32
	reader io.Reader
}

func (b *bpsImageCountingBody) Read(p []byte) (int, error) { b.reads.Add(1); return b.reader.Read(p) }
func (b *bpsImageCountingBody) Close() error               { return nil }

func bpsImageTestRouter(settings bpsImageTestSettings, next gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(string(ContextKeyAPIKey), &service.APIKey{Group: &service.Group{Platform: service.PlatformOpenAI}})
		c.Next()
	})
	r.Use(ExcelBPSImageAdmission(settings, 256<<20))
	for _, path := range []string{"/responses", "/v1/responses", "/backend-api/codex/responses", "/v1/chat/completions", "/chat/completions", "/v1/messages"} {
		r.POST(path, next)
	}
	r.POST("/v1/responses/*subpath", next)
	return r
}

func TestExcelBPSImageAdmission200ConcurrentRequests(t *testing.T) {
	for _, tt := range []struct {
		name     string
		length   int64
		encoding string
		allowed  int
	}{
		{"small", 1024, "", 32},
		{"large", 32 << 20, "", 2},
		{"chunked", -1, "", 1},
		{"compressed", 1024, "gzip", 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var reads atomic.Int32
			var active atomic.Int32
			var peak atomic.Int32
			entered := make(chan struct{}, 200)
			release := make(chan struct{})
			var once sync.Once
			var wg sync.WaitGroup
			t.Cleanup(func() { once.Do(func() { close(release) }); wg.Wait() })
			r := bpsImageTestRouter(bpsImageTestSettings{enabled: true}, func(c *gin.Context) {
				n := active.Add(1)
				for old := peak.Load(); n > old; old = peak.Load() {
					if peak.CompareAndSwap(old, n) {
						break
					}
				}
				_, err := io.Copy(io.Discard, c.Request.Body)
				if err != nil {
					t.Error(err)
				}
				entered <- struct{}{}
				<-release
				active.Add(-1)
				c.Status(http.StatusNoContent)
			})
			start := make(chan struct{})
			results := make(chan *httptest.ResponseRecorder, 200)
			paths := []string{"/responses", "/v1/responses/compact", "/backend-api/codex/responses", "/chat/completions", "/v1/messages"}
			for i := 0; i < 200; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					<-start
					req := httptest.NewRequest(http.MethodPost, paths[i%len(paths)], nil)
					req.ContentLength = tt.length
					req.Header.Set("Content-Encoding", tt.encoding)
					req.Body = &bpsImageCountingBody{reads: &reads, reader: strings.NewReader("test")}
					w := httptest.NewRecorder()
					r.ServeHTTP(w, req)
					results <- w
				}(i)
			}
			close(start)
			deadline := time.After(10 * time.Second)
			for i := 0; i < 200-tt.allowed; i++ {
				select {
				case w := <-results:
					require.Equal(t, http.StatusServiceUnavailable, w.Code)
					require.Equal(t, "1", w.Header().Get("Retry-After"))
					require.Contains(t, w.Body.String(), "basispoints_image_request_busy")
				case <-deadline:
					t.Fatal("rejected requests did not finish without reading their body")
				}
			}
			for i := 0; i < tt.allowed; i++ {
				select {
				case <-entered:
				case <-deadline:
					t.Fatal("admitted requests did not enter")
				}
			}
			require.Equal(t, int32(tt.allowed), peak.Load())
			require.Equal(t, int32(tt.allowed*2), reads.Load(), "only admitted bodies may be read")
			once.Do(func() { close(release) })
			for i := 0; i < tt.allowed; i++ {
				select {
				case w := <-results:
					require.Equal(t, http.StatusNoContent, w.Code)
				case <-deadline:
					t.Fatal("admitted request stuck")
				}
			}
			wg.Wait()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/responses", strings.NewReader("next"))
			r.ServeHTTP(w, req)
			require.Equal(t, http.StatusNoContent, w.Code, "completed requests must release their budget")
		})
	}
}

func TestExcelBPSImageAdmissionLimitsAndDisabled(t *testing.T) {
	for _, tt := range []struct {
		name     string
		settings bpsImageTestSettings
		length   int64
		body     string
		status   int
		wantRead bool
	}{
		{"oversized", bpsImageTestSettings{enabled: true}, 65 << 20, "body", 413, false},
		{"settings unavailable", bpsImageTestSettings{err: errors.New("private database error")}, 4, "body", 503, false},
		{"disabled", bpsImageTestSettings{}, 65 << 20, "body", 204, true},
		{"understated length", bpsImageTestSettings{enabled: true}, 1, "body", 413, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var reads atomic.Int32
			r := bpsImageTestRouter(tt.settings, func(c *gin.Context) {
				_, err := io.Copy(io.Discard, c.Request.Body)
				if err != nil {
					c.Status(413)
					return
				}
				c.Status(204)
			})
			req := httptest.NewRequest(http.MethodPost, "/responses", nil)
			req.ContentLength = tt.length
			req.Body = &bpsImageCountingBody{reads: &reads, reader: strings.NewReader(tt.body)}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			require.Equal(t, tt.status, w.Code)
			require.Equal(t, tt.wantRead, reads.Load() > 0)
			require.NotContains(t, w.Body.String(), "private database error")
		})
	}
}

func TestExcelBPSImageAdmissionReleaseIsIdempotent(t *testing.T) {
	budget := &bpsImageAdmissionBudget{}
	release, ok := budget.acquire(bpsImageBudgetBytes, 0)
	require.True(t, ok)
	_, ok = budget.acquire(1, 0)
	require.False(t, ok)
	release()
	release()
	require.Zero(t, budget.bytes)
	require.Zero(t, budget.requests)
}

func TestExcelBPSImageAdmissionReleasesAfterCancellation(t *testing.T) {
	entered := make(chan struct{})
	done := make(chan struct{})
	r := bpsImageTestRouter(bpsImageTestSettings{enabled: true}, func(c *gin.Context) {
		if c.GetHeader("Hold") == "true" {
			close(entered)
			<-c.Request.Context().Done()
		}
		c.Status(http.StatusNoContent)
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest(http.MethodPost, "/responses", nil).WithContext(ctx)
	req.ContentLength = -1
	req.Header.Set("Hold", "true")
	go func() {
		r.ServeHTTP(httptest.NewRecorder(), req)
		close(done)
	}()
	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("request did not enter")
	}
	cancel()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("cancelled request did not finish")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/responses", nil))
	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestExcelBPSImageAdmissionConfiguredMaxRequests(t *testing.T) {
	const configured = 2
	const excess = 3
	var peak atomic.Int32
	var active atomic.Int32
	entered := make(chan struct{}, configured)
	release := make(chan struct{})
	var wg sync.WaitGroup
	t.Cleanup(wg.Wait)
	r := bpsImageTestRouter(bpsImageTestSettings{enabled: true, maxRequests: configured}, func(c *gin.Context) {
		n := active.Add(1)
		for old := peak.Load(); n > old; old = peak.Load() {
			if peak.CompareAndSwap(old, n) {
				break
			}
		}
		entered <- struct{}{}
		<-release
		active.Add(-1)
		c.Status(http.StatusNoContent)
	})

	// Phase A: fill the gate with `configured` held requests.
	for i := 0; i < configured; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/responses", strings.NewReader("{}")))
		}()
	}
	for i := 0; i < configured; i++ {
		select {
		case <-entered:
		case <-time.After(5 * time.Second):
			t.Fatal("configured number of requests did not enter")
		}
	}
	require.Equal(t, int32(configured), peak.Load(), "in-flight peak must respect the configured max")

	// Phase B: while the gate is full, synchronous attempts must be rejected
	// immediately with 503 (never queued).
	rejected := 0
	for i := 0; i < excess; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/responses", strings.NewReader("{}")))
		if w.Code == http.StatusServiceUnavailable {
			rejected++
		}
	}
	require.Equal(t, excess, rejected, "excess requests must be rejected with 503")

	close(release)
	wg.Wait()
}
