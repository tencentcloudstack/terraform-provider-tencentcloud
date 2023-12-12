Use this data source to query detailed information of antiddos bgp_biz_trend

Example Usage

```hcl
data "tencentcloud_antiddos_bgp_biz_trend" "bgp_biz_trend" {
  business    = "bgp-multip"
  start_time  = "2023-11-22 09:25:00"
  end_time    = "2023-11-22 10:25:00"
  metric_name = "intraffic"
  instance_id = "bgp-00000ry7"
  flag        = 0
}
```