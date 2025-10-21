---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_stream_monitor"
sidebar_current: "docs-tencentcloud-resource-css_stream_monitor"
description: |-
  Provides a resource to create a css stream_monitor
---

# tencentcloud_css_stream_monitor

Provides a resource to create a css stream_monitor

## Example Usage

```hcl
resource "tencentcloud_css_stream_monitor" "stream_monitor" {
  ai_asr_input_index_list = [
    1,
  ]
  ai_format_diagnose = 1
  ai_ocr_input_index_list = [
    1,
  ]
  allow_monitor_report        = 1
  asr_language                = 1
  check_stream_broken         = 1
  check_stream_low_frame_rate = 1
  monitor_name                = "test"
  ocr_language                = 1

  input_list {
    input_app         = "live"
    input_domain      = "177154.push.tlivecloud.com"
    input_stream_name = "ppp"
  }

  notify_policy {
    callback_url       = "http://example.com/test"
    notify_policy_type = 1
  }

  output_info {
    output_domain        = "test122.jingxhu.top"
    output_stream_height = 1080
    output_stream_name   = "afc7847d-1fe1-43bc-b1e4-20d86303c393"
    output_stream_width  = 1920
  }
}
```

## Argument Reference

The following arguments are supported:

* `input_list` - (Required, List) Wait monitor input info list.
* `output_info` - (Required, List) Monitor task output info.
* `ai_asr_input_index_list` - (Optional, Set: [`Int`]) AI asr input index list.(first input index is 1.).
* `ai_format_diagnose` - (Optional, Int) If enable format diagnose.
* `ai_ocr_input_index_list` - (Optional, Set: [`Int`]) Ai ocr input index list(first input index is 1.).
* `allow_monitor_report` - (Optional, Int) If store monitor event.
* `asr_language` - (Optional, Int) Asr language.0: close.1: Chinese2: English3: Japanese4: Korean.
* `check_stream_broken` - (Optional, Int) If enable stream broken check.
* `check_stream_low_frame_rate` - (Optional, Int) If enable low frame rate check.
* `monitor_name` - (Optional, String) Monitor task name.
* `notify_policy` - (Optional, List) Monitor event notify policy.
* `ocr_language` - (Optional, Int) Intelligent text recognition language settings: ocr language.0: close.1. Chinese,English.

The `input_list` object supports the following:

* `input_stream_name` - (Required, String) Wait monitor input stream name.limit 256 bytes.
* `description` - (Optional, String) Description content.limit 256 bytes.
* `input_app` - (Optional, String) Wait monitor input push path.limit 32 bytes.
* `input_domain` - (Optional, String) Wait monitor input push domain.limit 128 bytes.
* `input_url` - (Optional, String) Wait monitor input stream push url.

The `notify_policy` object supports the following:

* `callback_url` - (Optional, String) Callback url.limit [0,512].only http or https.
* `notify_policy_type` - (Optional, Int) Notify policy type.0: not notify.1: use global policy.

The `output_info` object supports the following:

* `output_stream_height` - (Required, Int) Monitor task output height, limit[1, 1080].
* `output_stream_width` - (Required, Int) Output stream width, limit[1, 1920].
* `output_app` - (Optional, String) Monitor task play path.limit 32 bytes.
* `output_domain` - (Optional, String) Monitor task output play domain.limit 128 bytes.
* `output_stream_name` - (Optional, String) Monitor task output stream name.limit 256 bytes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

css stream_monitor can be imported using the id, e.g.

```
terraform import tencentcloud_css_stream_monitor.stream_monitor stream_monitor_id
```

