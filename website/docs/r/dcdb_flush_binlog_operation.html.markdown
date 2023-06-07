---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_flush_binlog_operation"
sidebar_current: "docs-tencentcloud-resource-dcdb_flush_binlog_operation"
description: |-
  Provides a resource to create a dcdb flush_binlog_operation
---

# tencentcloud_dcdb_flush_binlog_operation

Provides a resource to create a dcdb flush_binlog_operation

## Example Usage

```hcl
resource "tencentcloud_dcdb_flush_binlog_operation" "flush_binlog_operation" {
  instance_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



