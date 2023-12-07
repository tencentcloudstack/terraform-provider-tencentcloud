---
subcategory: "TDMQ for CMQ(tcmq)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcmq_topic"
sidebar_current: "docs-tencentcloud-resource-tcmq_topic"
description: |-
  Provides a resource to create a tcmq topic
---

# tencentcloud_tcmq_topic

Provides a resource to create a tcmq topic

## Example Usage

```hcl
resource "tencentcloud_tcmq_topic" "topic" {
  topic_name = "topic_name"
}
```

## Argument Reference

The following arguments are supported:

* `topic_name` - (Required, String) Topic name, which must be unique in the same topic under the same account in the same region. It can contain up to 64 letters, digits, and hyphens and must begin with a letter.
* `filter_type` - (Optional, Int) Used to specify the message match policy for the topic. `1`: tag match policy (default value); `2`: routing match policy.
* `max_msg_size` - (Optional, Int) Maximum message length. Value range: 1024-65536 bytes (i.e., 1-64 KB). Default value: 65536.
* `msg_retention_seconds` - (Optional, Int) Message retention period. Value range: 60-86400 seconds (i.e., 1 minute-1 day). Default value: 86400.
* `trace` - (Optional, Bool) Whether to enable message trace. true: yes; false: no. If this field is left empty, the feature will not be enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tcmq topic can be imported using the id, e.g.

```
terraform import tencentcloud_tcmq_topic.topic topic_id
```

