Provides a resource to create a as load_balancer

~> **NOTE:** `load_balancer_ids` A list of traditional load balancer IDs, with a maximum of 20 traditional load balancers bound to each scaling group. Only one LoadBalancerIds and ForwardLoadBalancers can be specified simultaneously.
~> **NOTE:** `forward_load_balancers` List of application type load balancers, with a maximum of 100 bound application type load balancers for each scaling group. Only one LoadBalancerIds and ForwardLoadBalancers can be specified simultaneously.

Example Usage

If use `load_balancer_ids`

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "as"
}

data "tencentcloud_images" "image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
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

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "tf-example"
  image_id           = data.tencentcloud_images.image.images.0.image_id
  instance_types     = ["SA1.SMALL1", "SA2.SMALL1", "SA2.SMALL2", "SA2.SMALL4"]
  instance_name_settings {
    instance_name = "test-ins-name"
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
  clb_name     = "clb-example"
  project_id   = 0
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_clb_listener" "example" {
  clb_id        = tencentcloud_clb_instance.example.id
  listener_name = "listener-example"
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
  load_balancer_ids     = [tencentcloud_clb_instance.example.id]
}
```

If use `forward_load_balancers`

```hcl
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

as load_balancer can be imported using the id, e.g.

```
terraform import tencentcloud_as_load_balancer.load_balancer auto_scaling_group_id
```
