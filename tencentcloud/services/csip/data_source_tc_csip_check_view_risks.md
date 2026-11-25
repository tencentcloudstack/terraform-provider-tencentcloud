Use this data source to query the cloud asset configuration risks under the check item view.

Example Usage

```hcl
data "tencentcloud_csip_check_view_risks" "example" {
  member_id = ["mem-6wfo0fzks3"]
  filters {
    name   = "Status"
    values = ["CHECKING"]
  }
  order = "desc"
  by    = "update_time"
}
```
