Provides details about a specific VPC.

This resource can prove useful when a module accepts a vpc id as an input variable and needs to, for example, determine the CIDR block of that VPC.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_vpc_instances.

Example Usage

```hcl
variable "vpc_id" {}

data "tencentcloud_vpc" "selected" {
  id = var.vpc_id
}

resource "tencentcloud_subnet" "main" {
  name              = "my test subnet"
  cidr_block        = cidrsubnet(data.tencentcloud_vpc.selected.cidr_block, 4, 1)
  availability_zone = "eu-frankfurt-1"
  vpc_id            = data.tencentcloud_vpc.selected.id
}
```