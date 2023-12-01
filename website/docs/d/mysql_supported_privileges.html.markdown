---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_supported_privileges"
sidebar_current: "docs-tencentcloud-datasource-mysql_supported_privileges"
description: |-
  Use this data source to query detailed information of mysql supported_privileges
---

# tencentcloud_mysql_supported_privileges

Use this data source to query detailed information of mysql supported_privileges

## Example Usage

```hcl
data "tencentcloud_mysql_supported_privileges" "supported_privileges" {
  instance_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `column_supported_privileges` - The database column permissions supported by the instance.
* `database_supported_privileges` - Database permissions supported by the instance.
* `global_supported_privileges` - Global permissions supported by the instance.
* `table_supported_privileges` - Database table permissions supported by the instance.


