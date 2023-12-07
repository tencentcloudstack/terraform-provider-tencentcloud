Use this data source to query detailed information of tsf application_config

Example Usage

```hcl
data "tencentcloud_tsf_application_config" "application_config" {
  application_id = "app-123456"
  config_id = "config-123456"
  config_id_list =
  config_name = "test-config"
  config_version = "1.0"
}
```