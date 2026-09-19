package service

import "context"

const SettingKeyOpenAICodexTicketStrategy = "openai_codex_ticket_strategy"

func NormalizeCodexTicketStrategy(value string) string {
	if value == "fixed" {
		return "fixed"
	}
	return "standby"
}

// Read once per harvest round; missing settings retain early refresh.
func (s *SettingService) GetCodexTicketStrategy(ctx context.Context) string {
	if s == nil || s.settingRepo == nil {
		return "standby"
	}
	value, err := s.settingRepo.GetValue(ctx, SettingKeyOpenAICodexTicketStrategy)
	if err != nil {
		return "standby"
	}
	return NormalizeCodexTicketStrategy(value)
}
