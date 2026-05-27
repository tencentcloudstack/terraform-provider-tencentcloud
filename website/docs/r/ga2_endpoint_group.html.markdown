---
subcategory: "Global Accelerator 2(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_endpoint_group"
sidebar_current: "docs-tencentcloud-resource-ga2_endpoint_group"
description: |-
  Provides a resource to create a GA2 (Global Accelerator 2) endpoint group.
---

# tencentcloud_ga2_endpoint_group

Provides a resource to create a GA2 (Global Accelerator 2) endpoint group.

## Example Usage

```hcl
resource "tencentcloud_ga2_endpoint_group" "example" {
  global_accelerator_id = "ga2-xxxxxxxx"
  listener_id           = "lis-xxxxxxxx"
  endpoint_group_type   = "DEFAULT"

  endpoint_group_configuration {
    name                  = "tf-example"
    endpoint_group_region = "ap-guangzhou"
    description           = "tf example endpoint group"
    enable_health_check   = true
    check_type            = "HTTP"
    check_port            = "80"
    check_path            = "/"
    check_method          = "GET"
    connect_timeout       = 5000
    health_check_interval = 30
    healthy_threshold     = 3
    unhealthy_threshold   = 3
    forward_protocol      = "HTTP"

    endpoint_configurations {
      endpoint_type    = "PublicIp"
      endpoint_service = "1.1.1.1"
      weight           = 10
    }

    endpoint_configurations {
      endpoint_type    = "Domain"
      endpoint_service = "example.com"
      weight           = 20
    }
  }
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
* `endpoint_type` - (Optional, String) Endpoint type. Valid values: `Domain`, `PublicIp`.
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
* `endpoint_configurations` - (Optional, List) Endpoint configurations under this group.
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

* `create` - (Defaults to `20m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

## Import

GA2 endpoint group can be imported using the composite ID `<global_accelerator_id>#<listener_id>#<endpoint_group_id>`, e.g.

```
terraform import tencentcloud_ga2_endpoint_group.example ga2-xxxxxxxx#lis-xxxxxxxx#eg-xxxxxxxx
```

