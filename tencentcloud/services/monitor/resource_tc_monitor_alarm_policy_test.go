package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TestAccTencentCloudMonitorAlarmPolicyResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmPolicy,
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckResourceAttrSet("tencentcloud_monitor_alarm_policy.policy", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_alarm_policy.policy", "policy_name", "terraform"),
				),
			},
		},
	})
}

const testAccMonitorAlarmPolicy string = `
resource "tencentcloud_monitor_alarm_policy" "policy" {
  enable       = 1
  monitor_type = "MT_QCE"
  namespace    = "cvm_device"
  notice_ids   = [
    "notice-f2svbu3w",
  ]
  policy_name  = "terraform"
  project_id   = 0

  conditions {
    is_union_rule = 0

    rules {
      continue_period  = 5
      description      = "CPUUtilization"
      is_power_notice  = 0
      metric_name      = "CpuUsage"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
    rules {
      continue_period  = 5
      description      = "PublicBandwidthUtilization"
      is_power_notice  = 0
      metric_name      = "Outratio"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
    rules {
      continue_period  = 5
      description      = "MemoryUtilization"
      is_power_notice  = 0
      metric_name      = "MemUsage"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
    rules {
      continue_period  = 5
      description      = "DiskUtilization"
      is_power_notice  = 0
      metric_name      = "CvmDiskUsage"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
  }

  event_conditions {
    continue_period  = 0
    description      = "DiskReadonly"
    is_power_notice  = 0
    metric_name      = "disk_readonly"
    notice_frequency = 0
    period           = 0
  }

  policy_tag {
    key   = "test-tag"
    value = "unit-test"
  }
}
`

// mockMetaForMonitorAlarmPolicy implements tccommon.ProviderMeta
type mockMetaForMonitorAlarmPolicy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForMonitorAlarmPolicy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForMonitorAlarmPolicy{}

func newMockMetaForMonitorAlarmPolicy() *mockMetaForMonitorAlarmPolicy {
	return &mockMetaForMonitorAlarmPolicy{client: &connectivity.TencentCloudClient{}}
}

// buildMockAlarmPolicy builds an AlarmPolicy struct with the minimum fields
// required by resourceTencentMonitorAlarmPolicyRead to avoid nil-pointer panics
// when iterating Condition.Rules and EventCondition.Rules.
func buildMockAlarmPolicy(isBindAll *int64) *monitor.AlarmPolicy {
	return &monitor.AlarmPolicy{
		PolicyId:    helper.String("policy-fake-id"),
		PolicyName:  helper.String("terraform"),
		MonitorType: helper.String("MT_QCE"),
		Namespace:   helper.String("cvm_device"),
		Remark:      helper.String(""),
		Enable:      helper.IntInt64(1),
		ProjectId:   helper.IntInt64(0),
		IsBindAll:   isBindAll,
		Condition: &monitor.AlarmPolicyCondition{
			IsUnionRule: helper.IntInt64(0),
			Rules:       []*monitor.AlarmPolicyRule{},
		},
		EventCondition: &monitor.AlarmPolicyEventCondition{
			Rules: []*monitor.AlarmPolicyRule{},
		},
		NoticeIds:    []*string{},
		TriggerTasks: []*monitor.AlarmPolicyTriggerTask{},
		TagInstances: []*monitor.TagInstance{},
	}
}

// go test ./tencentcloud/services/monitor/ -run "TestMonitorAlarmPolicy_Read_IsBindAll" -v -count=1 -gcflags="all=-l"

// TestMonitorAlarmPolicy_Read_IsBindAll tests Read populates is_bind_all from Policy.IsBindAll
func TestMonitorAlarmPolicy_Read_IsBindAll(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	monitorClient := &monitor.Client{}
	patches.ApplyMethodReturn(newMockMetaForMonitorAlarmPolicy().client, "UseMonitorClient", monitorClient)

	patches.ApplyMethodFunc(monitorClient, "DescribeAlarmPolicy", func(request *monitor.DescribeAlarmPolicyRequest) (*monitor.DescribeAlarmPolicyResponse, error) {
		resp := monitor.NewDescribeAlarmPolicyResponse()
		resp.Response = &monitor.DescribeAlarmPolicyResponseParams{
			Policy:    buildMockAlarmPolicy(helper.IntInt64(1)),
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMonitorAlarmPolicy()
	res := svcmonitor.ResourceTencentCloudMonitorAlarmPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"policy_name":  "terraform",
		"monitor_type": "MT_QCE",
		"namespace":    "cvm_device",
	})
	d.SetId("policy-fake-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 1, d.Get("is_bind_all"))
}

// TestMonitorAlarmPolicy_Read_IsBindAllNil tests Read does not set is_bind_all when IsBindAll is nil
func TestMonitorAlarmPolicy_Read_IsBindAllNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	monitorClient := &monitor.Client{}
	patches.ApplyMethodReturn(newMockMetaForMonitorAlarmPolicy().client, "UseMonitorClient", monitorClient)

	patches.ApplyMethodFunc(monitorClient, "DescribeAlarmPolicy", func(request *monitor.DescribeAlarmPolicyRequest) (*monitor.DescribeAlarmPolicyResponse, error) {
		resp := monitor.NewDescribeAlarmPolicyResponse()
		resp.Response = &monitor.DescribeAlarmPolicyResponseParams{
			Policy:    buildMockAlarmPolicy(nil),
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMonitorAlarmPolicy()
	res := svcmonitor.ResourceTencentCloudMonitorAlarmPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"policy_name":  "terraform",
		"monitor_type": "MT_QCE",
		"namespace":    "cvm_device",
	})
	d.SetId("policy-fake-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// is_bind_all should remain zero value because it was not set
	assert.Equal(t, 0, d.Get("is_bind_all"))
}

// TestMonitorAlarmPolicy_Read_NilPolicy tests Read handles nil policy
func TestMonitorAlarmPolicy_Read_NilPolicy(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	monitorClient := &monitor.Client{}
	patches.ApplyMethodReturn(newMockMetaForMonitorAlarmPolicy().client, "UseMonitorClient", monitorClient)

	patches.ApplyMethodFunc(monitorClient, "DescribeAlarmPolicy", func(request *monitor.DescribeAlarmPolicyRequest) (*monitor.DescribeAlarmPolicyResponse, error) {
		resp := monitor.NewDescribeAlarmPolicyResponse()
		resp.Response = &monitor.DescribeAlarmPolicyResponseParams{
			Policy:    nil,
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMonitorAlarmPolicy()
	res := svcmonitor.ResourceTencentCloudMonitorAlarmPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"policy_name":  "terraform",
		"monitor_type": "MT_QCE",
		"namespace":    "cvm_device",
	})
	d.SetId("policy-fake-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestMonitorAlarmPolicy_Create_WithIsBindAll tests Create passes IsBindAll to the CreateAlarmPolicy request
func TestMonitorAlarmPolicy_Create_WithIsBindAll(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	monitorClient := &monitor.Client{}
	patches.ApplyMethodReturn(newMockMetaForMonitorAlarmPolicy().client, "UseMonitorClient", monitorClient)

	var capturedRequest *monitor.CreateAlarmPolicyRequest
	patches.ApplyMethodFunc(monitorClient, "CreateAlarmPolicy", func(request *monitor.CreateAlarmPolicyRequest) (*monitor.CreateAlarmPolicyResponse, error) {
		capturedRequest = request
		resp := monitor.NewCreateAlarmPolicyResponse()
		resp.Response = &monitor.CreateAlarmPolicyResponseParams{
			PolicyId:  helper.String("policy-fake-id"),
			OriginId:  helper.String("origin-fake-id"),
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(monitorClient, "DescribeAlarmPolicy", func(request *monitor.DescribeAlarmPolicyRequest) (*monitor.DescribeAlarmPolicyResponse, error) {
		resp := monitor.NewDescribeAlarmPolicyResponse()
		resp.Response = &monitor.DescribeAlarmPolicyResponseParams{
			Policy:    buildMockAlarmPolicy(helper.IntInt64(1)),
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMonitorAlarmPolicy()
	res := svcmonitor.ResourceTencentCloudMonitorAlarmPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"policy_name":  "terraform",
		"monitor_type": "MT_QCE",
		"namespace":    "cvm_device",
		"is_bind_all":  1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "policy-fake-id", d.Id())
	// verify IsBindAll was passed to the CreateAlarmPolicy request
	assert.NotNil(t, capturedRequest.IsBindAll)
	assert.Equal(t, int64(1), *capturedRequest.IsBindAll)
}

// TestMonitorAlarmPolicy_Create_WithoutIsBindAll tests Create does not pass IsBindAll when not set
func TestMonitorAlarmPolicy_Create_WithoutIsBindAll(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	monitorClient := &monitor.Client{}
	patches.ApplyMethodReturn(newMockMetaForMonitorAlarmPolicy().client, "UseMonitorClient", monitorClient)

	var capturedRequest *monitor.CreateAlarmPolicyRequest
	patches.ApplyMethodFunc(monitorClient, "CreateAlarmPolicy", func(request *monitor.CreateAlarmPolicyRequest) (*monitor.CreateAlarmPolicyResponse, error) {
		capturedRequest = request
		resp := monitor.NewCreateAlarmPolicyResponse()
		resp.Response = &monitor.CreateAlarmPolicyResponseParams{
			PolicyId:  helper.String("policy-fake-id"),
			OriginId:  helper.String("origin-fake-id"),
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(monitorClient, "DescribeAlarmPolicy", func(request *monitor.DescribeAlarmPolicyRequest) (*monitor.DescribeAlarmPolicyResponse, error) {
		resp := monitor.NewDescribeAlarmPolicyResponse()
		resp.Response = &monitor.DescribeAlarmPolicyResponseParams{
			Policy:    buildMockAlarmPolicy(nil),
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMonitorAlarmPolicy()
	res := svcmonitor.ResourceTencentCloudMonitorAlarmPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"policy_name":  "terraform",
		"monitor_type": "MT_QCE",
		"namespace":    "cvm_device",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "policy-fake-id", d.Id())
	// verify IsBindAll was NOT passed to the CreateAlarmPolicy request
	assert.Nil(t, capturedRequest.IsBindAll)
}

// TestMonitorAlarmPolicy_Schema_IsBindAll tests the is_bind_all schema definition
func TestMonitorAlarmPolicy_Schema_IsBindAll(t *testing.T) {
	res := svcmonitor.ResourceTencentCloudMonitorAlarmPolicy()

	assert.Contains(t, res.Schema, "is_bind_all")
	field := res.Schema["is_bind_all"]
	assert.NotNil(t, field)
	assert.True(t, field.Optional)
	assert.True(t, field.ForceNew)
	assert.Equal(t, schema.TypeInt, field.Type)
	assert.NotNil(t, field.ValidateFunc)
}
