---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_invite_organization_member_operation"
sidebar_current: "docs-tencentcloud-resource-invite_organization_member_operation"
description: |-
  Provides a resource to create a invite organization member
---

# tencentcloud_invite_organization_member_operation

Provides a resource to create a invite organization member

## Example Usage

```hcl
resource "tencentcloud_invite_organization_member_operation" "example" {
  member_uin     = "100040906211"
  name           = "tf-example"
  policy_type    = "Financial"
  node_id        = 2014419
  is_allow_quit  = "Allow"
  permission_ids = [1, 2, 4]
  remark         = "Remarks."
  tags {
    tag_key   = "CreateBy"
    tag_value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `member_uin` - (Required, Int, ForceNew) Invited account Uin.
* `name` - (Required, String, ForceNew) Member name. The maximum length is 25 characters and supports English letters, numbers, Chinese characters, symbols `+`, `@`, `&`, `.`, `[`, `]`, `-`, `:`, `,` and enumeration comma.
* `node_id` - (Required, Int, ForceNew) Node ID of the member's department.
* `permission_ids` - (Required, Set: [`Int`], ForceNew) List of member financial authority IDs. Values: 1-View bill, 2-View balance, 3-Fund transfer, 4-Consolidated disbursement, 5-Invoice, 6-Benefit inheritance, 7-Proxy payment, 1 and 2 must be default.
* `policy_type` - (Required, String, ForceNew) Relationship strategies. Value taken: Financial.
* `auth_file` - (Optional, List, ForceNew) List of supporting documents of mutual trust entities.
* `is_allow_quit` - (Optional, String, ForceNew) Whether to allow members to withdraw. Allow: Allow, Disallow: Denied.
* `pay_uin` - (Optional, String, ForceNew) Payer Uin. Member needs to pay on behalf of.
* `relation_auth_name` - (Optional, String, ForceNew) Name of the real-name subject of mutual trust.
* `remark` - (Optional, String, ForceNew) Remark.
* `tags` - (Optional, List, ForceNew) List of member tags. Maximum 10.

The `auth_file` object supports the following:

* `name` - (Required, String) File name.
* `url` - (Required, String) File path.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



