package wedata_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataBaselineResource_basic -v
func TestAccTencentCloudNeedFixWedataBaselineResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataBaseline,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "baseline_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "baseline_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "create_uin"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "create_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "in_charge_uin"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "in_charge_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "promise_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "warning_margin"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "is_new_alarm"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_baseline.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataBaselineUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "baseline_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "baseline_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "create_uin"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "create_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "in_charge_uin"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "in_charge_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "promise_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "warning_margin"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline.example", "is_new_alarm"),
				),
			},
		},
	})
}

const testAccWedataBaseline = `
resource "tencentcloud_wedata_baseline" "example" {
  project_id     = "1927766435649077248"
  baseline_name  = "tf_example"
  baseline_type  = "D"
  create_uin     = "tf_user"
  create_name    = "100028439226"
  in_charge_uin  = "tf_user"
  in_charge_name = "100028439226"
  promise_tasks {
    project_id          = "1927766435649077248"
    task_name           = "tf_demo_task"
    task_id             = "20231030145334153"
    task_cycle          = "D"
    workflow_name       = "交易"
    workflow_id         = "e4dafb2e-76eb-11ee-bfeb-b8cef68a6637"
    task_in_charge_name = ";tf_user;"
  }
  promise_time   = "00:00:00"
  warning_margin = 30
  is_new_alarm   = true
  baseline_create_alarm_rule_request {
    alarm_types = [
      "baseLineBroken",
      "baseLineWarning",
      "baseLineTaskFailure"
    ]
    alarm_level = 2
    alarm_ways  = [
      "email",
      "sms"
    ]
    alarm_recipient_type = 1
    alarm_recipients     = [
      "tf_user"
    ]
    alarm_recipient_ids = [
      "100028439226"
    ]
  }
}
`

const testAccWedataBaselineUpdate = `
resource "tencentcloud_wedata_baseline" "example" {
  project_id     = "1927766435649077248"
  baseline_name  = "tf_example_update"
  baseline_type  = "D"
  create_uin     = "tf_user"
  create_name    = "100028439226"
  in_charge_uin  = "tf_user"
  in_charge_name = "100028439226"
  promise_tasks {
    project_id          = "1927766435649077248"
    task_name           = "tf_demo_task"
    task_id             = "20231030145334153"
    task_cycle          = "D"
    workflow_name       = "交易"
    workflow_id         = "e4dafb2e-76eb-11ee-bfeb-b8cef68a6637"
    task_in_charge_name = ";tf_user;"
  }
  promise_time   = "00:00:00"
  warning_margin = 30
  is_new_alarm   = true
  baseline_create_alarm_rule_request {
    alarm_types = [
      "baseLineBroken",
      "baseLineWarning",
      "baseLineTaskFailure"
    ]
    alarm_level = 2
    alarm_ways  = [
      "email",
      "sms"
    ]
    alarm_recipient_type = 1
    alarm_recipients     = [
      "tf_user"
    ]
    alarm_recipient_ids = [
      "100028439226"
    ]
  }
}
`
