Use this data source to query detailed information of DLC describe updatable data engines

Example Usage

```hcl
data "tencentcloud_dlc_describe_updatable_data_engines" "example" {
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
}
```

Or

```hcl
data "tencentcloud_dlc_describe_updatable_data_engines" "example" {
  data_engine_config_command = "UpdateSparkSQLResultPath"
}
```
