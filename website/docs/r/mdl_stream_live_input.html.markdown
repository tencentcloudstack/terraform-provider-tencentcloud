---
subcategory: "StreamLive(MDL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mdl_stream_live_input"
sidebar_current: "docs-tencentcloud-resource-mdl_stream_live_input"
description: |-
  Provides a resource to create a mdl streamlive_input
---

# tencentcloud_mdl_stream_live_input

Provides a resource to create a mdl streamlive_input

## Example Usage

```hcl
resource "tencentcloud_mdl_stream_live_input" "stream_live_input" {
  name = "terraform_test"
  type = "RTP_PUSH"
  security_group_ids = [
    "6405DF9D000007DFB4EC"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Input name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level.
* `type` - (Required, String, ForceNew) Input typeValid values: `RTMP_PUSH`, `RTP_PUSH`, `UDP_PUSH`, `RTMP_PULL`, `HLS_PULL`, `MP4_PULL`.
* `input_settings` - (Optional, List) Input settings. For the type `RTMP_PUSH`, `RTMP_PULL`, `HLS_PULL`, or `MP4_PULL`, 1 or 2 inputs of the corresponding type can be configured.
* `security_group_ids` - (Optional, Set: [`String`]) ID of the input security group to attachYou can attach only one security group to an input.

The `input_settings` object supports the following:

* `app_name` - (Optional, String) Application name, which is valid if `Type` is `RTMP_PUSH` and can contain 1-32 letters and digitsNote: This field may return `null`, indicating that no valid value was found.
* `delay_time` - (Optional, Int) Delayed time (ms) for playback, which is valid if `Type` is `RTMP_PUSH`Value range: 0 (default) or 10000-600000The value must be a multiple of 1,000.Note: This field may return `null`, indicating that no valid value was found.
* `input_address` - (Optional, String) RTP/UDP input address, which does not need to be entered for the input parameter.Note: this field may return null, indicating that no valid values can be obtained.
* `input_domain` - (Optional, String) The domain of an SRT_PUSH address. If this is a request parameter, you do not need to specify it.Note: This field may return `null`, indicating that no valid value was found.
* `password` - (Optional, String) The password, which is used for authentication.Note: This field may return `null`, indicating that no valid value was found.
* `source_type` - (Optional, String) Source type for stream pulling and relaying. To pull content from private-read COS buckets under the current account, set this parameter to `TencentCOS`; otherwise, leave it empty.Note: this field may return `null`, indicating that no valid value was found.
* `source_url` - (Optional, String) Source URL, which is valid if `Type` is `RTMP_PULL`, `HLS_PULL`, or `MP4_PULL` and can contain 1-512 charactersNote: This field may return `null`, indicating that no valid value was found.
* `stream_name` - (Optional, String) Stream name, which is valid if `Type` is `RTMP_PUSH` and can contain 1-32 letters and digitsNote: This field may return `null`, indicating that no valid value was found.
* `user_name` - (Optional, String) The username, which is used for authentication.Note: This field may return `null`, indicating that no valid value was found.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mdl stream_live_input can be imported using the id, e.g.

```
terraform import tencentcloud_mdl_stream_live_input.stream_live_input id
```

