Use this resource to create dayu layer 4 rule

~> **NOTE:** This resource only support resource Anti-DDoS of type `bgpip` and `net`

Example Usage

```hcl
resource "tencentcloud_dayu_l4_rule" "test_rule" {
  resource_type             = "bgpip"
  resource_id               = "bgpip-00000294"
  name                      = "rule_test"
  protocol                  = "TCP"
  s_port                    = 80
  d_port                    = 60
  source_type               = 2
  health_check_switch       = true
  health_check_timeout      = 30
  health_check_interval     = 35
  health_check_health_num   = 5
  health_check_unhealth_num = 10
  session_switch            = false
  session_time              = 300

  source_list {
    source = "1.1.1.1"
    weight = 100
  }
  source_list {
    source = "2.2.2.2"
    weight = 50
  }
}
```