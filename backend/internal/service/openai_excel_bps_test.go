package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service/basispoints"
	"github.com/stretchr/testify/require"
)

func testExcelBPSAccessToken(t *testing.T, accountID string) string {
	t.Helper()
	payload, err := json.Marshal(map[string]any{
		"https://api.openai.com/auth": map[string]string{"chatgpt_account_id": accountID},
	})
	require.NoError(t, err)
	return "header." + base64.RawURLEncoding.EncodeToString(payload) + ".signature"
}

func TestExcelBPSAccountID(t *testing.T) {
	account := &Account{
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Credentials: map[string]any{"chatgpt_account_id": "stored-account"},
	}
	require.Equal(t, "stored-account", excelBPSAccountID(account, testExcelBPSAccessToken(t, "jwt-account")))

	delete(account.Credentials, "chatgpt_account_id")
	require.Equal(t, "jwt-account", excelBPSAccountID(account, testExcelBPSAccessToken(t, "jwt-account")))
	require.Empty(t, excelBPSAccountID(account, "not-a-jwt"))
	require.Empty(t, excelBPSAccountID(account, testExcelBPSAccessToken(t, "")))
}

func TestNewExcelBPSRequestHeaders(t *testing.T) {
	token := testExcelBPSAccessToken(t, "jwt-account")
	req, err := newExcelBPSRequest(context.Background(), []byte(`{"model":"gpt-5.6-sol"}`), token, "jwt-account")
	require.NoError(t, err)
	require.Equal(t, http.MethodPost, req.Method)
	require.Equal(t, basispoints.ResponsesURL, req.URL.String())
	require.Equal(t, "Bearer "+token, req.Header.Get("authorization"))
	require.Equal(t, "jwt-account", req.Header.Get("chatgpt-account-id"))
	require.Equal(t, "jwt-account", req.Header.Get("x-openai-account-id"))
	require.Equal(t, "chatgpt", req.Header.Get("x-basispoints-auth-mode"))
}
