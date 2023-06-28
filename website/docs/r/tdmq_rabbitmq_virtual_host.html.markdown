---
subcategory: "TDMQ for RabbitMQ(trabbit)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_virtual_host"
sidebar_current: "docs-tencentcloud-resource-tdmq_rabbitmq_virtual_host"
description: |-
  Provides a resource to create a tdmq rabbitmq_virtual_host
---

# tencentcloud_tdmq_rabbitmq_virtual_host

Provides a resource to create a tdmq rabbitmq_virtual_host

## Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_virtual_host" "rabbitmq_virtual_host" {
  instance_id  = "amqp-kzbe8p3n"
  virtual_host = "vh-test-1"
  description  = "desc"
  trace_flag   = false
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Cluster instance ID.
* `virtual_host` - (Required, String) vhost name.
* `description` - (Optional, String) describe.
* `trace_flag` - (Optional, Bool) Message track switch, true is on, false is off, default is off.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



