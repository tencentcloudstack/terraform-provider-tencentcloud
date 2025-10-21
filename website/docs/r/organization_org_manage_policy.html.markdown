---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_manage_policy"
sidebar_current: "docs-tencentcloud-resource-organization_org_manage_policy"
description: |-
  Provides a resource to create a organization org_manage_policy
---

# tencentcloud_organization_org_manage_policy

Provides a resource to create a organization org_manage_policy

## Example Usage

```hcl
resource "tencentcloud_organization_org_manage_policy" "org_manage_policy" {
  name        = "FullAccessPolicy"
  content     = "{\"version\":\"2.0\",\"statement\":[{\"effect\":\"allow\",\"action\":\"*\",\"resource\":\"*\"}]}"
  type        = "SERVICE_CONTROL_POLICY"
  description = "Full access policy"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Policy content. Refer to the CAM policy syntax.
* `name` - (Required, String) Policy name.
The length is 1~128 characters, which can include Chinese characters, English letters, numbers, and underscores.
* `description` - (Optional, String) Policy description.
* `type` - (Optional, String) Policy type. Default value is SERVICE_CONTROL_POLICY.
Valid values:
  - `SERVICE_CONTROL_POLICY`: Service control policy.
  - `TAG_POLICY`: Tag policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `policy_id` - Policy Id.


## Import

organization org_manage_policy can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_manage_policy.org_manage_policy policy_id#type
```

