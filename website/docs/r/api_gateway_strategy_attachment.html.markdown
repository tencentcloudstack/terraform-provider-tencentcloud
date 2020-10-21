---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_strategy_attachment"
sidebar_current: "docs-tencentcloud-resource-api_gateway_strategy_attachment"
description: |-
  Use this resource to create IP strategy attachment of API gateway.
---

# tencentcloud_api_gateway_strategy_attachment

Use this resource to create IP strategy attachment of API gateway.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_ip_strategy" "test" {
  service_id    = tencentcloud_api_gateway_service.service.id
  strategy_name = "tf_test"
  strategy_type = "BLACK"
  strategy_data = "9.9.9.9"
}

resource "tencentcloud_api_gateway_api" "api" {
  service_id            = tencentcloud_api_gateway_service.service.id
  api_name              = "hello_update"
  api_desc              = "my hello api update"
  auth_type             = "SECRET"
  protocol              = "HTTP"
  enable_cors           = true
  request_config_path   = "/user/info"
  request_config_method = "POST"
  request_parameters {
    name          = "email"
    position      = "QUERY"
    type          = "string"
    desc          = "your email please?"
    default_value = "tom@qq.com"
    required      = true
  }
  service_config_type      = "HTTP"
  service_config_timeout   = 10
  service_config_url       = "http://www.tencent.com"
  service_config_path      = "/user"
  service_config_method    = "POST"
  response_type            = "XML"
  response_success_example = "<note>success</note>"
  response_fail_example    = "<note>fail</note>"
  response_error_codes {
    code           = 10
    msg            = "system error"
    desc           = "system error code"
    converted_code = -10
    need_convert   = true
  }
}

resource "tencentcloud_api_gateway_throttling_api" "foo" {
  service_id       = tencentcloud_api_gateway_service.service.id
  strategy         = "400"
  environment_name = "test"
  api_ids          = [tencentcloud_api_gateway_api.api.id]
}

resource "tencentcloud_api_gateway_service_release" "service" {
  service_id       = tencentcloud_api_gateway_throttling_api.foo.service_id
  environment_name = "release"
  release_desc     = "test service release"
}

resource "tencentcloud_api_gateway_strategy_attachment" "test" {
  service_id       = tencentcloud_api_gateway_service_release.service.service_id
  strategy_id      = tencentcloud_api_gateway_ip_strategy.test.strategy_id
  environment_name = "release"
  bind_api_id      = tencentcloud_api_gateway_api.api.id
}
```

## Argument Reference

The following arguments are supported:

* `bind_api_id` - (Required, ForceNew) The API that needs to be bound.
* `environment_name` - (Required, ForceNew) The environment of the strategy association. Valid values: `test`, `release`, `prepub`.
* `service_id` - (Required, ForceNew) The ID of the API gateway service.
* `strategy_id` - (Required, ForceNew) The ID of the API gateway strategy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

IP strategy attachment of API gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_strategy_attachment.test service-pk2r6bcc#IPStrategy-4kz2ljfi#api-h3wc5r0s#release
```

