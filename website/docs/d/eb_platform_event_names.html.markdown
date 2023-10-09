---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_platform_event_names"
sidebar_current: "docs-tencentcloud-datasource-eb_platform_event_names"
description: |-
  Use this data source to query detailed information of eb platform_event_names
---

# tencentcloud_eb_platform_event_names

Use this data source to query detailed information of eb platform_event_names

## Example Usage

```hcl
data "tencentcloud_eb_platform_event_names" "platform_event_names" {
  product_type = ""
}
```

## Argument Reference

The following arguments are supported:

* `product_type` - (Required, String) Platform product event type.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `event_names` - Platform product list.
  * `event_name` - Event name.Note: This field may return null, indicating that no valid value can be obtained.
  * `event_type` - Event type.Note: This field may return null, indicating that no valid value can be obtained.


