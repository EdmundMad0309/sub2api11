package mihomo

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const maxDynamicProxies = 256

// dynamicProxy is kept private so credentials cannot accidentally escape in a
// status DTO or an error value.
type dynamicProxy struct {
	Scheme   string
	Host     string
	Port     int
	Username string
	Password string
}

func normalizeDynamicProxies(raw []string) ([]string, error) {
	if len(raw) > maxDynamicProxies {
		return nil, fmt.Errorf("at most %d dynamic proxies", maxDynamicProxies)
	}
	result := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for _, value := range raw {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		_, normalized, err := parseDynamicProxy(value)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	return result, nil
}

func parseDynamicProxy(raw string) (dynamicProxy, string, error) {
	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() == "" || u.User == nil || u.Fragment != "" || u.Path != "" || u.RawQuery != "" {
		return dynamicProxy{}, "", errors.New("invalid dynamic proxy URL")
	}
	scheme := strings.ToLower(u.Scheme)
	switch scheme {
	case "http", "https", "socks5", "socks5h":
	default:
		return dynamicProxy{}, "", errors.New("dynamic proxy must use http, https, socks5, or socks5h")
	}
	portText := u.Port()
	if portText == "" {
		return dynamicProxy{}, "", errors.New("dynamic proxy port is required")
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port < 1 || port > 65535 {
		return dynamicProxy{}, "", errors.New("dynamic proxy port is invalid")
	}
	password, hasPassword := u.User.Password()
	if u.User.Username() == "" || !hasPassword || password == "" {
		return dynamicProxy{}, "", errors.New("dynamic proxy username and password are required")
	}
	proxy := dynamicProxy{
		Scheme:   scheme,
		Host:     u.Hostname(),
		Port:     port,
		Username: u.User.Username(),
		Password: password,
	}
	return proxy, u.String(), nil
}

func dynamicProxyNodes(raw []string) ([]map[string]any, map[string]string, error) {
	normalized, err := normalizeDynamicProxies(raw)
	if err != nil {
		return nil, nil, err
	}
	nodes := make([]map[string]any, 0, len(normalized))
	names := make(map[string]string, len(normalized))
	for index, value := range normalized {
		proxy, _, err := parseDynamicProxy(value)
		if err != nil {
			return nil, nil, err
		}
		hash := sha256.Sum256([]byte(value))
		name := "DYNAMIC-" + hex.EncodeToString(hash[:])[:16]
		nodeType := proxy.Scheme
		if nodeType == "https" {
			nodeType = "http"
		} else if nodeType == "socks5h" {
			// Mihomo's socks5 outbound performs remote DNS for proxy hosts;
			// its YAML type is still "socks5".
			nodeType = "socks5"
		}
		node := map[string]any{
			"name":     name,
			"type":     nodeType,
			"server":   proxy.Host,
			"port":     proxy.Port,
			"username": proxy.Username,
			"password": proxy.Password,
		}
		if proxy.Scheme == "https" {
			node["tls"] = true
		}
		nodes = append(nodes, node)
		names[name] = fmt.Sprintf("dynamic-%02d", index+1)
	}
	return nodes, names, nil
}

// buildDynamicClashSubscription produces a provider-compatible YAML document.
// The manager consumes the node representation directly, while this function
// keeps the conversion reusable for diagnostics and future provider endpoints.
func buildDynamicClashSubscription(raw []string) ([]byte, error) {
	nodes, _, err := dynamicProxyNodes(raw)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(map[string]any{"proxies": nodes})
}
