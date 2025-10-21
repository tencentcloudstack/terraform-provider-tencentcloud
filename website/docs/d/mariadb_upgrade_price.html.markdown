---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_upgrade_price"
sidebar_current: "docs-tencentcloud-datasource-mariadb_upgrade_price"
description: |-
  Use this data source to query detailed information of mariadb upgrade_price
---

# tencentcloud_mariadb_upgrade_price

Use this data source to query detailed information of mariadb upgrade_price

## Example Usage

```hcl
data "tencentcloud_mariadb_upgrade_price" "upgrade_price" {
  instance_id = "tdsql-9vqvls95"
  memory      = 4
  storage     = 40
  node_count  = 2
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `memory` - (Required, Int) Memory size in GB, which can be obtained by querying the instance specification through the `DescribeDBInstanceSpecs` API.
* `storage` - (Required, Int) Storage capacity in GB. The maximum and minimum storage space can be obtained by querying instance specification through the `DescribeDBInstanceSpecs` API.
* `amount_unit` - (Optional, String) Price unit. Valid values: `* pent` (cent), `* microPent` (microcent).
* `node_count` - (Optional, Int) New instance nodes, zero means not change.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `formula` - Price calculation formula.
* `original_price` - Original price * Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. * Currency: CNY (Chinese site), USD (international site).
* `price` - The actual price may be different from the original price due to discounts. * Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. * Currency: CNY (Chinese site), USD (international site).


