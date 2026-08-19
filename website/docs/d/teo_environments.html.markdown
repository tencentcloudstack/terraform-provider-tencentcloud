---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_environments"
sidebar_current: "docs-tencentcloud-datasource-teo_environments"
description: |-
  Use this data source to query detailed information of TEO environments, including the environment list and the source version (`source_version`) of each effective config group version.
---

# tencentcloud_teo_environments

Use this data source to query detailed information of TEO environments, including the environment list and the source version (`source_version`) of each effective config group version.

## Example Usage

```hcl
data "tencentcloud_teo_environments" "teo_environments" {
  zone_id = "zone-2qtuhspy7cr6"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Zone ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `env_infos` - Environment list.
  * `current_config_group_version_infos` - For the effective versions of each configuration group in the current environment, there are two possible scenarios based on the value of Status: When Status is set to version_deploying, the returned value of this field represents the previously effective version. In other words, during the deployment of the new version, the effective version is the one that was in effect before any changes were made. When Status is set to running, the value returned by this field is the currently effective version.
    * `create_time` - Version creation time. The time format follows the ISO 8601 standard and is represented in Coordinated Universal Time (UTC).
    * `description` - Version description.
    * `group_id` - Configuraration group ID.
    * `group_type` - Configuration group type. Valid values: l7_acceleration (L7 acceleration configuration group), edge_functions (Edge function configuration group).
    * `source_version` - The source version ID that the config group version was derived from.
    * `status` - Version status. Valid values: creating (Being created), inactive (Not effective), active (Effective).
    * `version_id` - Version ID.
    * `version_number` - Version No.
  * `env_id` - Environment ID.
  * `env_type` - Environment type. Valid values: production (Production environment), staging (Test environment).
  * `scope` - Effective scope of the configuration in the current environment. Valid values: ALL (It takes effect on the entire network when EnvType is set to production), It returns the IP address of the test node for host binding during testing when EnvType is set to staging.
  * `status` - Environment status. Valid values: creating (Being created), running (The environment is stable, with version changes allowed), version_deploying (The version is currently being deployed, with no more changes allowed).


