---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_project"
sidebar_current: "docs-tencentcloud-datasource-rum_project"
description: |-
  Use this data source to query detailed information of rum project
---

# tencentcloud_rum_project

Use this data source to query detailed information of rum project

## Example Usage

```hcl
data "tencentcloud_rum_project" "project" {
  instance_id = "rum-pasZKEI3RLgakj"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `project_set` - Project list.
  * `create_time` - CreateTime.
  * `creator` - Creator ID.
  * `desc` - Project description.
  * `enable_url_group` - Whether to enable URL aggregation.
  * `instance_id` - Instance ID.
  * `instance_key` - Instance key.
  * `instance_name` - Instance name.
  * `is_star` - Starred status. `1`: yes; `0`: no.
  * `key` - Unique project key (12 characters).
  * `name` - Project name.
  * `pid` - Project ID.
  * `project_status` - Project status (`1`: Creating; `2`: Running; `3`: Abnormal; `4`: Restarting; `5`: Stopping; `6`: Stopped; `7`: Terminating; `8`: Terminated).
  * `rate` - Project sample rate.
  * `repo` - Project repository address.
  * `type` - Project type.
  * `url` - Project URL.


