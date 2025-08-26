---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_member_auth_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-organization_member_auth_policy_attachment"
description: |-
  Provides a resource to create a Organization member auth policy attachment
---

# tencentcloud_organization_member_auth_policy_attachment

Provides a resource to create a Organization member auth policy attachment

## Example Usage

```hcl
resource "tencentcloud_organization_member_auth_policy_attachment" "example" {
  policy_id           = 252421751
  org_sub_account_uin = 100037718939
}
```

## Argument Reference

The following arguments are supported:

* `org_sub_account_uin` - (Required, Int, ForceNew) Organization administrator sub-account Uin.
* `policy_id` - (Required, Int, ForceNew) Pilicy ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `bind_type` - Bind type. 1-Subaccount, 2-User Group.
* `create_time` - Create time.
* `identity_id` - Identity ID.
* `identity_role_alias_name` - Identity role alias name.
* `identity_role_name` - Identity role name.
* `member_name` - Member name.
* `member_uin` - Member UIN.
* `org_sub_account_name` - Org sub account name.
* `policy_name` - Policy name.


## Import

Organization member auth policy attachment can be imported using the id, e.g.

```
terraform import tencentcloud_organization_member_auth_policy_attachment.example 252421751#100037718939
```

