---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_resource_group_to_project_attachment"
sidebar_current: "docs-tencentcloud-resource-wedata_resource_group_to_project_attachment"
description: |-
  Provides a resource to create a WeData resource group to project attachment
---

# tencentcloud_wedata_resource_group_to_project_attachment

Provides a resource to create a WeData resource group to project attachment

## Example Usage

```hcl
resource "tencentcloud_wedata_resource_group_to_project_attachment" "example" {
  resource_group_id = "20250909161820129828"
  project_id        = "2983848457986924544"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String, ForceNew) Project ID.
* `resource_group_id` - (Required, String, ForceNew) Resource group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AssociateResourceGroupToProject](https://cloud.tencent.com/document/api/1267/123985) - 关联资源组到某个项目
* Read: [ListResourceGroups](https://cloud.tencent.com/document/api/1267/123980) - 查看资源组列表
* Delete: [DissociateResourceGroupFromProject](https://cloud.tencent.com/document/api/1267/123982) - 解绑项目关联的资源组

## Import

WeData resource group to project attachment can be imported using the resourceGroupId#projectId, e.g.

```
terraform import tencentcloud_wedata_resource_group_to_project_attachment.example 20250909161820129828#2983848457986924544
```

