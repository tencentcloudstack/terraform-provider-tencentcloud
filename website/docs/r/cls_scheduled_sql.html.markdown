---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_scheduled_sql"
sidebar_current: "docs-tencentcloud-resource-cls_scheduled_sql"
description: |-
  Provides a resource to create a CLS scheduled sql
---

# tencentcloud_cls_scheduled_sql

Provides a resource to create a CLS scheduled sql

## Example Usage

```hcl
resource "tencentcloud_cls_scheduled_sql" "example" {
  src_topic_id     = "2360e4a2-2b9e-4640-9829-72bc5052c9dd"
  src_topic_region = "ap-guangzhou"
  name             = "tf-example"
  enable_flag      = 1
  dst_resource {
    topic_id      = "3f5cc19d-f601-48b6-9671-9019b7b09bd0"
    region        = "ap-guangzhou"
    biz_type      = 1
    metric_names  = ["name1", "name2", "name3"]
    metric_labels = ["lable1", "lable2", "lable3"]
    custom_time   = "__WindowStartTime__"
    custom_metric_labels {
      key   = "cmlKey1"
      value = "cmlValue1"
    }

    custom_metric_labels {
      key   = "cmlKey2"
      value = "cmlValue2"
    }

    custom_metric_labels {
      key   = "cmlKey3"
      value = "cmlValue3"
    }
  }

  scheduled_sql_content = "verb:get AND responseStatus.code>=400\n| select stageTimestamp, responseStatus.code, objectRef.resource, objectRef.name, objectRef.namespace, user.username, \"annotations.authorization.k8s.io/decision\" as auth_decision\n  order by stageTimestamp desc\n  limit 50"
  process_start_time    = 1788417300000
  process_type          = 1
  process_period        = 1
  process_time_window   = "@m-1m,@m"
  process_delay         = 60
  syntax_rule           = 1
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

The `custom_metric_labels` object of `dst_resource` supports the following:

* `key` - (Required, String) custom metric label key.
* `value` - (Required, String) custom metric label value.

The `dst_resource` object supports the following:

* `topic_id` - (Required, String) dst topic id.
* `biz_type` - (Optional, Int) topic type.
* `custom_metric_labels` - (Optional, List) custom metric labels, used to add static dimensions to metrics.
* `custom_time` - (Optional, String) metric timestamp field, the default value is the left boundary time of the SQL query range.
* `metric_labels` - (Optional, List) metric dimensions, time type is not accepted.
* `metric_name` - (Optional, String) metric name.
* `metric_names` - (Optional, List) metric names, used when biz_type is 1 (metric topic) for multi-metric scenarios.
* `region` - (Optional, String) topic region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `task_id` - task id.


## Import

CLS scheduled sql can be imported using the id, e.g.

```
terraform import tencentcloud_cls_scheduled_sql.example aebbad1b-4228-4ccc-8d62-0333c9739452
```

