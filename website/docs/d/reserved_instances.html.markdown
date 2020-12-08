---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_reserved_instances"
sidebar_current: "docs-tencentcloud-datasource-reserved_instances"
description: |-
  Use this data source to query reserved instances.
---

# tencentcloud_reserved_instances

Use this data source to query reserved instances.

## Example Usage

```hcl
data "tencentcloud_reserved_instances" "instances" {
  availability_zone = "na-siliconvalley-1"
  instance_type     = "S2.MEDIUM8"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The available zone that the reserved instance locates at.
* `instance_type` - (Optional) The type of reserved instance.
* `reserved_instance_id` - (Optional) ID of the reserved instance to be query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `reserved_instance_list` - An information list of reserved instance. Each element contains the following attributes:
  * `availability_zone` - Availability zone of the reserved instance.
  * `end_time` - Expiry time of the reserved instance.
  * `instance_count` - Number of reserved instance.
  * `instance_type` - The type of reserved instance.
  * `reserved_instance_id` - ID of the reserved instance.
  * `start_time` - Start time of the reserved instance.
  * `status` - Status of the reserved instance.


