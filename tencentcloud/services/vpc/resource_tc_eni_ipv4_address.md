Provides a resource to create a vpc eni ipv4 address

Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-6"
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_eni" "example" {
  name        = "tf-example"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni desc."
  ipv4_count  = 1
  security_groups = [
    tencentcloud_security_group.example.id,
  ]
}

resource "tencentcloud_eni_ipv4_address" "example" {
  network_interface_id               = tencentcloud_eni.example.id
  qos_level                          = "DEFAULT"
  secondary_private_ip_address_count = 3
}
```

Or

```hcl
resource "tencentcloud_eni_ipv4_address" "example" {
  network_interface_id = tencentcloud_eni.example.id
  private_ip_addresses {
    is_wan_ip_blocked  = false
    private_ip_address = "10.0.0.15"
    qos_level          = "DEFAULT"
  }

  private_ip_addresses {
    is_wan_ip_blocked  = false
    private_ip_address = "10.0.0.4"
    qos_level          = "DEFAULT"
  }
}
```

Import

vpc eni ipv4 address can be imported using the id, e.g.

```
terraform import tencentcloud_eni_ipv4_address.example eni-65369ozn
```