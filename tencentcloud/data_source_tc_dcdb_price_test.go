package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbPriceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbPriceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_price.price"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_price.price", "zone"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_price.price", "instance_count", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_price.price", "period", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_price.price", "shard_node_count", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_price.price", "shard_memory", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_price.price", "shard_storage", "10"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_price.price", "shard_count", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_price.price", "paymode", "postpaid"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_price.price", "amount_unit", "pent"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_price.price", "original_price"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_price.price", "price"),
				),
			},
		},
	})
}

const testAccDcdbPriceDataSource = defaultAzVariable + `

data "tencentcloud_dcdb_price" "price" {
	instance_count   = 1
	zone             = var.default_az
	period           = 1
	shard_node_count = 2
	shard_memory     = 2
	shard_storage    = 10
	shard_count      = 2
	paymode          = "postpaid"
	amount_unit      = "pent"
}

`
