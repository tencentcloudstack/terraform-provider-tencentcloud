---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_inst_tables"
sidebar_current: "docs-tencentcloud-datasource-mysql_inst_tables"
description: |-
  Use this data source to query detailed information of mysql inst_tables
---

# tencentcloud_mysql_inst_tables

Use this data source to query detailed information of mysql inst_tables

## Example Usage

```hcl
data "tencentcloud_mysql_inst_tables" "inst_tables" {
  instance_id = "cdb-fitq5t9h"
  database    = "tf_ci_test"
  # table_regexp = ""
}
```

## Argument Reference

The following arguments are supported:

* `database` - (Required, String) The name of the database.
* `instance_id` - (Required, String) The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.
* `result_output_file` - (Optional, String) Used to save results.
* `table_regexp` - (Optional, String) Match the regular expression of the database table name, the rules are the same as MySQL official website.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - The returned database table information.


