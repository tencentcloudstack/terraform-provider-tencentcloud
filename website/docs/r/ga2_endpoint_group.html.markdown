---
subcategory: "Global Accelerator(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_endpoint_group"
sidebar_current: "docs-tencentcloud-resource-ga2_endpoint_group"
description: |-
  Provides a resource to create a GA2 (Global Accelerator 2) endpoint group.
---

# tencentcloud_ga2_endpoint_group

Provides a resource to create a GA2 (Global Accelerator 2) endpoint group.

## Example Usage

### If enable_health_check is false

```hcl
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                 = "tf-example"
  instance_charge_type = "POSTPAID"
  description          = "tf example global accelerator"

  tags = {
    createdBy = "Terraform"
  }
}

resource "tencentcloud_ga2_listener" "example1" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
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
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
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

## Argument Reference

The following arguments are supported:

* `endpoint_group_configuration` - (Required, List) Endpoint group configuration.
* `endpoint_group_type` - (Required, String, ForceNew) Endpoint group type. Valid values: `VIRTUAL`, `DEFAULT`.
* `global_accelerator_id` - (Required, String, ForceNew) Global accelerator instance ID.
* `listener_id` - (Required, String, ForceNew) Listener ID.

The `endpoint_configurations` object of `endpoint_group_configuration` supports the following:

* `endpoint_service` - (Optional, String) Endpoint domain or IP.
* `endpoint_type` - (Optional, String) Endpoint type. Valid values: `CustomDomain`, `CustomPublicIp`.
* `weight` - (Optional, Int) Endpoint weight.

The `endpoint_group_configuration` object supports the following:

* `check_domain` - (Optional, String) Health check domain.
* `check_method` - (Optional, String) Health check request method.
* `check_path` - (Optional, String) Health check URL path.
* `check_port` - (Optional, String) Health check port.
* `check_recv_context` - (Optional, String) Health check expected response.
* `check_send_context` - (Optional, String) Health check request payload.
* `check_type` - (Optional, String) Health check protocol. Valid values: `TCP`, `HTTP`, `HTTPS`, `PING`, `CUSTOM`.
* `cipher_policy_id` - (Optional, String) HTTPS cipher policy ID.
* `connect_timeout` - (Optional, Int) Response timeout in milliseconds.
* `context_type` - (Optional, String) Health check content type.
* `description` - (Optional, String) Description. Maximum length is 100 bytes.
* `enable_health_check` - (Optional, Bool) Whether to enable health check.
* `endpoint_configurations` - (Optional, Set) Endpoint configurations under this group. This is an unordered set; element order in HCL has no semantic meaning.
* `endpoint_group_region` - (Optional, String) Region of the endpoint group.
* `forward_protocol` - (Optional, String) Forward protocol back to origin.
* `health_check_interval` - (Optional, Int) Health check interval in seconds.
* `healthy_threshold` - (Optional, Int) Healthy threshold count.
* `isp_type` - (Optional, String) ISP type.
* `name` - (Optional, String) Name. Maximum length is 60 bytes.
* `port_overrides` - (Optional, List) Port overrides for the endpoint group.
* `status_mask` - (Optional, List) Status code masks for health check.
* `unhealthy_threshold` - (Optional, Int) Unhealthy threshold count.

The `port_overrides` object of `endpoint_group_configuration` supports the following:

* `endpoint_port` - (Optional, Int) Endpoint port.
* `listener_port` - (Optional, Int) Listener port.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `endpoint_group_id` - Endpoint group instance ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `5m`) Used when creating the resource.
* `update` - (Defaults to `5m`) Used when updating the resource.
* `delete` - (Defaults to `5m`) Used when deleting the resource.

## Import

GA2 endpoint group can be imported using the composite ID `<global_accelerator_id>#<listener_id>#<endpoint_group_id>`, e.g.

```
terraform import tencentcloud_ga2_endpoint_group.example ga-4mredmiu#lsr-1vd1fdwf#epg-h0ebutmo
```

