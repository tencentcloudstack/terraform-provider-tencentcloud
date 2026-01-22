Provides a resource to create a cdwpg cdwpg_dbconfig

Example Usage

```hcl
resource "tencentcloud_cdwpg_dbconfig" "cdwpg_dbconfig" {
  instance_id = "cdwpg-ua8wkqrt"
  node_config_params {
	  node_type = "cn"
	  parameter_name = "log_min_duration_statement"
	  parameter_value = "10001"
  }
}
```
