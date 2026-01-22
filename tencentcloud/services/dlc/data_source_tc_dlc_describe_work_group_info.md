Use this data source to query detailed information of DLC describe work group info

Example Usage

```hcl
data "tencentcloud_dlc_describe_work_group_info" "example" {
  work_group_id = 70220
  type          = "User"
  sort_by       = "create-time"
  sorting       = "desc"
}
```
