Use this data source to query detailed information of DLC describe user info

Example Usage

```hcl
data "tencentcloud_dlc_describe_user_info" "example" {
  user_id = "100021240189"
  type    = "Group"
  sort_by = "create-time"
  sorting = "desc"
}
```
