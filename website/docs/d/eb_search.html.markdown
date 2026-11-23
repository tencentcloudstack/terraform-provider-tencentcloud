---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_search"
sidebar_current: "docs-tencentcloud-datasource-eb_search"
description: |-
  Use this data source to query detailed information of eb eb_search
---

# tencentcloud_eb_search

Use this data source to query detailed information of eb eb_search

## Example Usage

```hcl
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_put_events" "put_events" {
  event_list {
    source = "ckafka.cloud.tencent"
    data = jsonencode(
      {
        "topic" : "test-topic",
        "Partition" : 1,
        "offset" : 37,
        "msgKey" : "test",
        "msgBody" : "Hello from Ckafka again!"
      }
    )
    type    = "connector:ckafka"
    subject = "qcs::ckafka:ap-guangzhou:uin/1250000000:ckafkaId/uin/1250000000/ckafka-123456"
    time    = 1691572461939

  }
  event_bus_id = tencentcloud_eb_event_bus.foo.id
}

data "tencentcloud_eb_search" "eb_search" {
  start_time   = 1691637288422
  end_time     = 1691648088422
  event_bus_id = "eb-jzytzr4e"
  group_field  = "RuleIds"
  filter {
    type = "OR"
    filters {
      key      = "status"
      operator = "eq"
      value    = "1"
    }
  }

  filter {
    type = "OR"
    filters {
      key      = "type"
      operator = "eq"
      value    = "connector:ckafka"
    }
  }
  # order_fields = [""]
  order_by = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, Int) end time.
* `event_bus_id` - (Required, String) event bus Id.
* `start_time` - (Required, Int) start time.
* `filter` - (Optional, List) filter criteria.
* `group_field` - (Optional, String) aggregate field, When querying the log index dimension value, you must enter.
* `order_by` - (Optional, String) Sort by, asc from old to new, desc from new to old, take effect when the log is retrieved.
* `order_fields` - (Optional, Set: [`String`]) sort array, take effect when the log is retrieved.
* `result_output_file` - (Optional, String) Used to save results.

The `filter` object supports the following:

* `filters` - (Optional, List) LogFilters array.
* `key` - (Optional, String) filter field name.
* `operator` - (Optional, String) operator, congruent eq, not equal neq, similar like, exclude similar not like, less than lt, less than and equal to lte, greater than gt, greater than and equal to gte, in range range, not in range norange.
* `type` - (Optional, String) The logical relationship of the level filters, the value AND or OR.
* `value` - (Optional, String) Filter value, range operation needs to enter two values at the same time, separated by commas.

The `filters` object of `filter` supports the following:

* `key` - (Required, String) filter field name.
* `operator` - (Required, String) operator, congruent eq, not equal neq, similar like, exclude similar not like, less than lt, less than and equal to lte, greater than gt, greater than and equal to gte, within range range, not within range norange.
* `value` - (Required, String) Filter values, range operations need to enter two values at the same time, separated by commas.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dimension_values` - Index retrieves dimension values.
* `results` - Log search results, note: this field may return null, indicating that no valid value can be obtained.
  * `message` - Log content details, note: this field may return null, indicating that no valid value can be obtained.
  * `region` - Region, Note: This field may return null, indicating that no valid value can be obtained.
  * `rule_ids` - Event matching rules, note: this field may return null, indicating that no valid value can be obtained.
  * `source` - Event source, note: this field may return null, indicating that no valid value can be obtained.
  * `status` - Event status, note: this field may return null, indicating that no valid value can be obtained.
  * `subject` - Instance ID, note: this field may return null, indicating that no valid value can be obtained.
  * `timestamp` - The reporting time of a single log, note: this field may return null, indicating that no valid value can be obtained.
  * `type` - Event type, note: this field may return null, indicating that no valid value can be obtained.

## API Reference

This data source queries the following TencentCloud APIs:

* Read:
  - [DescribeLogTagValue](https://www.tencentcloud.com/document/api/1108/57987)
  - [SearchLog](https://www.tencentcloud.com/document/api/1108/57986)

