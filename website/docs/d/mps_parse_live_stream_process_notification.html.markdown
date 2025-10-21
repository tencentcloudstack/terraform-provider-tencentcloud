---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_parse_live_stream_process_notification"
sidebar_current: "docs-tencentcloud-datasource-mps_parse_live_stream_process_notification"
description: |-
  Use this data source to query detailed information of mps parse_live_stream_process_notification
---

# tencentcloud_mps_parse_live_stream_process_notification

Use this data source to query detailed information of mps parse_live_stream_process_notification

## Example Usage

```hcl
data "tencentcloud_mps_parse_live_stream_process_notification" "parse_live_stream_process_notification" {
  content = "your_content"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Live stream event notification obtained from CMQ.
* `result_output_file` - (Optional, String) Used to save results.


