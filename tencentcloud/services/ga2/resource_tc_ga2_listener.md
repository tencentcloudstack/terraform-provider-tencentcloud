Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) listener.

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

resource "tencentcloud_ga2_listener" "example1" {
  global_accelerator_id = tencentcloud_ga2_accelerate_area.example.global_accelerator_id
  name                  = "tf-example-tcp"
  protocol              = "TCP"

  port_ranges {
    from_port = 80
    to_port   = 80
  }

  description      = "tf example listener"
  get_real_ip_type = "ProxyProtocol"
  client_affinity  = "Open"
  listener_type    = "Standard"
  idle_timeout     = 900
}

resource "tencentcloud_ga2_listener" "example2" {
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

  depends_on = [tencentcloud_ga2_listener.example1]
}
```

Import

GA2 listener can be imported using the composite id `<global_accelerator_id>#<listener_id>`, e.g.

```
terraform import tencentcloud_ga2_listener.example ga-4mredmiu#lsr-llr0dng1
```
