Use this data source to query security policy rule.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_security_policy" "foo" {
  proxy_id = tencentcloud_gaap_proxy.foo.id
  action   = "ACCEPT"
}

resource "tencentcloud_gaap_security_rule" "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

data "tencentcloud_gaap_security_rules" "protocol" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  protocol  = tencentcloud_gaap_security_rule.foo.protocol
}
```