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

The project id is returned by `CreateProject` and used as the resource `id`.
`ModifyProject` only supports updating `name`, `root_dir`, `output_dir`,
`build_cmd`, `install_cmd`, `framework` and `node_js_version`; every other
configurable field (zone_id, channel, provider, source, area, template_id,
repo_*, env_vars, namespace, coding_project_id) forces resource replacement
when changed.

Note: the cloud API echoes `node_js_version` back as `nodejs_version` and only
returns `key`/`value` for `env_vars` entries; the provider preserves the
configured `id`, `comment` and `env` values of each environment variable, so
drift detection for those sub-fields is limited to the local state.
Removing an updatable attribute from the configuration is rejected because the
cloud API does not document whether `ModifyProject` keeps or clears omitted
fields.

## Example Usage

```hcl
resource "tencentcloud_makers_project" "example" {
  zone_id         = "zone-12345"
  name            = "tf-example-project"
  channel         = "Custom"
  provider        = "Github"
  area            = "mainland"
  node_js_version = "20.18.0"
  framework       = "react"
  root_dir        = "./"
  output_dir      = "dist"
  build_cmd       = "npm run build"
  install_cmd     = "npm install"
  repo_url        = "https://github.com/example/tf-example"
  repo_owner      = "example"
  repo_name       = "tf-example"
  repo_branch     = "main"

  env_vars {
    key     = "API_KEY"
    value   = "secret-value"
    comment = "api key for production"
    env     = ["Production"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `channel` - (Required, String) Creation channel. Valid values: `Template`, `Custom`, `CLI`, `Guest`.
* `name` - (Required, String) Project name.
* `zone_id` - (Required, String) Makers zone ID.
* `area` - (Optional, String) Acceleration area. Valid values: `mainland`, `overseas`, `global`.
* `build_cmd` - (Optional, String) Build command.
* `coding_project_id` - (Optional, Int) Coding project ID. Used when `provider` is `Coding`.
* `env_vars` - (Optional, List) Initial environment variables.
* `framework` - (Optional, String) Framework.
* `install_cmd` - (Optional, String) Install command.
* `namespace` - (Optional, String) Coding namespace. Used when `provider` is `Coding`.
* `node_js_version` - (Optional, String) Node.js version, for example `20.18.0`. The value is echoed back by the cloud API as `nodejs_version`.
* `output_dir` - (Optional, String) Build output directory.
* `provider` - (Optional, String) Code source. Valid values: `Upload`, `Github`, `Gitee`, `Gitlab`, `Coding`.
* `repo_branch` - (Optional, String) Production branch.
* `repo_name` - (Optional, String) Repository name.
* `repo_owner` - (Optional, String) Repository owner.
* `repo_url` - (Optional, String) Repository URL.
* `root_dir` - (Optional, String) Project root directory.
* `source` - (Optional, String) Integrator identifier.
* `template_id` - (Optional, String) Template ID.

The `env_vars` object supports the following:

* `key` - (Required, String) Variable name, up to 255 characters.
* `value` - (Required, String) Variable value, up to 1000 characters.
* `comment` - (Optional, String) Remark, up to 255 characters.
* `env` - (Optional, Set) Environment slugs the variable takes effect on, for example `Production`, `Preview`. Multiple rows of the same key cannot overlap environments.
* `id` - (Optional, Int) Variable ID. It is not sent on create.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_on` - Creation time of the project.
* `id` - ID of the resource, the Makers project ID.
* `modified_on` - Last modification time of the project.
* `preset_domain` - Preset domain of the project.
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

Makers project can be imported using the project id, e.g.

```
terraform import tencentcloud_makers_project.example proj-67890
```

