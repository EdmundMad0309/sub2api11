package service

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCodexTicketHarvestCookiesFollowFreshnessWindow(t *testing.T) {
	now := time.Now()
	ticket := &openAICodexTicket{CapturedAt: now, HarvestCookies: []string{"__cflb=secret", "__oai=secret2"}}
	require.True(t, codexTicketCookiesFresh(ticket, now.Add(239*time.Second)))
	require.False(t, codexTicketCookiesFresh(ticket, now.Add(240*time.Second)))

	headers := http.Header{}
	restoreBoundCodexTicketHarvestIdentity(headers, &openAICodexTicket{
		HarvestSessionID: "fresh-session",
		CapturedAt:       now,
		HarvestCookies:   ticket.HarvestCookies,
		Model:            "gpt-6-astra",
	})
	require.Equal(t, "__cflb=secret; __oai=secret2", headers.Get("Cookie"))
}
