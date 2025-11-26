Provides a resource to create a IGTM monitor

Example Usage

```hcl
resource "tencentcloud_igtm_monitor" "example" {
  monitor_name          = "tf-example"
  check_protocol        = "PING"
  check_interval        = 60
  timeout               = 5
  fail_times            = 1
  fail_rate             = 50
  detector_style        = "INTERNAL"
  detector_group_ids    = [30, 31, 32, 34, 37, 38, 39, 1, 2, 3, 7, 8, 9, 10, 11, 12]
  ping_num              = 20
  tcp_port              = 443
  path                  = "/"
  return_code_threshold = 500
  enable_redirect       = "DISABLED"
  enable_sni            = "DISABLED"
  packet_loss_rate      = 90
  continue_period       = 1
}
```

Import

IGTM monitor can be imported using the id, e.g.

```
terraform import tencentcloud_igtm_monitor.example 12355
```
