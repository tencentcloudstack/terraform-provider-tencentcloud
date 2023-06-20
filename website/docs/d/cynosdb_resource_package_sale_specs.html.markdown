---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_resource_package_sale_specs"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_resource_package_sale_specs"
description: |-
  Use this data source to query detailed information of cynosdb resource_package_sale_specs
---

# tencentcloud_cynosdb_resource_package_sale_specs

Use this data source to query detailed information of cynosdb resource_package_sale_specs

## Example Usage

```hcl
data "tencentcloud_cynosdb_resource_package_sale_specs" "resource_package_sale_specs" {
  instance_type  = "cynosdb-serverless"
  package_region = "china"
  package_type   = "CCU"
}
```

## Argument Reference

The following arguments are supported:

* `instance_type` - (Required, String) Instance Type. Value range: cynosdb-serverless, cynosdb, cdb.
* `package_region` - (Required, String) Resource package usage region China - common in mainland China, overseas - common in Hong Kong, Macao, Taiwan, and overseas.
* `package_type` - (Required, String) Resource package type CCU - Computing resource package DISK - Storage resource package.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `detail` - Resource package details note: This field may return null, indicating that a valid value cannot be obtained.
  * `expire_day` - Resource package validity period, in days. Note: This field may return null, indicating that a valid value cannot be obtained.
  * `max_package_spec` - The maximum number of resources in the current version of the resource package, calculated in units of resources; Storage resource: GB Note: This field may return null, indicating that a valid value cannot be obtained.
  * `min_package_spec` - The minimum number of resources in the current version of the resource package, calculated in units of resources; Storage resource: GB Note: This field may return null, indicating that a valid value cannot be obtained.
  * `package_region` - Note: This field may return null, indicating that a valid value cannot be obtained.
  * `package_type` - Resource package type CCU - Compute resource package DISK - Store resource package Note: This field may return null, indicating that a valid value cannot be obtained.
  * `package_version` - Resource package version base basic version, common general version, enterprise enterprise version Note: This field may return null, indicating that a valid value cannot be obtained.


