---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_scheduled_sql"
sidebar_current: "docs-tencentcloud-resource-cls_scheduled_sql"
description: |-
  Provides a resource to create a cls scheduled_sql
---

# tencentcloud_cls_scheduled_sql

Provides a resource to create a cls scheduled_sql

## Example Usage

```hcl
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-example-logset"
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_cls_topic" "topic" {
  topic_name           = "tf-example-topic"
  logset_id            = tencentcloud_cls_logset.logset.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags = {
    "test" = "test",
  }
}
resource "tencentcloud_cls_scheduled_sql" "scheduled_sql" {
  src_topic_id = tencentcloud_cls_topic.topic.id
  name         = "tf-example-task"
  enable_flag  = 1
  dst_resource {
    topic_id    = tencentcloud_cls_topic.topic.id
    region      = "ap-guangzhou"
    biz_type    = 0
    metric_name = "test"

  }
  scheduled_sql_content = "xxx"
  process_start_time    = 1690515360000
  process_type          = 1
  process_period        = 10
  process_time_window   = "@m-15m,@m"
  process_delay         = 5
  src_topic_region      = "ap-guangzhou"
  process_end_time      = 1690515360000
  syntax_rule           = 0
}
```

## Argument Reference

The following arguments are supported:

* `dst_resource` - (Required, List) scheduled slq dst resource.
* `enable_flag` - (Required, Int) task enable flag.
* `name` - (Required, String) task name.
* `process_delay` - (Required, Int) process delay.
* `process_period` - (Required, Int) process period.
* `process_start_time` - (Required, Int) process start timestamp.
* `process_time_window` - (Required, String) process time window.
* `process_type` - (Required, Int) process type.
* `scheduled_sql_content` - (Required, String) scheduled sql content.
* `src_topic_id` - (Required, String) src topic id.
* `src_topic_region` - (Required, String) src topic region.
* `process_end_time` - (Optional, Int) process end timestamp.
* `syntax_rule` - (Optional, Int) syntax rule.

The `dst_resource` object supports the following:

* `topic_id` - (Required, String) dst topic id.
* `biz_type` - (Optional, Int) topic type.
* `metric_name` - (Optional, String) metric name.
* `region` - (Optional, String) topic region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls scheduled_sql can be imported using the id, e.g.

```
terraform import tencentcloud_cls_scheduled_sql.scheduled_sql scheduled_sql_id
```

