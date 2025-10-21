---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_topic_with_full_id"
sidebar_current: "docs-tencentcloud-resource-tdmq_topic_with_full_id"
description: |-
  Provide a resource to create a TDMQ topic with full id.
---

# tencentcloud_tdmq_topic_with_full_id

Provide a resource to create a TDMQ topic with full id.

## Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_namespace" "example" {
  environ_name = "tf_example"
  msg_ttl      = 300
  cluster_id   = tencentcloud_tdmq_instance.example.id
  retention_policy {
    time_in_minutes = 60
    size_in_mb      = 10
  }
  remark = "remark."
}

resource "tencentcloud_tdmq_topic_with_full_id" "example" {
  environ_id        = tencentcloud_tdmq_namespace.example.environ_name
  cluster_id        = tencentcloud_tdmq_instance.example.id
  topic_name        = "tf-example-topic"
  partitions        = 6
  pulsar_topic_type = 3
  remark            = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) The Dedicated Cluster Id.
* `environ_id` - (Required, String, ForceNew) The name of tdmq namespace.
* `partitions` - (Required, Int) The partitions of topic.
* `topic_name` - (Required, String, ForceNew) The name of topic to be created.
* `pulsar_topic_type` - (Optional, Int) Pulsar Topic Type 0: Non-persistent non-partitioned 1: Non-persistent partitioned 2: Persistent non-partitioned 3: Persistent partitioned.
* `remark` - (Optional, String) Description of the namespace.
* `topic_type` - (Optional, Int, **Deprecated**) This input will be gradually discarded and can be switched to PulsarTopicType parameter 0: Normal message; 1: Global sequential messages; 2: Local sequential messages; 3: Retrying queue; 4: Dead letter queue. The type of topic.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of resource.


## Import

Tdmq Topic can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_topic_with_full_id.test ${cluster_id}#${environ_id}#${topic_name}
```

