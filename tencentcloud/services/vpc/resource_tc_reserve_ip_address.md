Provides a resource to create a VPC reserve ip address

Example Usage

```hcl
resource "tencentcloud_reserve_ip_address" "example" {
  vpc_id      = "vpc-i5yyodl9"
  subnet_id   = "subnet-hhi88a58"
  ip_address  = "10.0.30.27"
  name        = "tf-example"
  description = "remark."
  tags = {
    createBy = "Terraform"
  }
}
```

Import

VPC reserve ip address can be imported using the vpcId#reserveIpId, e.g.

```
terraform import tencentcloud_reserve_ip_address.example vpc-i5yyodl9#rsvip-cz3ces44
```
