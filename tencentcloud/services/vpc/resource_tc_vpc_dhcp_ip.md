Provides a resource to create a vpc dhcp_ip

Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
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
```

Import

vpc dhcp_ip can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_dhcp_ip.dhcp_ip dhcp_ip_id
```