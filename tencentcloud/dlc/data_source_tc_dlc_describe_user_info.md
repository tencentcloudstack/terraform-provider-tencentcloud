Use this data source to query detailed information of dlc describe_user_info

Example Usage

```hcl
data "tencentcloud_dlc_describe_user_info" "describe_user_info" {
  user_id = "100032772113"
  type = "Group"
  sort_by = "create-time"
  sorting = "desc"
}
```