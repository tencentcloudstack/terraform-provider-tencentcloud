Provides a resource to create a vpc bandwidth_package_attachment

Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
}

resource "tencentcloud_vpc_bandwidth_package" "example" {
  network_type           = "BGP"
  charge_type            = "TOP5_POSTPAID_BY_MONTH"
  bandwidth_package_name = "tf-example"
  tags                   = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_clb_instance" "example" {
  network_type = "INTERNAL"
  clb_name     = "tf-example"
  project_id   = 0
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id

  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_vpc_bandwidth_package_attachment" "attachment" {
  resource_id          = tencentcloud_clb_instance.example.id
  bandwidth_package_id = tencentcloud_vpc_bandwidth_package.example.id
  network_type         = "BGP"
  resource_type        = "LoadBalance"
}
```