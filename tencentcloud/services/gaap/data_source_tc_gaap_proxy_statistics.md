Use this data source to query detailed information of gaap proxy statistics

Example Usage

```hcl
data "tencentcloud_gaap_proxy_statistics" "proxy_statistics" {
  proxy_id = "link-m9t4yho9"
  start_time = "2024-05-20 00:00:00"
  end_time = "2024-05-20 23:59:59"
  metric_names = ["InBandwidth", "OutBandwidth", "InFlow", "OutFlow", "InPackets", "OutPackets", "Concurrent", "HttpQPS", "HttpsQPS", "Latency", "PacketLoss"]
  granularity = 300
}
```