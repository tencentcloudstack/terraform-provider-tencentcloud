Use this data source to query vpc subnets information.

Example Usage

Create subnet resource

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_availability_zones_by_product" "zones" {
  product = "vpc"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet1"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_subnet" "subnetCDC" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet2"
  cidr_block        = "10.0.0.0/16"
  cdc_id            = "cluster-lchwgxhs"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  is_multicast      = false
}
```

Query all subnets

```hcl
data "tencentcloud_vpc_subnets" "subnets" {}
```

Query subnets by filter

```hcl
data "tencentcloud_vpc_subnets" "by_vpc" {
  vpc_id = tencentcloud_vpc.vpc.id
}

data "tencentcloud_vpc_subnets" "by_subnet_id" {
  subnet_id = tencentcloud_subnet.subnet.id
}

data "tencentcloud_vpc_subnets" "by_name" {
  name = tencentcloud_subnet.subnet.name
}

data "tencentcloud_vpc_subnets" "by_tags" {
  tags = tencentcloud_subnet.subnet.tags
}

data "tencentcloud_vpc_subnets" "by_cdc_id" {
  cdc_id = tencentcloud_subnet.subnetCDC.cdc_id
}
```
