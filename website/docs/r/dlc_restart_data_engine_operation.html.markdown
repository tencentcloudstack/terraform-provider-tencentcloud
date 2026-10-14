---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_restart_data_engine_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_restart_data_engine_operation"
description: |-
  Provides a resource to create a DLC restart data engine
---

# tencentcloud_dlc_restart_data_engine_operation

Provides a resource to create a DLC restart data engine

## Example Usage

```hcl
resource "tencentcloud_dlc_restart_data_engine_operation" "example" {
  data_engine_id   = "DataEngine-g5ds87d8"
  forced_operation = false
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String, ForceNew) Engine ID.
* `forced_operation` - (Optional, Bool, ForceNew) Whether to restart by force and ignore tasks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeDataEngine](https://cloud.tencent.com/document/api/1342/99274) - This API is used to obtain detailed data engine information based on names.
* Create: [DescribeDataEngines](https://cloud.tencent.com/document/api/1342/86308)
* Create: [RestartDataEngine](https://cloud.tencent.com/document/api/1342/99268) - This API is used to restart engines.

