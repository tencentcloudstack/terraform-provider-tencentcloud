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

~> **NOTE** If you need to adjust the order of the `env_vars` entries, you must delete and recreate them.

## Example Usage

```hcl
resource "tencentcloud_makers_project_envs" "example" {
  zone_id    = "zone-3hzk0ndvnnyp"
  project_id = "makers-jw7jha6dktob"

  env_vars = [
    # Variables specific to the production environment
    {
      key     = "API_KEY"
      value   = "secret-value"
      comment = "production api key"
      env     = ["Production"]
    },

    # Variables specific to the preview environment (sets of environment variables with the same key cannot overlap)
    {
      key     = "API_KEY"
      value   = "preview-key"
      comment = "preview api key"
      env     = ["Preview"]
    },

    # All environment variables
    {
      key     = "DEBUG"
      value   = "1"
      comment = "debug switch"
      env     = ["Production", "Preview"]
    },
    {
      key   = "DEBUG1"
      value = "1"
    }
  ]
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
* `id` - (Optional, Int) Variable ID assigned by the cloud side. Not required when creating.

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
  - [DescribeProjectEnvs](https://www.tencentcloud.com/search/?wd=DescribeProjectEnvs&type=document&lang=en)
  - [ModifyProjectEnvs](https://www.tencentcloud.com/search/?wd=ModifyProjectEnvs&type=document&lang=en)
* Delete: [DeleteProjectEnvs](https://www.tencentcloud.com/search/?wd=DeleteProjectEnvs&type=document&lang=en)

## Import

tencentcloud_makers_project_envs can be imported using the `<zone_id>#<project_id>`, e.g.

```
$ terraform import tencentcloud_makers_project_envs.example zone-xxxx#proj-yyyy
```

