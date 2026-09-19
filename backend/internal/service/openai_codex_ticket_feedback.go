package service

import (
	"context"
	"net/http"
	"time"
)

// A mismatch invalidates only the ticket actually sent by this request. The
// current response continues through normal accounting; never replay it here.
func (s *OpenAIGatewayService) observeCodexTicketResponse(req *http.Request, resp *http.Response, account *Account) {
	if s == nil || req == nil || resp == nil || account == nil || resp.StatusCode != http.StatusOK {
		return
	}
	sent := req.Header.Get(openAICodexTurnStateHeader)
	returned := extractOpenAICodexTurnState(resp.Header)
	if sent == "" || returned == "" || !s.openAICodexTicketEnabledContext(req.Context()) {
		return
	}
	shape, err := parseOpenAICodexTicketShape(returned)
	if err == nil && shape.Blocks == openAICodexTicketExpectedBlocks(account) && !shape.IssuedAt.After(time.Now().Add(30*time.Second)) && time.Now().Before(shape.IssuedAt.Add(time.Hour-30*time.Second)) {
		return
	}
	s.openaiCodexTicketStateMu.Lock()
	defer s.openaiCodexTicketStateMu.Unlock()
	cfg := s.openAICodexTicketConfig()
	for _, model := range cfg.Models {
		current := s.lookupCodexTicketLocked(account, model)
		if current == nil || current.State != sent || current.Revoked {
			continue
		}
		next := *current
		if current.Standby.valid(time.Now(), openAICodexTicketTargetLength(account, cfg)) {
			next = *current.Standby
		} else {
			next.Revoked = true
		}
		next.CapturedAt = time.Now()
		s.openaiCodexTickets.Store(openAICodexTicketKey(account.ID, model), &next)
		if s.accountRepo != nil {
			ctx, cancel := context.WithTimeout(context.WithoutCancel(req.Context()), time.Second)
			_ = s.accountRepo.UpdateExtra(ctx, account.ID, map[string]any{openAICodexTicketExtraKey(model): &next})
			cancel()
		}
	}
}
