Provides a resource to create a VPC private nat gateway

Example Usage

```hcl
resource "tencentcloud_vpc_private_nat_gateway" "example" {
  nat_gateway_name = "tf-example"
  vpc_id           = "vpc-i5yyodl9"
  tags = {
    createBy = "Terraform"
  }
}
```

Import

VPC private nat gateway can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_private_nat_gateway.example intranat-ljdy849x
```
