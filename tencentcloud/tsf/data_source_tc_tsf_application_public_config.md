Use this data source to query detailed information of tsf application_public_config

Example Usage

```hcl
data "tencentcloud_tsf_application_public_config" "application_public_config" {
  config_id = "dcfg-p-evjrbgly"
  # config_id_list = [""]
  config_name = "dsadsa"
  config_version = "123"
}
```