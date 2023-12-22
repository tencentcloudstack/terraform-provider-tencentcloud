package dcdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbOrdersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbOrdersDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_orders.orders"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_orders.orders", "deal_names.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_orders.orders", "deals.#"),
				),
			},
		},
	})
}

const testAccDcdbOrdersDataSource = `

data "tencentcloud_dcdb_orders" "orders" {
  deal_names = ["20230612249034137670121"]
}

`
