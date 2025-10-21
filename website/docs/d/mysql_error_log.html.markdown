---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_error_log"
sidebar_current: "docs-tencentcloud-datasource-mysql_error_log"
description: |-
  Use this data source to query detailed information of mysql error_log
---

# tencentcloud_mysql_error_log

Use this data source to query detailed information of mysql error_log

## Example Usage

```hcl
data "tencentcloud_mysql_error_log" "error_log" {
  instance_id = "cdb-fitq5t9h"
  start_time  = 1683538307
  end_time    = 1686043908
  key_words   = ["Shutting"]
  inst_type   = "slave"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, Int) End timestamp. For example 1585142640.
* `instance_id` - (Required, String) instance id.
* `start_time` - (Required, Int) Start timestamp. For example 1585142640.
* `inst_type` - (Optional, String) Only valid when the instance is the master instance or disaster recovery instance, the optional value: slave, which means to pull the log of the slave machine.
* `key_words` - (Optional, Set: [`String`]) A list of keywords to match, up to 15 keywords are supported.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - The records returned.
  * `content` - error details.
  * `timestamp` - The time the error occurred.


