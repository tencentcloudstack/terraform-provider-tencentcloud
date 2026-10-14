---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_service"
sidebar_current: "docs-tencentcloud-resource-api_gateway_service"
description: |-
  Use this resource to create API gateway service.
---

# tencentcloud_api_gateway_service

Use this resource to create API gateway service.

~> **NOTE:** After setting `uniq_vpc_id`, it cannot be modified.

## Example Usage

### Shared Service

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "example-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_api_gateway_service" "example" {
  service_name = "tf-example"
  protocol     = "http&https"
  service_desc = "desc."
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
  uniq_vpc_id  = tencentcloud_vpc.vpc.id

  tags = {
    createdBy = "terraform"
  }

  release_limit = 500
  pre_limit     = 500
  test_limit    = 500
}
```

### Exclusive Service

```hcl
resource "tencentcloud_api_gateway_service" "example" {
  service_name = "tf-example"
  protocol     = "http&https"
  service_desc = "desc."
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
  uniq_vpc_id  = tencentcloud_vpc.vpc.id
  instance_id  = "instance-rc6fcv4e"

  tags = {
    createdBy = "terraform"
  }

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
* `instance_id` - (Optional, String, ForceNew) Exclusive instance ID.
* `ip_version` - (Optional, String, ForceNew) IP version number. Valid values: `IPv4`, `IPv6`. Default value: `IPv4`.
* `pre_limit` - (Optional, Int) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.
* `release_limit` - (Optional, Int) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.
* `service_desc` - (Optional, String) Custom service description.
* `tags` - (Optional, Map) Tag description list.
* `test_limit` - (Optional, Int) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.
* `uniq_vpc_id` - (Optional, String) VPC ID.

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

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateService](https://cloud.tencent.com/document/api/628/45201) - This API is used to create a service.
* Create: [DescribeService](https://cloud.tencent.com/document/api/628/45199) - This API is used to query the details of a service, such as its description, domain name, protocol, creation time, and releases.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Create: [ModifyServiceEnvironmentStrategy](https://cloud.tencent.com/document/api/628/45178) - This API is used to modify a service throttling policy.
* Read: [DescribeApiUsagePlan](https://cloud.tencent.com/document/api/628/45240) - This API is used to query the details of API usage plans in a service.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Read: [DescribeService](https://cloud.tencent.com/document/api/628/45199) - This API is used to query the details of a service, such as its description, domain name, protocol, creation time, and releases.
* Read: [DescribeServiceEnvironmentStrategy](https://cloud.tencent.com/document/api/628/45180) - This API is used to display a service throttling policy.
* Read: [DescribeServiceUsagePlan](https://cloud.tencent.com/document/api/628/45195) - This API is used to query the details of usage plans in a service.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [ModifyService](https://cloud.tencent.com/document/api/628/45193) - This API is used to modify the relevant information of a service. After a service is created, its name, description, and service type can be modified.
* Update: [ModifyServiceEnvironmentStrategy](https://cloud.tencent.com/document/api/628/45178) - This API is used to modify a service throttling policy.
* Delete: [DeleteService](https://cloud.tencent.com/document/api/628/45200) - This API is used to delete a service in API Gateway.
* Delete: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Delete: [UnReleaseService](https://cloud.tencent.com/document/api/628/45191) - This API is used to deactivate a service.

## Import

API gateway service can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_service.service service-pg6ud8pa
```

