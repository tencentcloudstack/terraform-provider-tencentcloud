Provides a resource to create a GA2 (Global Accelerator 2) endpoint group.

Example Usage

If enable_health_check is false

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
  idle_timeout     = 800
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
  idle_timeout            = 30
  request_timeout         = 60
  listener_type           = "Standard"
  x_forwarded_for_real_ip = true

  depends_on = [tencentcloud_ga2_listener.example1]
}

resource "tencentcloud_ga2_endpoint_group" "example1" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  listener_id           = tencentcloud_ga2_listener.example1.listener_id
  endpoint_group_type   = "DEFAULT"

  endpoint_group_configuration {
    name                  = "tf-example"
    endpoint_group_region = "ap-guangzhou"
    description           = "tf example endpoint group"
    enable_health_check   = true
    check_type            = "TCP"
    connect_timeout       = 2
    health_check_interval = 30
    healthy_threshold     = 3
    unhealthy_threshold   = 3

    endpoint_configurations {
      endpoint_type    = "CustomPublicIp"
      endpoint_service = "1.1.1.1"
      weight           = 50
    }

    endpoint_configurations {
      endpoint_type    = "CustomDomain"
      endpoint_service = "example.com"
      weight           = 90
    }

    port_overrides {
      listener_port = 80
      endpoint_port = 90
    }
  }
}

resource "tencentcloud_ga2_endpoint_group" "example2" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  listener_id           = tencentcloud_ga2_listener.example2.listener_id
  endpoint_group_type   = "VIRTUAL"

  endpoint_group_configuration {
    name                  = "tf-example"
    endpoint_group_region = "ap-guangzhou"
    description           = "tf example endpoint group"
    enable_health_check   = true
    forward_protocol      = "HTTP"
    check_type            = "HTTP"
    check_domain          = "check.com"
    check_method          = "GET"
    check_path            = "/path"
    connect_timeout       = 2
    health_check_interval = 30
    healthy_threshold     = 3
    unhealthy_threshold   = 3
    status_mask = [
      "http_2xx",
      "http_3xx",
      "http_4xx"
    ]

    endpoint_configurations {
      endpoint_type    = "CustomPublicIp"
      endpoint_service = "1.1.1.1"
      weight           = 10
    }

    endpoint_configurations {
      endpoint_type    = "CustomDomain"
      endpoint_service = "example.com"
      weight           = 20
    }

    port_overrides {
      listener_port = 90
      endpoint_port = 9090
    }
  }

  depends_on = [tencentcloud_ga2_endpoint_group.example1]
}
```

Import

GA2 endpoint group can be imported using the composite ID `<global_accelerator_id>#<listener_id>#<endpoint_group_id>`, e.g.

```
terraform import tencentcloud_ga2_endpoint_group.example ga-4mredmiu#lsr-1vd1fdwf#epg-h0ebutmo
```
