---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_sign"
sidebar_current: "docs-tencentcloud-datasource-rum_sign"
description: |-
  Use this data source to query detailed information of rum sign
---

# tencentcloud_rum_sign

Use this data source to query detailed information of rum sign

## Example Usage

```hcl
data "tencentcloud_rum_sign" "sign" {
  timeout   = 1800
  file_type = 1
}
```

## Argument Reference

The following arguments are supported:

* `file_type` - (Optional, Int) Bucket type. `1`:web project; `2`:app project.
* `result_output_file` - (Optional, String) Used to save results.
* `timeout` - (Optional, Int) Timeout duration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `expired_time` - Expiration timestamp.
* `secret_id` - Temporary access key ID.
* `secret_key` - Temporary access key.
* `session_token` - Temporary access key token.
* `start_time` - Start timestamp.


