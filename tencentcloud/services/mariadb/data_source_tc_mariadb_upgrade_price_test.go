package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixMariadbUpgradePriceDataSource_basic -v
func TestAccTencentCloudNeedFixMariadbUpgradePriceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbUpgradePriceDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_upgrade_price.upgrade_price"),
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
