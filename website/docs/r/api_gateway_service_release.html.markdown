---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_service_release"
sidebar_current: "docs-tencentcloud-resource-api_gateway_service_release"
description: |-
  Use this resource to create API gateway service release.
---

# tencentcloud_api_gateway_service_release

Use this resource to create API gateway service release.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "myservice"
  protocol     = "http"
  service_desc = "my nice service"
  net_type     = ["INNER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
  service_id            = tencentcloud_api_gateway_service.service.id
  api_name              = "tf_example"
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

resource "tencentcloud_api_gateway_service_release" "service" {
  service_id       = tencentcloud_api_gateway_api.api.service.id
  environment_name = "release"
  release_desc     = "test service release"
}
```

## Argument Reference

The following arguments are supported:

* `environment_name` - (Required, String, ForceNew) API gateway service environment name to be released. Valid values: `test`, `prepub`, `release`.
* `release_desc` - (Required, String, ForceNew) This release description of the API gateway service.
* `service_id` - (Required, String, ForceNew) ID of API gateway service.
* `release_version` - (Optional, String) The release version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

API gateway service release can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_service_release.service service-jjt3fs3s#release#20201015121916d85fb161-eaec-4dda-a7e0-659aa5f401be
```

