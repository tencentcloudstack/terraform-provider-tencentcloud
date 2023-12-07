Use this data source to query detailed information of mps schedules

Example Usage

Query the enabled schedules.

```hcl
data "tencentcloud_mps_schedules" "schedules" {
  status       = "Enabled"
}
```

Query the specified one.

```hcl
data "tencentcloud_mps_schedules" "schedules" {
  schedule_ids = [%d]
  trigger_type = "CosFileUpload"
  status       = "Enabled"
}
```