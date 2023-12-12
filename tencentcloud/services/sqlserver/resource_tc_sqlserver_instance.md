Use this resource to create SQL Server instance

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_sqlserver_instance" "example" {
  name              = "tf-example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  project_id        = 0
  memory            = 16
  storage           = 100
}
```

Import

SQL Server instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_instance.example mssql-3cdq7kx5
```