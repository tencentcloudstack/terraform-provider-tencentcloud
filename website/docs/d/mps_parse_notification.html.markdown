---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_parse_notification"
sidebar_current: "docs-tencentcloud-datasource-mps_parse_notification"
description: |-
  Use this data source to query detailed information of mps parse_notification
---

# tencentcloud_mps_parse_notification

Use this data source to query detailed information of mps parse_notification

## Example Usage

```hcl
data "tencentcloud_mps_parse_notification" "parse_notification" {
  content = "your_content"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Event notification obtained from CMQ.
* `result_output_file` - (Optional, String) Used to save results.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [ParseNotification](https://cloud.tencent.com/document/api/862/39228) - After a message is received from CMQ, parse out the content of the MPS event notification from the msgBody field.

