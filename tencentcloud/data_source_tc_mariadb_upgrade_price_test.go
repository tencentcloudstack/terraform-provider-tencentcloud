package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbUpgradePriceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbUpgradePriceDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_upgrade_price.upgrade_price")),
			},
		},
	})
}

const testAccMariadbUpgradePriceDataSource = `

data "tencentcloud_mariadb_upgrade_price" "upgrade_price" {
  instance_id = ""
  memory = 
  storage = 
  node_count = 
  amount_unit = ""
      }

`
