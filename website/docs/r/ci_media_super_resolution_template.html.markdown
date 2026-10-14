---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_super_resolution_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_super_resolution_template"
description: |-
  Provides a resource to create a ci media_super_resolution_template
---

# tencentcloud_ci_media_super_resolution_template

Provides a resource to create a ci media_super_resolution_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_super_resolution_template" "media_super_resolution_template" {
  bucket          = "terraform-ci-1308919341"
  name            = "super_resolution_template"
  resolution      = "sdtohd"
  enable_scale_up = "true"
  version         = "Enhance"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.
* `resolution` - (Required, String) Resolution Options sdtohd: Standard Definition to Ultra Definition, hdto4k: HD to 4K.
* `enable_scale_up` - (Optional, String) Auto scaling switch, off by default.
* `version` - (Optional, String) version, default value Base, Base: basic version, Enhance: enhanced version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci media_super_resolution_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_super_resolution_template.media_super_resolution_template terraform-ci-xxxxxx#t1d707eb2be3294e22b47123894f85cb8f
```

