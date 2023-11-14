package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTemScaleRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemScaleRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tem_scale_rule.scale_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_tem_scale_rule.scale_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTemScaleRule = `

resource "tencentcloud_tem_scale_rule" "scale_rule" {
  environment_id = "en-xxx"
  application_id = "app-xxx"
  autoscaler {
		autoscaler_name = "test"
		description = "test"
		enabled = true
		min_replicas = 1
		max_replicas = 2
		cron_horizontal_autoscaler {
			name = "test"
			period = "test"
			priority = 1
			enabled = true
			schedules {
				start_at = "03:00"
				target_replicas = 1
			}
		}
		horizontal_autoscaler {
			metrics = "test"
			enabled = true
			max_replicas = 2
			min_replicas = 1
			threshold = 60
		}

  }
}

`
