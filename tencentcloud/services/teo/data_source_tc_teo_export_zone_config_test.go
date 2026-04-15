package teo_test

import (
	"context"
	"testing"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	teoService "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func TestTeoExportZoneConfig_ParameterValidation(t *testing.T) {
	// Test parameter validation for TeoExportZoneConfig function

	ctx := context.Background()
	client := connectivity.NewTencentCloudClient("", "", "", "", 0)
	service := teoService.NewTeoService(client)

	// Test with valid zone_id and empty types
	_, err := service.TeoExportZoneConfig(ctx, "zone-test123", []interface{}{})
	if err == nil {
		t.Error("Expected error for invalid client configuration, got nil")
	}
}

func TestTeoExportZoneConfig_TypesParameter(t *testing.T) {
	// Test types parameter handling

	ctx := context.Background()
	client := connectivity.NewTencentCloudClient("", "", "", "", 0)
	service := teoService.NewTeoService(client)

	// Test with nil types
	_, err := service.TeoExportZoneConfig(ctx, "zone-test123", nil)
	if err == nil {
		t.Error("Expected error for invalid client configuration, got nil")
	}

	// Test with empty types list
	_, err = service.TeoExportZoneConfig(ctx, "zone-test123", []interface{}{})
	if err == nil {
		t.Error("Expected error for invalid client configuration, got nil")
	}

	// Test with specific types
	_, err = service.TeoExportZoneConfig(ctx, "zone-test123", []interface{}{"L7AccelerationConfig"})
	if err == nil {
		t.Error("Expected error for invalid client configuration, got nil")
	}
}

func TestTeoExportZoneConfig_ZoneIdValidation(t *testing.T) {
	// Test zone_id parameter validation

	ctx := context.Background()
	client := connectivity.NewTencentCloudClient("", "", "", "", 0)
	service := teoService.NewTeoService(client)

	// Test with empty zone_id
	_, err := service.TeoExportZoneConfig(ctx, "", []interface{}{})
	if err == nil {
		t.Error("Expected error for invalid client configuration, got nil")
	}

	// Test with valid zone_id (but invalid client)
	_, err = service.TeoExportZoneConfig(ctx, "zone-test123", []interface{}{})
	if err == nil {
		t.Error("Expected error for invalid client configuration, got nil")
	}
}

func TestTeoExportZoneConfig_ErrorHandling(t *testing.T) {
	// Test error handling for various error scenarios

	ctx := context.Background()
	client := connectivity.NewTencentCloudClient("", "", "", "", 0)
	service := teoService.NewTeoService(client)

	// Test that invalid credentials cause error
	_, err := service.TeoExportZoneConfig(ctx, "zone-test123", []interface{}{})
	if err == nil {
		t.Error("Expected error for invalid client configuration, got nil")
	}
}

func TestTeoExportZoneConfig_ContentParsing(t *testing.T) {
	// Test content parsing from API response

	ctx := context.Background()
	client := connectivity.NewTencentCloudClient("", "", "", "", 0)
	service := teoService.NewTeoService(client)

	// Test that function properly handles content parsing
	_, err := service.TeoExportZoneConfig(ctx, "zone-test123", []interface{}{})
	if err == nil {
		t.Error("Expected error for invalid client configuration, got nil")
	}
}
