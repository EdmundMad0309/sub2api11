package service

import (
	"context"
	"strconv"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/service/basispoints"
)

const (
	SettingKeyExcelBPSImageRelayEnabled      = "excel_bps_image_relay_enabled"
	SettingKeyExcelBPSImageBaseURL           = "excel_bps_image_base_url"
	SettingKeyExcelBPSImageRelayMaxRequests = "excel_bps_image_relay_max_requests"
)

// ExcelBPSImageRelayDefaultMaxRequests 是未配置时的在途请求上限，与图片中继
// 准入中间件的历史硬编码值保持一致。
const ExcelBPSImageRelayDefaultMaxRequests = 32

// ExcelBPSImageRelayMaxRequestsUpperBound 限制单实例图片暂存压力：
// 磁盘最多暂存 512 张图片，更高的在途并发没有意义。
const ExcelBPSImageRelayMaxRequestsUpperBound = 512

type ExcelBPSImageRelaySettings struct {
	Enabled     bool
	BaseURL     string
	MaxRequests int
}

func normalizeExcelBPSImageRelayMaxRequests(maxRequests int) int {
	if maxRequests <= 0 {
		return ExcelBPSImageRelayDefaultMaxRequests
	}
	if maxRequests > ExcelBPSImageRelayMaxRequestsUpperBound {
		return ExcelBPSImageRelayMaxRequestsUpperBound
	}
	return maxRequests
}

func normalizeExcelBPSImageRelaySettings(enabled bool, baseURL string, maxRequests int) (ExcelBPSImageRelaySettings, error) {
	baseURL = strings.TrimSpace(baseURL)
	if enabled || baseURL != "" {
		if err := basispoints.ValidateImageRelayOrigin(baseURL); err != nil {
			return ExcelBPSImageRelaySettings{}, infraerrors.BadRequest("INVALID_EXCEL_BPS_IMAGE_BASE_URL", err.Error())
		}
	}
	return ExcelBPSImageRelaySettings{
		Enabled:     enabled,
		BaseURL:     strings.TrimRight(baseURL, "/"),
		MaxRequests: normalizeExcelBPSImageRelayMaxRequests(maxRequests),
	}, nil
}

// Read only these settings for each request so saves take effect immediately,
// including on other instances sharing the settings database.
func (s *SettingService) GetExcelBPSImageRelaySettings(ctx context.Context) (ExcelBPSImageRelaySettings, error) {
	if s == nil || s.settingRepo == nil {
		return ExcelBPSImageRelaySettings{}, nil
	}
	dbCtx, cancel := context.WithTimeout(ctx, gatewayForwardingDBTimeout)
	defer cancel()
	values, err := s.settingRepo.GetMultiple(dbCtx, []string{
		SettingKeyExcelBPSImageRelayEnabled,
		SettingKeyExcelBPSImageBaseURL,
		SettingKeyExcelBPSImageRelayMaxRequests,
	})
	if err != nil {
		return ExcelBPSImageRelaySettings{}, infraerrors.ServiceUnavailable("EXCEL_BPS_IMAGE_SETTINGS_UNAVAILABLE", "Excel BPS image settings are unavailable")
	}
	maxRequests, _ := strconv.Atoi(strings.TrimSpace(values[SettingKeyExcelBPSImageRelayMaxRequests]))
	return normalizeExcelBPSImageRelaySettings(
		values[SettingKeyExcelBPSImageRelayEnabled] == "true",
		values[SettingKeyExcelBPSImageBaseURL],
		maxRequests,
	)
}
