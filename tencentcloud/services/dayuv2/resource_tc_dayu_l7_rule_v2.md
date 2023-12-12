Use this resource to create dayu new layer 7 rule

~> **NOTE:** This resource only support resource Anti-DDoS of type `bgpip`

Example Usage

```hcl
resource "tencentcloud_dayu_l7_rule_v2" "tencentcloud_dayu_l7_rule_v2" {
  resource_type="bgpip"
  resource_id="bgpip-000004xe"
  resource_ip="119.28.217.162"
  rule {
    keep_enable=false
    keeptime=0
    source_list {
      source="1.2.3.5"
      weight=100
    }
	source_list {
      source="1.2.3.6"
      weight=100
    }
    lb_type=1
    protocol="http"
    source_type=2
    domain="github.com"
  }
}
```