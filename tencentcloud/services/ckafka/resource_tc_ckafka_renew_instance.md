Provides a resource to create a ckafka renew_instance

Example Usage

```hcl
resource "tencentcloud_ckafka_renew_instance" "renew_ckafka_instance" {
  instance_id = "InstanceId"
  time_span = 1
}
```