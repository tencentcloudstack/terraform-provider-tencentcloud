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

```hcl
resource "tencentcloud_cls_alarm" "alarm" {
  name = "terraform-alarm-test"
  alarm_notice_ids = [
    "notice-0850756b-245d-4bc7-bb27-2a58fffc780b",
  ]
  alarm_period     = 15
  condition        = "test"
  message_template = "{{.Label}}"
  status           = true
  tags = {
    "createdBy" = "terraform"
  }
  trigger_count = 1

  alarm_targets {
    end_time_offset   = 0
    logset_id         = "33aaf0ae-6163-411b-a415-9f27450f68db"
    number            = 1
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    topic_id          = "88735a07-bea4-4985-8763-e9deb6da4fad"
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
}
```

## Argument Reference

The following arguments are supported:

* `alarm_notice_ids` - (Required, Set: [`String`]) list of alarm notice id.
* `alarm_period` - (Required, Int) alarm repeat cycle.
* `alarm_targets` - (Required, List) list of alarm target.
* `condition` - (Required, String) triggering conditions.
* `monitor_time` - (Required, List) monitor task execution time.
* `name` - (Required, String) log alarm name.
* `trigger_count` - (Required, Int) continuous cycle.
* `analysis` - (Optional, List) multidimensional analysis.
* `call_back` - (Optional, List) user define callback.
* `message_template` - (Optional, String) user define alarm notice.
* `status` - (Optional, Bool) whether to enable the alarm policy.
* `tags` - (Optional, Map) Tag description list.

The `alarm_targets` object supports the following:

* `end_time_offset` - (Required, Int) search end time of offset.
* `logset_id` - (Required, String) logset id.
* `number` - (Required, Int) the number of alarm object.
* `query` - (Required, String) query rules.
* `start_time_offset` - (Required, Int) search start time of offset.
* `topic_id` - (Required, String) topic id.

The `analysis` object supports the following:

* `content` - (Required, String) analysis content.
* `name` - (Required, String) analysis name.
* `type` - (Required, String) analysis type.
* `config_info` - (Optional, List) configuration.

The `call_back` object supports the following:

* `body` - (Required, String) callback body.
* `headers` - (Optional, Set) callback headers.

The `config_info` object supports the following:

* `key` - (Required, String) key.
* `value` - (Required, String) value.

The `monitor_time` object supports the following:

* `time` - (Required, Int) time period or point in time.
* `type` - (Required, String) Period for periodic execution, Fixed for regular execution.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls alarm can be imported using the id, e.g.

```
terraform import tencentcloud_cls_alarm.alarm alarm_id
```

