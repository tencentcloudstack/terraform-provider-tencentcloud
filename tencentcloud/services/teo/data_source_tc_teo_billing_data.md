Use this data source to query TEO billing data, such as traffic, bandwidth, and request metrics.

Example Usage

```hcl
data "tencentcloud_teo_billing_data" "example" {
  start_time  = "2025-01-01T00:00:00+08:00"
  end_time    = "2025-01-02T00:00:00+08:00"
  zone_ids    = ["zone-2qtuhspy7cr6"]
  metric_name = "acc_flux"
  interval    = "hour"
  filters {
    type  = "host"
    value = "test.example.com"
  }
  group_by = ["zone-id", "host"]
}
```
