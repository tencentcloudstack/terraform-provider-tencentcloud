---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_sub_accounts"
sidebar_current: "docs-tencentcloud-datasource-cam_sub_accounts"
description: |-
  Use this data source to query detailed information of cam sub accounts
---

# tencentcloud_cam_sub_accounts

Use this data source to query detailed information of cam sub accounts

## Example Usage

```hcl
data "tencentcloud_cam_sub_accounts" "example" {
  filter_sub_account_uin = ["100037718139"]
}
```

## Argument Reference

The following arguments are supported:

* `filter_sub_account_uin` - (Required, Set: [`Int`]) List of sub-user UINs. Up to 50 UINs are supported.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `sub_accounts` - Sub-user list.


