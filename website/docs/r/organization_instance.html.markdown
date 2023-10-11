---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_instance"
sidebar_current: "docs-tencentcloud-resource-organization_instance"
description: |-
  Provides a resource to create a organization organization
---

# tencentcloud_organization_instance

Provides a resource to create a organization organization

## Example Usage

```hcl
resource "tencentcloud_organization_instance" "organization" {
}
```

## Argument Reference

The following arguments are supported:



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Organize the creation time.Note: This field may return NULL, indicating that the valid value cannot be obtained.
* `host_uin` - Creator Uin.Note: This field may return NULL, indicating that the valid value cannot be obtained.
* `is_allow_quit` - Whether the members are allowed to withdraw.Allow: Allow, not allowed: DENIEDNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `is_assign_manager` - Whether a trusted service administrator.Yes: true, no: falseNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `is_auth_manager` - Whether the real -name subject administrator.Yes: true, no: falseNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `is_manager` - Whether to organize an administrator.Yes: true, no: falseNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `join_time` - Members join time.Note: This field may return NULL, indicating that the valid value cannot be obtained.
* `nick_name` - Creator nickname.Note: This field may return NULL, indicating that the valid value cannot be obtained.
* `org_id` - Enterprise organization ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.
* `org_permission` - List of membership authority of members.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `id` - Permissions ID.
  * `name` - Permission name.
* `org_policy_name` - Strategic name.Note: This field may return NULL, indicating that the valid value cannot be obtained.
* `org_policy_type` - Strategy type.Financial Management: FinancialNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `org_type` - Enterprise organization type.Note: This field may return NULL, indicating that the valid value cannot be obtained.
* `pay_name` - The name of the payment.Note: This field may return NULL, indicating that the valid value cannot be obtained.
* `pay_uin` - UIN on behalf of the payer.Note: This field may return NULL, indicating that the valid value cannot be obtained.
* `root_node_id` - Organize the root node ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.


## Import

organization organization can be imported using the id, e.g.

```
terraform import tencentcloud_organization_instance.organization organization_id
```

