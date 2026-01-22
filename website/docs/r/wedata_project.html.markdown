---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_project"
sidebar_current: "docs-tencentcloud-resource-wedata_project"
description: |-
  Provides a resource to create a WeData project
---

# tencentcloud_wedata_project

Provides a resource to create a WeData project

## Example Usage

```hcl
resource "tencentcloud_wedata_project" "example" {
  project {
    project_name  = "tf_example"
    display_name  = "display_name"
    project_model = "SIMPLE"
  }

  dlc_info {
    compute_resources     = ["svmgao_stability"]
    region                = "ap-guangzhou"
    default_database      = "db_name"
    standard_mode_env_tag = "Dev"
    access_account        = "OWNER"
  }

  resource_ids = [
    "20250909193110713075",
    "20250820215449817917"
  ]

  status = 1
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, List) Project basic information.
* `dlc_info` - (Optional, List) DLC binding cluster information.
* `resource_ids` - (Optional, Set: [`String`]) List of bound resource group IDs.
* `status` - (Optional, Int) Item status: 0: disabled, 1: enabled.

The `dlc_info` object supports the following:

* `compute_resources` - (Required, Set) DLC resource name (need to add role Uin to DLC, otherwise may not be able to obtain resources).
* `default_database` - (Required, String) Specify the default database for DLC cluster.
* `region` - (Required, String) DLC region.
* `access_account` - (Optional, String) Access account (only effective for standard mode projects and required for standard mode), used to submit DLC tasks.
It is recommended to use a specified sub-account and set corresponding database table permissions for the sub-account; task runner mode may cause task failure when the responsible person leaves; main account mode is not easy for permission control when multiple projects have different permissions.

Enum values:
- TASK_RUNNER (Task Runner)
- OWNER (Main Account Mode)
- SUB (Sub Account Mode).
* `standard_mode_env_tag` - (Optional, String) Cluster configuration tag (only effective for standard mode projects and required for standard mode). Enum values:
- Prod  (Production environment)
- Dev  (Development environment).
* `sub_account_uin` - (Optional, String) Sub-account ID (only effective for standard mode projects), when AccessAccount is in sub-account mode, the sub-account ID information needs to be specified, other modes do not need to be specified.

The `project` object supports the following:

* `display_name` - (Required, String, ForceNew) Project display name, can be Chinese name starting with a letter, can contain letters, numbers, and underscores, cannot exceed 32 characters.
* `project_name` - (Required, String) Project identifier, English name starting with a letter, can contain letters, numbers, and underscores, cannot exceed 32 characters.
* `project_model` - (Optional, String, ForceNew) Project mode, SIMPLE (default): Simple mode STANDARD: Standard mode.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `project_id` - Project ID.


