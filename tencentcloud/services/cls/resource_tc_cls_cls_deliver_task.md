Provides a resource to create a CLS deliver task

Example Usage

```hcl
resource "tencentcloud_cls_logset" "source_logset" {
  logset_name = "tf-example-source-logset"
}

resource "tencentcloud_cls_topic" "source_topic" {
  topic_name   = "tf-example-source-topic"
  logset_id    = tencentcloud_cls_logset.source_logset.id
  period       = 10
  storage_type = "hot"
}

resource "tencentcloud_cls_logset" "target_logset" {
  logset_name = "tf-example-target-logset"
}

resource "tencentcloud_cls_topic" "target_topic" {
  topic_name   = "tf-example-target-topic"
  logset_id    = tencentcloud_cls_logset.target_logset.id
  period       = 10
  storage_type = "hot"
}

resource "tencentcloud_cls_cls_deliver_task" "deliver_task" {
  task_name   = "tf-example-deliver-task"
  compliance  = 1

  source_topic_config {
    topic_filter_type = 1
    logset_id         = tencentcloud_cls_logset.source_logset.id
    topics {
      topic_id = tencentcloud_cls_topic.source_topic.id
    }
  }

  target_topic_config {
    account_type = 1
    region       = "ap-guangzhou"
    logset_id    = tencentcloud_cls_logset.target_logset.id
    topic_id     = tencentcloud_cls_topic.target_topic.id
  }

  deliver_rule {
    data_scope = 3
  }
}
```

Import

CLS deliver task can be imported using the id, e.g.

```
terraform import tencentcloud_cls_cls_deliver_task.deliver_task deliver_task_id
```