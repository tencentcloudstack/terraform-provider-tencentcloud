---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_upgrade_price"
sidebar_current: "docs-tencentcloud-datasource-dcdb_upgrade_price"
description: |-
  Use this data source to query detailed information of dcdb upgrade_price
---

# tencentcloud_dcdb_upgrade_price

Use this data source to query detailed information of dcdb upgrade_price

## Example Usage

```hcl
data "tencentcloud_dcdb_upgrade_price" "add_upgrade_price" {
  instance_id  = local.dcdb_id
  upgrade_type = "ADD"
  add_shard_config {
    shard_count   = 2
    shard_memory  = 2
    shard_storage = 100
  }
  amount_unit = "pent"
}

data "tencentcloud_dcdb_upgrade_price" "expand_upgrade_price" {
  instance_id  = local.dcdb_id
  upgrade_type = "EXPAND"

  expand_shard_config {
    shard_instance_ids = ["shard-1b5r04az"]
    shard_memory       = 2
    shard_storage      = 40
    shard_node_count   = 2
  }
  amount_unit = "pent"
}

data "tencentcloud_dcdb_upgrade_price" "split_upgrade_price" {
  instance_id  = local.dcdb_id
  upgrade_type = "SPLIT"

  split_shard_config {
    shard_instance_ids = ["shard-1b5r04az"]
    split_rate         = 50
    shard_memory       = 2
    shard_storage      = 100
  }
  amount_unit = "pent"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `upgrade_type` - (Required, String) Upgrade type, ADD: add new shard, EXPAND: upgrade the existing shard, SPLIT: split existing shard.
* `add_shard_config` - (Optional, List) Config for adding new shard.
* `amount_unit` - (Optional, String) Price unit. Valid values: `pent` (cent), `microPent` (microcent).
* `expand_shard_config` - (Optional, List) Config for expanding existing shard.
* `result_output_file` - (Optional, String) Used to save results.
* `split_shard_config` - (Optional, List) Config for splitting existing shard.

The `add_shard_config` object supports the following:

* `shard_count` - (Required, Int) The number of new shards.
* `shard_memory` - (Required, Int) Shard memory size in GB.
* `shard_storage` - (Required, Int) Shard storage capacity in GB.

The `expand_shard_config` object supports the following:

* `shard_instance_ids` - (Required, Set) List of shard ID.
* `shard_memory` - (Required, Int) Shard memory size in GB.
* `shard_storage` - (Required, Int) Shard storage capacity in GB.
* `shard_node_count` - (Optional, Int) Shard node count.

The `split_shard_config` object supports the following:

* `shard_instance_ids` - (Required, Set) List of shard ID.
* `shard_memory` - (Required, Int) Shard memory size in GB.
* `shard_storage` - (Required, Int) Shard storage capacity in GB.
* `split_rate` - (Required, Int) Data split ratio, fixed at 50%.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `formula` - Price calculation formula.
* `original_price` - Original price. Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. Currency: CNY (Chinese site), USD (international site).
* `price` - The actual price may be different from the original price due to discounts. Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. Currency: CNY (Chinese site), USD (international site).


