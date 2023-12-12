Provides a resource to create a mysql ro_instance_ip

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_mysql_ro_instance_ip" "example" {
  instance_id    = "cdbro-bdlvcfpj"
  uniq_subnet_id = tencentcloud_subnet.subnet.id
  uniq_vpc_id    = tencentcloud_vpc.vpc.id
}
```