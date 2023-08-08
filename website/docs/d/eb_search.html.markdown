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
data "tencentcloud_eb_search" "eb_search" {
  start_time   =
  end_time     =
  event_bus_id = ""
  group_field  = ""
  filter {
    key      = ""
    operator = ""
    value    = ""
    type     = ""
    filters {
      key      = ""
      operator = ""
      value    = ""
    }

  }
  order_fields =
  order_by     = ""
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, Int) end time.
* `event_bus_id` - (Required, String) event bus Id.
* `group_field` - (Required, String) aggregate field.
* `start_time` - (Required, Int) start time.
* `filter` - (Optional, List) filter criteria.
* `order_by` - (Optional, String) Sort by, asc from old to new, desc from new to old.
* `order_fields` - (Optional, Set: [`String`]) sort array.
* `result_output_file` - (Optional, String) Used to save results.

The `filter` object supports the following:

* `filters` - (Optional, List) LogFilters array.
* `key` - (Optional, String) filter field name.
* `operator` - (Optional, String) operator, congruent eq, not equal neq, similar like, exclude similar not like, less than lt, less than and equal to lte, greater than gt, greater than and equal to gte, in range range, not in range norange.
* `type` - (Optional, String) The logical relationship of the level filters, the value AND or OR.
* `value` - (Optional, String) Filter value, range operation needs to enter two values at the same time, separated by commas.

The `filters` object supports the following:

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


