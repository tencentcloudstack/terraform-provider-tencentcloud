---
subcategory: "Web & Agents Development Platform(Makers)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_makers_project"
sidebar_current: "docs-tencentcloud-resource-makers_project"
description: |-
  Provides a resource to create a Makers project.
---

# tencentcloud_makers_project

Provides a resource to create a Makers project.

## Example Usage

```hcl
resource "tencentcloud_makers_project" "example" {
  zone_id         = "zone-3hzk0ndvnnyp"
  name            = "tf-test5"
  channel         = "Custom"
  area            = "global"
  code_provider   = "Upload"
  node_js_version = "22.11.0"
}
```

## Argument Reference

The following arguments are supported:

* `channel` - (Required, String) Creation channel. Valid values: `Template`, `Custom`, `CLI`, `Guest`.
* `name` - (Required, String) Project name.
* `zone_id` - (Required, String) Makers zone ID.
* `area` - (Optional, String) Acceleration area. Valid values: `mainland`, `overseas`, `global`.
* `build_cmd` - (Optional, String) Build command.
* `code_provider` - (Optional, String) Code source. Valid values: `Upload`, `Github`, `Gitee`, `Gitlab`, `Coding`.
* `coding_project_id` - (Optional, Int) Coding project ID. Used when `provider` is `Coding`.
* `framework` - (Optional, String) Framework.
* `install_cmd` - (Optional, String) Install command.
* `namespace` - (Optional, String) Coding namespace. Used when `provider` is `Coding`.
* `node_js_version` - (Optional, String) Node.js version, for example `20.18.0`. The value is echoed back by the cloud API as `nodejs_version`.
* `output_dir` - (Optional, String) Build output directory.
* `repo_branch` - (Optional, String) Production branch.
* `repo_name` - (Optional, String) Repository name.
* `repo_owner` - (Optional, String) Repository owner.
* `repo_url` - (Optional, String) Repository URL.
* `root_dir` - (Optional, String) Project root directory.
* `source` - (Optional, String) Integrator identifier.
* `template_id` - (Optional, String) Template ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_on` - Creation time of the project.
* `id` - ID of the resource, the Makers project ID.
* `modified_on` - Last modification time of the project.
* `preset_domain` - Preset domain of the project.
* `project_id` - Project ID.
* `status` - Project status. Valid values: `Normal`, `Ban`, `Unlink`, `Process`, `Pending`, `Offline`.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateProject](https://www.tencentcloud.com/search/?wd=CreateProject&type=document&lang=en)
  - [DescribeProjects](https://www.tencentcloud.com/search/?wd=DescribeProjects&type=document&lang=en)
* Read: [DescribeProjects](https://www.tencentcloud.com/search/?wd=DescribeProjects&type=document&lang=en)
* Update:
  - [DescribeProjects](https://www.tencentcloud.com/search/?wd=DescribeProjects&type=document&lang=en)
  - [ModifyProject](https://www.tencentcloud.com/search/?wd=ModifyProject&type=document&lang=en)
* Delete: [DeleteProject](https://www.tencentcloud.com/search/?wd=DeleteProject&type=document&lang=en)

## Import

Makers project can be imported using the composite id `zone_id#project_id`, e.g.

```
terraform import tencentcloud_makers_project.example zone-12345#proj-67890
```

