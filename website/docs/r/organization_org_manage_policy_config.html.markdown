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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [EnablePolicyType](https://cloud.tencent.com/document/api/850/105312) - This API is used to enable a policy type.
* Read: [DescribePolicyConfig](https://cloud.tencent.com/document/api/850/105315) - This API is used to query enterprise organization policy configurations.
* Delete: [DisablePolicyType](https://cloud.tencent.com/document/api/850/105313) - This API is used to disable a policy type.

## Import

organization org_manage_policy_config can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_manage_policy_config.org_manage_policy_config organization_id#policy_type
```

