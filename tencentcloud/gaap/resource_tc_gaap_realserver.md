Provides a resource to create a GAAP realserver.

Example Usage

```hcl
resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"

  tags = {
    test = "test"
  }
}
```

Import

GAAP realserver can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_realserver.foo rs-4ftghy6
```