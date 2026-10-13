---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_application_file_config_release"
sidebar_current: "docs-tencentcloud-resource-tsf_application_file_config_release"
description: |-
  Provides a resource to create a tsf application_file_config_release
---

# tencentcloud_tsf_application_file_config_release

Provides a resource to create a tsf application_file_config_release

## Example Usage

```hcl
resource "tencentcloud_tsf_application_file_config_release" "application_file_config_release" {
  config_id    = "dcfg-f-123456"
  group_id     = "group-123456"
  release_desc = "product release"
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required, String, ForceNew) File config id.
* `group_id` - (Required, String, ForceNew) Group Id.
* `release_desc` - (Optional, String, ForceNew) release Description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ReleaseFileConfig](https://cloud.tencent.com/document/api/649/58592) - 发布文件配置
* Read: [DescribeFileConfigReleases](https://cloud.tencent.com/document/api/649/85855) - 查询文件配置项发布信息
* Delete: [DescribeFileConfigReleases](https://cloud.tencent.com/document/api/649/85855) - 查询文件配置项发布信息
* Delete: [RevokeFileConfig](https://cloud.tencent.com/document/api/649/85850) - 撤回已发布的文件配置

## Import

tsf applicationfile_config_release can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_file_config_release.application_file_config_release application_file_config_release_id
```

