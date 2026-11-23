---
subcategory: "TDSQL Boundless(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_specs"
description: |-
  Use this data source to query TDSQL Boundless (TEO Tdmysql) salable component specs.
---

# tencentcloud_tdmysql_specs

Use this data source to query TDSQL Boundless (TEO Tdmysql) salable component specs.

## Example Usage

```hcl
data "tencentcloud_tdmysql_specs" "example" {
  full_replications      = 2
  is_exclusive_instance  = 1
  zones                  = ["ap-guangzhou-1"]
  instance_id            = "tdsql-o2mtest"
}

output "first_hybrid_spec_code" {
  value = data.tencentcloud_tdmysql_specs.example.hybrid_node_specs[0].spec_code
}
```
