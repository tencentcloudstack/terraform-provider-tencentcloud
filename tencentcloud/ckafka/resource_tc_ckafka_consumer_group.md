Provides a resource to create a ckafka consumer_group

Example Usage

```hcl
resource "tencentcloud_ckafka_consumer_group" "consumer_group" {
  instance_id = "InstanceId"
  group_name = "GroupName"
  topic_name_list = ["xxxxxx"]
}
```

Import

ckafka consumer_group can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_consumer_group.consumer_group consumer_group_id
```