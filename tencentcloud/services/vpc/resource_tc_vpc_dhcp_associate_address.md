Provides a resource to create a vpc dhcp associate address

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "vpc"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_vpc_dhcp_ip" "example" {
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id
  dhcp_ip_name = "tf-example"
}

resource "tencentcloud_eip" "eip" {
  name = "example-eip"
}

resource "tencentcloud_vpc_dhcp_associate_address" "example" {
  dhcp_ip_id = tencentcloud_vpc_dhcp_ip.example.id
  address_ip = tencentcloud_eip.eip.public_ip
}
```

Import

vpc dhcp_associate_address can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_dhcp_associate_address.dhcp_associate_address dhcp_associate_address_id
```
