---
subcategory: "Web & Agents Development Platform(Makers)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_makers_deployments"
sidebar_current: "docs-tencentcloud-data_source-makers_deployments"
description: |-
  Use this data source to query the list of Makers deployment records (build history) of a site project.
---

# tencentcloud_makers_deployments

Use this data source to query the list of Makers deployment records (build history) of a site project.

## Example Usage

### Query all deployment records of a Makers project

```hcl
data "tencentcloud_makers_deployments" "example" {
  zone_id    = "zone-abcdefg"
  project_id = "prj-123456"
}
```

### Query deployment records with filters and sorting

```hcl
data "tencentcloud_makers_deployments" "example" {
  zone_id    = "zone-abcdefg"
  project_id = "prj-123456"
  filters {
    name   = "Status"
    values = ["Success"]
  }
  order_by   = "CreatedOn"
  order_type = "Desc"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `zone_id` - (Required, String) Makers site ID.
* `deployment_ids` - (Optional, List) List of deployment IDs to query.
* `end_time` - (Optional, String) End time.
* `filters` - (Optional, List) Filter conditions. Supported filter names: ProjectId, DeploymentId, Name, Status, StartTime, EndTime, RepoBranch.
* `order_by` - (Optional, String) Sort field. Valid values: CreatedOn, ModifiedOn.
* `order_type` - (Optional, String) Sort type. Valid values: Asc, Desc.
* `start_time` - (Optional, String) Start time.

The `filters` object supports the following:

* `name` - (Required, String) Filter field. Valid values: ProjectId, DeploymentId, Name, Status, StartTime, EndTime, RepoBranch.
* `values` - (Required, List) Filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `build_cost` - Build cost time, in milliseconds.
  * `code` - Build exit code, 0 means success.
  * `cover_url` - Cover image URL.
  * `created_on` - Creation time.
  * `deployment_id` - Deployment ID.
  * `dpl_repo_url` - Repository URL used by this deployment.
  * `env` - Environment slug.
  * `meta_data` - Route configuration, JSON string.
  * `modified_on` - Modification time.
  * `original_status` - Original status of the deployment before it was removed.
  * `preview_url` - Access URL of this deployment.
  * `project_id` - Project ID.
  * `project_url` - Default domain URL of the project.
  * `repo_branch` - Repository branch.
  * `repo_commit_hash` - Commit hash.
  * `repo_commit_msg` - Commit message.
  * `status` - Deployment status. Valid values: Success, Failed, Process, Invalid, Pending, Timeout, Cancelled.
  * `used_in_prod` - Whether this is the current production version.
  * `via_meta` - Trigger source. Valid values: Github, ReDeploy, Api, Coding, Gitee, Upload, Gitlab.
* `deployments` - Deployment list.
* `total_count` - Total number of deployment records.


