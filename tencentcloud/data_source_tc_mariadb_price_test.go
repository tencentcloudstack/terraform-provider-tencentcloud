package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbPriceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbPriceDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_price.price")),
			},
		},
	})
}

const testAccMariadbPriceDataSource = `

data "tencentcloud_mariadb_price" "price" {
  zone = ""
  node_count = 
  memory = 
  storage = 
  period = 
  paymode = ""
  amount_unit = ""
    }

`
