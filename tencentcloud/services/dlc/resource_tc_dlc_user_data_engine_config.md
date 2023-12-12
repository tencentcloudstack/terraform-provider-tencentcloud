Provides a resource to create a dlc user_data_engine_config

Example Usage

```hcl
resource "tencentcloud_dlc_user_data_engine_config" "user_data_engine_config" {
  data_engine_id = "DataEngine-cgkvbas6"
  data_engine_config_pairs {
		config_item = "qq"
		config_value = "ff"
  }
  session_resource_template {
		driver_size = "small"
		executor_size = "small"
		executor_nums = 1
		executor_max_numbers = 1
  }
}
```

Import

dlc user_data_engine_config can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_user_data_engine_config.user_data_engine_config user_data_engine_config_id
```