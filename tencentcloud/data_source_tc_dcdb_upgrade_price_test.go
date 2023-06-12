package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbUpgradePriceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbUpgradePriceDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_upgrade_price.upgrade_price")),
			},
		},
	})
}

const testAccDcdbUpgradePriceDataSource = `

data "tencentcloud_dcdb_upgrade_price" "upgrade_price" {
  instance_id = ""
  upgrade_type = ""
  add_shard_config {
		shard_count = 
		shard_memory = 
		shard_storage = 

  }
  expand_shard_config {
		shard_instance_ids = 
		shard_memory = 
		shard_storage = 
		shard_node_count = 

  }
  split_shard_config {
		shard_instance_ids = 
		split_rate = 
		shard_memory = 
		shard_storage = 

  }
  amount_unit = ""
      }

`
