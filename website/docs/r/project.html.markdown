---
subcategory: "Project"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_project"
sidebar_current: "docs-tencentcloud-resource-project"
description: |-
  Provides a resource to create a project
---

# tencentcloud_project

Provides a resource to create a project

~> **NOTE:** Project can not be destroyed. If run `terraform destroy`, project will be set invisible.

## Example Usage

```hcl
resource "tencentcloud_project" "project" {
  project_name = "terraform-test"
  info         = "for terraform test"
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required, String) Name of project.
* `disable` - (Optional, Int) If disable project. 1 means disable, 0 means enable. Default 0.
* `info` - (Optional, String) Description of project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `creator_uin` - Uin of creator.


## Import

tag project can be imported using the id, e.g.

```
terraform import tencentcloud_project.project project_id
```

