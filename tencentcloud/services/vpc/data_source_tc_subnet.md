Provides details about a specific VPC subnet.

This resource can prove useful when a module accepts a subnet id as an input variable and needs to, for example, determine the id of the VPC that the subnet belongs to.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_vpc_subnets.

Example Usage

```hcl
variable "subnet_id" {}
variable "vpc_id" {}

data "tencentcloud_subnet" "selected" {
  vpc_id    = var.vpc_id
  subnet_id = var.subnet_id
}

resource "tencentcloud_security_group" "default" {
  name        = "test subnet data"
  description = "test subnet data description"
}

resource "tencentcloud_security_group_rule" "subnet" {
  security_group_id = tencentcloud_security_group.default.id
  type              = "ingress"
  cidr_ip           = data.tencentcloud_subnet.selected.cidr_block
  ip_protocol       = "tcp"
  port_range        = "80,8080"
  policy            = "accept"
}
```