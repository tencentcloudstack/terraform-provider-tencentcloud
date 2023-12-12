Provides a resource to create a cls topic.

Example Usage

```hcl
resource "tencentcloud_cls_topic" "topic" {
  topic_name           = "topic"
  logset_id            = "5cd3a17e-fb0b-418c-afd7-77b365397426"
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test",
  }
}
```

Import

cls topic can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_topic.topic 2f5764c1-c833-44c5-84c7-950979b2a278
```