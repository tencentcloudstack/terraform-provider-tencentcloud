---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_sync_check_job_operation"
sidebar_current: "docs-tencentcloud-resource-dts_sync_check_job_operation"
description: |-
  Provides a resource to create a dts sync_check_job_operation
---

# tencentcloud_dts_sync_check_job_operation

Provides a resource to create a dts sync_check_job_operation

## Example Usage

```hcl
resource "tencentcloud_dts_sync_check_job_operation" "sync_check_job_operation" {
  job_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String, ForceNew) Sync job id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateCheckSyncJob](https://www.tencentcloud.com/document/product/571/51796) - This API is used to verify a sync task by checking required parameters and peripheral configuration.
* Create: [DescribeCheckSyncJobResult](https://www.tencentcloud.com/document/product/571/51794) - This API is used to query the result of the sync check task and check the required parameters and peripheral configurations.

