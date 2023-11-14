---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_start_stream_monitor"
sidebar_current: "docs-tencentcloud-resource-css_start_stream_monitor"
description: |-
  Provides a resource to create a css start_stream_monitor
---

# tencentcloud_css_start_stream_monitor

Provides a resource to create a css start_stream_monitor

## Example Usage

```hcl
resource "tencentcloud_css_start_stream_monitor" "start_stream_monitor" {
  monitor_id               = "3d5738dd-1ca2-4601-a6e9-004c5ec75c0b"
  audible_input_index_list = [1]
}
```

## Argument Reference

The following arguments are supported:

* `monitor_id` - (Required, String, ForceNew) Monitor id.
* `audible_input_index_list` - (Optional, Set: [`Int`], ForceNew) The input index for monitoring the screen audio, supports multiple input audio sources.The valid range for InputIndex is that it must already exist.If left blank, there will be no audio output by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

css start_stream_monitor can be imported using the id, e.g.

```
terraform import tencentcloud_css_start_stream_monitor.start_stream_monitor start_stream_monitor_id
```

