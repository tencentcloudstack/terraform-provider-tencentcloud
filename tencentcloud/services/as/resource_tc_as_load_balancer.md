Provides a resource to create a as load balancer

~> **NOTE:** This resource must exclusive in one auto scaling group, do not declare additional rule resources of this auto scaling group elsewhere.

~> **NOTE:** If the `auto_scaling_group_id` field of this resource comes from the `tencentcloud_as_scaling_group` resource, then the `forward_balancer_ids` field of the `tencentcloud_as_scaling_group` resource cannot be set simultaneously with this resource, which may result in conflicts

~> **NOTE:** `forward_load_balancers` List of application type load balancers, with a maximum of 100 bound application type load balancers for each scaling group.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
  is_multicast      = false
}


resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "tf-example"
  image_id           = "img-eb30mz89"
  instance_types     = ["S6.MEDIUM4"]
  instance_name_settings {
    instance_name = "demo-ins-name"
  }
}

resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name = "tf-example"
  configuration_id   = tencentcloud_as_scaling_config.example.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_clb_instance" "example" {
  network_type = "INTERNAL"
  clb_name     = "tf-example"
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id
  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_clb_listener" "example" {
  clb_id        = tencentcloud_clb_instance.example.id
  listener_name = "tf-example"
  port          = 80
  protocol      = "HTTP"
}

resource "tencentcloud_clb_listener_rule" "example" {
  listener_id = tencentcloud_clb_listener.example.listener_id
  clb_id      = tencentcloud_clb_instance.example.id
  domain      = "foo.net"
  url         = "/bar"
}

resource "tencentcloud_as_load_balancer" "example" {
  auto_scaling_group_id = tencentcloud_as_scaling_group.example.id

  forward_load_balancers {
    load_balancer_id = tencentcloud_clb_instance.example.id
    listener_id      = tencentcloud_clb_listener.example.listener_id
    location_id      = tencentcloud_clb_listener_rule.example.rule_id

    target_attributes {
      port   = 8080
      weight = 20
    }
  }
}
```

Import

as load balancer can be imported using the id, e.g.

```
terraform import tencentcloud_as_load_balancer.example asg-bpp4uol2
```
