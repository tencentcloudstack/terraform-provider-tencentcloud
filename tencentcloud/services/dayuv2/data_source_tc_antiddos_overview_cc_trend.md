Use this data source to query detailed information of antiddos overview_cc_trend

Example Usage

```hcl
data "tencentcloud_antiddos_overview_cc_trend" "overview_cc_trend" {
  period = 300
  start_time = "2023-11-20 00:00:00"
  end_time = "2023-11-21 00:00:00"
  metric_name = "inqps"
  business = "bgpip"
}
```