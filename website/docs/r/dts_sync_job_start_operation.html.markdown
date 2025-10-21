---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_sync_job_start_operation"
sidebar_current: "docs-tencentcloud-resource-dts_sync_job_start_operation"
description: |-
  Provides a resource to create a dts sync_job_start_operation
---

# tencentcloud_dts_sync_job_start_operation

Provides a resource to create a dts sync_job_start_operation

## Example Usage

```hcl
resource "tencentcloud_dts_sync_job_start_operation" "sync_job_start_operation" {
  job_id = "sync-werwfs23"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String, ForceNew) Synchronization instance id (i.e. identifies a synchronization job).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



