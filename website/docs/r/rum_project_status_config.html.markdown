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

* Read: [DescribeProjects](https://cloud.tencent.com/document/api/1464/61475)
* Update: [ResumeProject](https://cloud.tencent.com/document/api/1464/86367)
* Update: [StopProject](https://cloud.tencent.com/document/api/1464/84458)

## Import

rum project_status_config can be imported using the id, e.g.

```
terraform import tencentcloud_rum_project_status_config.project_status_config project_id
```

