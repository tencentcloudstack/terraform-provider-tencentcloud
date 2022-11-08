---
subcategory: "TDMQ for RabbitMQ(RabbitMQ)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_route_relation_attachment"
sidebar_current: "docs-tencentcloud-resource-tdmq_rabbitmq_route_relation_attachment"
description: |-
  Provides a resource to create a tdmq rabbitmq_route_relation_attachment
---

# tencentcloud_tdmq_rabbitmq_route_relation_attachment

Provides a resource to create a tdmq rabbitmq_route_relation_attachment

## Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_route_relation_attachment" "rabbitmq_route_relation_attachment" {
  cluster_id      = ""
  vhost_id        = ""
  source_exchange = ""
  dest_type       = ""
  dest_value      = ""
  remark          = ""
  routing_key     = ""
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) cluster id.
* `dest_type` - (Required, String, ForceNew) destination type, the optional value is Queue or Exchange.
* `dest_value` - (Required, String, ForceNew) destination value.
* `source_exchange` - (Required, String, ForceNew) source exchange name.
* `vhost_id` - (Required, String, ForceNew) vhost id.
* `remark` - (Optional, String, ForceNew) route relation comment.
* `routing_key` - (Optional, String, ForceNew) route key, default value is `default`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdmq rabbitmq_route_relation_attachment can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rabbitmq_route_relation_attachment.rabbitmq_route_relation_attachment rabbitmqRouteRelationAttachment_id
```

