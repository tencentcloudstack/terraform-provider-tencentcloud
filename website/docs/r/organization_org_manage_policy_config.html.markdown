---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_manage_policy_config"
sidebar_current: "docs-tencentcloud-resource-organization_org_manage_policy_config"
description: |-
  Provides a resource to create a organization org_manage_policy_config
---

# tencentcloud_organization_org_manage_policy_config

Provides a resource to create a organization org_manage_policy_config

## Example Usage

```hcl
resource "tencentcloud_organization_org_manage_policy_config" "org_manage_policy_config" {
  organization_id = 80001
  policy_type     = "SERVICE_CONTROL_POLICY"
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (Required, Int, ForceNew) Organization ID.
* `policy_type` - (Optional, String, ForceNew) Policy type. Default value is SERVICE_CONTROL_POLICY.
Valid values:
  - `SERVICE_CONTROL_POLICY`: Service control policy.
  - `TAG_POLICY`: Tag policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

organization org_manage_policy_config can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_manage_policy_config.org_manage_policy_config organization_id#policy_type
```

