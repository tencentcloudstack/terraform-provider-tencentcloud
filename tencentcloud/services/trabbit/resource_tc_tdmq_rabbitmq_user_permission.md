Provides a resource to create a tdmq rabbitmq_user_permission

Example Usage

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
  instance_id    = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  user           = tencentcloud_tdmq_rabbitmq_user.example.user
  virtual_host   = tencentcloud_tdmq_rabbitmq_virtual_host.example.virtual_host
  config_regexp  = ".*"
  write_regexp   = ".*"
  read_regexp    = ".*"
}
```

Import

tdmq rabbitmq_user_permission can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rabbitmq_user_permission.example amqp-xxxxxxxx#user#vhost
```
