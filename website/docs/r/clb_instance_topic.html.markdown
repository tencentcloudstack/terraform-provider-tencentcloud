---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instance_topic"
sidebar_current: "docs-tencentcloud-resource-clb_instance_topic"
description: |-
  Provides a resource to create a CLB instance topic.
---

# tencentcloud_clb_instance_topic

Provides a resource to create a CLB instance topic.

## Example Usage

```hcl
resource "tencentcloud_clb_instances_topic" "topic" {
  topic_name      = "clb-topic"
  partition_count = 3
}
```

## Argument Reference

The following arguments are supported:

* `topic_name` - (Required) Log topic of CLB instance.
* `partition_count` - (Optional) Topic partition count of CLB instance.(Default 1).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



