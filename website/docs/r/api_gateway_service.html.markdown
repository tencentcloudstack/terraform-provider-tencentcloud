---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_service"
sidebar_current: "docs-tencentcloud-resource-api_gateway_service"
description: |-
  Use this resource to create API gateway service.
---

# tencentcloud_api_gateway_service

Use this resource to create API gateway service.

## Example Usage

Shared Service

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
  tags = {
    test-key1 = "test-value1"
    test-key2 = "test-value2"
  }
  release_limit = 500
  pre_limit     = 500
  test_limit    = 500
}
```

Exclusive Service

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "service"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
  tags = {
    test-key1 = "test-value1"
  }
  instance_id   = "instance-rc6fcv4e"
  release_limit = 500
  pre_limit     = 500
  test_limit    = 500
}
```

## Argument Reference

The following arguments are supported:

* `net_type` - (Required, Set: [`String`]) Network type list, which is used to specify the supported network types. Valid values: `INNER`, `OUTER`. `INNER` indicates access over private network, and `OUTER` indicates access over public network.
* `protocol` - (Required, String) Service frontend request type. Valid values: `http`, `https`, `http&https`.
* `service_name` - (Required, String) Custom service name.
* `exclusive_set_name` - (Optional, String, ForceNew, **Deprecated**) It has been deprecated from version 1.81.9. Self-deployed cluster name, which is used to specify the self-deployed cluster where the service is to be created.
* `instance_id` - (Optional, String) Exclusive instance ID.
* `ip_version` - (Optional, String, ForceNew) IP version number. Valid values: `IPv4`, `IPv6`. Default value: `IPv4`.
* `pre_limit` - (Optional, Int) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.
* `release_limit` - (Optional, Int) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.
* `service_desc` - (Optional, String) Custom service description.
* `tags` - (Optional, Map) Tag description list.
* `test_limit` - (Optional, Int) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_list` - A list of APIs.
  * `api_desc` - Description of the API.
  * `api_id` - ID of the API.
  * `api_name` - Name of the API.
  * `method` - Method of the API.
  * `path` - Path of the API.
* `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
* `inner_http_port` - Port number for http access over private network.
* `inner_https_port` - Port number for https access over private network.
* `internal_sub_domain` - Private network access subdomain name.
* `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
* `outer_sub_domain` - Public network access subdomain name.
* `usage_plan_list` - A list of attach usage plans.
  * `api_id` - ID of the API.
  * `bind_type` - Binding type.
  * `usage_plan_id` - ID of the usage plan.
  * `usage_plan_name` - Name of the usage plan.


## Import

API gateway service can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_service.service service-pg6ud8pa
```

