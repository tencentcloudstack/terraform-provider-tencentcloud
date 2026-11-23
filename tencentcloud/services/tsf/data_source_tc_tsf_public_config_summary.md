Use this data source to query detailed information of tsf public config summary

Example Usage

```hcl
data "tencentcloud_tsf_public_config_summary" "example" {
  search_word                = "test"
  order_by                   = "last_update_time"
  order_type                 = 0
  disable_program_auth_check = true
  config_id_list             = ["dcfg-p-ygbdw5mv"]
}
```