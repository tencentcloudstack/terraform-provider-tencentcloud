package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbShardSpecDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbShardSpecDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_shard_spec.shard_spec"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shard_spec.shard_spec", "spec_config.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shard_spec.shard_spec", "spec_config.0.machine"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shard_spec.shard_spec", "spec_config.0.spec_config_infos.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shard_spec.shard_spec", "spec_config.0.spec_config_infos.0.node_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shard_spec.shard_spec", "spec_config.0.spec_config_infos.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shard_spec.shard_spec", "spec_config.0.spec_config_infos.0.min_storage"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shard_spec.shard_spec", "spec_config.0.spec_config_infos.0.max_storage"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shard_spec.shard_spec", "spec_config.0.spec_config_infos.0.qps"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shard_spec.shard_spec", "spec_config.0.spec_config_infos.0.suit_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shard_spec.shard_spec", "spec_config.0.spec_config_infos.0.pid"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shard_spec.shard_spec", "spec_config.0.spec_config_infos.0.cpu"),
				),
			},
		},
	})
}

const testAccDcdbShardSpecDataSource = `

data "tencentcloud_dcdb_shard_spec" "shard_spec" {}

`
