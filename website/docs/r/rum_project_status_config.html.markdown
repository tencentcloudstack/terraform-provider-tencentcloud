---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_project_status_config"
sidebar_current: "docs-tencentcloud-resource-rum_project_status_config"
description: |-
  Provides a resource to create a rum project_status_config
---

# tencentcloud_rum_project_status_config

Provides a resource to create a rum project_status_config

## Example Usage

```hcl
resource "tencentcloud_rum_project_status_config" "project_status_config" {
  project_id = 131407
  operate    = "stop"
}
```

## Argument Reference

The following arguments are supported:

* `operate` - (Required, String) `resume`, `stop`.
* `project_id` - (Required, Int, ForceNew) Project ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Read: [DescribeProjects](https://www.tencentcloud.com/document/api/1131/45308)
* Update:
  - [ResumeProject](https://www.tencentcloud.com/document/api/1131/53811)
  - [StopProject](https://www.tencentcloud.com/document/api/1131/53261)

## Import

rum project_status_config can be imported using the id, e.g.

```
terraform import tencentcloud_rum_project_status_config.project_status_config project_id
```

