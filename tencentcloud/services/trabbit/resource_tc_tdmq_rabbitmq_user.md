Provides a resource to create a tdmq rabbitmq_user

Example Usage

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