---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_log_files"
sidebar_current: "docs-tencentcloud-datasource-mariadb_log_files"
description: |-
  Use this data source to query detailed information of mariadb log_files
---

# tencentcloud_mariadb_log_files

Use this data source to query detailed information of mariadb log_files

## Example Usage

```hcl
data "tencentcloud_mariadb_log_files" "log_files" {
  instance_id = "tdsql-9vqvls95"
  type        = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID in the format of `tdsql-ow728lmc`.
* `type` - (Required, Int) Requested log type. Valid values: 1 (binlog), 2 (cold backup), 3 (errlog), 4 (slowlog).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `files` - Information such as `uri`, `length`, and `mtime` (modification time).
  * `file_name` - Filename.
  * `length` - File length.
  * `mtime` - Last modified time of log.
  * `uri` - Uniform resource identifier (URI) used during log download.
* `normal_prefix` - For an instance in a common network, this prefix plus URI can be used as the download address.
* `vpc_prefix` - For an instance in a VPC, this prefix plus URI can be used as the download address.


