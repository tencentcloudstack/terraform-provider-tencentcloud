---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_plateform_event_template"
sidebar_current: "docs-tencentcloud-datasource-eb_plateform_event_template"
description: |-
  Use this data source to query detailed information of eb plateform_event_template
---

# tencentcloud_eb_plateform_event_template

Use this data source to query detailed information of eb plateform_event_template

## Example Usage

```hcl
data "tencentcloud_eb_plateform_event_template" "plateform_event_template" {
  event_type = "eb_platform_test:TEST:ALL"
}
```

## Argument Reference

The following arguments are supported:

* `event_type` - (Required, String) Platform product event type.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `event_template` - Platform product event template.


