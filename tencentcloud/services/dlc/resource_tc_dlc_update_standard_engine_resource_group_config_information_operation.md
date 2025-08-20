Provides a resource to create a DLC update standard engine resource group config information operation

Example Usage

```hcl
resource "tencentcloud_dlc_update_standard_engine_resource_group_config_information_operation" "example" {
  engine_resource_group_name = "tf-example"
  update_conf_context {
    config_type = "StaticConfigType"
    params {
      config_item  = "spark.sql.shuffle.partitions"
      config_value = "300"
      operate      = "ADD"
    }
  }
}
```
