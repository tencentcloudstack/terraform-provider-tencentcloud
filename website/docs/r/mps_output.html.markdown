---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_output"
sidebar_current: "docs-tencentcloud-resource-mps_output"
description: |-
  Provides a resource to create a mps output
---

# tencentcloud_mps_output

Provides a resource to create a mps output

## Example Usage

### Create a output group with RTP

```hcl
resource "tencentcloud_mps_output" "output" {
  flow_id = "your_flow_id"
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
```

## Argument Reference

The following arguments are supported:

* `flow_id` - (Required, String) Flow ID.
* `output` - (Required, List) Output configuration of the transport stream.

The `destinations` object supports the following:

* `ip` - (Required, String) Output IP.
* `port` - (Required, Int) output port.

The `destinations` object supports the following:

* `ip` - (Required, String) The target IP of the relay.
* `port` - (Required, Int) Destination port for relays.

The `destinations` object supports the following:

* `stream_key` - (Required, String) relayed StreamKey, in the format: stream?key=value.
* `url` - (Required, String) relayed URL, the format is: rtmp://domain/live.

The `output` object supports the following:

* `description` - (Required, String) Output description.
* `output_name` - (Required, String) The name of the output.
* `output_region` - (Required, String) Output region.
* `protocol` - (Required, String) Output protocol, optional [SRT|RTP|RTMP|RTMP_PULL].
* `allow_ip_list` - (Optional, Set) IP whitelist list, the format is CIDR, such as 0.0.0.0/0. When the Protocol is RTMP_PULL, it is valid, and if it is empty, it means that the client IP is not limited.
* `max_concurrent` - (Optional, Int) The maximum number of concurrent pull streams, the maximum is 4, and the default is 4. Only SRT or RTMP_PULL can set this parameter.
* `rtmp_settings` - (Optional, List) Output RTMP configuration.
* `rtp_settings` - (Optional, List) Output RTP configuration.
* `srt_settings` - (Optional, List) configuration of the output SRT.

The `rtmp_settings` object supports the following:

* `destinations` - (Required, List) The target address of the relay can be filled in 1~2.
* `chunk_size` - (Optional, Int) RTMP Chunk size, range is [4096, 40960].

The `rtp_settings` object supports the following:

* `destinations` - (Required, List) The target address of the relay can be filled in 1~2.
* `fec` - (Required, String) You can only fill in none.
* `idle_timeout` - (Required, Int) Idle timeout, unit ms.

The `srt_settings` object supports the following:

* `destinations` - (Required, List) The target address of the relay is required when Mode is CALLER, and only one group can be filled in.
* `latency` - (Optional, Int) The total delay of relaying SRT, the default is 0, the unit is ms, and the range is [0, 3000].
* `mode` - (Optional, String) SRT mode, optional [LISTENER|CALLER], default is CALLER.
* `passphrase` - (Optional, String) The encryption key for relaying SRT, which is empty by default, indicating no encryption. Only ascii code values can be filled in, and the length is [10, 79].
* `pb_key_len` - (Optional, Int) The key length of relay SRT, the default is 0, optional [0|16|24|32].
* `peer_idle_timeout` - (Optional, Int) The peer idle timeout for relaying SRT, the default is 5000, the unit is ms, and the range is [1000, 10000].
* `peer_latency` - (Optional, Int) The peer delay of relaying SRT, the default is 0, the unit is ms, and the range is [0, 3000].
* `recv_latency` - (Optional, Int) The reception delay of relay SRT, the default is 120, the unit is ms, the range is [0, 3000].
* `stream_id` - (Optional, String) relay the stream ID of SRT. You can choose uppercase and lowercase letters, numbers and special characters (.#!:&amp;,=_-). The length is 0~512.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps output can be imported using the id, e.g.

```
terraform import tencentcloud_mps_output.output flow_id#output_id
```

