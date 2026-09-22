package mihomo

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestDynamicProxyNodesAcceptBareHTTPURLsAndRedactStatus(t *testing.T) {
	raw := []string{
		"user-one:pass-one@proxy.example:2000",
		"socks5h://user-two:pass-two@proxy.example:2001",
		"user-one:pass-one@proxy.example:2000",
	}
	nodes, names, err := dynamicProxyNodes(raw)
	require.NoError(t, err)
	require.Len(t, nodes, 2)
	require.Len(t, names, 2)
	require.Equal(t, "http", nodes[0]["type"])
	require.Equal(t, "socks5", nodes[1]["type"])
	require.Equal(t, "proxy.example", nodes[0]["server"])
	require.Equal(t, 2000, nodes[0]["port"])

	m := New(t.TempDir())
	m.saved.DynamicProxies = []string{raw[0]}
	status, err := json.Marshal(m.Status())
	require.NoError(t, err)
	require.Contains(t, string(status), `"dynamic_proxies":1`)
	require.NotContains(t, string(status), "pass-one")
}

func TestDynamicProxyValidationRejectsUnsupportedOrCredentiallessInput(t *testing.T) {
	for _, raw := range []string{
		"https://proxy.example:2000",
		"ftp://user:pass@proxy.example:2000",
		"http://user:pass@proxy.example:0",
		"http://user:pass@proxy.example:2000/path",
	} {
		_, err := normalizeDynamicProxies([]string{raw})
		require.Error(t, err, raw)
	}
}

func TestDynamicClashSubscriptionContainsOnlyProxyNodes(t *testing.T) {
	b, err := buildDynamicClashSubscription([]string{"user:pass@proxy.example:2000"})
	require.NoError(t, err)
	var document struct {
		Proxies []map[string]any `yaml:"proxies"`
	}
	require.NoError(t, yaml.Unmarshal(b, &document))
	require.Len(t, document.Proxies, 1)
	require.Equal(t, "http", document.Proxies[0]["type"])
	require.NotContains(t, strings.TrimSpace(string(b)), "external-controller")
}
