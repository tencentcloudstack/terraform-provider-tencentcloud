---
subcategory: "TDMQ for RabbitMQ(trabbit)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_user_permission"
sidebar_current: "docs-tencentcloud-resource-tdmq_rabbitmq_user_permission"
description: |-
  Provides a resource to create a tdmq rabbitmq_user_permission
---

# tencentcloud_tdmq_rabbitmq_user_permission

Provides a resource to create a tdmq rabbitmq_user_permission

## Example Usage

```hcl
# create rabbitmq instance
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [var.availability_zone]
  vpc_id                                = var.vpc_id
  subnet_id                             = var.subnet_id
  cluster_name                          = "tf-example-rabbitmq"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}

# create rabbitmq user
resource "tencentcloud_tdmq_rabbitmq_user" "example" {
  instance_id = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  user        = "tf-example-user"
  password    = "Password@123"
  description = "test user"
  tags        = ["management"]
}

# create virtual host
resource "tencentcloud_tdmq_rabbitmq_virtual_host" "example" {
  instance_id  = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  virtual_host = "tf-example-vhost"
  description  = "test virtual host"
  trace_flag   = false
}

# create user permission
resource "tencentcloud_tdmq_rabbitmq_user_permission" "example" {
  instance_id   = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  user          = tencentcloud_tdmq_rabbitmq_user.example.user
  virtual_host  = tencentcloud_tdmq_rabbitmq_virtual_host.example.virtual_host
  config_regexp = ".*"
  write_regexp  = ".*"
  read_regexp   = ".*"
}
```

## Argument Reference

The following arguments are supported:

* `config_regexp` - (Required, String) Configure permission regexp, controls which resources can be declared.
* `instance_id` - (Required, String, ForceNew) Cluster instance ID.
* `read_regexp` - (Required, String) Read permission regexp, controls which resources can be read.
* `user` - (Required, String, ForceNew) Username.
* `virtual_host` - (Required, String, ForceNew) VirtualHost name.
* `write_regexp` - (Required, String) Write permission regexp, controls which resources can be written.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdmq rabbitmq_user_permission can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rabbitmq_user_permission.example amqp-xxxxxxxx#user#vhost
```

