---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_compare_task_stop_operation"
sidebar_current: "docs-tencentcloud-resource-dts_compare_task_stop_operation"
description: |-
  Provides a resource to create a dts compare_task_stop_operation
---

# tencentcloud_dts_compare_task_stop_operation

Provides a resource to create a dts compare_task_stop_operation

## Example Usage

```hcl
resource "tencentcloud_dts_compare_task_stop_operation" "compare_task_stop_operation" {
  job_id          = "dts-8yv4w2i1"
  compare_task_id = "dts-8yv4w2i1-cmp-37skmii9"
}
```

## Argument Reference

The following arguments are supported:

* `compare_task_id` - (Required, String, ForceNew) Compare task id.
* `job_id` - (Required, String, ForceNew) job id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [StopCompare](https://www.tencentcloud.com/document/product/571/51800) - This API is used to stop a data consistency check task.

