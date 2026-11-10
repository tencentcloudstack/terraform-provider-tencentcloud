---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_provision_role_configuration"
sidebar_current: "docs-tencentcloud-action-organization_provision_role_configuration"
description: |-
  Provides an action to provision (deploy) a permission configuration to a target account of the Tencent Cloud Organization via the ProvisionRoleConfiguration API. The action waits for the asynchronous task returned by the API to finish. This is a one-time operation; no cloud-side state is persisted after the action completes.
---

# tencentcloud_organization_provision_role_configuration

Provides an action to provision (deploy) a permission configuration to a target account of the Tencent Cloud Organization via the ProvisionRoleConfiguration API. The action waits for the asynchronous task returned by the API to finish. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

## Example Usage

```hcl
action "tencentcloud_organization_provision_role_configuration" "example" {
  config {
    zone_id               = "z-1os7c9znogct"
    role_configuration_id = "rc-ihogrs0e6ceg"
    target_type           = "MemberUin"
    target_uin            = 100043 * * * * 45

    timeouts {
      invoke = "3m"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `role_configuration_id` - (Required, String) Permission configuration ID.
* `target_type` - (Required, String) Type of the synchronized target account of the Tencent Cloud Organization. Valid values: `ManagerUin` (management account), `MemberUin` (member account).
* `target_uin` - (Required, Int) UIN of the target account of the Tencent Cloud Organization.
* `zone_id` - (Required, String) Space ID.
* `timeouts` - (Optional, Object) The timeouts block allows you to specify the timeout for the invoke operation.

## API Reference

This action invokes the following TencentCloud APIs:

* Invoke:
  - [GetTaskStatus](https://www.tencentcloud.com/document/api/1031/63608)
  - [ProvisionRoleConfiguration](https://www.tencentcloud.com/document/api/1031/63605)

