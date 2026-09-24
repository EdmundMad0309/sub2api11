package service

import (
	"context"
	"fmt"
	"strings"
)

// QualityPolicy is opt-in. Legacy connectivity/HTML tests never modify membership.
type QualityPolicy struct {
	ExpectedAnswer string  `json:"expected_answer"`
	Action         string  `json:"action"`
	RemoveGroupIDs []int64 `json:"remove_group_ids"`
	AutoRestore    bool    `json:"auto_restore"`
}

func validateQualityPolicy(plan *ScheduledTestPlan) error {
	q := plan.PelicanConfig.Quality
	if q == nil {
		return nil
	}
	if plan.AutoRecover {
		return fmt.Errorf("quality plans use auto_restore, not connectivity auto_recover")
	}
	if plan.PelicanConfig.QuestionKind != "candy" {
		return fmt.Errorf("quality plans require a text answer question")
	}
	if strings.TrimSpace(q.ExpectedAnswer) == "" || len(q.ExpectedAnswer) > 4000 {
		return fmt.Errorf("expected answer must be 1–4000 bytes")
	}
	if q.Action != "remove_groups" && q.Action != "disable_scheduling" {
		return fmt.Errorf("invalid quality action")
	}
	if q.Action == "remove_groups" && len(q.RemoveGroupIDs) == 0 {
		return fmt.Errorf("select at least one group to remove")
	}
	if len(q.RemoveGroupIDs) > 100 {
		return fmt.Errorf("select at most 100 groups")
	}
	seen := map[int64]bool{}
	for _, id := range q.RemoveGroupIDs {
		if id <= 0 || seen[id] {
			return fmt.Errorf("invalid or duplicate group ID")
		}
		seen[id] = true
	}
	return nil
}

// One completed wrong answer quarantines; restoration requires every probe to pass.
// Transport errors alone are inconclusive, never evidence of degradation.
func qualityOutcome(results []*ScheduledTestResult) string {
	allPassed := len(results) > 0
	for _, r := range results {
		if r != nil && r.Status == "failed" && r.ErrorMessage == "answer_mismatch" {
			return "failed"
		}
		if r == nil || r.Status != "success" {
			allPassed = false
		}
	}
	if allPassed {
		return "passed"
	}
	return "inconclusive"
}

func (s *ScheduledTestService) ListQualityPlans(ctx context.Context) ([]*ScheduledTestPlan, error) {
	return s.planRepo.ListQualityPlans(ctx)
}
func (s *ScheduledTestService) TriggerQuality(ctx context.Context, id int64) error {
	return s.planRepo.TriggerQuality(ctx, id)
}
