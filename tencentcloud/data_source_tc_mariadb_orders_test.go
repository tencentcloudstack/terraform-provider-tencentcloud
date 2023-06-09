package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbOrdersDataSource_basic -v
func TestAccTencentCloudMariadbOrdersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbOrdersDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_orders.orders"),
				),
			},
		},
	})
}

const testAccMariadbOrdersDataSource = `
data "tencentcloud_mariadb_orders" "orders" {
  deal_name = "20230607164033835942781"
}
`
