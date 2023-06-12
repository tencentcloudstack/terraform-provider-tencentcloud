---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_renewal_price"
sidebar_current: "docs-tencentcloud-datasource-mariadb_renewal_price"
description: |-
  Use this data source to query detailed information of mariadb renewal_price
---

# tencentcloud_mariadb_renewal_price

Use this data source to query detailed information of mariadb renewal_price

## Example Usage

```hcl
data "tencentcloud_mariadb_renewal_price" "renewal_price" {
  instance_id = "tdsql-9vqvls95"
  period      = 2
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `amount_unit` - (Optional, String) Price unit. Valid values: `* pent` (cent), `* microPent` (microcent).
* `period` - (Optional, Int) Renewal duration, default: 1 month.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `original_price` - Original price * Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. * Currency: CNY (Chinese site), USD (international site).
* `price` - The actual price may be different from the original price due to discounts. * Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. * Currency: CNY (Chinese site), USD (international site).


