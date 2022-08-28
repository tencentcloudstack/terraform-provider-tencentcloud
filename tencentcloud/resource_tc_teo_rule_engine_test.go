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
					resource.TestCheckResourceAttrSet("tencentcloud_teo_rule_engine.ruleEngine", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_rule_engine.ruleEngine",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoRuleEngine = `

resource "tencentcloud_teo_rule_engine" "ruleEngine" {
  zone_id   = ""
  rule_name = ""
  status    = ""
  rules {
    conditions {
      conditions {
        operator = ""
        target   = ""
        values   = ""
      }
    }
    actions {
      normal_action {
        action = ""
        parameters {
          name   = ""
          values = ""
        }
      }
      rewrite_action {
        action = ""
        parameters {
          action = ""
          name   = ""
          values = ""
        }
      }
      code_action {
        action = ""
        parameters {
          name        = ""
          values      = ""
          status_code = ""
        }
      }
    }

  }
  tags      = {
    "createdBy" = "terraform"
  }
}

`
