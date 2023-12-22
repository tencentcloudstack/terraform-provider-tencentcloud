package dayuv2_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataDayuL4RulesV2Name = "data.tencentcloud_dayu_l4_rules_v2.test"

func TestAccTencentCloudDataDayuL4RulesV2(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_INTERNATIONAL) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDayuL4RuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataDayuL4RulesV2Baic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDayuL4RuleV2Exists(testDayuL4RuleV2ResourceKeyTCP),
					resource.TestCheckResourceAttr(testDataDayuL4RulesV2Name, "list.#", "1"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesV2Name, "list.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesV2Name, "list.0.source_list.#", "1"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesV2Name, "list.0.virtual_port", "2020"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesV2Name, "list.0.source_type", "2"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesV2Name, "list.0.source_port", "20"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesV2Name, "list.0.rule_name", "test"),
				),
			},
		},
	})
}

const testAccTencentCloudDataDayuL4RulesV2Baic = `
resource "tencentcloud_dayu_l4_rule_v2" "test" {
	business     = "bgpip"
	resource_id  = "bgpip-000004xg"
	vpn          = "162.62.163.50"
	virtual_port = 2020
  
	rules {
	  keep_enable   = false
	  keeptime      = 0
	  lb_type       = 1
	  protocol      = "TCP"
	  region        = 5
	  remove_switch = false
	  rule_name     = "test"
	  source_list {
		source = "1.2.3.9"
		weight = 0
	  }
	  source_port  = 20
	  source_type  = 2
	  virtual_port = 2020
	}
  }

data "tencentcloud_dayu_l4_rules_v2" "test" {
    business = tencentcloud_dayu_l4_rule_v2.test.business
	ip = tencentcloud_dayu_l4_rule_v2.test.vpn
	virtual_port = tencentcloud_dayu_l4_rule_v2.test.virtual_port
}
`
