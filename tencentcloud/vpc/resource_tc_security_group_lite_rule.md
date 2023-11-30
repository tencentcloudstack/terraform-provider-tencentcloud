Provide a resource to create security group some lite rules quickly.

-> **NOTE:** It can't be used with tencentcloud_security_group_rule, and don't create multiple tencentcloud_security_group_rule resources, otherwise it may cause problems.

Example Usage

```hcl
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"
}

resource "tencentcloud_security_group_lite_rule" "foo" {
  security_group_id = tencentcloud_security_group.foo.id

  ingress = [
    "ACCEPT#192.168.1.0/24#80#TCP",
    "DROP#8.8.8.8#80,90#UDP",
    "ACCEPT#0.0.0.0/0#80-90#TCP",
    "ACCEPT#sg-7ixn3foj#80-90#TCP",
    "ACCEPT#ipm-epjq5kn0#80-90#TCP",
    "ACCEPT#ipmg-3loavam6#80-90#TCP",
    "ACCEPT#0.0.0.0/0##ppm-xxxxxxxx"
    "ACCEPT#0.0.0.0/0##ppmg-xxxxxxxx"
  ]

  egress = [
    "ACCEPT#192.168.0.0/16#ALL#TCP",
    "ACCEPT#10.0.0.0/8#ALL#ICMP",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]
}
```

Import

Security group lite rule can be imported using the id, e.g.

```
  $ terraform import tencentcloud_security_group_lite_rule.foo sg-ey3wmiz1
```