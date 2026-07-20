package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestAccTencentCloudTeoPlanForZone" -v -count=1 -gcflags="all=-l"
// TestAccTencentCloudTeoPlanForZoneOperation_Success tests successful plan purchase
func TestAccTencentCloudTeoPlanForZoneOperation_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseTeoV20220901Client to return a non-nil client
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Patch CreatePlanForZone to return success
	patches.ApplyMethodFunc(teoClient, "CreatePlanForZone", func(request *teov20220901.CreatePlanForZoneRequest) (*teov20220901.CreatePlanForZoneResponse, error) {
		assert.Equal(t, "zone-27h0vbm5w1e", *request.ZoneId)
		assert.Equal(t, "sta_global", *request.PlanType)

		resp := teov20220901.NewCreatePlanForZoneResponse()
		resp.Response = &teov20220901.CreatePlanForZoneResponseParams{
			ResourceNames: []*string{
				ptrString("plan-res-123456"),
			},
			DealNames: []*string{
				ptrString("2025122000000123"),
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoPlanForZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-27h0vbm5w1e",
		"plan_type": "sta_global",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	// Verify resource_names was set
	resourceNames := d.Get("resource_names").([]interface{})
	assert.Len(t, resourceNames, 1)
	assert.Equal(t, "plan-res-123456", resourceNames[0].(string))

	// Verify deal_names was set
	dealNames := d.Get("deal_names").([]interface{})
	assert.Len(t, dealNames, 1)
	assert.Equal(t, "2025122000000123", dealNames[0].(string))
}

// TestAccTencentCloudTeoPlanForZoneOperation_MissingZoneId tests missing required parameter zone_id
func TestAccTencentCloudTeoPlanForZoneOperation_MissingZoneId(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoPlanForZone()

	// Create resource data without zone_id (empty string)
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "",
		"plan_type": "sta_global",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "zone_id is required")
}

// TestAccTencentCloudTeoPlanForZoneOperation_MissingPlanType tests missing required parameter plan_type
func TestAccTencentCloudTeoPlanForZoneOperation_MissingPlanType(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoPlanForZone()

	// Create resource data without plan_type (empty string)
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-27h0vbm5w1e",
		"plan_type": "",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plan_type is required")
}

// TestAccTencentCloudTeoPlanForZoneOperation_APIError tests API error handling
func TestAccTencentCloudTeoPlanForZoneOperation_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreatePlanForZone", func(request *teov20220901.CreatePlanForZoneRequest) (*teov20220901.CreatePlanForZoneResponse, error) {
		assert.Equal(t, "zone-invalid", *request.ZoneId)
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter.ZoneHasBeenBound, Message=Zone has been bound")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoPlanForZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-invalid",
		"plan_type": "sta_global",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ZoneHasBeenBound")
}

// TestAccTencentCloudTeoPlanForZoneOperation_EmptyResponse tests nil response handling
func TestAccTencentCloudTeoPlanForZoneOperation_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreatePlanForZone", func(request *teov20220901.CreatePlanForZoneRequest) (*teov20220901.CreatePlanForZoneResponse, error) {
		resp := teov20220901.NewCreatePlanForZoneResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoPlanForZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-27h0vbm5w1e",
		"plan_type": "sta_global",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

// TestAccTencentCloudTeoPlanForZoneOperation_Read tests Read is no-op
func TestAccTencentCloudTeoPlanForZoneOperation_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoPlanForZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-27h0vbm5w1e",
		"plan_type": "sta_global",
	})
	d.SetId("fake-token-id")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestAccTencentCloudTeoPlanForZoneOperation_Delete tests Delete is no-op
func TestAccTencentCloudTeoPlanForZoneOperation_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoPlanForZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-27h0vbm5w1e",
		"plan_type": "sta_global",
	})
	d.SetId("fake-token-id")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestAccTencentCloudTeoPlanForZoneOperation_Schema validates schema definition
func TestAccTencentCloudTeoPlanForZoneOperation_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoPlanForZone()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "plan_type")
	assert.Contains(t, res.Schema, "resource_names")
	assert.Contains(t, res.Schema, "deal_names")

	// Check zone_id
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	// Check plan_type
	planType := res.Schema["plan_type"]
	assert.Equal(t, schema.TypeString, planType.Type)
	assert.True(t, planType.Required)
	assert.True(t, planType.ForceNew)

	// Check resource_names
	resourceNames := res.Schema["resource_names"]
	assert.Equal(t, schema.TypeList, resourceNames.Type)
	assert.True(t, resourceNames.Computed)

	// Check deal_names
	dealNames := res.Schema["deal_names"]
	assert.Equal(t, schema.TypeList, dealNames.Type)
	assert.True(t, dealNames.Computed)
}
