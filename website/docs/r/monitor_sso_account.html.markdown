---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_sso_account"
sidebar_current: "docs-tencentcloud-resource-monitor_sso_account"
description: |-
  Provides a resource to create a monitor ssoAccount
---

# tencentcloud_monitor_sso_account

Provides a resource to create a monitor ssoAccount

## Example Usage

```hcl
resource "tencentcloud_monitor_sso_account" "ssoAccount" {
  instance_id = "grafana-50nj6v00"
  user_id     = "111"
  notes       = "desc12222"
  role {
    organization = "Main Org."
    role         = "Admin"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) grafana instance id.
* `user_id` - (Required, String) sub account uin of specific user.
* `notes` - (Optional, String) account related description.
* `role` - (Optional, List) grafana role.

The `role` object supports the following:

* `organization` - (Required, String) Grafana organization id string.
* `role` - (Required, String) Grafana role, one of {Admin,Editor,Viewer}.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor ssoAccount can be imported using the instance_id#user_id, e.g.
```
$ terraform import tencentcloud_monitor_sso_account.ssoAccount grafana-50nj6v00#111
```

