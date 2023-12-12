Use this data source to query detailed information of gaap listener statistics

Example Usage

```hcl
data "tencentcloud_gaap_listener_statistics" "listener_statistics" {
  listener_id = "listener-xxxxxx"
  start_time = "2023-10-19 00:00:00"
  end_time = "2023-10-19 23:59:59"
  metric_names = ["InBandwidth", "OutBandwidth", "InPackets", "OutPackets", "Concurrent"]
  granularity = 300
}
```