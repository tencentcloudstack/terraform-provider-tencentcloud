Provides a resource to create a vpc private nat gateway

Example Usage

```hcl
resource "tencentcloud_vpc_private_nat_gateway" "private_nat_gateway" {
	nat_gateway_name = "xxx"
	vpc_id = "xxx"
}
```

Import

vpc private_nat_gateway can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_private_nat_gateway.private_nat_gateway private_nat_gateway_id
```