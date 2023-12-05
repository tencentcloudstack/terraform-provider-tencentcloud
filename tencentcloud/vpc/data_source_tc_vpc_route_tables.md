Use this data source to query vpc route tables information.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "guagua-ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "route_table" {
  vpc_id = tencentcloud_vpc.foo.id
  name   = "ci-temp-test-rt"

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_vpc_route_tables" "id_instances" {
  route_table_id = tencentcloud_route_table.route_table.id
}

data "tencentcloud_vpc_route_tables" "name_instances" {
  name = tencentcloud_route_table.route_table.name
}

data "tencentcloud_vpc_route_tables" "vpc_default_instance" {
  vpc_id           = tencentcloud_vpc.foo.id
  association_main = true
}

data "tencentcloud_vpc_route_tables" "tags_instances" {
  tags = tencentcloud_route_table.route_table.tags
}
```