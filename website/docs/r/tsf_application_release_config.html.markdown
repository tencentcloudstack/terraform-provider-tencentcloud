---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_application_release_config"
sidebar_current: "docs-tencentcloud-resource-tsf_application_release_config"
description: |-
  Provides a resource to create a tsf application_release_config
---

# tencentcloud_tsf_application_release_config

Provides a resource to create a tsf application_release_config

## Example Usage

```hcl
resource "tencentcloud_tsf_application_release_config" "application_release_config" {
  config_id    = "dcfg-nalqbqwv"
  group_id     = "group-yxmz72gv"
  release_desc = "terraform-test"
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required, String, ForceNew) Configuration ID.
* `group_id` - (Required, String, ForceNew) deployment group ID.
* `release_desc` - (Optional, String, ForceNew) release description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `application_id` - Application ID.
* `cluster_id` - cluster ID.
* `cluster_name` - cluster name.
* `config_name` - configuration item name.
* `config_release_id` - configuration item release ID.
* `config_version` - configuration item version.
* `group_name` - deployment group name.
* `namespace_id` - Namespace ID.
* `namespace_name` - namespace name.
* `release_time` - release time.


## Import

tsf application_release_config can be imported using the configId#groupId#configReleaseId, e.g.

```
terraform import tencentcloud_tsf_application_release_config.application_release_config dcfg-nalqbqwv#group-yxmz72gv#dcfgr-maeeq2ea
```

