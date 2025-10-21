---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_tts_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_tts_template"
description: |-
  Provides a resource to create a ci media_tts_template
---

# tencentcloud_ci_media_tts_template

Provides a resource to create a ci media_tts_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_tts_template" "media_tts_template" {
  bucket     = "terraform-ci-xxxxxx"
  name       = "tts_template"
  mode       = "Asyc"
  codec      = "pcm"
  voice_type = "ruxue"
  volume     = "0"
  speed      = "100"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.
* `codec` - (Optional, String) Audio format, default wav (synchronous)/pcm (asynchronous, wav, mp3, pcm.
* `mode` - (Optional, String) Processing mode, default value Asyc, Asyc (asynchronous composition), Sync (synchronous composition), When Asyc is selected, the codec only supports pcm.
* `speed` - (Optional, String) Speech rate, the default value is 100, [50,200].
* `voice_type` - (Optional, String) Timbre, the default value is ruxue.
* `volume` - (Optional, String) Volume, default value 0, [-10,10].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci media_tts_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_tts_template.media_tts_template terraform-ci-xxxxxx#t1ed421df8bd2140b6b73474f70f99b0f8
```

