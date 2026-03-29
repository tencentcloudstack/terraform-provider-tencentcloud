package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoL7AccRuleDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoL7AccRuleDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_l7_acc_rule.example"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.example", "zone_id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.example", "total_count"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.example", "rules"),
			),
		}},
	})
}

func TestAccTencentCloudTeoL7AccRuleDataSource_withRuleId(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoL7AccRuleDataSourceWithRuleId,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_l7_acc_rule.example"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.example", "zone_id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.example", "rule_id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.example", "total_count"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.example", "rules"),
			),
		}},
	})
}

const testAccTeoL7AccRuleDataSource = `
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
`

const testAccTeoL7AccRuleDataSourceWithRuleId = `
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"
  rule_id = "rule-xxxxx"
}
`
