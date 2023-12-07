Provides a resource to create a tdmq rabbitmq_virtual_host

Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_virtual_host" "rabbitmq_virtual_host" {
  instance_id  = "amqp-kzbe8p3n"
  virtual_host = "vh-test-1"
  description  = "desc"
  trace_flag   = false
}
```