Use this data source to query detailed information of tsf public_config_summary

Example Usage

```hcl
data "tencentcloud_tsf_describe_public_config_summary" "describe_public_config_summary" {
  search_word = "test"
  order_by = "last_update_time"
  order_type = 0
  # config_tag_list = [""]
  disable_program_auth_check = true
  config_id_list = ["dcfg-p-ygbdw5mv"]
}
```