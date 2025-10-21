---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_export_instance_error_logs"
sidebar_current: "docs-tencentcloud-resource-cynosdb_export_instance_error_logs"
description: |-
  Provides a resource to create a cynosdb export_instance_error_logs
---

# tencentcloud_cynosdb_export_instance_error_logs

Provides a resource to create a cynosdb export_instance_error_logs

## Example Usage

```hcl
resource "tencentcloud_cynosdb_export_instance_error_logs" "export_instance_error_logs" {
  instance_id   = "cynosdbmysql-ins-afqx1hy0"
  start_time    = "2022-01-01 12:00:00"
  end_time      = "2022-01-01 14:00:00"
  log_levels    = ["note"]
  key_words     = ["content"]
  file_type     = "csv"
  order_by      = "Timestamp"
  order_by_type = "ASC"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `end_time` - (Optional, String, ForceNew) Latest log time.
* `file_type` - (Optional, String, ForceNew) File type, optional values: csv, original.
* `key_words` - (Optional, Set: [`String`], ForceNew) keyword.
* `log_levels` - (Optional, Set: [`String`], ForceNew) Log level.
* `order_by_type` - (Optional, String, ForceNew) ASC or DESC.
* `order_by` - (Optional, String, ForceNew) Optional value Timestamp.
* `start_time` - (Optional, String, ForceNew) Log earliest time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `error_log_item_export` - List of instances in the read-write instance group.
  * `content` - log content.
  * `level` - Log level, optional values note, warning, error.
  * `timestamp` - time.


