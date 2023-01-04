---
subcategory: "Private Link(PLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_end_point_service"
sidebar_current: "docs-tencentcloud-resource-vpc_end_point_service"
description: |-
  Provides a resource to create a vpc end_point_service
---

# tencentcloud_vpc_end_point_service

Provides a resource to create a vpc end_point_service

## Example Usage

```hcl
resource "tencentcloud_vpc_end_point_service" "end_point_service" {
  vpc_id                 = "vpc-391sv4w3"
  end_point_service_name = "terraform-endpoint-service"
  auto_accept_flag       = false
  service_instance_id    = "lb-o5f6x7ke"
  service_type           = "CLB"
}
```

## Argument Reference

The following arguments are supported:

* `auto_accept_flag` - (Required, Bool) Whether to automatically accept.
* `end_point_service_name` - (Required, String) Name of end point service.
* `service_instance_id` - (Required, String) Id of service instance, like lb-xxx.
* `vpc_id` - (Required, String) ID of vpc instance.
* `service_type` - (Optional, String) Type of service instance, like `CLB`, `CDB`, `CRS`, default is `CLB`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create Time.
* `end_point_count` - Count of end point.
* `service_owner` - APPID.
* `service_vip` - VIP of backend service.


## Import

vpc end_point_service can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_end_point_service.end_point_service end_point_service_id
```

