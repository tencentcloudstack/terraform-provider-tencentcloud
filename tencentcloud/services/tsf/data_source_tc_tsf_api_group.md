Use this data source to query detailed information of tsf api_group

Example Usage

```hcl
data "tencentcloud_tsf_api_group" "api_group" {
  search_word = "xxx01"
  group_type = "ms"
  auth_type = "none"
  status = "released"
  order_by = "created_time"
  order_type = 0
  gateway_instance_id = "gw-ins-lvdypq5k"
}
```