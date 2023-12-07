Provides a resource to create a eni_sg_attachment

-> **Note:** If this resource is used to bind security groups to eni, it cannot be linked to `tentcloud_eni` binding security group for simultaneous use.

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

resource "tencentcloud_security_group" "example1" {
  name        = "tf-example-sg1"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_security_group" "example2" {
  name        = "tf-example-sg2"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_eni" "example" {
  name        = "tf-example-eni"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni desc."
  ipv4_count  = 1
}

resource "tencentcloud_eni_sg_attachment" "eni_sg_attachment" {
  network_interface_ids = [tencentcloud_eni.example.id]
  security_group_ids    = [
    tencentcloud_security_group.example1.id,
    tencentcloud_security_group.example2.id
  ]
}
```

Import

vpc eni_sg_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_eni_sg_attachment.eni_sg_attachment eni_sg_attachment_id
```