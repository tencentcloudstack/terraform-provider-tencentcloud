---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_standard_engine_resource_group"
sidebar_current: "docs-tencentcloud-resource-dlc_standard_engine_resource_group"
description: |-
  Provides a resource to create a DLC standard engine resource group
---

# tencentcloud_dlc_standard_engine_resource_group

Provides a resource to create a DLC standard engine resource group

~> **NOTE:** If you are creating a machine learning resource group for the first time, you need to contact DLC product for whitelisting.

~> **NOTE:** Field `auto_pause_time` is meaningful only when the values ​​of fields `auto_launch` and `auto_pause` are 0.

~> **NOTE:** If you need to set the `static_config_pairs` or `dynamic_config_pairs`, it is recommended to use resource `tencentcloud_dlc_standard_engine_resource_group_config_info`.

## Example Usage

### Only SQL analysis resource group

```hcl
resource "tencentcloud_dlc_standard_engine_resource_group" "example" {
  engine_resource_group_name = "tf-example"
  data_engine_name           = "tf-engine"
  auto_launch                = 0
  auto_pause                 = 0
  auto_pause_time            = 10
  static_config_pairs {
    config_item  = "key"
    config_value = "value"
  }

  dynamic_config_pairs {
    config_item  = "key"
    config_value = "value"
  }
  max_concurrency      = 5
  resource_group_scene = "SparkSQL"
  spark_spec_mode      = "fast"
  spark_size           = 16
  running_state        = true
}
```

### Machine learning resource group

```hcl
resource "tencentcloud_dlc_standard_engine_resource_group" "example" {
  engine_resource_group_name = "tf-example"
  data_engine_name           = "tf-engine"
  max_concurrency            = 5
  resource_group_scene       = "Artificial-Intelligence"
  spark_spec_mode            = "fast"
  spark_size                 = 16
  frame_type                 = "machine-learning"
  size                       = 16
  python_cu_spec             = "large"
  image_type                 = "built-in"
  image_version              = "97319759-0b80-48b4-a7a7-436d9ef3b666"
  image_name                 = "pytorch-v2.5.1"
  running_state              = false
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_name` - (Required, String) Standard engine name.
* `engine_resource_group_name` - (Required, String) Standard engine resource group name.
* `auto_launch` - (Optional, Int) Automatic start (task submission automatically pulls up the resource group) 0-automatic start, 1-not automatic start.
* `auto_pause_time` - (Optional, Int) Automatic suspension time, in minutes, with a value range of 1-999 (after no tasks have reached AutoPauseTime, the resource group will automatically suspend).
* `auto_pause` - (Optional, Int) Automatically suspend resource groups. 0 - Automatically suspend, 1 - Not automatically suspend.
* `driver_cu_spec` - (Optional, String) Driver CU specifications: Currently supported: small (default, 1 CU), medium (2 CU), large (4 CU), xlarge (8 CU). Memory CUs are CPUs with a ratio of 1:8, m.small (1 CU memory), m.medium (2 CU memory), m.large (4 CU memory), and m.xlarge (8 CU memory).
* `dynamic_config_pairs` - (Optional, List) Dynamic parameters of the resource group, effective in the next task.
* `executor_cu_spec` - (Optional, String) Executor CU specifications: Currently supported: small (default, 1 CU), medium (2 CU), large (4 CU), xlarge (8 CU). Memory CUs are CPUs with a ratio of 1:8, m.small (1 CU memory), m.medium (2 CU memory), m.large (4 CU memory), and m.xlarge (8 CU memory).
* `frame_type` - (Optional, String) The framework type of the AI type resource group, machine-learning, python, spark-ml, if not filled in, the default is machine-learning.
* `image_name` - (Optional, String) Image Name. 
Example value: image-xxx. If using a built-in image (ImageType is built-in), the ImageName for different frameworks is: machine-learning: pytorch-v2.5.1, scikit-learn-v1.6.0, tensorflow-v2.18.0, python: python-v3.10, spark-m: Standard-S 1.1.
* `image_type` - (Optional, String) Image type, build-in: built-in, custom: custom, if not filled in, the default is build-in.
* `image_version` - (Optional, String) Image ID.
* `max_concurrency` - (Optional, Int) The number of concurrent tasks is 5 by default.
* `max_executor_nums` - (Optional, Int) Maximum number of executors.
* `min_executor_nums` - (Optional, Int) Minimum number of executors.
* `network_config_names` - (Optional, Set: [`String`]) Network configuration name.
* `public_domain` - (Optional, String) Customized mirror domain name.
* `python_cu_spec` - (Optional, String) The resource limit for a Python stand-alone node in a Python resource group must be smaller than the resource limit for the resource group. Small: 1cu Medium: 2cu Large: 4cu Xlarge: 8cu 4xlarge: 16cu 8xlarge: 32cu 16xlarge: 64cu. If the resource type is high memory, add m before the type.
* `region_name` - (Optional, String) Custom image location.
* `registry_id` - (Optional, String) Custom image instance ID.
* `resource_group_scene` - (Optional, String) Resource group scenario.
* `running_state` - (Optional, Bool) The state of the resource group. true: launch standard engine resource group; false: pause standard engine resource group. Default is true.
* `size` - (Optional, Int) The AI resource group is valid, and the upper limit of available resources in the resource group must be less than the upper limit of engine resources.
* `spark_size` - (Optional, Int) Only the SQL resource group resource limit, only used for the express module.
* `spark_spec_mode` - (Optional, String) Only SQL resource group resource configuration mode, fast: fast mode, custom: custom mode.
* `static_config_pairs` - (Optional, List) Static parameters of the resource group, which require restarting the resource group to take effect.

The `dynamic_config_pairs` object supports the following:

* `config_item` - (Required, String) Configuration items.
* `config_value` - (Required, String) Configuration values.

The `static_config_pairs` object supports the following:

* `config_item` - (Required, String) Configuration items.
* `config_value` - (Required, String) Configuration values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `engine_resource_group_id` - Standard engine resource group ID.


