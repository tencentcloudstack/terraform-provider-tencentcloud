---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_provision_role_configuration_operation"
sidebar_current: "docs-tencentcloud-resource-provision_role_configuration_operation"
description: |-
  Provides a resource to create a Organization provision role configuration
---

# tencentcloud_provision_role_configuration_operation

Provides a resource to create a Organization provision role configuration

~> **NOTE:** This is a one-time provisioning operation rather than a continuously reconciled cloud resource. It is recommended to use the `tencentcloud_organization_provision_role_configuration` **action** resource instead, which is purpose-built for one-off operations like deploying a permission configuration to a target account. Actions avoid persisting unnecessary state (this resource performs no-op `Read` and `Delete`), natively support invocation timeout configuration via the `timeouts` block, and are available in Terraform 1.14+. See [tencentcloud_organization_provision_role_configuration](https://registry.terraform.io/providers/tencentcloudstack/tencentcloud/latest/docs/actions/organization_provision_role_configuration) for details.

## Example Usage

```hcl
resource "tencentcloud_provision_role_configuration_operation" "example" {
  zone_id               = "z-1os7c9znogct"
  role_configuration_id = "rc-ihogrs0e6ceg"
  target_type           = "MemberUin"
  target_uin            = 100043 * * * * 45
}
```

## Argument Reference

The following arguments are supported:

* `role_configuration_id` - (Required, String, ForceNew) Permission configuration ID.
* `target_type` - (Required, String, ForceNew) Type of the synchronized target account of the Tencent Cloud Organization. ManagerUin: admin account; MemberUin: member account.
* `target_uin` - (Required, Int, ForceNew) UIN of the target account of the Tencent Cloud Organization.
* `zone_id` - (Required, String, ForceNew) Space ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [GetTaskStatus](https://www.tencentcloud.com/document/api/1031/63608)
  - [ProvisionRoleConfiguration](https://www.tencentcloud.com/document/api/1031/63605)

