Provides a resource to create a GA2 (Global Accelerator 2) endpoint group.

Example Usage

If listener protocol is TCP

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

resource "tencentcloud_ga2_endpoint_group" "example" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  listener_id           = tencentcloud_ga2_listener.example.listener_id
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
```

If listener protocol is UDP

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
  name                  = "tf-example-udp"
  protocol              = "UDP"

  port_ranges {
    from_port = 70
    to_port   = 70
  }

  description     = "tf example listener"
  client_affinity = "Open"
  listener_type   = "Standard"
  idle_timeout    = 20
}


resource "tencentcloud_ga2_endpoint_group" "example" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  listener_id           = tencentcloud_ga2_listener.example.listener_id
  endpoint_group_type   = "DEFAULT"

  endpoint_group_configuration {
    name                  = "tf-example"
    endpoint_group_region = "ap-guangzhou"
    description           = "tf example endpoint group"
    enable_health_check   = false
    check_type            = "PING"
    connect_timeout       = 2
    health_check_interval = 30
    healthy_threshold     = 3
    unhealthy_threshold   = 3

    endpoint_configurations {
      endpoint_type    = "CustomDomain"
      endpoint_service = "example.com"
      weight           = 100
    }

    port_overrides {
      listener_port = 70
      endpoint_port = 9090
    }
  }
}
```

If listener protocol is HTTP

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

resource "tencentcloud_ga2_endpoint_group" "example" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  listener_id           = tencentcloud_ga2_listener.example.listener_id
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
}
```

If listener protocol is HTTPS

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
  name                  = "tf-example-https"
  protocol              = "HTTPS"

  port_ranges {
    from_port = 9090
    to_port   = 9090
  }

  listener_type           = "Standard"
  idle_timeout            = 38
  request_timeout         = 60
  x_forwarded_for_real_ip = true
  certification_type      = "MUTUAL"
  cipher_policy_id        = "tls_policy_1.2_strict-1.3"
  server_certificates     = ["ZDxux2I9"]
  client_ca_certificates  = ["UL45hi9B"]
  http_version            = "HTTP/2"
}

resource "tencentcloud_ga2_endpoint_group" "example" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  listener_id           = tencentcloud_ga2_listener.example.listener_id
  endpoint_group_type   = "VIRTUAL"

  endpoint_group_configuration {
    name                  = "tf-example"
    endpoint_group_region = "ap-guangzhou"
    enable_health_check   = true
    forward_protocol      = "HTTPS"
    check_type            = "HTTP"
    check_port            = "8080"
    check_domain          = "www.tencent.com"
    check_path            = "/path"
    check_method          = "GET"
    connect_timeout       = 2
    health_check_interval = 30
    healthy_threshold     = 3
    unhealthy_threshold   = 3
    cipher_policy_id      = "tls_policy_1.2_strict-1.3"
    http_version          = "HTTP/1.1"
    status_mask = [
      "http_2xx",
      "http_3xx",
      "http_4xx",
      "http_5xx"
    ]

    endpoint_configurations {
      endpoint_type    = "CustomPublicIp"
      endpoint_service = "12.15.13.16"
      weight           = 100
    }

    port_overrides {
      listener_port = 9090
      endpoint_port = 80
    }
  }
}
```

Import

GA2 endpoint group can be imported using the composite ID `<global_accelerator_id>#<listener_id>#<endpoint_group_id>`, e.g.

```
terraform import tencentcloud_ga2_endpoint_group.example ga-4mredmiu#lsr-1vd1fdwf#epg-h0ebutmo
```
