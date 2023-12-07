Provides a resource to create a monitor grafana_sso_cam_config

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_sso_cam_config" "grafana_sso_cam_config" {
  instance_id          = "grafana-dp2hnnfa"
  enable_sso_cam_check = false
}
```

Import

monitor grafana_sso_cam_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_sso_cam_config.grafana_sso_cam_config instance_id
```