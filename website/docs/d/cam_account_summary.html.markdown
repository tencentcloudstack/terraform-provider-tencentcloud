---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_account_summary"
sidebar_current: "docs-tencentcloud-datasource-cam_account_summary"
description: |-
  Use this data source to query detailed information of cam account_summary
---

# tencentcloud_cam_account_summary

Use this data source to query detailed information of cam account_summary

## Example Usage

```hcl
data "tencentcloud_cam_account_summary" "account_summary" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `group` - The number of Group.
* `identity_providers` - The number of identity provider.
* `member` - The number of grouped users.
* `policies` - The number of policy.
* `roles` - The number of role.
* `user` - The number of Sub-user.


