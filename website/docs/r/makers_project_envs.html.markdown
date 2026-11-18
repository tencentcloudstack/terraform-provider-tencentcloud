---
subcategory: "Makers"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_makers_project_envs"
sidebar_current: "docs-tencentcloud-resource-makers_project_envs"
description: |-
  Provides a resource to manage environment variables of a Makers (TEO makers) project.
---

# tencentcloud_makers_project_envs

Provides a resource to manage environment variables of a Makers (TEO makers) project.

~> **NOTE** The cloud API `ModifyProjectEnvs` writes the full desired set of environment
variables for the project, and `DeleteProjectEnvs` deletes variables by key. Variables
removed from the `env_vars` configuration are deleted explicitly on update; variables
modified outside of Terraform are detected only through the read-back of the variable
list. The `zone_id` and `project_id` locating fields are not echoed back by the describe
API, so drift detection covers the variable list only.

## Example Usage

```hcl
resource "tencentcloud_makers_project_envs" "example" {
  zone_id    = "zone-xxxx"
  project_id = "proj-yyyy"
  env_vars {
    key     = "API_KEY"
    value   = "secret-value"
    comment = "production key"
    env     = ["Production"]
  }
  env_vars {
    key   = "DEBUG"
    value = jsonencode({ level = "info" })
  }
}
```

## Argument Reference

The following arguments are supported:

* `env_vars` - (Required, List) The environment variables to write.
* `project_id` - (Required, String) Project ID.
* `zone_id` - (Required, String) Makers site ID.

The `env_vars` object supports the following:

* `key` - (Required, String) Variable name, at most 255 characters.
* `value` - (Required, String) Variable value, at most 1000 characters.
* `comment` - (Optional, String) Comment, at most 255 characters.
* `env` - (Optional, Set) Slugs of the environments the variable applies to, e.g. `Production`, `Preview`. Env sets of variables sharing the same key must not overlap.
* `id` - (Optional, Int) Variable ID. Not required when creating.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource in the form `<zone_id>#<project_id>`.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [DescribeProjectEnvs](https://www.tencentcloud.com/search/?wd=DescribeProjectEnvs&type=document&lang=en)
  - [ModifyProjectEnvs](https://www.tencentcloud.com/search/?wd=ModifyProjectEnvs&type=document&lang=en)
* Read: [DescribeProjectEnvs](https://www.tencentcloud.com/search/?wd=DescribeProjectEnvs&type=document&lang=en)
* Update:
  - [DeleteProjectEnvs](https://www.tencentcloud.com/search/?wd=DeleteProjectEnvs&type=document&lang=en)
  - [DescribeProjectEnvs](https://www.tencentcloud.com/search/?wd=DescribeProjectEnvs&type=document&lang=en)
  - [ModifyProjectEnvs](https://www.tencentcloud.com/search/?wd=ModifyProjectEnvs&type=document&lang=en)
* Delete: [DeleteProjectEnvs](https://www.tencentcloud.com/search/?wd=DeleteProjectEnvs&type=document&lang=en)

## Import

tencentcloud_makers_project_envs can be imported using the `<zone_id>#<project_id>`, e.g.

```
$ terraform import tencentcloud_makers_project_envs.example zone-xxxx#proj-yyyy
```

