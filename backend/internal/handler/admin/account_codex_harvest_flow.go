package admin

import (
	"net/http"
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type setCodexSkipHarvestRequest struct {
	SkipHarvest *bool `json:"skip_harvest" binding:"required"`
}

// GetCodexHarvestFlow returns the live harvest pipeline for the admin flow page.
// GET /api/v1/admin/accounts/codex-harvest-flow
func (h *AccountHandler) GetCodexHarvestFlow(c *gin.Context) {
	if h == nil {
		response.Error(c, http.StatusServiceUnavailable, "Account handler not available")
		return
	}
	ctx := c.Request.Context()
	var accounts []service.Account
	if h.adminService != nil {
		listed, err := h.adminService.ListAccountsForSchedulerScoreFilter(ctx, service.PlatformOpenAI, "", "", "", 0, "")
		if err != nil {
			response.ErrorFrom(c, err)
			return
		}
		accounts = listed
	}
	response.Success(c, service.BuildCodexHarvestFlow(ctx, h.cfg, h.codexTicketSettings, accounts))
}

// SetCodexSkipHarvest marks an account as harvest-excluded. Fail-closed gated
// routing then refuses leftover tickets on that account.
// PUT /api/v1/admin/accounts/:id/codex-skip-harvest
func (h *AccountHandler) SetCodexSkipHarvest(c *gin.Context) {
	if h == nil || h.adminService == nil {
		response.Error(c, http.StatusServiceUnavailable, "Account handler not available")
		return
	}
	accountID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || accountID <= 0 {
		response.BadRequest(c, "Invalid account ID")
		return
	}
	var req setCodexSkipHarvestRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.SkipHarvest == nil {
		response.BadRequest(c, "skip_harvest is required")
		return
	}
	if err := h.adminService.UpdateAccountExtra(c.Request.Context(), accountID, map[string]any{
		service.OpenAICodexSkipHarvestExtraKey: *req.SkipHarvest,
	}); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"account_id": accountID, "skip_harvest": *req.SkipHarvest})
}
