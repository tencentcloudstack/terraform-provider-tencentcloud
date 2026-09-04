---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_log_topic"
sidebar_current: "docs-tencentcloud-resource-clb_log_topic"
description: |-
  Provides a resource to create a CLB log topic.
---

# tencentcloud_clb_log_topic

Provides a resource to create a CLB log topic.

## Example Usage

```hcl
resource "tencentcloud_clb_log_topic" "example" {
  log_set_id = "2ed70190-bf06-4777-980d-2d8a327a2554"
  topic_name = "tf-example"
  status     = true
  period     = 30
  tags = {
    env = "prod"
  }
}
```

## Argument Reference

The following arguments are supported:

* `log_set_id` - (Required, String, ForceNew) Log topic of CLB instance.
* `topic_name` - (Required, String, ForceNew) Log topic of CLB instance.
* `period` - (Optional, Int) Log storage lifecycle in days. Standard storage supports 1-3600; 3640 means permanent retention. Defaults to 30 when unset.
* `status` - (Optional, Bool) The status of log topic. true: enable; false: disable. Default is true.
* `tags` - (Optional, Map) Tags of clb log topic.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Log topic creation time.


## Import

CLB log topic can be imported using the id, e.g.

```
terraform import tencentcloud_clb_log_topic.example be1a83dd-04b4-4807-89bf-8daddce0df71
```

