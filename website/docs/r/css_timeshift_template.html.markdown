---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_timeshift_template"
sidebar_current: "docs-tencentcloud-resource-css_timeshift_template"
description: |-
  Provides a resource to create a css timeshift_template
---

# tencentcloud_css_timeshift_template

Provides a resource to create a css timeshift_template

## Example Usage

```hcl
resource "tencentcloud_css_timeshift_template" "timeshift_template" {
  area                   = "Mainland"
  description            = "timeshift template"
  duration               = 604800
  item_duration          = 5
  remove_watermark       = true
  template_name          = "tf-test"
  transcode_template_ids = []
}
```

## Argument Reference

The following arguments are supported:

* `duration` - (Required, Int) The time shifting duration.Unit: Second.
* `template_name` - (Required, String) The template name.Maximum length: 255 bytes.Only letters, numbers, underscores, and hyphens are supported.
* `area` - (Optional, String) The region.`Mainland`: The Chinese mainland.`Overseas`: Outside the Chinese mainland.Default value: `Mainland`.
* `description` - (Optional, String) The template description.Only letters, numbers, underscores, and hyphens are supported.
* `item_duration` - (Optional, Int) The segment size.Value range: 3-10.Unit: Second.Default value: 5.
* `remove_watermark` - (Optional, Bool) Whether to remove watermarks.If you pass in `true`, the original stream will be recorded.Default value: `false`.
* `transcode_template_ids` - (Optional, Set: [`Int`]) The transcoding template IDs.This API works only if `RemoveWatermark` is `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

css timeshift_template can be imported using the id, e.g.

```
terraform import tencentcloud_css_timeshift_template.timeshift_template templateId
```

