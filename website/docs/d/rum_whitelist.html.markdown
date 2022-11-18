---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_whitelist"
sidebar_current: "docs-tencentcloud-datasource-rum_whitelist"
description: |-
  Use this data source to query detailed information of rum whitelist
---

# tencentcloud_rum_whitelist

Use this data source to query detailed information of rum whitelist

## Example Usage

```hcl
data "tencentcloud_rum_whitelist" "whitelist" {
  instance_id = "rum-pasZKEI3RLgakj"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, such as taw-123.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `whitelist_set` - While list.
  * `aid` - Business identifier.
  * `create_time` - Creation time.
  * `create_user` - Creator ID.
  * `remark` - Remarks.
  * `ttl` - End time.
  * `whitelist_uin` - uin: business identifier.
  * `wid` - Auto-Increment allowlist ID.


