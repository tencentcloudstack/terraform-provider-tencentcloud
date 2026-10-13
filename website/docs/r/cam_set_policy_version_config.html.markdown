---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_set_policy_version_config"
sidebar_current: "docs-tencentcloud-resource-cam_set_policy_version_config"
description: |-
  Provides a resource to create a CAM set policy version config
---

# tencentcloud_cam_set_policy_version_config

Provides a resource to create a CAM set policy version config

## Example Usage

```hcl
resource "tencentcloud_cam_set_policy_version_config" "example" {
  policy_id  = 234290251
  version_id = 3
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, Int) Policy ID.
* `version_id` - (Required, Int) The policy version number, which can be obtained from ListPolicyVersions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [SetDefaultPolicyVersion](https://cloud.tencent.com/document/api/598/43838) - 设置生效的策略版本
* Read: [ListPolicyVersions](https://cloud.tencent.com/document/api/598/43839) - 获取策略版本列表
* Update: [SetDefaultPolicyVersion](https://cloud.tencent.com/document/api/598/43838) - 设置生效的策略版本

## Import

CAM set policy version config can be imported using the id, e.g.

```
terraform import tencentcloud_cam_set_policy_version_config.example 234290251#3
```

