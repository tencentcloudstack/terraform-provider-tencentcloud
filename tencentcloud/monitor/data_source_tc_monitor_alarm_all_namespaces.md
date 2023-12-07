Use this data source to query detailed information of monitor alarm_all_namespaces

Example Usage

```hcl
data "tencentcloud_monitor_alarm_all_namespaces" "alarm_all_namespaces" {
  scene_type    = "ST_ALARM"
  module        = "monitor"
  monitor_types = ["MT_QCE"]
  ids           = ["qaap_tunnel_l4_listeners"]
}
```