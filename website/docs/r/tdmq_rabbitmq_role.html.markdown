---
subcategory: "TDMQ for RabbitMQ(RabbitMQ)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_role"
sidebar_current: "docs-tencentcloud-resource-tdmq_rabbitmq_role"
description: |-
  Provides a resource to create a tdmq rabbitmq_role
---

# tencentcloud_tdmq_rabbitmq_role

Provides a resource to create a tdmq rabbitmq_role

## Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_role" "rabbitmq_role" {
  role_name  = ""
  cluster_id = ""
  remark     = ""
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) cluster id.
* `role_name` - (Required, String) role name.
* `remark` - (Optional, String) role description, 128 characters or less.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdmq rabbitmq_role can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rabbitmq_role.rabbitmq_role rabbitmqRole_id
```

