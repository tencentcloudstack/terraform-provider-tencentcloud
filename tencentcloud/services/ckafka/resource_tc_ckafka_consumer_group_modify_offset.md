Provides a resource to create a ckafka consumer_group_modify_offset

Example Usage

```hcl
resource "tencentcloud_ckafka_consumer_group_modify_offset" "consumer_group_modify_offset" {
  instance_id = "ckafka-xxxxxx"
  group = "xxxxxx"
  offset = 0
  strategy = 2
  topics = ["xxxxxx"]
}
```