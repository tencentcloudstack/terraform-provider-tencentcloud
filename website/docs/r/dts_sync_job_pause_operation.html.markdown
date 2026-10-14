---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_sync_job_pause_operation"
sidebar_current: "docs-tencentcloud-resource-dts_sync_job_pause_operation"
description: |-
  Provides a resource to create a dts sync_job_pause_operation
---

# tencentcloud_dts_sync_job_pause_operation

Provides a resource to create a dts sync_job_pause_operation

## Example Usage

```hcl
resource "tencentcloud_dts_sync_job_pause_operation" "sync_job_pause_operation" {
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

* Create: [DescribeSyncJobs](https://www.tencentcloud.com/document/product/571/51793) - This API is used to query the information of a sync task.
* Create: [PauseSyncJob](https://www.tencentcloud.com/document/product/571/53979) - This API is used to pause a data sync task.

