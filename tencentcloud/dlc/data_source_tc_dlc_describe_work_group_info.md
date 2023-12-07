Use this data source to query detailed information of dlc describe_work_group_info

Example Usage

```hcl
data "tencentcloud_dlc_describe_work_group_info" "describe_work_group_info" {
  work_group_id = 23181
  type = "User"
  sort_by = "create-time"
  sorting = "desc"
  }
```