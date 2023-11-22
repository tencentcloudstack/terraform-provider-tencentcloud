package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const defaultTestingPrometheusId = "prom-3oq89zvg"

func TestAccTencentCloudTestingMonitorRecordingRule_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTestingRecordingRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordingRuleExists("tencentcloud_monitor_tmp_recording_rule.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_recording_rule.basic", "name", "recording_rule-test"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_recording_rule.basic", "rule_state", "2"),
				),
			},
			{
				Config: testTestingRecordingRule_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordingRuleExists("tencentcloud_monitor_tmp_recording_rule.update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_recording_rule.update", "name", "recording_rule-update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_recording_rule.update", "rule_state", "3"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_recording_rule.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testTestingRecordingRuleVar = `
variable "prometheus_id" {
  default = "` + defaultTestingPrometheusId + `"
}
`
const testTestingRecordingRule_basic = testTestingRecordingRuleVar + `
resource "tencentcloud_monitor_tmp_recording_rule" "basic" {
  name			= "recording_rule-test"
  instance_id	= var.prometheus_id
  rule_state	= 2
  group			= <<EOF
---
name: example-test
rules:
  - record: job:http_inprogress_requests:sum
    expr: sum by (job) (http_inprogress_requests)
EOF
}`

const testTestingRecordingRule_update = testTestingRecordingRuleVar + `
resource "tencentcloud_monitor_tmp_recording_rule" "update" {
  name			= "recording_rule-update"
  instance_id	= var.prometheus_id
  rule_state	= 3
  group			= <<EOF
---
name: example-test-update
rules:
  - record: job:http_inprogress_requests:sum
    expr: sum by (job) (http_inprogress_requests)
EOF
}`
