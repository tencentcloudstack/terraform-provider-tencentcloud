---
subcategory: "TDMQ for RabbitMQ(RabbitMQ)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_exchange"
sidebar_current: "docs-tencentcloud-resource-tdmq_rabbitmq_exchange"
description: |-
  Provides a resource to create a tdmq rabbitmq_exchange
---

# tencentcloud_tdmq_rabbitmq_exchange

Provides a resource to create a tdmq rabbitmq_exchange

## Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_exchange" "rabbitmq_exchange" {
  exchange           = ""
  vhost_id           = ""
  type               = ""
  cluster_id         = ""
  remark             = ""
  alternate_exchange = ""
  delayed_type       = ""
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) cluster id.
* `exchange` - (Required, String) exchange name.
* `type` - (Required, String) exchange type.
* `vhost_id` - (Required, String) vhost.
* `alternate_exchange` - (Optional, String) alternate exchange name.
* `delayed_type` - (Optional, String) delayed exchange type, the value must be one of Direct, Fanout, Topic.
* `remark` - (Optional, String) exchange comment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdmq rabbitmq_exchange can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rabbitmq_exchange.rabbitmq_exchange rabbitmqExchange_id
```

