---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_whitelist"
sidebar_current: "docs-tencentcloud-resource-rum_whitelist"
description: |-
  Provides a resource to create a rum whitelist
---

# tencentcloud_rum_whitelist

Provides a resource to create a rum whitelist

## Example Usage

```hcl
resource "tencentcloud_rum_whitelist" "whitelist" {
  instance_id   = "rum-pasZKEI3RLgakj"
  remark        = "white list remark"
  whitelist_uin = "20221122"
  # aid = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, such as taw-123.
* `remark` - (Required, String) Remarks.
* `whitelist_uin` - (Required, String) uin: business identifier.
* `aid` - (Optional, String) Business identifier.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time.
* `create_user` - Creator ID.
* `ttl` - End time.
* `wid` - Auto-Increment allowlist ID.


## Import

rum whitelist can be imported using the id, e.g.
```
$ terraform import tencentcloud_rum_whitelist.whitelist whitelist_id
```

