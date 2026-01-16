---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_deploy_config_group_version"
sidebar_current: "docs-tencentcloud-resource-teo_deploy_config_group_version"
description: |-
  Provides a resource to create a teo deploy config group version
---

# tencentcloud_teo_deploy_config_group_version

Provides a resource to create a teo deploy config group version

## Example Usage

```hcl
resource "tencentcloud_teo_deploy_config_group_version" "teo_deploy_config_group_version" {
  zone_id     = "zone-2xkazzl8yf6k"
  env_id      = "env-3lchxiq1h855"
  description = "Deploy config group version for production"
  # l7_acceleration
  config_group_version_infos {
    version_id = "ver-3lchxizh2mqn"
  }
  # edge_functions
  config_group_version_infos {
    version_id = "ver-3lchxjdciuzx"
  }
}
```

## Argument Reference

The following arguments are supported:

* `config_group_version_infos` - (Required, Set, ForceNew) Version information required for release. Multiple versions of different configuration groups can be modified simultaneously, while each group allows modifying only one version at a time.
* `description` - (Required, String, ForceNew) Change description. It is used to describe the content and reasons for this change. A maximum of 100 characters are supported.
* `env_id` - (Required, String, ForceNew) Environment ID. Please specify the environment ID to which the version should be released.
* `zone_id` - (Required, String, ForceNew) Zone ID.

The `config_group_version_infos` object supports the following:

* `version_id` - (Required, String, ForceNew) Version ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `deploy_time` - Deploy time. The time follows the ISO 8601 standard in the date and time format.
* `message` - Deploy result message.
* `record_id` - Deploy record ID.
* `status` - Deploy status. Valid values: deploying (Deploying), failure (Deploy failed), success (Deploy successful).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.

