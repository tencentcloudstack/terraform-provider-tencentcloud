---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_logset"
sidebar_current: "docs-tencentcloud-resource-cls_logset"
description: |-
  Provides a resource to create a CLS logset
---

# tencentcloud_cls_logset

Provides a resource to create a CLS logset

## Example Usage

```hcl
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf-example"
  tag_list {
    key   = "createdBy"
    value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `logset_name` - (Required, String) Logset name, which must be unique.
* `tag_list` - (Optional, List) Tag description list. The CLS API supports up to 10 tag key-value pairs, and duplicate keys are not allowed.
* `tags` - (Optional, Map, **Deprecated**) It is recommended to use `tag_list` because the current `tags` field is binding resources by calling the tag API. Tag description list.

The `tag_list` object supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time.
* `role_name` - If assumer_uin is not empty, it indicates the service provider who creates the logset.
* `topic_count` - Number of log topics in logset.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateLogset](https://cloud.tencent.com/document/api/614/58626) - 创建日志集
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Read: [DescribeLogsets](https://cloud.tencent.com/document/api/614/58624) - 获取日志集列表
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - 批量查看资源关联的标签
* Update: [ModifyLogset](https://cloud.tencent.com/document/api/614/58623) - 修改日志集
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Delete: [DeleteLogset](https://cloud.tencent.com/document/api/614/58625) - 删除日志集

## Import

CLS logset can be imported using the id, e.g.
```
$ terraform import tencentcloud_cls_logset.example 698902ff-8b5a-4c65-824b-d8956f366351
```

