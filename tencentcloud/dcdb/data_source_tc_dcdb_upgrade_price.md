Use this data source to query detailed information of dcdb upgrade_price

Example Usage

```hcl
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

data "tencentcloud_dcdb_upgrade_price" "expand_upgrade_price" {
  instance_id = local.dcdb_id
  upgrade_type = "EXPAND"

  expand_shard_config {
		shard_instance_ids = ["shard-1b5r04az"]
		shard_memory = 2
		shard_storage = 40
		shard_node_count = 2
  }
  amount_unit = "pent"
}

data "tencentcloud_dcdb_upgrade_price" "split_upgrade_price" {
  instance_id = local.dcdb_id
  upgrade_type = "SPLIT"

  split_shard_config {
		shard_instance_ids = ["shard-1b5r04az"]
		split_rate = 50
		shard_memory = 2
		shard_storage = 100
  }
  amount_unit = "pent"
}
```