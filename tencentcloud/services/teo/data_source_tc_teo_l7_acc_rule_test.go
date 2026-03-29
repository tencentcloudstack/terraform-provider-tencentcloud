package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudTeoL7AccRuleDataSource_basic -v
func TestAccTencentCloudTeoL7AccRuleDataSource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL7AccRuleDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_l7_acc_rule.test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.test", "total_count"),
				),
			},
		},
	})
}

const testAccTeoL7AccRuleDataSource = `
data "tencentcloud_teo_l7_acc_rule" "test" {
  zone_id = "zone-39quuimqg8r6"
}
`
