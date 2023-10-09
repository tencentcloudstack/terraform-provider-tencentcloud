---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_custom_header"
sidebar_current: "docs-tencentcloud-datasource-gaap_custom_header"
description: |-
  Use this data source to query detailed information of gaap custom header
---

# tencentcloud_gaap_custom_header

Use this data source to query detailed information of gaap custom header

## Example Usage

```hcl
data "tencentcloud_gaap_custom_header" "custom_header" {
  rule_id = "rule-9sdhv655"
}
```

## Argument Reference

The following arguments are supported:

* `rule_id` - (Required, String) Rule IdNote: This field may return null, indicating that a valid value cannot be obtained.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `headers` - HeadersNote: This field may return null, indicating that a valid value cannot be obtained.
  * `header_name` - Header Name.
  * `header_value` - Header Value.


