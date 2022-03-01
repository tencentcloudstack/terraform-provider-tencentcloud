---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_user_info"
sidebar_current: "docs-tencentcloud-datasource-user_info"
description: |-
  Use this data source to query user appid, uin and ownerUin.
---

# tencentcloud_user_info

Use this data source to query user appid, uin and ownerUin.

## Example Usage

```hcl
data "tencentcloud_user_info" "foo" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional) Used for save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `app_id` - Current account App ID.
* `owner_uin` - Current account OwnerUIN.
* `uin` - Current account UIN.


