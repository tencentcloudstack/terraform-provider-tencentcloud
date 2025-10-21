---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_domain_error_page_infos"
sidebar_current: "docs-tencentcloud-datasource-gaap_domain_error_page_infos"
description: |-
  Use this data source to query detailed information of gaap domain error page infos
---

# tencentcloud_gaap_domain_error_page_infos

Use this data source to query detailed information of gaap domain error page infos

## Example Usage

```hcl
data "tencentcloud_gaap_domain_error_page_infos" "domain_error_page_infos" {
  error_page_ids = ["errorPage-xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `error_page_ids` - (Required, Set: [`String`]) Customized error ID list, supporting up to 10.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `error_page_set` - Custom error response configuration setNote: This field may return null, indicating that a valid value cannot be obtained.
  * `body` - Response body set (excluding HTTP header)Note: This field may return null, indicating that a valid value cannot be obtained.
  * `clear_headers` - Response headers that need to be cleanedNote: This field may return null, indicating that a valid value cannot be obtained.
  * `domain` - domain name.
  * `error_nos` - Original error code.
  * `error_page_id` - Configuration ID for error customization response.
  * `listener_id` - Listener ID.
  * `new_error_no` - New error codeNote: This field may return null, indicating that a valid value cannot be obtained.
  * `set_headers` - Response header to be setNote: This field may return null, indicating that a valid value cannot be obtained.
    * `header_name` - HTTP header name.
    * `header_value` - HTTP header value.
  * `status` - Rule status, 0 indicates successNote: This field may return null, indicating that a valid value cannot be obtained.


