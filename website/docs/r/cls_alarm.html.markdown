---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_alarm"
sidebar_current: "docs-tencentcloud-resource-cls_alarm"
description: |-
  Provides a resource to create a cls alarm
---

# tencentcloud_cls_alarm

Provides a resource to create a cls alarm

## Example Usage

### Use single condition

```hcl
resource "tencentcloud_cls_alarm" "example" {
  name = "tf-example"
  alarm_notice_ids = [
    "notice-c2af43ee-1a4b-4c4a-ae3e-f81481280101",
  ]
  alarm_period     = 15
  condition        = "$1.source='10.0.0.1'"
  alarm_level      = 1
  message_template = "{{.Label}}"
  status           = true
  trigger_count    = 1

  alarm_targets {
    logset_id         = "e74efb8e-f647-48b2-a725-43f11b122081"
    topic_id          = "59cf3ec0-1612-4157-be3f-341b2e7a53cb"
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    end_time_offset   = 0
    number            = 1
    syntax_rule       = 1
  }

  analysis {
    content = "__FILENAME__"
    name    = "terraform"
    type    = "field"

    config_info {
      key   = "QueryIndex"
      value = "1"
    }
  }

  monitor_time {
    time = 1
    type = "Period"
  }

  tags = {
    createdBy = "terraform"
  }
}
```

### Use multi conditions

```hcl
resource "tencentcloud_cls_alarm" "example" {
  name = "tf-example"
  alarm_notice_ids = [
    "notice-c2af43ee-1a4b-4c4a-ae3e-f81481280101",
  ]
  alarm_period     = 15
  message_template = "{{.Label}}"
  status           = true
  trigger_count    = 1

  alarm_targets {
    logset_id         = "e74efb8e-f647-48b2-a725-43f11b122081"
    topic_id          = "59cf3ec0-1612-4157-be3f-341b2e7a53cb"
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    end_time_offset   = 0
    number            = 1
    syntax_rule       = 1
  }

  analysis {
    content = "__FILENAME__"
    name    = "terraform"
    type    = "field"

    config_info {
      key   = "QueryIndex"
      value = "1"
    }
  }

  multi_conditions {
    condition   = "[$1.__QUERYCOUNT__]> 0"
    alarm_level = 1
  }

  multi_conditions {
    condition   = "$1.source='10.0.0.1'"
    alarm_level = 2
  }

  monitor_time {
    time = 1
    type = "Period"
  }

  tags = {
    createdBy = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `alarm_notice_ids` - (Required, Set: [`String`]) list of alarm notice id.
* `alarm_period` - (Required, Int) alarm repeat cycle.
* `alarm_targets` - (Required, List) list of alarm target.
* `monitor_time` - (Required, List) monitor task execution time.
* `name` - (Required, String) log alarm name.
* `trigger_count` - (Required, Int) continuous cycle.
* `alarm_level` - (Optional, Int) Alarm level. 0: Warning; 1: Info; 2: Critical. Default is 0.
* `analysis` - (Optional, List) multidimensional analysis.
* `call_back` - (Optional, List) user define callback.
* `condition` - (Optional, String) Trigger condition.
* `message_template` - (Optional, String) user define alarm notice.
* `multi_conditions` - (Optional, List) Multiple triggering conditions.
* `status` - (Optional, Bool) whether to enable the alarm policy.
* `tags` - (Optional, Map) Tag description list.

The `alarm_targets` object supports the following:

* `end_time_offset` - (Required, Int) search end time of offset.
* `logset_id` - (Required, String) logset id.
* `number` - (Required, Int) the number of alarm object.
* `query` - (Required, String) query rules.
* `start_time_offset` - (Required, Int) search start time of offset.
* `topic_id` - (Required, String) topic id.
* `syntax_rule` - (Optional, Int) Retrieve grammar rules, 0: Lucene syntax, 1: CQL syntax, Default value is 0.

The `analysis` object supports the following:

* `content` - (Required, String) analysis content.
* `name` - (Required, String) analysis name.
* `type` - (Required, String) analysis type.
* `config_info` - (Optional, List) configuration.

The `call_back` object supports the following:

* `body` - (Required, String) callback body.
* `headers` - (Optional, Set) callback headers.

The `config_info` object of `analysis` supports the following:

* `key` - (Required, String) key.
* `value` - (Required, String) value.

The `monitor_time` object supports the following:

* `time` - (Required, Int) time period or point in time.
* `type` - (Required, String) Period for periodic execution, Fixed for regular execution.

The `multi_conditions` object supports the following:

* `alarm_level` - (Optional, Int) Alarm level. 0: Warning; 1: Info; 2: Critical. Default is 0.
* `condition` - (Optional, String) Trigger condition.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls alarm can be imported using the id, e.g.

```
terraform import tencentcloud_cls_alarm.example alarm-d8529662-e10f-440c-ba80-50f3dcf215a3
```

