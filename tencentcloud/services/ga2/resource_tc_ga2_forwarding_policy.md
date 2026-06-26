Provides a resource to create a GA2 forwarding policy

Example Usage

```hcl
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                 = "tf-example"
  instance_charge_type = "POSTPAID"
  description          = "tf example global accelerator"

  tags = {
    createdBy = "Terraform"
  }
}

resource "tencentcloud_ga2_accelerate_area" "example" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  accelerate_region     = "ap-guangzhou"
  bandwidth             = 10
  isp_type              = "BGP"
  ip_version            = "IPv4"
}

resource "tencentcloud_ga2_listener" "example" {
  global_accelerator_id = tencentcloud_ga2_accelerate_area.example.global_accelerator_id
  name                  = "tf-example-http"
  protocol              = "HTTP"

  port_ranges {
    from_port = 90
    to_port   = 90
  }

  description             = "tf example listener"
  idle_timeout            = 15
  request_timeout         = 60
  listener_type           = "Standard"
  x_forwarded_for_real_ip = true
}

resource "tencentcloud_ga2_forwarding_policy" "example" {
  global_accelerator_id = tencentcloud_ga2_accelerate_area.example.global_accelerator_id
  listener_id           = tencentcloud_ga2_listener.example.listener_id
  host                  = "example.com"
}
```

Import

GA2 forwarding policy can be imported using the composite id `<global_accelerator_id>#<listener_id>#<forwarding_policy_id>`, e.g.

```
terraform import tencentcloud_ga2_forwarding_policy.example ga-jnyfyyss#lsr-hzht200v#dm-kvassops
```