Provides a resource to create a CLB listener rule.

-> **NOTE:** This resource only be applied to the HTTP or HTTPS listeners.

Example Usage

```hcl
resource "tencentcloud_clb_listener_rule" "foo" {
  listener_id                = "lbl-hh141sn9"
  clb_id                     = "lb-k2zjp9lv"
  domain                     = "foo.net"
  url                        = "/bar"
  health_check_switch        = true
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  health_check_http_code     = 2
  health_check_http_path     = "Default Path"
  health_check_http_domain   = "Default Domain"
  health_check_http_method   = "GET"
  certificate_ssl_mode       = "MUTUAL"
  certificate_id             = "VjANRdz8"
  certificate_ca_id          = "VfqO4zkB"
  session_expire_time        = 30
  scheduler                  = "WRR"
}
```
Import

CLB listener rule can be imported using the id (version >= 1.47.0), e.g.

```
$ terraform import tencentcloud_clb_listener_rule.foo lb-7a0t6zqb#lbl-hh141sn9#loc-agg236ys
```