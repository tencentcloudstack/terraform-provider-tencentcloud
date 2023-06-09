---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_price"
sidebar_current: "docs-tencentcloud-datasource-mariadb_price"
description: |-
  Use this data source to query detailed information of mariadb price
---

# tencentcloud_mariadb_price

Use this data source to query detailed information of mariadb price

## Example Usage

```hcl
data "tencentcloud_mariadb_price" "price" {
  zone       = "ap-guangzhou-3"
  node_count = 2
  memory     = 2
  storage    = 20
  buy_count  = 1
  period     = 1
  paymode    = "prepaid"
}
```

## Argument Reference

The following arguments are supported:

* `buy_count` - (Required, Int) The quantity you want to purchase is queried by default for the price of purchasing 1 instance.
* `memory` - (Required, Int) Memory size in GB, which can be obtained by querying the instance specification through the `DescribeDBInstanceSpecs` API.
* `node_count` - (Required, Int) Number of instance nodes, which can be obtained by querying the instance specification through the `DescribeDBInstanceSpecs` API.
* `storage` - (Required, Int) Storage capacity in GB. The maximum and minimum storage space can be obtained by querying instance specification through the `DescribeDBInstanceSpecs` API.
* `zone` - (Required, String) AZ ID of the purchased instance.
* `amount_unit` - (Optional, String) Price unit. Valid values: `* pent` (cent), `* microPent` (microcent).
* `paymode` - (Optional, String) Billing type. Valid values: `postpaid` (pay-as-you-go), `prepaid` (monthly subscription).
* `period` - (Optional, Int) Purchase period in months.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `original_price` - Original price * Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. * Currency: CNY (Chinese site), USD (international site).
* `price` - The actual price may be different from the original price due to discounts. * Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. * Currency: CNY (Chinese site), USD (international site).


