---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_role_configuration_provisionings"
sidebar_current: "docs-tencentcloud-datasource-role_configuration_provisionings"
description: |-
  Use this data source to query detailed information of organization role_configuration_provisionings
---

# tencentcloud_role_configuration_provisionings

Use this data source to query detailed information of organization role_configuration_provisionings

## Example Usage

```hcl
data "tencentcloud_role_configuration_provisionings" "role_configuration_provisionings" {
  zone_id               = "xxxxxx"
  role_configuration_id = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Space ID.
* `deployment_status` - (Optional, String) Deployed: Deployment succeeded; DeployedRequired: Redeployment required; DeployFailed: Deployment failed.
* `filter` - (Optional, String) Search by configuration name is supported.
* `result_output_file` - (Optional, String) Used to save results.
* `role_configuration_id` - (Optional, String) Permission configuration ID.
* `target_type` - (Optional, String) Type of the synchronized target account of the Tencent Cloud Organization. ManagerUin: admin account; MemberUin: member account.
* `target_uin` - (Optional, Int) UIN of the synchronized target account of the Tencent Cloud Organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `role_configuration_provisionings` - Department member account list.


