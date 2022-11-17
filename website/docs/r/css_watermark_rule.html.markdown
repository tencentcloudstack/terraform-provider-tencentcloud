---
subcategory: "css"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_watermark_rule"
sidebar_current: "docs-tencentcloud-resource-css_watermark_rule"
description: |-
  Provides a resource to create a css watermark_rule
---

# tencentcloud_css_watermark_rule

Provides a resource to create a css watermark_rule

## Example Usage

```hcl
resource "tencentcloud_css_watermark_rule" "watermark_rule" {
  domain_name  = ""
  app_name     = ""
  stream_name  = ""
  watermark_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required, String) rule app name.
* `domain_name` - (Required, String) rule domain name.
* `stream_name` - (Required, String) rule stream name.
* `watermark_id` - (Required, Int) watermark id created by AddLiveWatermark.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - create time.
* `update_time` - update time.


## Import

css watermark_rule can be imported using the id, e.g.
```
$ terraform import tencentcloud_css_watermark_rule.watermark_rule watermarkRule_id
```

