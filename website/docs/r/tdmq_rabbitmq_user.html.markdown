---
subcategory: "TDMQ for RabbitMQ(trabbit)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_user"
sidebar_current: "docs-tencentcloud-resource-tdmq_rabbitmq_user"
description: |-
  Provides a resource to create a tdmq rabbitmq_user
---

# tencentcloud_tdmq_rabbitmq_user

Provides a resource to create a tdmq rabbitmq_user

## Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_user" "rabbitmq_user" {
  instance_id     = "amqp-kzbe8p3n"
  user            = "keep-user"
  password        = "asdf1234"
  description     = "test user"
  tags            = ["management", "monitoring"]
  max_connections = 3
  max_channels    = 3
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Cluster instance ID.
* `password` - (Required, String) Password, used when logging in.
* `user` - (Required, String) Username, used when logging in.
* `description` - (Optional, String) Describe.
* `max_channels` - (Optional, Int) The maximum number of channels for this user, if not filled in, there is no limit.
* `max_connections` - (Optional, Int) The maximum number of connections for this user, if not filled in, there is no limit.
* `tags` - (Optional, List: [`String`]) User tag, used to determine the permission range for changing user access to RabbitMQ Management. Management: regular console user, monitoring: management console user, other values: non console user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



