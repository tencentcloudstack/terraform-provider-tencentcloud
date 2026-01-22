Provides a resource to create a DLC standard engine resource group

~> **NOTE:** If you are creating a machine learning resource group for the first time, you need to contact DLC product for whitelisting.

~> **NOTE:** Field `auto_pause_time` is meaningful only when the values ​​of fields `auto_launch` and `auto_pause` are 0.

~> **NOTE:** If you need to set the `static_config_pairs` or `dynamic_config_pairs`, it is recommended to use resource `tencentcloud_dlc_standard_engine_resource_group_config_info`.

Example Usage

Only SQL analysis resource group 

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

Machine learning resource group

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
