Use this data source to query detailed information of tsf application_file_config

Example Usage

```hcl
data "tencentcloud_tsf_application_file_config" "application_file_config" {
  config_id = "dcfg-f-4y4ekzqv"
  # config_id_list = [""]
  config_name = "file-log1"
  application_id = "application-2vzk6n3v"
  config_version = "1.2"
}
```