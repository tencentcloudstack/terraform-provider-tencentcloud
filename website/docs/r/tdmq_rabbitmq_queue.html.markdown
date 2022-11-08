---
subcategory: "TDMQ for RabbitMQ(RabbitMQ)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_queue"
sidebar_current: "docs-tencentcloud-resource-tdmq_rabbitmq_queue"
description: |-
  Provides a resource to create a tdmq rabbitmq_queue
---

# tencentcloud_tdmq_rabbitmq_queue

Provides a resource to create a tdmq rabbitmq_queue

## Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_queue" "rabbitmq_queue" {
  queue                   = ""
  cluster_id              = ""
  vhost_id                = ""
  auto_delete             = ""
  remark                  = ""
  dead_letter_exchange    = ""
  dead_letter_routing_key = ""
}
```

## Argument Reference

The following arguments are supported:

* `auto_delete` - (Required, Bool) auto delete.
* `cluster_id` - (Required, String) cluster id.
* `queue` - (Required, String) queue name, 3~64 characters.
* `vhost_id` - (Required, String) vhost name.
* `dead_letter_exchange` - (Optional, String) dead letter exchange.
* `dead_letter_routing_key` - (Optional, String) dead letter routing key.
* `remark` - (Optional, String) queue description, 128 characters or less.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdmq rabbitmq_queue can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rabbitmq_queue.rabbitmq_queue rabbitmqQueue_id
```

