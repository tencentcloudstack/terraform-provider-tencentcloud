---
subcategory: "DNSPod"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_record_group"
sidebar_current: "docs-tencentcloud-resource-dnspod_record_group"
description: |-
  Provides a resource to create a dnspod record_group
---

# tencentcloud_dnspod_record_group

Provides a resource to create a dnspod record_group

## Example Usage

```hcl
resource "tencentcloud_dnspod_record_group" "record_group" {
  domain     = "dnspod.cn"
  group_name = "group_demo"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain.
* `group_name` - (Required, String) Record Group Name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `group_id` - Group ID.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateRecordGroup](https://cloud.tencent.com/document/api/1427/83227) - 添加记录分组
* Read: [DescribeRecordGroupList](https://cloud.tencent.com/document/api/1427/83225) - 查询解析记录分组列表
* Update: [ModifyRecordGroup](https://cloud.tencent.com/document/api/1427/83224) - 修改记录分组
* Delete: [DeleteRecordGroup](https://cloud.tencent.com/document/api/1427/83226) - 删除记录分组

## Import

dnspod record_group can be imported using the domain#groupId, e.g.

```
terraform import tencentcloud_dnspod_record_group.record_group domain#groupId
```

