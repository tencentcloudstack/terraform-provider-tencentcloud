Provides a SQL Server instance resource to create read-only database instances.

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

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_instance" "example" {
  name              = "tf-example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  multi_zones       = true
  charge_type       = "POSTPAID_BY_HOUR"
  engine_version    = "2019"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  project_id        = 0
  memory            = 16
  storage           = 20
  security_groups   = [tencentcloud_security_group.security_group.id]
}

resource "tencentcloud_sqlserver_readonly_instance" "example" {
  name                = "tf_example"
  availability_zone   = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type         = "POSTPAID_BY_HOUR"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  memory              = 4
  storage             = 20
  master_instance_id  = tencentcloud_sqlserver_instance.example.id
  readonly_group_type = 1
  force_upgrade       = true
  tags = {
    CreateBy = "Terraform"
  }
}
```

Import

SQL Server readonly instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_readonly_instance.example mssqlro-3cdq7kx5
```