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
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL7AccRuleDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "total_count"),
				),
			},
		},
	})
}

const testAccTeoL7AccRuleDataSource = `

resource "tencentcloud_teo_zone" "test" {
  zone_name = "test.example.com"
  type       = "full"
}

data "tencentcloud_teo_l7_acc_rule" "teo_l7_acc_rule" {
  zone_id = tencentcloud_teo_zone.test.id
}
`
