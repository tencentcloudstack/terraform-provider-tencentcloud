package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataDayuL7RulesV2Name = "data.tencentcloud_dayu_l7_rules_v2.test"

func TestAccTencentCloudDataDayuL7V2Rules(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_INTERNATIONAL) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuL7RuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataDayuL7RulesV2Baic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDayuL7RuleV2Exists("tencentcloud_dayu_l7_rule_v2.test_rule"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesV2Name, "list.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudDataDayuL7RulesV2Baic = `
resource "tencentcloud_dayu_l7_rule_v2" "test_rule" {
	resource_type="bgpip"
	resource_id="bgpip-000004xe"
	resource_ip="119.28.217.162"
	rule {
	  keep_enable=0
	  keeptime=0
	  source_list {
		source="1.2.3.5"
		weight=100
	  }
	  source_list {
		source="1.2.3.6"
		weight=100
	  }
	  lb_type=1
	  protocol="http"
	  source_type=2
	  domain="github.com"
	}
  }
data "tencentcloud_dayu_l7_rules_v2" "test" {
	business = "bgpip"
	domain   = tencentcloud_dayu_l7_rule_v2.test_rule.rule.0.domain
	offset   = 0
	limit    = 10
}
`
