---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_export_instance_slow_queries"
sidebar_current: "docs-tencentcloud-resource-cynosdb_export_instance_slow_queries"
description: |-
  Provides a resource to create a cynosdb export_instance_slow_queries
---

# tencentcloud_cynosdb_export_instance_slow_queries

Provides a resource to create a cynosdb export_instance_slow_queries

## Example Usage

```hcl
resource "tencentcloud_cynosdb_export_instance_slow_queries" "export_instance_slow_queries" {
  instance_id = "cynosdbmysql-ins-123"
  start_time  = "2022-01-01 12:00:00"
  end_time    = "2022-01-01 14:00:00"
  username    = "root"
  host        = "10.10.10.10"
  database    = "db1"
  file_type   = "csv"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `database` - (Optional, String, ForceNew) Database name.
* `end_time` - (Optional, String, ForceNew) Latest transaction start time.
* `file_type` - (Optional, String, ForceNew) File type, optional values: csv, original.
* `host` - (Optional, String, ForceNew) Client host.
* `start_time` - (Optional, String, ForceNew) Earliest transaction start time.
* `username` - (Optional, String, ForceNew) user name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `file_content` - Slow query export content.


