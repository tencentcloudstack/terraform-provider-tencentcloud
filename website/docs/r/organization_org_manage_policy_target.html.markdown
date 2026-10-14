---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_manage_policy_target"
sidebar_current: "docs-tencentcloud-resource-organization_org_manage_policy_target"
description: |-
  Provides a resource to create a organization org_manage_policy_target
---

# tencentcloud_organization_org_manage_policy_target

Provides a resource to create a organization org_manage_policy_target

## Example Usage

```hcl
resource "tencentcloud_organization_org_manage_policy_target" "org_manage_policy_target" {
  target_id   = 10001
  target_type = "NODE"
  policy_id   = 100001
  policy_type = "SERVICE_CONTROL_POLICY"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, Int, ForceNew) Policy Id.
* `target_id` - (Required, Int, ForceNew) Binding target ID of the policy. Member Uin or Department ID.
* `target_type` - (Required, String, ForceNew) Target type.
Valid values:
  - `NODE`: Department.
  - `MEMBER`: Check Member.
* `policy_type` - (Optional, String, ForceNew) Policy type. Default value is SERVICE_CONTROL_POLICY.
Valid values:
  - `SERVICE_CONTROL_POLICY`: Service control policy.
  - `TAG_POLICY`: Tag policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AttachPolicy](https://cloud.tencent.com/document/api/850/105319)
* Read: [ListTargetsForPolicy](https://cloud.tencent.com/document/api/850/105309) - This API is used to query the list of targets associated with a specified policy.
* Delete: [DetachPolicy](https://cloud.tencent.com/document/api/850/105314)

## Import

organization org_manage_policy_target can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_manage_policy_target.org_manage_policy_target policy_type#policy_id#target_type#target_id
```

