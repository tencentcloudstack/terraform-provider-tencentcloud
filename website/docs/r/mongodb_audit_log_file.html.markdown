---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_audit_log_file"
sidebar_current: "docs-tencentcloud-resource-mongodb_audit_log_file"
description: |-
  Provides a resource to create a MongoDB audit log file
---

# tencentcloud_mongodb_audit_log_file

Provides a resource to create a MongoDB audit log file

## Example Usage

```hcl
resource "tencentcloud_mongodb_audit_log_file" "example" {
  instance_id = "cmgo-5aqo4yf7"
  start_time  = "2026-06-01 10:29:20"
  end_time    = "2026-06-01 10:39:20"
  order       = "ASC"
  order_by    = "timestamp"

  filter {
    host        = ["10.0.0.1"]
    user        = ["admin"]
    exec_time   = 100
    affect_rows = 10
    atype       = ["insert", "update"]
    result      = ["ok"]
    param       = ["keyword"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time, format: "2021-07-12 10:39:20".
* `instance_id` - (Required, String, ForceNew) Instance ID, the format is: cmgo-xfts****.
* `start_time` - (Required, String) Start time, format: "2021-07-12 10:29:20".
* `filter` - (Optional, List) Filter conditions.
* `order_by` - (Optional, String) Sort field. Valid values: `timestamp`, `affectRows`, `execTime`.
* `order` - (Optional, String) Sort order. Valid values: `ASC`, `DESC`.

The `filter` object supports the following:

* `affect_rows` - (Optional, Int) Minimum affected rows.
* `atype` - (Optional, List) Operation types.
* `exec_time` - (Optional, Int) Minimum execution time in ms.
* `host` - (Optional, List) Client addresses.
* `param` - (Optional, List) Keywords to filter logs.
* `result` - (Optional, List) Execution results.
* `user` - (Optional, List) Usernames.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `file_name` - The generated audit log file name.
* `items` - Audit log file details.
  * `create_time` - Creation time.
  * `download_url` - Download URL.
  * `err_msg` - Error message.
  * `file_name` - File name.
  * `file_size` - File size in KB.
  * `progress_rate` - Download progress.
  * `status` - File status. Valid values: `creating`, `failed`, `success`.


## Import

mongodb audit_log_file can be imported using the composite instance_id#file_name, e.g.

```
terraform import tencentcloud_mongodb_audit_log_file.example cmgo-5aqo4yf7#1309118522_cmgo-5aqo4yf7_1780474413_109642711.csv
```

