Provides a resource to create a DLC standard engine resource group config info

~> **NOTE:** This resource must exclusive in one engine resource group, do not declare additional config resources of this conf context elsewhere.

~> **NOTE:** If you use the `tencentcloud_dlc_standard_engine_resource_group_config_info`. Please do not set `static_config_pairs` or `dynamic_config_pairs` in resource `tencentcloud_dlc_standard_engine_resource_group` simultaneously.

Example Usage

```hcl
resource "tencentcloud_dlc_standard_engine_resource_group_config_info" "example" {
  engine_resource_group_name = "tf-example"
  static_conf_context {
    params {
      config_item  = "item1"
      config_value = "value1"
    }

    params {
      config_item  = "item2"
      config_value = "value2"
    }
  }

  dynamic_conf_context {
    params {
      config_item  = "item3"
      config_value = "value3"
    }
  }
}
```

Import

DLC standard engine resource group config info can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_standard_engine_resource_group_config_info.example tf-example
```
