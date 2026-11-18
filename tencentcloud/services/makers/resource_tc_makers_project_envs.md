subcategory: ""
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_makers_project_envs"
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

## Import

tencentcloud_makers_project_envs can be imported using the `<zone_id>#<project_id>`, e.g.

```
$ terraform import tencentcloud_makers_project_envs.example zone-xxxx#proj-yyyy
```
