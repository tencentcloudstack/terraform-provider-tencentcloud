---
subcategory: "TDSQL Boundless(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_db_engines"
sidebar_current: "docs-tencentcloud-data_source-tdmysql_db_engines"
description: |-
  Use this data source to query the list of TDSQL Boundless (tdmysql) DB engine versions.
---

# tencentcloud_tdmysql_db_engines

Use this data source to query the list of TDSQL Boundless (tdmysql) DB engine versions.

## Example Usage

```hcl
data "tencentcloud_tdmysql_db_engines" "example" {}
```

## Argument Reference

The following arguments are supported:



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `description` - Engine description.
  * `is_support_param_template` - Whether the parameter template is supported.
  * `is_support_password_complexity` - Whether password complexity is supported.
  * `is_support_serverless` - Whether the Serverless mode is supported.
  * `is_support_tde` - Whether transparent encryption (TDE) is supported.
  * `name` - Engine name. Note: this field may return null, indicating that no valid values can be obtained.
  * `new` - Whether it is the latest version. Note: this field may return null, indicating that no valid values can be obtained.
  * `sql_mode` - Supported compatibility modes, separated by ,.
  * `type` - Engine type. Note: this field may return null, indicating that no valid values can be obtained.
  * `version` - Engine version. Note: this field may return null, indicating that no valid values can be obtained.
* `id` - ID of the data source, a random token.
* `items` - List of DB engine information.


