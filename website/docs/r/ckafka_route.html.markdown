---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_route"
sidebar_current: "docs-tencentcloud-resource-ckafka_route"
description: |-
  Provides a resource to create a ckafka route
---

# tencentcloud_ckafka_route

Provides a resource to create a ckafka route

## Example Usage

```hcl
resource "tencentcloud_ckafka_route" "example" {
  instance_id    = "ckafka-8j4ra928"
  vip_type       = 3
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-hhi88a58"
  access_type    = 0
  public_network = 3
  security_group_ids = [
    "sg-8ztqghor",
    "sg-rc9644tr",
    "sg-n8zf5ry9"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `vip_type` - (Required, Int) Routing network type (3:vpc routing; 4: standard support routing; 7: professional support routing).
* `access_type` - (Optional, Int) Access type. Valid values:
- 0: PLAINTEXT (in clear text, supported by both the old version and the community version without user information)
- 1: SASL_PLAINTEXT (in clear text, but at the beginning of the data, authentication will be logged in through SASL, which is only supported by the community version)
- 2: SSL (SSL encrypted communication without user information, supported by both older and community versions)
- 3: SASL_SSL (SSL encrypted communication. When the data starts, authentication will be logged in through SASL. Only the community version supports it).
* `auth_flag` - (Optional, Int) Auth flag.
* `caller_appid` - (Optional, Int) Caller appid.
* `ip` - (Optional, String) Ip.
* `public_network` - (Optional, Int) Public network.
* `security_group_ids` - (Optional, List: [`String`]) List of associated security group IDs in order.
* `subnet_id` - (Optional, String) Subnet id.
* `vpc_id` - (Optional, String) Vpc id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `broker_vip_list` - Virtual IP list (1 to 1 broker nodes).
  * `vip` - Virtual IP.
  * `vport` - Virtual port.
* `route_id` - Route ID.
* `vip_list` - Virtual IP list.
  * `vip` - Virtual IP.
  * `vport` - Virtual port.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateRoute](https://www.tencentcloud.com/document/api/597/74258)
  - [DescribeTaskStatus](https://www.tencentcloud.com/document/api/597/55881)
* Read:
  - [DescribeRoute](https://www.tencentcloud.com/document/api/597/37090)
  - [DescribeSecurityGroupRoutes](https://www.tencentcloud.com/document/api/597/74257)
* Update: [ModifyRouteSecurityGroups](https://www.tencentcloud.com/search/?wd=ModifyRouteSecurityGroups&type=document&lang=en)
* Delete: [DeleteRoute](https://www.tencentcloud.com/document/api/597/47140)

## Import

ckafka route can be imported using the instanceId#routeId, e.g.

```
terraform import tencentcloud_ckafka_route.example ckafka-8j4ra928#135912
```

