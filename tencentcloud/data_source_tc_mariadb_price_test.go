package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbPriceDataSource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_price.price"),
				),
			},
		},
	})
}

const testAccMariadbPriceDataSource = `
data "tencentcloud_mariadb_price" "price" {
  zone       = "ap-guangzhou-3"
  node_count = 2
  memory     = 2
  storage    = 20
  buy_count  = 1
  period     = 1
  paymode    = "prepaid"
}
`
