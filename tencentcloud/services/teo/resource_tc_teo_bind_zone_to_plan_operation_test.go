package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMetaBindZoneToPlan implements tccommon.ProviderMeta
type mockMetaBindZoneToPlan struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaBindZoneToPlan) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaBindZoneToPlan{}

func newMockMetaBindZoneToPlan() *mockMetaBindZoneToPlan {
	return &mockMetaBindZoneToPlan{client: &connectivity.TencentCloudClient{}}
}

func ptrStringBindZoneToPlan(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestBindZoneToPlanOperation" -v -count=1 -gcflags="all=-l"

// TestBindZoneToPlanOperation_Success tests successful binding
func TestBindZoneToPlanOperation_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBindZoneToPlan().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "BindZoneToPlanWithContext", func(_ context.Context, request *teov20220901.BindZoneToPlanRequest) (*teov20220901.BindZoneToPlanResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.NotNil(t, request.PlanId)
		assert.Equal(t, "edgeone-12345678", *request.PlanId)

		resp := teov20220901.NewBindZoneToPlanResponse()
		resp.Response = &teov20220901.BindZoneToPlanResponseParams{
			RequestId: ptrStringBindZoneToPlan("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBindZoneToPlan()
	res := teo.ResourceTencentCloudTeoBindZoneToPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"plan_id": "edgeone-12345678",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
}

// TestBindZoneToPlanOperation_APIError tests API error handling
func TestBindZoneToPlanOperation_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBindZoneToPlan().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "BindZoneToPlanWithContext", func(_ context.Context, request *teov20220901.BindZoneToPlanRequest) (*teov20220901.BindZoneToPlanResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newMockMetaBindZoneToPlan()
	res := teo.ResourceTencentCloudTeoBindZoneToPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"plan_id": "edgeone-invalid",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestBindZoneToPlanOperation_Read tests Read is no-op
func TestBindZoneToPlanOperation_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoBindZoneToPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"plan_id": "edgeone-12345678",
	})
	d.SetId("test-id")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestBindZoneToPlanOperation_Delete tests Delete is no-op
func TestBindZoneToPlanOperation_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoBindZoneToPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"plan_id": "edgeone-12345678",
	})
	d.SetId("test-id")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestBindZoneToPlanOperation_Schema validates schema definition
func TestBindZoneToPlanOperation_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoBindZoneToPlan()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "plan_id")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	planId := res.Schema["plan_id"]
	assert.Equal(t, schema.TypeString, planId.Type)
	assert.True(t, planId.Required)
	assert.True(t, planId.ForceNew)
}
