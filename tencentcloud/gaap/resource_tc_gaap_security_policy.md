Provides a resource to create a security policy of GAAP proxy.

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
  action   = "DROP"
}
```

Import

GAAP security policy can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_security_policy.foo pl-xxxx
```