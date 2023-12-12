Use this data source to query detailed information of antiddos overview ddos trend

Example Usage

```hcl
data "tencentcloud_antiddos_overview_ddos_trend" "overview_ddos_trend" {
  period = 300
  start_time = "2023-11-20 14:16:23"
  end_time = "2023-11-21 14:16:23"
  metric_name = "bps"
  business = "bgpip"
}
```