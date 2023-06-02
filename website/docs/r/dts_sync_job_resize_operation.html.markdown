---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_sync_job_resize_operation"
sidebar_current: "docs-tencentcloud-resource-dts_sync_job_resize_operation"
description: |-
  Provides a resource to create a dts sync_job_resize_operation
---

# tencentcloud_dts_sync_job_resize_operation

Provides a resource to create a dts sync_job_resize_operation

## Example Usage

```hcl
resource "tencentcloud_dts_sync_job_resize_operation" "sync_job_resize_operation" {
  job_id             = "sync-werwfs23"
  new_instance_class = "large"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String, ForceNew) Synchronization instance id (i.e. identifies a synchronization job).
* `new_instance_class` - (Required, String, ForceNew) Task specification.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



