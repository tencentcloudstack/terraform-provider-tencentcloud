Provides a resource to config a NAT gateway flow monitor.

Example Usage

```hcl
resource "tencentcloud_nat_gateway_flow_monitor" "example" {
  gateway_id = "nat-e6u6axsm"
  enable     = true
}
```

Import

NAT gateway flow monitor can be imported using the id, e.g.

```
$ terraform import tencentcloud_nat_gateway_flow_monitor.example nat-e6u6axsm
```
