---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_input"
sidebar_current: "docs-tencentcloud-resource-mps_input"
description: |-
  Provides a resource to create a mps input
---

# tencentcloud_mps_input

Provides a resource to create a mps input

## Example Usage

### Create mps input group with SRT

```hcl
resource "tencentcloud_mps_input" "input" {
  flow_id = tencentcloud_mps_flow.flow.id
  input_group {
    input_name    = "your_input_name"
    protocol      = "SRT"
    description   = "input name Description"
    allow_ip_list = ["0.0.0.0/0"]
    srt_settings {
      mode              = "LISTENER"
      stream_id         = "#!::u=johnny,r=resource,h=xxx.com,t=stream,m=play"
      latency           = 1000
      recv_latency      = 1000
      peer_latency      = 1000
      peer_idle_timeout = 1000
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `flow_id` - (Required, String) Flow ID.
* `input_group` - (Optional, List) The input group for the input. Only support one group for one `tencentcloud_mps_input`. Use `for_each` to create multiple inputs Scenario.

The `hls_pull_settings` object supports the following:

* `source_addresses` - (Required, List) There is only one origin address of the HLS origin station.

The `input_group` object supports the following:

* `input_name` - (Required, String) The input name, you can fill in uppercase and lowercase letters, numbers and underscores, and the length is [1, 32].
* `protocol` - (Required, String) Input protocol, optional [SRT|RTP|RTMP|RTMP_PULL].
* `allow_ip_list` - (Optional, Set) The input IP whitelist, the format is CIDR.
* `description` - (Optional, String) The input description with a length of [0, 255].
* `fail_over` - (Optional, String) The active/standby switch of the input, [OPEN|CLOSE] is optional, and the default is CLOSE.
* `hls_pull_settings` - (Optional, List) Input HLS_PULL configuration information.
* `resilient_stream` - (Optional, List) Delay broadcast smooth streaming configuration information.
* `rtmp_pull_settings` - (Optional, List) Input RTMP_PULL configuration information.
* `rtp_settings` - (Optional, List) Input RTP configuration information.
* `rtsp_pull_settings` - (Optional, List) Input RTSP_PULL configuration information.
* `srt_settings` - (Optional, List) The input SRT configuration information.

The `resilient_stream` object supports the following:

* `buffer_time` - (Optional, Int) Delay time, in seconds, currently supports a range of 10 to 300 seconds. Note: This field may return null, indicating that no valid value can be obtained.
* `enable` - (Optional, Bool) Whether to enable the delayed broadcast smooth spit stream, true is enabled, false is not enabled, and the default is not enabled. Note: This field may return null, indicating that no valid value can be obtained.

The `rtmp_pull_settings` object supports the following:

* `source_addresses` - (Required, List) The source site address of the RTMP source site, there can only be one.

The `rtp_settings` object supports the following:

* `fec` - (Optional, String) Defaults to &#39;none&#39;, optional values[&#39;none&#39;].
* `idle_timeout` - (Optional, Int) Idle timeout, the default is 5000, the unit is ms, and the range is [1000, 10000].

The `rtsp_pull_settings` object supports the following:

* `source_addresses` - (Required, List) The source site address of the RTSP source site, there can only be one.

The `source_addresses` object supports the following:

* `ip` - (Required, String) Peer IP.
* `port` - (Required, Int) Peer port.

The `source_addresses` object supports the following:

* `stream_key` - (Required, String) StreamKey information of the RTMP source site.
* `tc_url` - (Required, String) TcUrl address of the RTMP source server.

The `source_addresses` object supports the following:

* `url` - (Required, String) The URL address of the HLS origin site.

The `source_addresses` object supports the following:

* `url` - (Required, String) The URL address of the RTSP source site.

The `srt_settings` object supports the following:

* `latency` - (Optional, Int) Delay, default 0, unit ms, range [0, 3000].
* `mode` - (Optional, String) SRT mode, optional [LISTENER|CALLER], default is LISTENER.
* `passphrase` - (Optional, String) The decryption key, which is empty by default, means no encryption. Only ascii code values can be filled in, and the length is [10, 79].
* `pb_key_len` - (Optional, Int) Key length, default is 0, optional [0|16|24|32].
* `peer_idle_timeout` - (Optional, Int) Peer timeout, default is 5000, unit ms, range is [1000, 10000].
* `peer_latency` - (Optional, Int) Peer delay, the default is 0, the unit is ms, and the range is [0, 3000].
* `recv_latency` - (Optional, Int) Receiving delay, default is 120, unit ms, range is [0, 3000].
* `source_addresses` - (Optional, List) SRT peer address, required when Mode is CALLER, and only 1 set can be filled in.
* `stream_id` - (Optional, String) Stream ID, optional uppercase and lowercase letters, numbers and special characters (.#!:&amp;,=_-), length 0~512. Specific format can refer to:https://github.com/Haivision/srt/blob/master/docs/features/access-control.md#standard-keys.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps input can be imported using the id, e.g.

```
terraform import tencentcloud_mps_input.input input_id
```

