Provides a resource to create as instance refresh

Example Usage

```hcl
resource "tencentcloud_as_start_instance_refresh" "example" {
  auto_scaling_group_id = "asg-8n7fdm28"
  refresh_mode          = "ROLLING_UPDATE_RESET"
  refresh_settings {
    check_instance_target_health = false
    rolling_update_settings {
      batch_number = 1
      batch_pause  = "AUTOMATIC"
      max_surge    = 1
      fail_process = "AUTO_PAUSE"
    }
    check_instance_target_health_timeout = 1800
  }

  timeouts {
    create = "10m"
  }
}
```
