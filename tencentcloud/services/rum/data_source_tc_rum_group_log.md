Use this data source to query detailed information of rum group_log

Example Usage

```hcl
data "tencentcloud_rum_group_log" "group_log" {
  order_by    = "desc"
  start_time  = 1625444040000
  query       = "id:123 AND type:\"log\""
  end_time    = 1625454840000
  project_id  = 1
  group_field = "level"
}
```