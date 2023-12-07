Use this data source to query detailed information of emr auto_scale_records

Example Usage

```hcl
data "tencentcloud_emr_auto_scale_records" "auto_scale_records" {
  instance_id = "emr-bpum4pad"
  filters {
    key   = "StartTime"
    value = "2006-01-02 15:04:05"
  }
}
```
