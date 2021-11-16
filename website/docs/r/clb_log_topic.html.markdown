---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_log_topic"
sidebar_current: "docs-tencentcloud-resource-clb_log_topic"
description: |-
  Provides a resource to create a CLB instance topic.
---

# tencentcloud_clb_log_topic

Provides a resource to create a CLB instance topic.

## Example Usage

```hcl
resource "tencentcloud_clb_log_topic" "topic" {
  log_set_id = "${tencentcloud_clb_log_set.set.id}"
  topic_name = "clb-topic"
}
```

## Argument Reference

The following arguments are supported:

* `log_set_id` - (Required, ForceNew) Log topic of CLB instance.
* `topic_name` - (Required, ForceNew) Log topic of CLB instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Log topic creation time.
* `status` - The status of log topic.


## Import

CLB log topic can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_log_topic.topic lb-7a0t6zqb

