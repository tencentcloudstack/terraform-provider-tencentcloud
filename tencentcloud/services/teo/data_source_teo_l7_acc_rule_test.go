package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentcloudTeoL7AccRuleDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL7AccRuleDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "total_count"),
				),
			},
		},
	})
}

const testAccTeoL7AccRuleDataSourceBasic = `
resource "tencentcloud_teo_zone" "zone" {
  zone_name = "test-zone-acc-l7-rule"
  type = "full"
}

data "tencentcloud_teo_l7_acc_rule" "teo_l7_acc_rule" {
  zone_id = tencentcloud_teo_zone.zone.id
}
`
