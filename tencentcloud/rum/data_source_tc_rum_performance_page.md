Use this data source to query detailed information of rum performance_page

Example Usage

```hcl
data "tencentcloud_rum_performance_page" "performance_page" {
  project_id = 1
  start_time = 1625444040
  end_time   = 1625454840
  type       = "pagepv"
  level      = "1"
}
```