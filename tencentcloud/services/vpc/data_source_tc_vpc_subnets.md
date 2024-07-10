Use this data source to query vpc subnets information.

Example Usage

Create subnet resource

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
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
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  is_multicast      = false
}
```

Query all subnets

```hcl
data "tencentcloud_vpc_subnets" "subnets" {}
```

Query subnets by filter

```hcl
data "tencentcloud_vpc_subnets" "subnets" {
  vpc_id = tencentcloud_vpc.vpc.id
}

data "tencentcloud_vpc_subnets" "subnets" {
  subnet_id = tencentcloud_subnet.subnet.id
}

data "tencentcloud_vpc_subnets" "subnets" {
  name = tencentcloud_subnet.subnet.name
}

data "tencentcloud_vpc_subnets" "subnets" {
  tags = tencentcloud_subnet.subnet.tags
}

data "tencentcloud_vpc_subnets" "subnets" {
  cdc_id = tencentcloud_subnet.subnetCDC.cdc_id
}
```
