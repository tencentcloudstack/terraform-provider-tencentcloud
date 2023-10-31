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


