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

resource "tencentcloud_ga2_listener" "example3" {
  global_accelerator_id = tencentcloud_ga2_accelerate_area.example.global_accelerator_id
  name                  = "tf-example-https"
  protocol              = "HTTPS"

  port_ranges {
    from_port = 9090
    to_port   = 9090
  }

  listener_type       = "Standard"
  idle_timeout        = 60
  request_timeout     = 60
  certification_type  = "SVR"
  cipher_policy_id    = "tls_policy_1.2_strict-1.3"
  server_certificates = ["Yj6CmODs"]

  depends_on = [tencentcloud_ga2_listener.example2]
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

resource "tencentcloud_ga2_endpoint_group" "example3" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  listener_id           = tencentcloud_ga2_listener.example3.listener_id
  endpoint_group_type   = "DEFAULT"

  endpoint_group_configuration {
    name                  = "tf-example"
    endpoint_group_region = "ap-guangzhou"
    enable_health_check   = true
    forward_protocol      = "HTTPS"
    check_type            = "HTTP"
    check_port            = "8080"
    check_domain          = "www.tencent.com"
    check_path            = "/test"
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

  depends_on = [tencentcloud_ga2_endpoint_group.example2]
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_group_configuration` - (Required, List) Endpoint group configuration.
* `endpoint_group_type` - (Required, String, ForceNew) Endpoint group type. Valid values: `VIRTUAL`, `DEFAULT`.
* `global_accelerator_id` - (Required, String, ForceNew) Global accelerator instance ID.
* `listener_id` - (Required, String, ForceNew) Listener ID.

The `endpoint_configurations` object of `endpoint_group_configuration` supports the following:

* `endpoint_service` - (Optional, String) Endpoint domain name or IP address.
* `endpoint_type` - (Optional, String) Endpoint type. Valid values: `CustomDomain`, `CustomPublicIp`.
* `weight` - (Optional, Int) Endpoint weight.

The `endpoint_group_configuration` object supports the following:

* `check_domain` - (Optional, String) Health check domain. Length must be between 3 and 80 bytes. Required when `check_type` is `HTTP`.
* `check_method` - (Optional, String) Health check request method. Valid values: `GET`, `HEAD`. Required when `check_type` is `HTTP`.
* `check_path` - (Optional, String) Health check URL path. Must match the regular expression `^[a-zA-Z0-9_.-/]{1,80}$`. Required when `check_type` is `HTTP`.
* `check_port` - (Optional, String) Health check port. Valid range: [1, 65535]. Required when `check_type` is `CUSTOM`.
* `check_recv_context` - (Optional, String) Expected health check response content. Length must be between 1 and 500 bytes. Required when `check_type` is `CUSTOM`.
* `check_send_context` - (Optional, String) Health check request payload. Length must be between 1 and 500 bytes. Required when `check_type` is `CUSTOM`.
* `check_type` - (Optional, String) Health check protocol. Valid values: `TCP` (only when the endpoint group's listener protocol is TCP), `HTTP` (only when the listener protocol is HTTP/HTTPS), `PING` (only when the listener protocol is UDP), `CUSTOM` (only when the listener protocol is TCP/UDP). Required when `enable_health_check` is `true`.
* `cipher_policy_id` - (Optional, String) HTTPS cipher suite policy. Valid values: `tls_policy_1.0-2`, `tls_policy_1.1-2`, `tls_policy_1.2`, `tls_policy_1.2_strict`, `tls_policy_1.2_strict-1.3`. Required when `forward_protocol` is `HTTPS`.
* `connect_timeout` - (Optional, Int) Response timeout in seconds. Valid range: [1, 100]. Default: `2`. Required when `enable_health_check` is `true`.
* `context_type` - (Optional, String) Health check content type. Valid values: `TEXT` (plain text content). Required when `check_type` is `CUSTOM`.
* `description` - (Optional, String) Description. Default is empty (no description configured). Maximum length is 100 bytes.
* `enable_health_check` - (Optional, Bool) Whether to enable health check. Default: `false`.
* `endpoint_configurations` - (Optional, Set) Endpoint configurations under this group. This is an unordered set; element order in HCL has no semantic meaning.
* `endpoint_group_region` - (Optional, String) Region where the endpoint group resides.
* `forward_protocol` - (Optional, String) Protocol used to forward traffic to the origin. Valid values: `HTTP` (available when the endpoint group's listener protocol is HTTP/HTTPS), `HTTPS` (available when the listener protocol is HTTPS). Required when the listener protocol is HTTP or HTTPS.
* `health_check_interval` - (Optional, Int) Health check interval in seconds. Valid range: [5, 300]. Default: `30`. Required when `enable_health_check` is `true`.
* `healthy_threshold` - (Optional, Int) Healthy threshold count. Valid range: [1, 10]. Default: `3`. Required when `enable_health_check` is `true`.
* `http_version` - (Optional, String) Origin-pull protocol. Supports configuration of 'HTTP/1.1' or 'HTTP/2'. Enum values: HTTP/1.1: HTTP/1.1 version, HTTP/2: HTTP/2 version. This field is mandatory when the origin-pull protocol is HTTPS.
* `isp_type` - (Optional, String) ISP type. Valid values: `CMCC` (China Mobile), `CUCC` (China Unicom), `CTCC` (China Telecom). Required when the endpoint group region is a multi-ISP (three-network) region.
* `name` - (Optional, String) Endpoint group name. Maximum length is 128 bytes. Must start with a letter (a-z, A-Z) or a Chinese character.
* `port_overrides` - (Optional, List) Port mapping rules for the endpoint group. Layer-7 endpoint groups support at most 1 port mapping; layer-4 endpoint groups support at most 30 port mappings.
* `status_mask` - (Optional, List) Health check status code masks. Valid values: `http_2xx`, `http_3xx`, `http_4xx`, `http_5xx`. Required when `check_type` is `HTTP`.
* `unhealthy_threshold` - (Optional, Int) Unhealthy threshold count. Valid range: [1, 10]. Default: `3`. Required when `enable_health_check` is `true`.

The `port_overrides` object of `endpoint_group_configuration` supports the following:

* `endpoint_port` - (Optional, Int) Mapped endpoint port.
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

