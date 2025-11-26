---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_monitor"
sidebar_current: "docs-tencentcloud-resource-igtm_monitor"
description: |-
  Provides a resource to create a IGTM monitor
---

# tencentcloud_igtm_monitor

Provides a resource to create a IGTM monitor

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `check_interval` - (Required, Int) Check interval (seconds), optional values 15 60 120 300.
* `check_protocol` - (Required, String) Detection protocol, optional values `PING`, `TCP`, `HTTP`, `HTTPS`.
* `detector_group_ids` - (Required, Set: [`Int`]) Detector group ID list separated by commas.
* `detector_style` - (Required, String) Monitoring node type, optional values AUTO INTERNAL OVERSEAS IPV6 ALL.
* `fail_rate` - (Required, Int) Failure rate, values 20 30 40 50 60 70 80 100, default value 50.
* `fail_times` - (Required, Int) Retry count, optional values 0, 1, 2.
* `monitor_name` - (Required, String) Monitor name.
* `timeout` - (Required, Int) Timeout time, unit seconds, optional values 2 3 5 10.
* `continue_period` - (Optional, Int) Continuous period count, optional values 1-5.
* `enable_redirect` - (Optional, String) Follow 3XX redirect, DISABLED for disabled, ENABLED for enabled, default disabled.
* `enable_sni` - (Optional, String) Enable SNI, DISABLED for disabled, ENABLED for enabled, default disabled.
* `host` - (Optional, String) Host setting, default is business domain name.
* `packet_loss_rate` - (Optional, Int) Packet loss rate alarm threshold, required when CheckProtocol=ping, values 10 30 50 80 90 100.
* `path` - (Optional, String) URL path, default is "/".
* `ping_num` - (Optional, Int) PING packet count, required when CheckProtocol=ping, optional values 20 50 100.
* `return_code_threshold` - (Optional, Int) Return error code threshold, optional values 400 and 500, default value 500.
* `tcp_port` - (Optional, Int) Check port, optional values between 1-65535.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `monitor_id` - Monitor ID.


## Import

IGTM monitor can be imported using the id, e.g.

```
terraform import tencentcloud_igtm_monitor.example 12355
```

