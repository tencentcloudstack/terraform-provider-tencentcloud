package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_price.price")),
			},
		},
	})
}

const testAccDcdbPriceDataSource = `

data "tencentcloud_dcdb_price" "price" {
  zone = ""
  period = 
  shard_node_count = 
  shard_memory = 
  shard_storage = 
  shard_count = 
  paymode = ""
  amount_unit = ""
    }

`
