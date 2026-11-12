Provides a resource to create a NAT Gateway SNat rule.

Example Usage

If resource_type is SUBNET

```hcl
resource "tencentcloud_nat_gateway_snat" "example" {
  nat_gateway_id    = "nat-jmf9955g"
  resource_type     = "SUBNET"
  subnet_id         = "subnet-hhi88a58"
  subnet_cidr_block = "10.0.30.0/24"
  description       = "remark."
  public_ip_addr = [
    "106.55.167.123",
  ]
}
```

If resource_type is NETWORKINTERFACE 

```hcl
resource "tencentcloud_nat_gateway_snat" "example" {
  nat_gateway_id           = "nat-jmf9955g"
  resource_type            = "NETWORKINTERFACE"
  instance_id              = "ins-p97oyxt6"
  instance_private_ip_addr = "10.0.30.62"
  description              = "remark."
  public_ip_addr = [
    "106.55.167.123",
  ]
}
```

If resource_type is USERDEFINED 

```hcl
resource "tencentcloud_nat_gateway_snat" "example" {
  nat_gateway_id           = "nat-jmf9955g"
  resource_type            = "USERDEFINED"
  instance_private_ip_addr = "10.0.40.0/24"
  description              = "remark."
  public_ip_addr = [
    "106.55.167.123",
  ]
}
```

Import

VPN gateway route can be imported using the id, the id format must be '{nat_gateway_id}#{resource_id}', resource_id range `subnet_id`, `instance_id` and `""`, e.g.

SUBNET SNat

```
terraform import tencentcloud_nat_gateway_snat.example nat-jmf9955g#subnet-hhi88a58
```

NETWORKINTERFACT SNat

```
terraform import tencentcloud_nat_gateway_snat.example nat-jmf9955g#ins-p97oyxt6
```

USERDEFINED SNat

```
terraform import tencentcloud_nat_gateway_snat.example nat-jmf9955g#
```
