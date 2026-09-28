package config_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	configsvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/config"
)

// ---- gomonkey-based unit tests for tencentcloud_config_update_config_deliver ----
// Run with: go test ./tencentcloud/services/config/ -run "TestConfigUpdateConfigDeliver" -v -count=1 -gcflags="all=-l"

type mockMetaConfigUpdateConfigDeliver struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaConfigUpdateConfigDeliver) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaConfigUpdateConfigDeliver{}

func newMockMetaConfigUpdateConfigDeliver() *mockMetaConfigUpdateConfigDeliver {
	return &mockMetaConfigUpdateConfigDeliver{client: &connectivity.TencentCloudClient{}}
}

func ptrUint64ConfigDeliver(v uint64) *uint64 {
	return &v
}

func ptrStringConfigDeliver(s string) *string {
	return &s
}

// TestConfigUpdateConfigDeliver_Schema validates the schema definition.
func TestConfigUpdateConfigDeliver_Schema(t *testing.T) {
	res := configsvc.ResourceTencentCloudConfigUpdateConfigDeliver()

	statusField, ok := res.Schema["status"]
	assert.True(t, ok, "status should exist in schema")
	assert.Equal(t, schema.TypeInt, statusField.Type)
	assert.True(t, statusField.Required)

	createTimeField, ok := res.Schema["create_time"]
	assert.True(t, ok, "create_time should exist in schema")
	assert.Equal(t, schema.TypeString, createTimeField.Type)
	assert.True(t, createTimeField.Computed)
}

// TestConfigUpdateConfigDeliver_Create verifies Create sets the token id and calls UpdateConfigDeliver.
func TestConfigUpdateConfigDeliver_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaConfigUpdateConfigDeliver().client, "UseConfigV20220802Client", configClient)

	var capturedRequest *configv20220802.UpdateConfigDeliverRequest
	patches.ApplyMethodFunc(configClient, "UpdateConfigDeliverWithContext", func(ctx context.Context, request *configv20220802.UpdateConfigDeliverRequest) (*configv20220802.UpdateConfigDeliverResponse, error) {
		capturedRequest = request
		resp := configv20220802.NewUpdateConfigDeliverResponse()
		resp.Response = &configv20220802.UpdateConfigDeliverResponseParams{
			RequestId: ptrStringConfigDeliver("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeConfigDeliver called by the read after update.
	patches.ApplyMethodFunc(configClient, "DescribeConfigDeliver", func(request *configv20220802.DescribeConfigDeliverRequest) (*configv20220802.DescribeConfigDeliverResponse, error) {
		resp := configv20220802.NewDescribeConfigDeliverResponse()
		resp.Response = &configv20220802.DescribeConfigDeliverResponseParams{
			DeliverName:        ptrStringConfigDeliver("tf-example-deliver"),
			TargetArn:          ptrStringConfigDeliver("qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"),
			Status:             ptrUint64ConfigDeliver(1),
			CreateTime:         ptrStringConfigDeliver("2024-01-01 00:00:00"),
			DeliverPrefix:      ptrStringConfigDeliver("config"),
			DeliverType:        ptrStringConfigDeliver("COS"),
			DeliverContentType: ptrUint64ConfigDeliver(3),
			RequestId:          ptrStringConfigDeliver("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfigUpdateConfigDeliver()
	res := configsvc.ResourceTencentCloudConfigUpdateConfigDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"status":               1,
		"deliver_name":         "tf-example-deliver",
		"target_arn":           "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket",
		"deliver_prefix":       "config",
		"deliver_type":         "COS",
		"deliver_content_type": 3,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	// Verify status was passed to the update request.
	assert.NotNil(t, capturedRequest.Status)
	assert.Equal(t, uint64(1), *capturedRequest.Status)

	// Verify status was read back into state.
	assert.Equal(t, 1, d.Get("status").(int))
	assert.Equal(t, "COS", d.Get("deliver_type").(string))
	assert.Equal(t, 3, d.Get("deliver_content_type").(int))
	assert.Equal(t, "2024-01-01 00:00:00", d.Get("create_time").(string))
}

// TestConfigUpdateConfigDeliver_Update verifies field change triggers UpdateConfigDeliver.
func TestConfigUpdateConfigDeliver_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaConfigUpdateConfigDeliver().client, "UseConfigV20220802Client", configClient)

	var capturedRequest *configv20220802.UpdateConfigDeliverRequest
	patches.ApplyMethodFunc(configClient, "UpdateConfigDeliverWithContext", func(ctx context.Context, request *configv20220802.UpdateConfigDeliverRequest) (*configv20220802.UpdateConfigDeliverResponse, error) {
		capturedRequest = request
		resp := configv20220802.NewUpdateConfigDeliverResponse()
		resp.Response = &configv20220802.UpdateConfigDeliverResponseParams{
			RequestId: ptrStringConfigDeliver("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(configClient, "DescribeConfigDeliver", func(request *configv20220802.DescribeConfigDeliverRequest) (*configv20220802.DescribeConfigDeliverResponse, error) {
		resp := configv20220802.NewDescribeConfigDeliverResponse()
		resp.Response = &configv20220802.DescribeConfigDeliverResponseParams{
			DeliverName:        ptrStringConfigDeliver("tf-example-deliver"),
			TargetArn:          ptrStringConfigDeliver("qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"),
			Status:             ptrUint64ConfigDeliver(1),
			CreateTime:         ptrStringConfigDeliver("2024-01-01 00:00:00"),
			DeliverPrefix:      ptrStringConfigDeliver("config"),
			DeliverType:        ptrStringConfigDeliver("COS"),
			DeliverContentType: ptrUint64ConfigDeliver(3),
			RequestId:          ptrStringConfigDeliver("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfigUpdateConfigDeliver()
	res := configsvc.ResourceTencentCloudConfigUpdateConfigDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"status":               1,
		"deliver_name":         "tf-example-deliver",
		"target_arn":           "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket",
		"deliver_prefix":       "config",
		"deliver_type":         "COS",
		"deliver_content_type": 3,
	})
	d.SetId("fake-token-id")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify deliver_content_type was passed to the update request.
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.DeliverContentType)
	assert.Equal(t, uint64(3), *capturedRequest.DeliverContentType)

	// Verify deliver_content_type was read back into state.
	assert.Equal(t, 3, d.Get("deliver_content_type").(int))
}

// TestConfigUpdateConfigDeliver_ReadValue verifies Read sets fields when API returns values.
func TestConfigUpdateConfigDeliver_ReadValue(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaConfigUpdateConfigDeliver().client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "DescribeConfigDeliver", func(request *configv20220802.DescribeConfigDeliverRequest) (*configv20220802.DescribeConfigDeliverResponse, error) {
		resp := configv20220802.NewDescribeConfigDeliverResponse()
		resp.Response = &configv20220802.DescribeConfigDeliverResponseParams{
			DeliverName:        ptrStringConfigDeliver("tf-example-deliver"),
			TargetArn:          ptrStringConfigDeliver("qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"),
			Status:             ptrUint64ConfigDeliver(1),
			CreateTime:         ptrStringConfigDeliver("2024-01-01 00:00:00"),
			DeliverPrefix:      ptrStringConfigDeliver("config"),
			DeliverType:        ptrStringConfigDeliver("COS"),
			DeliverContentType: ptrUint64ConfigDeliver(3),
			RequestId:          ptrStringConfigDeliver("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfigUpdateConfigDeliver()
	res := configsvc.ResourceTencentCloudConfigUpdateConfigDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"status": 0,
	})
	d.SetId("fake-token-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify the API values were read back into state.
	assert.Equal(t, 1, d.Get("status").(int))
	assert.Equal(t, "tf-example-deliver", d.Get("deliver_name").(string))
	assert.Equal(t, "COS", d.Get("deliver_type").(string))
	assert.Equal(t, 3, d.Get("deliver_content_type").(int))
	assert.Equal(t, "2024-01-01 00:00:00", d.Get("create_time").(string))
}

// TestConfigUpdateConfigDeliver_ReadNil verifies Read clears id when service returns nil response.
func TestConfigUpdateConfigDeliver_ReadNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch the service-layer DescribeConfigDeliver to return a nil response with no error,
	// which exercises the resource-level nil guard (d.SetId("")).
	patches.ApplyMethodFunc(&configsvc.ConfigService{}, "DescribeConfigDeliver", func(ctx context.Context) (*configv20220802.DescribeConfigDeliverResponseParams, error) {
		return nil, nil
	})

	meta := newMockMetaConfigUpdateConfigDeliver()
	res := configsvc.ResourceTencentCloudConfigUpdateConfigDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"status": 1,
	})
	d.SetId("fake-token-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// The id should be cleared when the response is nil.
	assert.Empty(t, d.Id())
}
