Provides a resource to create a dlc update_data_engine_config_operation

Example Usage

```hcl
resource "tencentcloud_dlc_update_data_engine_config_operation" "update_data_engine_config_operation" {
  data_engine_id = "DataEngine-o3lzpqpo"
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
}
```