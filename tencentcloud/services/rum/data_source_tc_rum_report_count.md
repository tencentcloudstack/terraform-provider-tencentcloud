Use this data source to query detailed information of rum report_count

Example Usage

```hcl
data "tencentcloud_rum_report_count" "report_count" {
  start_time  = 1625444040
  end_time    = 1625454840
  project_id  = 1
  report_type = "log"
}
```