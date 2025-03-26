package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -test.run TestAccTencentCloudTeoL7AccRulePriorityResource_basic -v
func TestAccTencentCloudTeoL7AccRulePriorityResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL7AccRulePriority,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule_priority.teo_l7_acc_rule_priority", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_priority.teo_l7_acc_rule_priority", "rule_ids.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_priority.teo_l7_acc_rule_priority", "rule_ids.0", "rule-39pkyiu08edu"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_priority.teo_l7_acc_rule_priority", "rule_ids.1", "rule-39pky6n21mkf"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_l7_acc_rule_priority.teo_l7_acc_rule_priority",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoL7AccRulePriorityUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule_priority.teo_l7_acc_rule_priority", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_priority.teo_l7_acc_rule_priority", "rule_ids.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_priority.teo_l7_acc_rule_priority", "rule_ids.0", "rule-39pky6n21mkf"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_priority.teo_l7_acc_rule_priority", "rule_ids.1", "rule-39pkyiu08edu"),
				),
			},
		},
	})
}

const testAccTeoL7AccRulePriority = `

resource "tencentcloud_teo_l7_acc_rule_priority" "teo_l7_acc_rule_priority" {
    zone_id       = "zone-36bjhygh1bxe"
    rule_ids = [
        "rule-39pkyiu08edu",
        "rule-39pky6n21mkf",
    ]
}
`

const testAccTeoL7AccRulePriorityUp = `

resource "tencentcloud_teo_l7_acc_rule_priority" "teo_l7_acc_rule_priority" {
    zone_id       = "zone-36bjhygh1bxe"
    rule_ids = [
        "rule-39pky6n21mkf",
        "rule-39pkyiu08edu",
    ]
}
`
