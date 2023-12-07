Use this data source to query vpc subnets information.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "guagua_vpc_instance_test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "guagua_vpc_subnet_test"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_vpc_subnets" "id_instances" {
  subnet_id = tencentcloud_subnet.subnet.id
}

data "tencentcloud_vpc_subnets" "name_instances" {
  name = tencentcloud_subnet.subnet.name
}

data "tencentcloud_vpc_subnets" "tags_instances" {
  tags = tencentcloud_subnet.subnet.tags
}
```