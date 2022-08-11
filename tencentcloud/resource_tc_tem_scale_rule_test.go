package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTemScaleRule_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemScaleRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tem_scale_rule.scaleRule", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tem_scale_rule.scaleRule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTemScaleRule = `

resource "tencentcloud_tem_scale_rule" "scaleRule" {
  environment_id = "en-853mggjm"
  application_id = "app-3j29aa2p"
  autoscaler {
    autoscaler_name = "test3123"
    description     = "test"
    enabled         = true
    min_replicas    = 1
    max_replicas    = 3
    cron_horizontal_autoscaler {
      name     = "test"
      period   = "* * *"
      priority = 1
      enabled  = true
      schedules {
        start_at        = "03:00"
        target_replicas = 1
      }
    }
    cron_horizontal_autoscaler {
      name     = "test123123"
      period   = "* * *"
      priority = 0
      enabled  = true
      schedules {
        start_at        = "04:13"
        target_replicas = 1
      }
    }
    horizontal_autoscaler {
      metrics      = "CPU"
      enabled      = true
      max_replicas = 3
      min_replicas = 1
      threshold    = 60
    }

  }
}

`
