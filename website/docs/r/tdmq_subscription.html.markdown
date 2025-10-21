---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_subscription"
sidebar_current: "docs-tencentcloud-resource-tdmq_subscription"
description: |-
  Provides a resource to create a tdmq subscription
---

# tencentcloud_tdmq_subscription

Provides a resource to create a tdmq subscription

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

resource "tencentcloud_tdmq_topic" "example" {
  cluster_id        = tencentcloud_tdmq_instance.example.id
  environ_id        = tencentcloud_tdmq_namespace.example.environ_name
  topic_name        = "tf-example-topic"
  partitions        = 1
  pulsar_topic_type = 3
  remark            = "remark."
}

resource "tencentcloud_tdmq_subscription" "example" {
  cluster_id               = tencentcloud_tdmq_instance.example.id
  environment_id           = tencentcloud_tdmq_namespace.example.environ_name
  topic_name               = tencentcloud_tdmq_topic.example.topic_name
  subscription_name        = "tf-example-subscription"
  remark                   = "remark."
  auto_create_policy_topic = true
  auto_delete_policy_topic = true
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Pulsar cluster ID.
* `environment_id` - (Required, String, ForceNew) Environment (namespace) name.
* `subscription_name` - (Required, String, ForceNew) Subscriber name, which can contain up to 128 characters.
* `topic_name` - (Required, String, ForceNew) Topic name.
* `auto_create_policy_topic` - (Optional, Bool, ForceNew) Whether to automatically create a dead letter topic and a retry letter topic. true: yes; false: no(default value).
* `auto_delete_policy_topic` - (Optional, Bool, ForceNew) Whether to automatically delete a dead letter topic and a retry letter topic. Setting is only allowed when `auto_create_policy_topic` is true. Default is false.
* `remark` - (Optional, String, ForceNew) Remarks (up to 128 characters).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdmq subscription can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_subscription.example pulsar-q4k5898krpqj#tf_example#tf-example-topic#tf-example-subscription#true
```

