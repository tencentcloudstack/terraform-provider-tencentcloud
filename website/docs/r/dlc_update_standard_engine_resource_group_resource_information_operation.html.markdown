---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_update_standard_engine_resource_group_resource_information_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_update_standard_engine_resource_group_resource_information_operation"
description: |-
  Provides a resource to create a DLC update standard engine resource group resource information operation
---

# tencentcloud_dlc_update_standard_engine_resource_group_resource_information_operation

Provides a resource to create a DLC update standard engine resource group resource information operation

## Example Usage

```hcl
resource "tencentcloud_dlc_update_standard_engine_resource_group_resource_information_operation" "example" {
  engine_resource_group_name = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `engine_resource_group_name` - (Required, String, ForceNew) Engine resource group name.
* `driver_cu_spec` - (Optional, String, ForceNew) Driver CU specifications:
Currently supported: small (default, 1 CU), medium (2 CU), large (4 CU), xlarge (8 CU). Memory CUs are CPUs with a ratio of 1:8, m.small (1 CU memory), m.medium (2 CU memory), m.large (4 CU memory), and m.xlarge (8 CU memory).
* `executor_cu_spec` - (Optional, String, ForceNew) Executor CU specifications:
Currently supported: small (default, 1 CU), medium (2 CU), large (4 CU), xlarge (8 CU). Memory CUs are CPUs with a ratio of 1:8, m.small (1 CU memory), m.medium (2 CU memory), m.large (4 CU memory), and m.xlarge (8 CU memory).
* `frame_type` - (Optional, String, ForceNew) Framework Type.
* `image_name` - (Optional, String, ForceNew) Image name.
* `image_type` - (Optional, String, ForceNew) Image type, built-in image: built-in, custom image: custom.
* `image_version` - (Optional, String, ForceNew) Image version, image id.
* `max_executor_nums` - (Optional, Int, ForceNew) Maximum number of executors.
* `min_executor_nums` - (Optional, Int, ForceNew) Minimum number of executors.
* `public_domain` - (Optional, String, ForceNew) Customized mirror domain name.
* `python_cu_spec` - (Optional, String, ForceNew) The resource limit for a Python stand-alone node in a Python resource group must be smaller than the resource limit for the resource group. Small: 1cu Medium: 2cu Large: 4cu Xlarge: 8cu 4xlarge: 16cu 8xlarge: 32cu 16xlarge: 64cu. If the resource type is high memory, add m before the type.
* `region_name` - (Optional, String, ForceNew) Customize the image region.
* `registry_id` - (Optional, String, ForceNew) Custom image instance id.
* `size` - (Optional, Int, ForceNew) AI resource group resource limit.
* `spark_size` - (Optional, Int, ForceNew) SQL resource group resource limit only, only used in fast mode.
* `spark_spec_mode` - (Optional, String, ForceNew) Only SQL resource group resource configuration mode, fast: fast mode, custom: custom mode.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



