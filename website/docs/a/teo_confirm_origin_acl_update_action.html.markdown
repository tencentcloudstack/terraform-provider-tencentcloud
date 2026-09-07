---
subcategory: "Provider Meta"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_confirm_origin_acl_update_action"
sidebar_current: "docs-tencentcloud-action-teo_confirm_origin_acl_update_action"
description: |-
  Provides an action to confirm TEO origin ACL update for a zone. When the origin IP ranges of TEO change, you can use this action to confirm that the latest origin IP ranges have been updated to the origin firewall, and the change notification will stop being pushed.
---

# tencentcloud_teo_confirm_origin_acl_update_action

Provides an action to confirm TEO origin ACL update for a zone. When the origin IP ranges of TEO change, you can use this action to confirm that the latest origin IP ranges have been updated to the origin firewall, and the change notification will stop being pushed.

~> **NOTE:** `Action` type resources require Terraform version to be updated to at least 1.14.x or higher

## Example Usage

```hcl
action "tencentcloud_teo_confirm_origin_acl_update_action" "example" {
  config {
    zone_id = "zone-3fkff38fyw8s"
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Zone ID.


