package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTeoRuleEngine_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoRuleEngine,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_rule_engine.rule_engine", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_rule_engine.rule_engine",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoRuleEngine = `

resource "tencentcloud_teo_rule_engine" "rule_engine" {
  zone_id   = tencentcloud_teo_zone.zone.id
  rule_name = "rule0"
  status    = "enable"

  rules {
    conditions {
      conditions {
        operator = "equal"
        target   = "host"
        values   = [
          "www.sfurnace.work",
        ]
      }
    }

    actions {
      normal_action {
        action = "MaxAge"

        parameters {
          name   = "FollowOrigin"
          values = [
            "on",
          ]
        }
        parameters {
          name   = "MaxAgeTime"
          values = [
            "0",
          ]
        }
      }
    }
  }
}


`
