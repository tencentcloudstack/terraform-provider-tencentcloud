---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_platform_event_patterns"
sidebar_current: "docs-tencentcloud-datasource-eb_platform_event_patterns"
description: |-
  Use this data source to query detailed information of eb platform_event_patterns
---

# tencentcloud_eb_platform_event_patterns

Use this data source to query detailed information of eb platform_event_patterns

## Example Usage

```hcl
data "tencentcloud_eb_platform_event_patterns" "platform_event_patterns" {
  product_type = ""
}
```

## Argument Reference

The following arguments are supported:

* `product_type` - (Required, String) Platform product type.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `event_patterns` - Platform product event matching rules.
  * `event_name` - Platform event name.Note: This field may return null, indicating that no valid value can be obtained.
  * `event_pattern` - Platform event matching rules.Note: This field may return null, indicating that no valid value can be obtained.


