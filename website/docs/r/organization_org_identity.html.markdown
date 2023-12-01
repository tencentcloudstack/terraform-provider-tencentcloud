---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_identity"
sidebar_current: "docs-tencentcloud-resource-organization_org_identity"
description: |-
  Provides a resource to create a organization org_identity
---

# tencentcloud_organization_org_identity

Provides a resource to create a organization org_identity

## Example Usage

```hcl
resource "tencentcloud_organization_org_identity" "org_identity" {
  identity_alias_name = "example-iac-test"
  identity_policy {
    policy_id   = 1
    policy_name = "AdministratorAccess"
    policy_type = 2
  }
  description = "iac-test"
}
```

## Argument Reference

The following arguments are supported:

* `identity_alias_name` - (Required, String) Identity name.Supports English letters and numbers, the length cannot exceed 40 characters.
* `identity_policy` - (Required, List) Identity policy list.
* `description` - (Optional, String) Identity description.

The `identity_policy` object supports the following:

* `policy_document` - (Optional, String) Customize policy content and follow CAM policy syntax. Valid and required when PolicyType is the 1-custom policy.
* `policy_id` - (Optional, Int) CAM default policy ID. Valid and required when PolicyType is the 2-preset policy.
* `policy_name` - (Optional, String) CAM default policy name. Valid and required when PolicyType is the 2-preset policy.
* `policy_type` - (Optional, Int) Policy type. Value 1-custom policy 2-preset policy; default value 2.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

organization org_identity can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_identity.org_identity org_identity_id
```

