---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_snapshot_template"
sidebar_current: "docs-tencentcloud-resource-css_snapshot_template"
description: |-
  Provides a resource to create a css snapshot_template
---

# tencentcloud_css_snapshot_template

Provides a resource to create a css snapshot_template

## Example Usage

```hcl
resource "tencentcloud_css_snapshot_template" "snapshot_template" {
  cos_app_id        = 1308919341
  cos_bucket        = "keep-bucket"
  cos_region        = "ap-guangzhou"
  description       = "snapshot template"
  height            = 0
  porn_flag         = 0
  snapshot_interval = 2
  template_name     = "tf-snapshot-template"
  width             = 0
}
```

## Argument Reference

The following arguments are supported:

* `cos_app_id` - (Required, Int) Cos application ID.
* `cos_bucket` - (Required, String) Cos bucket name. Note: The CosBucket parameter value cannot include the - [appid] part.
* `cos_region` - (Required, String) Cos region.
* `template_name` - (Required, String) Template name. Maximum length: 255 bytes. Only Chinese, English, numbers, `_`, `-` are supported.
* `cos_file_name` - (Optional, String) Cos file name. If it is empty, set according to the default value {StreamID}-screenshot-{Hour}-{Minute}-{Second}-{Width}x{Height}{Ext}.
* `cos_prefix` - (Optional, String) Cos Bucket folder prefix. If it is empty, set according to the default value /{Year}-{Month}-{Day}/.
* `description` - (Optional, String) Description information. Maximum length: 1024 bytes. Only `Chinese`, `English`, `numbers`, `_`, `-` are supported.
* `height` - (Optional, Int) Screenshot height. Default: 0 (original height). Range: 0-2000.
* `porn_flag` - (Optional, Int) Whether porn is enabled, 0: not enabled, 1: enabled. Default: 0.
* `snapshot_interval` - (Optional, Int) Screenshot interval, unit: s, default: 10s. Range: 2s~300s.
* `width` - (Optional, Int) Screenshot width. Default: 0 (original width). Range: 0-3000.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

css snapshot_template can be imported using the id, e.g.

```
terraform import tencentcloud_css_snapshot_template.snapshot_template templateId
```

