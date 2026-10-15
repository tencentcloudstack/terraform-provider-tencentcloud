---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_isolate_hour_instance_operation"
sidebar_current: "docs-tencentcloud-resource-dcdb_isolate_hour_instance_operation"
description: |-
  Provides a resource to create a dcdb isolate_hour_instance_operation
---

# tencentcloud_dcdb_isolate_hour_instance_operation

Provides a resource to create a dcdb isolate_hour_instance_operation

## Example Usage

```hcl
resource "tencentcloud_dcdb_isolate_hour_instance_operation" "isolate_hour_instance_operation" {
  instance_ids = local.dcdb_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [IsolateHourDCDBInstance](https://cloud.tencent.com/document/api/557/70420) - 隔离TDSQL按量计费实例

