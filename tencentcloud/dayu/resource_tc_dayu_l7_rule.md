Use this resource to create dayu layer 7 rule

~> **NOTE:** This resource only support resource Anti-DDoS of type `bgpip`

Example Usage

```hcl
resource "tencentcloud_dayu_l7_rule" "test_rule" {
  resource_type             = "bgpip"
  resource_id               = "bgpip-00000294"
  name                      = "rule_test"
  domain                    = "zhaoshaona.com"
  protocol                  = "https"
  switch                    = true
  source_type               = 2
  source_list               = ["1.1.1.1:80", "2.2.2.2"]
  ssl_id                    = "%s"
  health_check_switch       = true
  health_check_code         = 31
  health_check_interval     = 30
  health_check_method       = "GET"
  health_check_path         = "/"
  health_check_health_num   = 5
  health_check_unhealth_num = 10
}
```