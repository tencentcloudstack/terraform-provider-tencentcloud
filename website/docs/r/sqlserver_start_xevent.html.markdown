---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_start_xevent"
sidebar_current: "docs-tencentcloud-resource-sqlserver_start_xevent"
description: |-
  Provides a resource to create a sqlserver start_xevent
---

# tencentcloud_sqlserver_start_xevent

Provides a resource to create a sqlserver start_xevent

## Example Usage

```hcl
resource "tencentcloud_sqlserver_start_xevent" "start_xevent" {
  instance_id = "mssql-gyg9xycl"
  event_config {
    event_type = "slow"
    threshold  = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `event_config` - (Required, List, ForceNew) Whether to start or stop an extended event.
* `instance_id` - (Required, String, ForceNew) Instance ID.

The `event_config` object supports the following:

* `event_type` - (Required, String) Event type. Valid values: slow (set threshold for slow SQL ), blocked (set threshold for the blocking and deadlock).
* `threshold` - (Required, Int) Threshold in milliseconds. Valid values: 0(disable), non-zero (enable).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



