Provides a resource to create a DLC update data engine config operation

Example Usage

```hcl
resource "tencentcloud_dlc_update_data_engine_config_operation" "example" {
  data_engine_id             = "DataEngine-o3lzpqpo"
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
}
```