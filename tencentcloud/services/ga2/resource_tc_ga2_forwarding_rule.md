Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) layer-7 forwarding rule.

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

resource "tencentcloud_ga2_forwarding_policy" "example" {
  global_accelerator_id = tencentcloud_ga2_accelerate_area.example.global_accelerator_id
  listener_id           = tencentcloud_ga2_listener.example.listener_id
  host                  = "example.com"
}

resource "tencentcloud_ga2_forwarding_rule" "example" {
  global_accelerator_id = tencentcloud_ga2_accelerate_area.example.global_accelerator_id
  listener_id           = tencentcloud_ga2_listener.example.listener_id
  forwarding_policy_id  = tencentcloud_ga2_forwarding_policy.example.forwarding_policy_id

  rule_conditions {
    rule_condition_type  = "Path"
    rule_condition_value = ["/path"]
  }

  rule_actions {
    rule_action_type  = "ForwardGroup"
    rule_action_value = tencentcloud_ga2_endpoint_group.example.endpoint_group_id
  }
}
```

with origin headers

```hcl
resource "tencentcloud_ga2_forwarding_rule" "example" {
  global_accelerator_id = "ga-fhhs8w84"
  listener_id           = "lsr-dyy8jhzp"
  forwarding_policy_id  = "dm-rjssxr8k"
  origin_host           = "2.2.2.2"

  rule_conditions {
    rule_condition_type  = "Path"
    rule_condition_value = ["/path"]
  }

  rule_actions {
    rule_action_type  = "ForwardGroup"
    rule_action_value = "epg-nt4iwozo"
  }

  origin_headers {
    key   = "key"
    value = "value"
  }
}
```

Import

GA2 forwarding rule can be imported using the composite id `<global_accelerator_id>#<listener_id>#<forwarding_policy_id>#<forwarding_rule_id>`, e.g.

```
terraform import tencentcloud_ga2_forwarding_rule.example ga-fhhs8w84#lsr-dyy8jhzp#dm-rjssxr8k#rule-757r3bk2
```
