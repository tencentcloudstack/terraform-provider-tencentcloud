---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_watermarks"
sidebar_current: "docs-tencentcloud-datasource-css_watermarks"
description: |-
  Use this data source to query detailed information of css watermarks
---

# tencentcloud_css_watermarks

Use this data source to query detailed information of css watermarks

## Example Usage

```hcl
data "tencentcloud_css_watermarks" "watermarks" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `watermark_list` - Watermark information list.
  * `create_time` - The time when the watermark was added.Note: Beijing time (UTC+8) is used.
  * `height` - Watermark height.
  * `picture_url` - Watermark image URL.
  * `status` - Current status. 0: not used. 1: in use.
  * `watermark_id` - Watermark ID.
  * `watermark_name` - Watermark name.
  * `width` - Watermark width.
  * `x_position` - Display position: X-axis offset.
  * `y_position` - Display position: Y-axis offset.


