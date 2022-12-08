---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_topic"
sidebar_current: "docs-tencentcloud-resource-tdmq_topic"
description: |-
  Provide a resource to create a TDMQ topic.
---

# tencentcloud_tdmq_topic

Provide a resource to create a TDMQ topic.

## Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "foo" {
  cluster_name = "example"
  remark       = "this is description."
}

resource "tencentcloud_tdmq_namespace" "bar" {
  environ_name = "example"
  msg_ttl      = 300
  cluster_id   = "${tencentcloud_tdmq_instance.foo.id}"
  remark       = "this is description."
}

resource "tencentcloud_tdmq_topic" "bar" {
  environ_id = "${tencentcloud_tdmq_namespace.bar.id}"
  topic_name = "example"
  partitions = 6
  topic_type = 0
  cluster_id = "${tencentcloud_tdmq_instance.foo.id}"
  remark     = "this is description."
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) The Dedicated Cluster Id.
* `environ_id` - (Required, String, ForceNew) The name of tdmq namespace.
* `partitions` - (Required, Int) The partitions of topic.
* `topic_name` - (Required, String, ForceNew) The name of topic to be created.
* `topic_type` - (Required, Int, ForceNew) The type of topic.
* `remark` - (Optional, String) Description of the namespace.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of resource.


## Import

Tdmq Topic can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_topic.test topic_id
```

