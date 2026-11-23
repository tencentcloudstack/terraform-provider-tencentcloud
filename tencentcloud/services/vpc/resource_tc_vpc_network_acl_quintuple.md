Provides a resource to create a vpc network acl quintuple

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
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
}

resource "tencentcloud_vpc_acl" "example" {
  vpc_id  = tencentcloud_vpc.vpc.id
  name    = "tf-example"
  ingress = [
    "ACCEPT#192.168.1.0/24#800#TCP",
    "ACCEPT#192.168.1.0/24#800-900#TCP",
  ]
  egress = [
    "ACCEPT#192.168.1.0/24#800#TCP",
    "ACCEPT#192.168.1.0/24#800-900#TCP",
  ]
}

resource "tencentcloud_vpc_network_acl_quintuple" "network_acl_quintuple" {
  network_acl_id = tencentcloud_vpc_acl.example.id
  network_acl_quintuple_set {
    ingress {
      protocol                       = "TCP"
      description                    = "ingress desc."
      source_port                    = "80"
      source_cidr                    = "192.168.0.0/24"
      destination_port               = "8080"
      destination_cidr               = "192.168.0.0/24"
      action                         = "DROP"
      network_acl_quintuple_entry_id = "acli45-q1phngkz"
      priority                       = 1
      network_acl_direction          = "INGRESS"
    }
    egress {
      protocol                       = "TCP"
      description                    = "egress desc."
      source_port                    = "80"
      source_cidr                    = "192.168.0.0/24"
      destination_port               = "8080"
      destination_cidr               = "192.168.0.0/24"
      action                         = "DROP"
      network_acl_quintuple_entry_id = "acli45-q1phngkz"
      priority                       = 1
      network_acl_direction          = "EGRESS"
    }
  }
}
```

Import

vpc network_acl_quintuple can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_network_acl_quintuple.network_acl_quintuple network_acl_quintuple_id
```
