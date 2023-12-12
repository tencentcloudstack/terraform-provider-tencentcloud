Provides a resource to create a monitor grafana_dns_config

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_dns_config" "grafana_dns_config" {
  instance_id  = "grafana-dp2hnnfa"
  name_servers = ["10.1.2.1", "10.1.2.2", "10.1.2.3"]
}
```

Import

monitor grafana_dns_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_dns_config.grafana_dns_config instance_id
```