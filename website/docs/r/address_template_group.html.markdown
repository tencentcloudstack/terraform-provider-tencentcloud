---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_address_template_group"
sidebar_current: "docs-tencentcloud-resource-address_template_group"
description: |-
  Provides a resource to manage address template group.
---

# tencentcloud_address_template_group

Provides a resource to manage address template group.

## Example Usage

```hcl
resource "tencentcloud_address_template_group" "foo" {
  name         = "group-test"
  template_ids = ["ipl-axaf24151", "ipl-axaf24152"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the address template group.
* `template_ids` - (Required, Set: [`String`]) Template ID list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAddressTemplateGroup](https://cloud.tencent.com/document/api/215/16709) - 创建IP地址模板集合
* Read: [DescribeAddressTemplateGroups](https://cloud.tencent.com/document/api/215/16716) - 查询IP地址模板集合
* Update: [ModifyAddressTemplateGroupAttribute](https://cloud.tencent.com/document/api/215/16721) - 修改IP地址模板集合
* Delete: [DeleteAddressTemplateGroup](https://cloud.tencent.com/document/api/215/16713) - 删除IP地址模板集合
* Delete: [DescribeAddressTemplateGroups](https://cloud.tencent.com/document/api/215/16716) - 查询IP地址模板集合

## Import

Address template group can be imported using the id, e.g.

```
$ terraform import tencentcloud_address_template_group.foo ipmg-0np3u974
```

