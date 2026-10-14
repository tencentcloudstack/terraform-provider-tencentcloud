---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_sync_job_isolate_operation"
sidebar_current: "docs-tencentcloud-resource-dts_sync_job_isolate_operation"
description: |-
  Provides a resource to create a dts sync_job_isolate_operation
---

# tencentcloud_dts_sync_job_isolate_operation

Provides a resource to create a dts sync_job_isolate_operation

## Example Usage

```hcl
resource "tencentcloud_dts_sync_job_isolate_operation" "sync_job_isolate_operation" {
  job_id = "sync-werwfs23"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String, ForceNew) Synchronization instance id (i.e. identifies a synchronization job).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeSyncJobs](https://cloud.tencent.com/document/api/571/82103) - This API is used to query the information of a sync task.
* Create: [IsolateSyncJob](https://cloud.tencent.com/document/api/571/82101) - This API is used to isolate a sync task. After the task is isolated, you can call the `DescribeSyncJobs` API to query its status, call `RecoverSyncJob` to recover it, or directly delete it.

