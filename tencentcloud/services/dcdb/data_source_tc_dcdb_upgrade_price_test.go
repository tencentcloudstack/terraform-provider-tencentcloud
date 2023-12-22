package dcdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbUpgradePriceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config:    testAccDcdbUpgradePriceDataSourceAdd,
				PreConfig: func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_upgrade_price.add_upgrade_price"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.add_upgrade_price", "instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.add_upgrade_price", "upgrade_type", "ADD"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.add_upgrade_price", "add_shard_config.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.add_upgrade_price", "add_shard_config.0.shard_count", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.add_upgrade_price", "add_shard_config.0.shard_memory", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.add_upgrade_price", "add_shard_config.0.shard_storage", "100"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.add_upgrade_price", "amount_unit", "pent"),

					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.add_upgrade_price", "original_price"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.add_upgrade_price", "price"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.add_upgrade_price", "formula"),
				),
			},
			{
				Config:    testAccDcdbUpgradePriceDataSourceExpand,
				PreConfig: func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_upgrade_price.expand_upgrade_price"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.expand_upgrade_price", "instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.expand_upgrade_price", "upgrade_type", "EXPAND"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.expand_upgrade_price", "expand_shard_config.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.expand_upgrade_price", "expand_shard_config.0.shard_instance_ids.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.expand_upgrade_price", "expand_shard_config.0.shard_memory", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.expand_upgrade_price", "expand_shard_config.0.shard_storage", "40"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.expand_upgrade_price", "expand_shard_config.0.shard_node_count", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.expand_upgrade_price", "amount_unit", "pent"),

					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.expand_upgrade_price", "original_price"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.expand_upgrade_price", "price"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.expand_upgrade_price", "formula"),
				),
			},
			{
				Config:    testAccDcdbUpgradePriceDataSourceSplit,
				PreConfig: func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_upgrade_price.split_upgrade_price"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.split_upgrade_price", "instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.split_upgrade_price", "upgrade_type", "SPLIT"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.split_upgrade_price", "split_shard_config.#"),

					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.split_upgrade_price", "split_shard_config.0.shard_instance_ids.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.split_upgrade_price", "split_shard_config.0.split_rate", "50"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.split_upgrade_price", "split_shard_config.0.shard_memory", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.split_upgrade_price", "split_shard_config.0.shard_storage", "100"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_upgrade_price.split_upgrade_price", "amount_unit", "pent"),

					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.split_upgrade_price", "original_price"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.split_upgrade_price", "price"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_upgrade_price.split_upgrade_price", "formula"),
				),
			},
		},
	})
}

const testAccDcdbUpgradePriceDataSourceAdd = `
locals {
	dcdb_id = "tdsqlshard-2imgzk5l"
	shardA_id = "shard-bgng7aqt"
	shardB_id = "shard-d4wx22xv"
}

data "tencentcloud_dcdb_upgrade_price" "add_upgrade_price" {
  instance_id = local.dcdb_id
  upgrade_type = "ADD"

  add_shard_config {
		shard_count = 2
		shard_memory = 2
		shard_storage = 100
  }
  amount_unit = "pent"
}

`

const testAccDcdbUpgradePriceDataSourceExpand = `
locals {
	dcdb_id = "tdsqlshard-2imgzk5l"
	shardA_id = "shard-bgng7aqt"
	shardB_id = "shard-d4wx22xv"
}

data "tencentcloud_dcdb_upgrade_price" "expand_upgrade_price" {
  instance_id = local.dcdb_id
  upgrade_type = "EXPAND"

  expand_shard_config {
		shard_instance_ids = [local.shardA_id]
		shard_memory = 2
		shard_storage = 40
		shard_node_count = 2
  }
  amount_unit = "pent"
}

`

const testAccDcdbUpgradePriceDataSourceSplit = `
locals {
	dcdb_id = "tdsqlshard-2imgzk5l"
	shardA_id = "shard-bgng7aqt"
	shardB_id = "shard-d4wx22xv"
}

data "tencentcloud_dcdb_upgrade_price" "split_upgrade_price" {
  instance_id = local.dcdb_id
  upgrade_type = "SPLIT"

  split_shard_config {
	    shard_instance_ids = [local.shardA_id]
		split_rate = 50
		shard_memory = 2
		shard_storage = 100
  }
  amount_unit = "pent"
}

`
