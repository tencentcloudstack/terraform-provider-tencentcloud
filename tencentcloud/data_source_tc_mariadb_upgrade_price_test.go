package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixMariadbUpgradePriceDataSource_basic -v
func TestAccTencentCloudNeedFixMariadbUpgradePriceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbUpgradePriceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_upgrade_price.upgrade_price"),
				),
			},
		},
	})
}

const testAccMariadbUpgradePriceDataSource = `
data "tencentcloud_mariadb_upgrade_price" "upgrade_price" {
  instance_id = "tdsql-9vqvls95"
  memory      = 4
  storage     = 40
  node_count  = 2
}
`
