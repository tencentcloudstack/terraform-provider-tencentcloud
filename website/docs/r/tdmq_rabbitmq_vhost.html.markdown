---
subcategory: "TDMQ for RabbitMQ(RabbitMQ)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_vhost"
sidebar_current: "docs-tencentcloud-resource-tdmq_rabbitmq_vhost"
description: |-
  Provides a resource to create a tdmq rabbitmq_vhost
---

# tencentcloud_tdmq_rabbitmq_vhost

Provides a resource to create a tdmq rabbitmq_vhost

## Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_vhost" "rabbitmq_vhost" {
  cluster_id = ""
  vhost_id   = ""
  msg_ttl    = ""
  remark     = ""
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) cluster id.
* `msg_ttl` - (Required, Int) retention time for unconsumed messages, the unit is ms, range is 60s~15days.
* `vhost_id` - (Required, String) vhost name, can only contain letters, numbers, '-' and '_'.
* `remark` - (Optional, String) cluster description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdmq rabbitmq_vhost can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rabbitmq_vhost.rabbitmq_vhost rabbitmqVhost_id
```

