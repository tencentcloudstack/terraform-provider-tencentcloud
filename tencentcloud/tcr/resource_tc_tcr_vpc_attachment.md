Use this resource to attach tcr instance with the vpc and subnet network.

Example Usage

Attach a tcr instance with vpc resource

```hcl
locals {
  vpc_id    = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
  tcr_id = tencentcloud_tcr_instance.example.id
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

data "tencentcloud_security_groups" "sg" {
  name = "default"
}

resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "basic"
  delete_bucket = true
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_vpc_attachment" "foo" {
  instance_id		= local.tcr_id
  vpc_id			= local.vpc_id
  subnet_id		 	= local.subnet_id
}
```

Import

tcr vpc attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_vpc_attachment.foo instance_id#vpc_id#subnet_id
```