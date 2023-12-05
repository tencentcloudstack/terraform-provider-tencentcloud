Provides a resource to create a group of AS (Auto scaling) instances.

Example Usage

Create a basic Scaling Group

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
```

Create a complete Scaling Group

```hcl
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

resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name   = "tf-example"
  configuration_id     = tencentcloud_as_scaling_config.example.id
  max_size             = 1
  min_size             = 0
  vpc_id               = tencentcloud_vpc.vpc.id
  subnet_ids           = [tencentcloud_subnet.subnet.id]
  project_id           = 0
  default_cooldown     = 400
  desired_capacity     = 1
  termination_policies = ["NEWEST_INSTANCE"]
  retry_policy         = "INCREMENTAL_INTERVALS"

  forward_balancer_ids {
    load_balancer_id = tencentcloud_clb_instance.example.id
    listener_id      = tencentcloud_clb_listener.example.listener_id
    rule_id          = tencentcloud_clb_listener_rule.example.rule_id

    target_attribute {
      port   = 80
      weight = 90
    }
  }

  tags = {
    "createBy" = "tfExample"
  }
}
```

Import

AutoScaling Groups can be imported using the id, e.g.

```
$ terraform import tencentcloud_as_scaling_group.scaling_group asg-n32ymck2
```