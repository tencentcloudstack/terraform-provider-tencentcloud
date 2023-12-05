Provides a resource to create a monitor grafana_version_upgrade

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_version_upgrade" "grafana_version_upgrade" {
  instance_id = "grafana-dp2hnnfa"
  alias       = "v8.2.7"
}
```

Import

monitor grafana_version_upgrade can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_version_upgrade.grafana_version_upgrade instance_id
```