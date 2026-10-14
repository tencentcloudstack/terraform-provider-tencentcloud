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

* Create: [CreateLogset](https://www.tencentcloud.com/document/product/614/42780) - This API is used to create a logset. The ID of the created logset is returned.
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Read: [DescribeLogsets](https://www.tencentcloud.com/document/product/614/42778) - This API is used to get the list of logsets.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [ModifyLogset](https://www.tencentcloud.com/document/product/614/42777) - This API is used to modify a logset.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteLogset](https://www.tencentcloud.com/document/product/614/42779) - This API is used to delete a logset.

## Import

CLS logset can be imported using the id, e.g.
```
$ terraform import tencentcloud_cls_logset.example 698902ff-8b5a-4c65-824b-d8956f366351
```

