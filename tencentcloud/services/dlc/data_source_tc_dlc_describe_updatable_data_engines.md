Use this data source to query detailed information of dlc describe_updatable_data_engines

Example Usage

```hcl
data "tencentcloud_dlc_describe_updatable_data_engines" "describe_updatable_data_engines" {
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
  }
```