package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

const SettingKeyOpenAICodexTicketHarvestScope = "openai_codex_ticket_harvest_scope"

// CodexTicketHarvestScope affects background harvesting only, not request routing
// or existing tickets. Selected with no groups intentionally harvests nothing.
type CodexTicketHarvestScope struct {
	Mode     string  `json:"mode"`
	GroupIDs []int64 `json:"group_ids"`
}

func NormalizeCodexTicketHarvestScope(scope CodexTicketHarvestScope) (CodexTicketHarvestScope, error) {
	if scope.Mode == "" {
		scope.Mode = "all"
	}
	if scope.Mode != "all" && scope.Mode != "selected" {
		return scope, fmt.Errorf("harvest scope mode must be all or selected")
	}
	if len(scope.GroupIDs) > 1000 {
		return scope, fmt.Errorf("at most 1000 harvest groups are allowed")
	}
	ids := append([]int64{}, scope.GroupIDs...)
	for _, id := range ids {
		if id <= 0 {
			return scope, fmt.Errorf("harvest group IDs must be positive")
		}
	}
	slices.Sort(ids)
	scope.GroupIDs = slices.Compact(ids)
	return scope, nil
}

func parseCodexTicketHarvestScope(raw string) (CodexTicketHarvestScope, error) {
	if strings.TrimSpace(raw) == "" {
		return NormalizeCodexTicketHarvestScope(CodexTicketHarvestScope{})
	}
	var scope CodexTicketHarvestScope
	if err := json.Unmarshal([]byte(raw), &scope); err != nil {
		return scope, fmt.Errorf("invalid harvest scope JSON")
	}
	// Only a missing legacy setting defaults to all. Corrupt persisted settings
	// must never accidentally widen a previously selected scope.
	if scope.Mode == "" {
		return scope, fmt.Errorf("harvest scope mode is required")
	}
	return NormalizeCodexTicketHarvestScope(scope)
}

// Read the complete scope atomically once per round. Storage errors stop the
// round rather than silently falling back to harvesting all accounts.
func (s *SettingService) GetCodexTicketHarvestScope(ctx context.Context) (CodexTicketHarvestScope, error) {
	if s == nil || s.settingRepo == nil {
		return parseCodexTicketHarvestScope("")
	}
	raw, err := s.settingRepo.GetValue(ctx, SettingKeyOpenAICodexTicketHarvestScope)
	if errors.Is(err, ErrSettingNotFound) {
		return parseCodexTicketHarvestScope("")
	}
	if err != nil {
		return CodexTicketHarvestScope{}, err
	}
	return parseCodexTicketHarvestScope(raw)
}

func (scope CodexTicketHarvestScope) includes(account *Account) bool {
	if scope.Mode == "all" {
		return true
	}
	for _, id := range account.GroupIDs {
		if slices.Contains(scope.GroupIDs, id) {
			return true
		}
	}
	return false
}

// A shared account is harvested once, at its best priority among selected
// memberships. Unselected memberships must not boost its harvest priority.
// Legacy snapshots without membership priorities fall back to account priority.
func (scope CodexTicketHarvestScope) priority(account *Account) int {
	best, found := account.Priority, false
	for _, membership := range account.AccountGroups {
		if scope.Mode == "selected" && !slices.Contains(scope.GroupIDs, membership.GroupID) {
			continue
		}
		if !found || membership.Priority < best {
			best, found = membership.Priority, true
		}
	}
	return best
}

type codexHarvestTier struct {
	Schedulable     bool
	Priority        int
	AccountPriority int
}

func (scope CodexTicketHarvestScope) tier(account *Account) codexHarvestTier {
	return codexHarvestTier{Schedulable: account.Schedulable, Priority: scope.priority(account), AccountPriority: account.Priority}
}
