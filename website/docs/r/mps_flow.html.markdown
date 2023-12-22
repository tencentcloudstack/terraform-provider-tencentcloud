---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_flow"
sidebar_current: "docs-tencentcloud-resource-mps_flow"
description: |-
  Provides a resource to create a mps flow
---

# tencentcloud_mps_flow

Provides a resource to create a mps flow

## Example Usage

### Create a mps RTP flow

```hcl
resource "tencentcloud_mps_event" "event" {
  event_name  = "tf_test_event_srt_%d"
  description = "tf test mps event description"
}

resource "tencentcloud_mps_flow" "flow" {
  flow_name     = "tf_test_mps_flow_srt_%d"
  max_bandwidth = 10000000
  input_group {
    input_name    = "test_inputname"
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
  event_id = tencentcloud_mps_event.event.id
}
```

### Create a mps RTP flow

```hcl
resource "tencentcloud_mps_event" "event_rtp" {
  event_name  = "tf_test_event_rtp_%d"
  description = "tf test mps event description"
}

resource "tencentcloud_mps_flow" "flow_rtp" {
  flow_name     = "tf_test_mps_flow_rtp_%d"
  max_bandwidth = 10000000
  input_group {
    input_name    = "test_inputname"
    protocol      = "RTP"
    description   = "input name Description"
    allow_ip_list = ["0.0.0.0/0"]
    rtp_settings {
      fec          = "none"
      idle_timeout = 1000
    }
  }
  event_id = tencentcloud_mps_event.event_rtp.id
}
```

### Create a mps RTP flow and start it

Before you start a mps flow, you need to create a output first.

```hcl
resource "tencentcloud_mps_event" "event_rtp" {
  event_name  = "your_event_name"
  description = "tf test mps event description"
}

resource "tencentcloud_mps_flow" "flow_rtp" {
  flow_name     = "your_flow_name"
  max_bandwidth = 10000000
  input_group {
    input_name    = "test_inputname"
    protocol      = "RTP"
    description   = "input name Description"
    allow_ip_list = ["0.0.0.0/0"]
    rtp_settings {
      fec          = "none"
      idle_timeout = 1000
    }
  }
  event_id = tencentcloud_mps_event.event_rtp.id
}

resource "tencentcloud_mps_output" "output" {
  flow_id = tencentcloud_mps_flow.flow_rtp.id
  output {
    output_name   = "your_output_name"
    description   = "tf mps output group"
    protocol      = "RTP"
    output_region = "ap-guangzhou"
    rtp_settings {
      destinations {
        ip   = "203.205.141.84"
        port = 65535
      }
      fec          = "none"
      idle_timeout = 1000
    }
  }
}

resource "tencentcloud_mps_start_flow_operation" "operation" {
  flow_id    = tencentcloud_mps_flow.flow_rtp.id
  start      = true
  depends_on = [tencentcloud_mps_output.output]
}
```

## Argument Reference

The following arguments are supported:

* `flow_name` - (Required, String) Flow name.
* `max_bandwidth` - (Required, Int) Maximum bandwidth, unit bps, optional [10000000, 20000000, 50000000].
* `event_id` - (Optional, String) The event ID associated with this Flow. Each flow can only be associated with one Event.
* `input_group` - (Optional, List) The input group for the flow.

The `hls_pull_settings` object of `input_group` supports the following:

* `source_addresses` - (Required, List) There is only one origin address of the HLS origin station.

The `input_group` object supports the following:

* `input_name` - (Required, String) Input name, you can fill in uppercase and lowercase letters, numbers and underscores, and the length is [1, 32].
* `protocol` - (Required, String) Input protocol, optional [SRT|RTP|RTMP|RTMP_PULL].
* `allow_ip_list` - (Optional, Set) The input IP whitelist, the format is CIDR.
* `description` - (Optional, String) Input description with a length of [0, 255].
* `fail_over` - (Optional, String) The active/standby switch of the input, [OPEN|CLOSE] is optional, and the default is CLOSE.
* `hls_pull_settings` - (Optional, List) Input HLS_PULL configuration information.
* `resilient_stream` - (Optional, List) Delay broadcast smooth streaming configuration information.
* `rtmp_pull_settings` - (Optional, List) Input RTMP_PULL configuration information.
* `rtp_settings` - (Optional, List) RTP configuration information.
* `rtsp_pull_settings` - (Optional, List) Input RTSP_PULL configuration information.
* `srt_settings` - (Optional, List) The input SRT configuration information.

The `resilient_stream` object of `input_group` supports the following:

* `buffer_time` - (Optional, Int) Delay time, in seconds, currently supports a range of 10 to 300 seconds. Note: This field may return null, indicating that no valid value can be obtained.
* `enable` - (Optional, Bool) Whether to enable the delayed broadcast smooth spit stream, true is enabled, false is not enabled, and the default is not enabled. Note: This field may return null, indicating that no valid value can be obtained.

The `rtmp_pull_settings` object of `input_group` supports the following:

* `source_addresses` - (Required, List) The source site address of the RTMP source site, there can only be one.

The `rtp_settings` object of `input_group` supports the following:

* `fec` - (Optional, String) Defaults to none, optional values[none].
* `idle_timeout` - (Optional, Int) Idle timeout, the default is 5000, the unit is ms, and the range is [1000, 10000].

The `rtsp_pull_settings` object of `input_group` supports the following:

* `source_addresses` - (Required, List) The source site address of the RTSP source site, there can only be one.

The `source_addresses` object of `hls_pull_settings` supports the following:

* `url` - (Required, String) The URL address of the HLS origin site.

The `source_addresses` object of `rtmp_pull_settings` supports the following:

* `stream_key` - (Required, String) StreamKey information of the RTMP source site.
* `tc_url` - (Required, String) TcUrl address of the RTMP source server.

The `source_addresses` object of `rtsp_pull_settings` supports the following:

* `url` - (Required, String) The URL address of the RTSP source site.

The `source_addresses` object of `srt_settings` supports the following:

* `ip` - (Required, String) Peer IP.
* `port` - (Required, Int) Peer port.

The `srt_settings` object of `input_group` supports the following:

* `latency` - (Optional, Int) Delay, default 0, unit ms, range [0, 3000].
* `mode` - (Optional, String) SRT mode, optional [LISTENER|CALLER], default is LISTENER.
* `passphrase` - (Optional, String) The decryption key, which is empty by default, means no encryption. Only ascii code values can be filled in, and the length is [10, 79].
* `pb_key_len` - (Optional, Int) Key length, default is 0, optional [0|16|24|32].
* `peer_idle_timeout` - (Optional, Int) Peer timeout, default is 5000, unit ms, range is [1000, 10000].
* `peer_latency` - (Optional, Int) Peer delay, the default is 0, the unit is ms, and the range is [0, 3000].
* `recv_latency` - (Optional, Int) Receiving delay, default is 120, unit ms, range is [0, 3000].
* `source_addresses` - (Optional, List) SRT peer address, required when Mode is CALLER, and only 1 set can be filled in.
* `stream_id` - (Optional, String) Stream ID, optional uppercase and lowercase letters, numbers and special characters (.#!:&amp;,=_-), length 0~512. For specific format, please refer to:https://github.com/Haivision/srt/blob/master/docs/features/access-control.md#standard-keys.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps flow can be imported using the id, e.g.

```
terraform import tencentcloud_mps_flow.flow flow_id
```

