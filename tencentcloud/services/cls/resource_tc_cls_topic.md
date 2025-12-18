Provides a resource to create a cls topic.

~> **NOTE:** Field `encryption` can only be enabled, not disabled.

Example Usage

Create a standard cls topic

```hcl
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf_example"
  tags = {
    tagKey = "tagValue"
  }
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf_example"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 30
  storage_type         = "hot"
  describes            = "Test Demo."
  hot_period           = 10
  tags = {
    tagKey = "tagValue"
  }
}
```

Create a cls topic with web tracking

```hcl
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf_example"
  tags = {
    tagKey = "tagValue"
  }
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf_example"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 30
  storage_type         = "hot"
  describes            = "Test Demo."
  hot_period           = 10
  is_web_tracking      = true

  extends {
    anonymous_access {
      operations = ["trackLog", "realtimeProducer"]
      conditions {
        attributes      = "VpcID"
        rule            = 1
        condition_value = "vpc-ahr3xajx"
      }
    }
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

Import

cls topic can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_topic.example 2f5764c1-c833-44c5-84c7-950979b2a278
```