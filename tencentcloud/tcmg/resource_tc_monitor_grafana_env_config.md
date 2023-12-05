Provides a resource to create a monitor grafana_env_config

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_env_config" "grafana_env_config" {
  instance_id = "grafana-dp2hnnfa"
  envs = {
    "aaa" = "ccc"
    "bbb"  = "ccc"
  }
}
```

Import

monitor grafana_env_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_env_config.grafana_env_config instance_id
```