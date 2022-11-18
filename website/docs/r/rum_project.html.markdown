---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_project"
sidebar_current: "docs-tencentcloud-resource-rum_project"
description: |-
  Provides a resource to create a rum project
---

# tencentcloud_rum_project

Provides a resource to create a rum project

## Example Usage

```hcl
resource "tencentcloud_rum_project" "project" {
  name             = "projectName"
  instance_id      = "rum-pasZKEI3RLgakj"
  rate             = "100"
  enable_url_group = "0"
  type             = "web"
  repo             = ""
  url              = "iac-tf.com"
  desc             = "projectDesc-1"
}
```

## Argument Reference

The following arguments are supported:

* `enable_url_group` - (Required, Int) Whether to enable aggregation.
* `instance_id` - (Required, String) Business system ID.
* `name` - (Required, String) Name of the created project (required and up to 200 characters).
* `rate` - (Required, String) Project sampling rate (greater than or equal to 0).
* `type` - (Required, String) Project type (valid values: `web`, `mp`, `android`, `ios`, `node`, `hippy`, `weex`, `viola`, `rn`).
* `desc` - (Optional, String) 	Description of the created project (optional and up to 1,000 characters).
* `repo` - (Optional, String) Repository address of the project (optional and up to 256 characters).
* `url` - (Optional, String) Webpage address of the project (optional and up to 256 characters).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creata Time.
* `creator` - Creator ID.
* `instance_key` - Instance key.
* `instance_name` - Instance name.
* `is_star` - Starred status. `1`: yes; `0`: no.
* `key` - Unique project key (12 characters).
* `project_status` - Project status (`1`: Creating; `2`: Running; `3`: Abnormal; `4`: Restarting; `5`: Stopping; `6`: Stopped; `7`: Terminating; `8`: Terminated).


## Import

rum project can be imported using the id, e.g.
```
$ terraform import tencentcloud_rum_project.project project_id
```

