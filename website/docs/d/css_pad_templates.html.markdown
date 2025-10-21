---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_pad_templates"
sidebar_current: "docs-tencentcloud-datasource-css_pad_templates"
description: |-
  Use this data source to query detailed information of css pad_templates
---

# tencentcloud_css_pad_templates

Use this data source to query detailed information of css pad_templates

## Example Usage

```hcl
data "tencentcloud_css_pad_templates" "pad_templates" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `templates` - Live pad template information.
  * `create_time` - Template create time.
  * `description` - Description info.
  * `max_duration` - Maximum pad duration.Value range: 0 - positive infinity.Unit: milliseconds.
  * `template_id` - Template id.
  * `template_name` - Template name.
  * `type` - Pad content type: 1: Image, 2: Video. Default value: 1.
  * `update_time` - Template modify time.
  * `url` - Pad content.
  * `wait_duration` - Stream interruption waiting time.Value range: 0-30000.Unit: milliseconds.


