---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_node"
sidebar_current: "docs-tencentcloud-resource-organization_org_node"
description: |-
  Provides a resource to create a organization org_node
---

# tencentcloud_organization_org_node

Provides a resource to create a organization org_node

## Example Usage

```hcl
resource "tencentcloud_organization_org_node" "org_node" {
  name           = "terraform_test"
  parent_node_id = 2003721
  remark         = "for terraform test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Node name.
* `parent_node_id` - (Required, Int) Parent node ID.
* `remark` - (Optional, String) Notes.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Node creation time.
* `update_time` - Node update time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AddOrganizationNode](https://cloud.tencent.com/document/api/850/82930) - 添加企业组织节点
* Read: [DescribeOrganizationNodes](https://cloud.tencent.com/document/api/850/82926) - 获取组织节点列表
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Update: [UpdateOrganizationNode](https://cloud.tencent.com/document/api/850/82924) - 更新企业组织节点
* Delete: [DeleteOrganizationNodes](https://cloud.tencent.com/document/api/850/82928) - 批量删除企业组织节点

## Import

organization org_node can be imported using the id, e.g.
```
$ terraform import tencentcloud_organization_org_node.org_node orgNode_id
```

