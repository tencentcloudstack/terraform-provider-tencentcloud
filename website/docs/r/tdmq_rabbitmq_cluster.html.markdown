---
subcategory: "TDMQ for RabbitMQ(RabbitMQ)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_cluster"
sidebar_current: "docs-tencentcloud-resource-tdmq_rabbitmq_cluster"
description: |-
  Provides a resource to create a tdmq rabbitmq_cluster
---

# tencentcloud_tdmq_rabbitmq_cluster

Provides a resource to create a tdmq rabbitmq_cluster

## Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_cluster" "rabbitmq_cluster" {
  name   = ""
  remark = ""
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) cluster name.
* `remark` - (Optional, String) cluster description, 128 characters or less.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdmq rabbitmq_cluster can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rabbitmq_cluster.rabbitmq_cluster rabbitmqCluster_id
```

