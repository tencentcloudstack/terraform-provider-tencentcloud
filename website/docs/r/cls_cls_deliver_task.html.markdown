---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_cls_deliver_task"
sidebar_current: "docs-tencentcloud-resource-cls_cls_deliver_task"
description: |-
  Provides a resource to create a CLS deliver task
---

# tencentcloud_cls_cls_deliver_task

Provides a resource to create a CLS deliver task

## Example Usage

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
  task_name  = "tf-example-deliver-task"
  compliance = 1

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

## Argument Reference

The following arguments are supported:

* `compliance` - (Required, Int) Compliance. 1: agree cross-domain data transfer terms.
* `deliver_rule` - (Required, List) Deliver rule.
* `source_topic_config` - (Required, List) Source topic config.
* `target_topic_config` - (Required, List) Target topic config.
* `task_name` - (Required, String) Deliver task name.
* `enable` - (Optional, Int) Task status. 0: run, 1: pause.
* `has_services_log` - (Optional, Int) Whether to deliver service log. 1: off, 2: on.

The `deliver_rule` object supports the following:

* `data_scope` - (Required, Int) Data deliver scope. 1: history+new, 2: custom time range, 3: new only.

The `source_topic_config` object supports the following:

* `logset_id` - (Required, String) Source logset id.
* `topic_filter_type` - (Required, Int) Topic filter type. 1: static select.
* `topics` - (Optional, List) Source topic list.

The `target_topic_config` object supports the following:

* `account_type` - (Required, Int) Target account type. 1: current main account, 2: other main account.
* `logset_id` - (Required, String) Target logset id.
* `region` - (Required, String) Target region, e.g. ap-guangzhou.
* `topic_id` - (Required, String) Target topic id.
* `external_id` - (Optional, String) External ID, required when AccountType=2.
* `role_arn` - (Optional, String) Role ARN, required when AccountType=2.

The `topics` object of `source_topic_config` supports the following:

* `topic_id` - (Required, String) Log topic id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time (unix seconds).
* `progress` - Task progress percentage.
* `status` - Task status. 0: running, 1: paused, 2: completed, 3: abnormal.
* `task_id` - Deliver task id.
* `total` - Total count of matched tasks.
* `uin` - Main account id.
* `update_time` - Update time (unix seconds).


## Import

CLS deliver task can be imported using the id, e.g.

```
terraform import tencentcloud_cls_cls_deliver_task.deliver_task deliver_task_id
```

