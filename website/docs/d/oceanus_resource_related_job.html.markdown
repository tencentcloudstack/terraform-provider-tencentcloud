---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_resource_related_job"
sidebar_current: "docs-tencentcloud-datasource-oceanus_resource_related_job"
description: |-
  Use this data source to query detailed information of oceanus resource_related_job
---

# tencentcloud_oceanus_resource_related_job

Use this data source to query detailed information of oceanus resource_related_job

## Example Usage

```hcl
data "tencentcloud_oceanus_resource_related_job" "example" {
  resource_id                    = "resource-8y9lzcuz"
  desc_by_job_config_create_time = 0
  resource_config_version        = 1
  work_space_id                  = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String) Resource ID.
* `desc_by_job_config_create_time` - (Optional, Int) Default:0; 1:sort by job version creation time in descending order.
* `resource_config_version` - (Optional, Int) Resource version number.
* `result_output_file` - (Optional, String) Used to save results.
* `work_space_id` - (Optional, String) Workspace SerialId.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ref_job_infos` - Associated job information.
  * `job_config_version` - Job configuration version.
  * `job_id` - Job ID.
  * `resource_version` - Resource version.


