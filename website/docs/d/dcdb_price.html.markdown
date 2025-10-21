---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_price"
sidebar_current: "docs-tencentcloud-datasource-dcdb_price"
description: |-
  Use this data source to query detailed information of dcdb price
---

# tencentcloud_dcdb_price

Use this data source to query detailed information of dcdb price

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `instance_count` - (Required, Int) The count of instances wants to buy.
* `period` - (Required, Int) Purchase period in months.
* `shard_count` - (Required, Int) Number of instance shards.
* `shard_memory` - (Required, Int) Shard memory size in GB.
* `shard_node_count` - (Required, Int) Number of instance shard nodes.
* `shard_storage` - (Required, Int) Shard storage capacity in GB.
* `zone` - (Required, String) AZ ID of the purchased instance.
* `amount_unit` - (Optional, String) Price unit. Valid values: `pent` (cent), `microPent` (microcent).
* `paymode` - (Optional, String) Billing type. Valid values: `postpaid` (pay-as-you-go), `prepaid` (monthly subscription).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `original_price` - Original price. Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. Currency: CNY (Chinese site), USD (international site).
* `price` - The actual price may be different from the original price due to discounts. Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. Currency: CNY (Chinese site), USD (international site).


