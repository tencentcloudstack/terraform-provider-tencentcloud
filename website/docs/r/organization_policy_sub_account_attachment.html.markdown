---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_policy_sub_account_attachment"
sidebar_current: "docs-tencentcloud-resource-organization_policy_sub_account_attachment"
description: |-
  Provides a resource to create a organization policy_sub_account_attachment
---

# tencentcloud_organization_policy_sub_account_attachment

Provides a resource to create a organization policy_sub_account_attachment

## Example Usage

```hcl
resource "tencentcloud_organization_policy_sub_account_attachment" "policy_sub_account_attachment" {
  member_uin          = 100028582828
  org_sub_account_uin = 100028223737
  policy_id           = 144256499
}
```

## Argument Reference

The following arguments are supported:

* `member_uin` - (Required, Int, ForceNew) Organization member uin.
* `org_sub_account_uin` - (Required, Int, ForceNew) Organization administrator sub account uin list.
* `policy_id` - (Required, Int, ForceNew) Policy ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time.
* `identity_id` - Manage Identity ID.
* `identity_role_alias_name` - Identity role alias name.
* `identity_role_name` - Identity role name.
* `org_sub_account_name` - Organization administrator sub account name.
* `policy_name` - Policy name.
* `update_time` - Update time.


## Import

organization policy_sub_account_attachment can be imported using the id, e.g.
```
$ terraform import tencentcloud_organization_policy_sub_account_attachment.policy_sub_account_attachment policyId#memberUin#orgSubAccountUin
```

