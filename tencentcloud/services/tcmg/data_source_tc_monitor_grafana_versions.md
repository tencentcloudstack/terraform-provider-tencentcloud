Use this data source to query available Grafana versions of monitor

Example Usage

```hcl
data "tencentcloud_monitor_grafana_versions" "grafana_versions" {
}

output "available_versions" {
  value = data.tencentcloud_monitor_grafana_versions.grafana_versions.versions
}
```
