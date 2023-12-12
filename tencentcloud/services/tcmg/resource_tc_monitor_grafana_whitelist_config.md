Provides a resource to create a monitor grafana_whitelist_config

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_whitelist_config" "grafana_whitelist_config" {
  instance_id = "grafana-dp2hnnfa"
  whitelist   = ["10.1.1.1", "10.1.1.2", "10.1.1.3"]
}
```

Import

monitor grafana_whitelist_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_whitelist_config.grafana_whitelist_config instance_id
```