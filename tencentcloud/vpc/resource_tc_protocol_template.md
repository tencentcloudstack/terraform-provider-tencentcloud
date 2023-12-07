Provides a resource to manage protocol template.

Example Usage

```hcl
resource "tencentcloud_protocol_template" "foo" {
  name                = "protocol-template-test"
  protocols = ["tcp:80","udp:all","icmp:10-30"]
}
```

Import

Protocol template can be imported using the id, e.g.

```
$ terraform import tencentcloud_protocol_template.foo ppm-nwrggd14
```