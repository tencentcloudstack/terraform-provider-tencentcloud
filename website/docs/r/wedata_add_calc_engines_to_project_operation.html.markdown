---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_add_calc_engines_to_project_operation"
sidebar_current: "docs-tencentcloud-resource-wedata_add_calc_engines_to_project_operation"
description: |-
  Provides a resource to create a WeData add calc engines to project operation
---

# tencentcloud_wedata_add_calc_engines_to_project_operation

Provides a resource to create a WeData add calc engines to project operation

## Example Usage

```hcl
resource "tencentcloud_wedata_add_calc_engines_to_project_operation" "example" {
  project_id = "20241107221758402"
  dlc_info {
    compute_resources = [
      "dlc_linau6d4bu8bd5u52ffu52a8"
    ]
    region           = "ap-guangzhou"
    default_database = "default_db"
  }
}
```

## Argument Reference

The following arguments are supported:

* `dlc_info` - (Required, List, ForceNew) DLC cluster information.
* `project_id` - (Required, String, ForceNew) Project ID to be modified.

The `dlc_info` object supports the following:

* `compute_resources` - (Required, Set) DLC resource names (need to add role Uin to DLC, otherwise resources may not be available).
* `default_database` - (Required, String) Specify the default database for the DLC cluster.
* `region` - (Required, String) DLC region.
* `access_account` - (Optional, String) Access account (only effective for standard mode projects and required for standard mode), used to submit DLC tasks.
It is recommended to use a specified sub-account and set corresponding database table permissions for the sub-account; task runner mode may cause task failures when the responsible person leaves; main account mode is not easy for permission control when multiple projects have different permissions.

Enum values:
- TASK_RUNNER (Task Runner)
- OWNER (Main Account Mode)
- SUB (Sub-Account Mode).
* `standard_mode_env_tag` - (Optional, String) Cluster configuration tag (only effective for standard mode projects and required for standard mode). Enum values:
- Prod  (Production environment)
- Dev  (Development environment).
* `sub_account_uin` - (Optional, String) Sub-account ID (only effective for standard mode projects), when AccessAccount is in sub-account mode, the sub-account ID information needs to be specified, other modes do not need to be specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



