---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_member_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-organization_org_member_policy_attachment"
description: |-
  Provides a resource to create a organization org_member_policy_attachment
---

# tencentcloud_organization_org_member_policy_attachment

Provides a resource to create a organization org_member_policy_attachment

## Example Usage

```hcl
resource "tencentcloud_organization_org_member_policy_attachment" "org_member_policy_attachment" {
  member_uins = [100033905366, 100033905356]
  policy_name = "example-iac"
  identity_id = 1
}
```

## Argument Reference

The following arguments are supported:

* `identity_id` - (Required, Int, ForceNew) Organization identity ID.
* `member_uins` - (Required, Set: [`Int`], ForceNew) Member Uin list. Up to 10.
* `policy_name` - (Required, String, ForceNew) Policy name.The maximum length is 128 characters, supporting English letters, numbers, and symbols +=,.@_-.
* `description` - (Optional, String, ForceNew) Notes.The maximum length is 128 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

organization org_member_policy_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_member_policy_attachment.org_member_policy_attachment org_member_policy_attachment_id
```

