Use this data source to query detailed information of dbbrain diag_events

Example Usage

```hcl
data "tencentcloud_dbbrain_diag_events" "diag_events" {
  instance_ids = ["%s"]
  start_time = "%s"
  end_time = "%s"
  severities = [1,4,5]
}
```