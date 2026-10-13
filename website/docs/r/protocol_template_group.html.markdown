---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_protocol_template_group"
sidebar_current: "docs-tencentcloud-resource-protocol_template_group"
description: |-
  Provides a resource to manage protocol template group.
---

# tencentcloud_protocol_template_group

Provides a resource to manage protocol template group.

## Example Usage

```hcl
resource "tencentcloud_protocol_template_group" "foo" {
  name         = "group-test"
  template_ids = ["ipl-axaf24151", "ipl-axaf24152"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the protocol template group.
* `template_ids` - (Required, Set: [`String`]) Service template ID list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateServiceTemplateGroup](https://cloud.tencent.com/document/api/215/16711) - 创建协议端口模板集合
* Read: [DescribeServiceTemplateGroups](https://cloud.tencent.com/document/api/215/16718) - 查询协议端口模板集合
* Update: [ModifyServiceTemplateGroupAttribute](https://cloud.tencent.com/document/api/215/16723) - 修改协议端口模板集合
* Delete: [DeleteServiceTemplateGroup](https://cloud.tencent.com/document/api/215/16715) - 删除协议端口模板集合
* Delete: [DescribeServiceTemplateGroups](https://cloud.tencent.com/document/api/215/16718) - 查询协议端口模板集合

## Import

Protocol template group can be imported using the id, e.g.

```
$ terraform import tencentcloud_protocol_template_group.foo ppmg-0np3u974
```

