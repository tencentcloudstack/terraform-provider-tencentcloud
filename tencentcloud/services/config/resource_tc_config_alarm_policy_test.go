package config_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	config_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcconfig "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/config"
)

func TestAccTencentCloudConfigAlarmPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAlarmPolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_config_alarm_policy.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_config_alarm_policy.example", "alarm_policy_id"),
					resource.TestCheckResourceAttr("tencentcloud_config_alarm_policy.example", "name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_config_alarm_policy.example", "status", "1"),
				),
			},
			{
				Config: testAccConfigAlarmPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_config_alarm_policy.example", "name", "tf-example-update"),
					resource.TestCheckResourceAttr("tencentcloud_config_alarm_policy.example", "status", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_config_alarm_policy.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccConfigAlarmPolicy = `
resource "tencentcloud_config_alarm_policy" "example" {
  name                   = "tf-example"
  type                   = 1
  event_scope            = [1]
  risk_level             = [1, 2]
  notice_time            = "09:30:00~23:30:00"
  notification_mechanism = "实时发送"
  status                 = 1
  notice_period          = [1, 2, 3, 4, 5]
  description            = "tf example alarm policy"
}
`

const testAccConfigAlarmPolicyUpdate = `
resource "tencentcloud_config_alarm_policy" "example" {
  name                   = "tf-example-update"
  type                   = 1
  event_scope            = [1]
  risk_level             = [1, 2, 3]
  notice_time            = "00:00:00~23:59:59"
  notification_mechanism = "实时发送"
  status                 = 2
  notice_period          = [1, 2, 3, 4, 5, 6, 7]
  description            = "updated alarm policy"
}
`

type mockMetaForConfigAlarmPolicy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForConfigAlarmPolicy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForConfigAlarmPolicy{}

func newMockMetaForConfigAlarmPolicy() *mockMetaForConfigAlarmPolicy {
	return &mockMetaForConfigAlarmPolicy{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func buildConfigAlarmPolicyRsp() *config_sdk.AlarmPolicyRsp {
	return &config_sdk.AlarmPolicyRsp{
		AlarmPolicyId:         helper.IntUint64(123456),
		Name:                  helper.String("tf-example"),
		Type:                  helper.IntInt64(1),
		EventScope:            []*int64{helper.IntInt64(1)},
		RiskLevel:             []*int64{helper.IntInt64(1), helper.IntInt64(2)},
		NoticePeriod:          []*int64{helper.IntInt64(1), helper.IntInt64(2), helper.IntInt64(3), helper.IntInt64(4), helper.IntInt64(5)},
		NoticeTime:            helper.String("09:30:00~23:30:00"),
		NotificationMechanism: helper.String("实时发送"),
		Status:                helper.IntInt64(1),
		Description:           helper.String("tf example alarm policy"),
	}
}

func buildConfigAlarmPolicyListResponse() *config_sdk.ListAlarmPolicyResponse {
	resp := config_sdk.NewListAlarmPolicyResponse()
	resp.Response = &config_sdk.ListAlarmPolicyResponseParams{
		Total: helper.IntInt64(1),
		AlarmPolicyList: []*config_sdk.AlarmPolicyRsp{
			buildConfigAlarmPolicyRsp(),
		},
	}
	return resp
}

// TestConfigAlarmPolicy_Schema verifies the schema fields match the design.
func TestConfigAlarmPolicy_Schema(t *testing.T) {
	res := svcconfig.ResourceTencentCloudConfigAlarmPolicy()

	cases := []struct {
		field    string
		typ      schema.ValueType
		required bool
		optional bool
		computed bool
		forceNew bool
	}{
		{"name", schema.TypeString, true, false, false, false},
		{"type", schema.TypeInt, true, false, false, true},
		{"event_scope", schema.TypeList, true, false, false, false},
		{"risk_level", schema.TypeList, true, false, false, false},
		{"notice_time", schema.TypeString, true, false, false, false},
		{"notification_mechanism", schema.TypeString, true, false, false, false},
		{"status", schema.TypeInt, true, false, false, false},
		{"notice_period", schema.TypeList, true, false, false, false},
		{"description", schema.TypeString, false, true, false, false},
		{"alarm_policy_id", schema.TypeString, false, false, true, false},
	}

	for _, c := range cases {
		s, ok := res.Schema[c.field]
		assert.True(t, ok, "field %s should exist in schema", c.field)
		assert.Equal(t, c.typ, s.Type, "field %s type mismatch", c.field)
		if c.required {
			assert.True(t, s.Required, "field %s should be required", c.field)
		}
		if c.optional {
			assert.True(t, s.Optional, "field %s should be optional", c.field)
		}
		if c.computed {
			assert.True(t, s.Computed, "field %s should be computed", c.field)
		}
		if c.forceNew {
			assert.True(t, s.ForceNew, "field %s should be ForceNew", c.field)
		}
	}
}

// TestConfigAlarmPolicy_Create covers the Create business logic end-to-end with mocked cloud API.
func TestConfigAlarmPolicy_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAlarmPolicy().client, "UseConfigV20220802Client", configClient)

	var capturedRequest *config_sdk.AddAlarmPolicyRequest
	patches.ApplyMethodFunc(configClient, "AddAlarmPolicyWithContext", func(ctx context.Context, request *config_sdk.AddAlarmPolicyRequest) (*config_sdk.AddAlarmPolicyResponse, error) {
		capturedRequest = request
		resp := config_sdk.NewAddAlarmPolicyResponse()
		resp.Response = &config_sdk.AddAlarmPolicyResponseParams{
			AlarmPolicyId: helper.IntUint64(123456),
		}
		return resp, nil
	})
	// Read after create uses the service layer ListAlarmPolicy.
	patches.ApplyMethodFunc(configClient, "ListAlarmPolicy", func(request *config_sdk.ListAlarmPolicyRequest) (*config_sdk.ListAlarmPolicyResponse, error) {
		return buildConfigAlarmPolicyListResponse(), nil
	})

	meta := newMockMetaForConfigAlarmPolicy()
	res := svcconfig.ResourceTencentCloudConfigAlarmPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                   "tf-example",
		"type":                   1,
		"event_scope":            []interface{}{1},
		"risk_level":             []interface{}{1, 2},
		"notice_time":            "09:30:00~23:30:00",
		"notification_mechanism": "实时发送",
		"status":                 1,
		"notice_period":          []interface{}{1, 2, 3, 4, 5},
		"description":            "tf example alarm policy",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "123456", d.Id())
	assert.Equal(t, "123456", d.Get("alarm_policy_id"))

	assert.NotNil(t, capturedRequest, "AddAlarmPolicyWithContext should be called")
	assert.Equal(t, "tf-example", *capturedRequest.Name)
	assert.Equal(t, int64(1), *capturedRequest.Type)
	assert.Equal(t, "tf example alarm policy", *capturedRequest.Description)
	assert.Equal(t, int64(1), *capturedRequest.EventScope[0])
	assert.Equal(t, int64(2), *capturedRequest.RiskLevel[1])
	assert.Len(t, capturedRequest.NoticePeriod, 5)
	assert.Equal(t, "09:30:00~23:30:00", *capturedRequest.NoticeTime)
	assert.Equal(t, "实时发送", *capturedRequest.NotificationMechanism)
	assert.Equal(t, int64(1), *capturedRequest.Status)
}

// TestConfigAlarmPolicy_Read covers the Read business logic (populate state from cloud).
func TestConfigAlarmPolicy_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAlarmPolicy().client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "ListAlarmPolicy", func(request *config_sdk.ListAlarmPolicyRequest) (*config_sdk.ListAlarmPolicyResponse, error) {
		return buildConfigAlarmPolicyListResponse(), nil
	})

	meta := newMockMetaForConfigAlarmPolicy()
	res := svcconfig.ResourceTencentCloudConfigAlarmPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("123456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "123456", d.Id())
	assert.Equal(t, "tf-example", d.Get("name"))
	assert.Equal(t, 1, d.Get("type"))
	assert.Equal(t, "123456", d.Get("alarm_policy_id"))
	assert.Equal(t, "09:30:00~23:30:00", d.Get("notice_time"))
	assert.Equal(t, "实时发送", d.Get("notification_mechanism"))
	assert.Equal(t, 1, d.Get("status"))
	assert.Equal(t, "tf example alarm policy", d.Get("description"))

	eventScope := d.Get("event_scope").([]interface{})
	assert.Len(t, eventScope, 1)
	assert.Equal(t, 1, eventScope[0])

	riskLevel := d.Get("risk_level").([]interface{})
	assert.Len(t, riskLevel, 2)
	assert.Equal(t, 2, riskLevel[1])

	noticePeriod := d.Get("notice_period").([]interface{})
	assert.Len(t, noticePeriod, 5)
	assert.Equal(t, 5, noticePeriod[4])
}

// TestConfigAlarmPolicy_Update covers the Update business logic (field change).
func TestConfigAlarmPolicy_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAlarmPolicy().client, "UseConfigV20220802Client", configClient)

	var capturedRequest *config_sdk.UpdateAlarmPolicyRequest
	patches.ApplyMethodFunc(configClient, "UpdateAlarmPolicyWithContext", func(ctx context.Context, request *config_sdk.UpdateAlarmPolicyRequest) (*config_sdk.UpdateAlarmPolicyResponse, error) {
		capturedRequest = request
		resp := config_sdk.NewUpdateAlarmPolicyResponse()
		resp.Response = &config_sdk.UpdateAlarmPolicyResponseParams{}
		return resp, nil
	})
	patches.ApplyMethodFunc(configClient, "ListAlarmPolicy", func(request *config_sdk.ListAlarmPolicyRequest) (*config_sdk.ListAlarmPolicyResponse, error) {
		return buildConfigAlarmPolicyListResponse(), nil
	})

	meta := newMockMetaForConfigAlarmPolicy()
	res := svcconfig.ResourceTencentCloudConfigAlarmPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                   "tf-example",
		"type":                   1,
		"event_scope":            []interface{}{1},
		"risk_level":             []interface{}{1, 2},
		"notice_time":            "09:30:00~23:30:00",
		"notification_mechanism": "实时发送",
		"status":                 1,
		"notice_period":          []interface{}{1, 2, 3, 4, 5},
		"description":            "tf example alarm policy",
	})
	d.SetId("123456")

	// Simulate a name + status change.
	assert.NoError(t, d.Set("name", "tf-example-update"))
	assert.NoError(t, d.Set("status", 2))

	err := res.Update(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest, "UpdateAlarmPolicyWithContext should be called")
	assert.Equal(t, uint64(123456), *capturedRequest.AlarmPolicyId)
	assert.Equal(t, "tf-example-update", *capturedRequest.Name)
	assert.Equal(t, int64(2), *capturedRequest.Status)
	assert.Equal(t, int64(1), *capturedRequest.EventScope[0])
	assert.Len(t, capturedRequest.RiskLevel, 2)
	assert.Len(t, capturedRequest.NoticePeriod, 5)
	assert.Equal(t, "tf example alarm policy", *capturedRequest.Description)
}

// TestConfigAlarmPolicy_Delete covers the Delete business logic.
func TestConfigAlarmPolicy_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAlarmPolicy().client, "UseConfigV20220802Client", configClient)

	var capturedRequest *config_sdk.DeleteAlarmPolicyRequest
	patches.ApplyMethodFunc(configClient, "DeleteAlarmPolicyWithContext", func(ctx context.Context, request *config_sdk.DeleteAlarmPolicyRequest) (*config_sdk.DeleteAlarmPolicyResponse, error) {
		capturedRequest = request
		resp := config_sdk.NewDeleteAlarmPolicyResponse()
		resp.Response = &config_sdk.DeleteAlarmPolicyResponseParams{}
		return resp, nil
	})

	meta := newMockMetaForConfigAlarmPolicy()
	res := svcconfig.ResourceTencentCloudConfigAlarmPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("123456")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest, "DeleteAlarmPolicyWithContext should be called")
	assert.Equal(t, uint64(123456), *capturedRequest.AlarmPolicyId)
}

// TestConfigAlarmPolicy_Create_NilResponse verifies the create nil-id guard returns an error.
func TestConfigAlarmPolicy_Create_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAlarmPolicy().client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "AddAlarmPolicyWithContext", func(ctx context.Context, request *config_sdk.AddAlarmPolicyRequest) (*config_sdk.AddAlarmPolicyResponse, error) {
		resp := config_sdk.NewAddAlarmPolicyResponse()
		resp.Response = &config_sdk.AddAlarmPolicyResponseParams{}
		return resp, nil
	})

	meta := newMockMetaForConfigAlarmPolicy()
	res := svcconfig.ResourceTencentCloudConfigAlarmPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                   "tf-example",
		"type":                   1,
		"event_scope":            []interface{}{1},
		"risk_level":             []interface{}{1},
		"notice_time":            "09:30:00~23:30:00",
		"notification_mechanism": "实时发送",
		"status":                 1,
		"notice_period":          []interface{}{1},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

// TestConfigAlarmPolicy_Read_NotFound verifies that a nil match clears the state id.
func TestConfigAlarmPolicy_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAlarmPolicy().client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "ListAlarmPolicy", func(request *config_sdk.ListAlarmPolicyRequest) (*config_sdk.ListAlarmPolicyResponse, error) {
		resp := config_sdk.NewListAlarmPolicyResponse()
		resp.Response = &config_sdk.ListAlarmPolicyResponseParams{
			Total:           helper.IntInt64(0),
			AlarmPolicyList: []*config_sdk.AlarmPolicyRsp{},
		}
		return resp, nil
	})

	meta := newMockMetaForConfigAlarmPolicy()
	res := svcconfig.ResourceTencentCloudConfigAlarmPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("999999")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestConfigAlarmPolicy_Create_APIError verifies create propagates API errors.
func TestConfigAlarmPolicy_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAlarmPolicy().client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "AddAlarmPolicyWithContext", func(ctx context.Context, request *config_sdk.AddAlarmPolicyRequest) (*config_sdk.AddAlarmPolicyResponse, error) {
		return nil, context.DeadlineExceeded
	})

	meta := newMockMetaForConfigAlarmPolicy()
	res := svcconfig.ResourceTencentCloudConfigAlarmPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                   "tf-example",
		"type":                   1,
		"event_scope":            []interface{}{1},
		"risk_level":             []interface{}{1},
		"notice_time":            "09:30:00~23:30:00",
		"notification_mechanism": "实时发送",
		"status":                 1,
		"notice_period":          []interface{}{1},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}
