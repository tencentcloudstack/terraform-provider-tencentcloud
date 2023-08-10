---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_data_transform"
sidebar_current: "docs-tencentcloud-resource-cls_data_transform"
description: |-
  Provides a resource to create a cls data_transform
---

# tencentcloud_cls_data_transform

Provides a resource to create a cls data_transform

## Example Usage

```hcl
resource "tencentcloud_cls_logset" "logset_src" {
  logset_name = "tf-example-src"
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_cls_topic" "topic_src" {
  topic_name           = "tf-example_src"
  logset_id            = tencentcloud_cls_logset.logset_src.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags = {
    "test" = "test",
  }
}
resource "tencentcloud_cls_logset" "logset_dst" {
  logset_name = "tf-example-dst"
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_cls_topic" "topic_dst" {
  topic_name           = "tf-example-dst"
  logset_id            = tencentcloud_cls_logset.logset_dst.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags = {
    "test" = "test",
  }
}
resource "tencentcloud_cls_data_transform" "data_transform" {
  func_type    = 1
  src_topic_id = tencentcloud_cls_topic.topic_src.id
  name         = "iac-test-src"
  etl_content  = "ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")"
  task_type    = 3
  enable_flag  = 1
  dst_resources {
    topic_id = tencentcloud_cls_topic.topic_dst.id
    alias    = "iac-test-dst"

  }
}
```

## Argument Reference

The following arguments are supported:

* `etl_content` - (Required, String) data transform content.
* `func_type` - (Required, Int) task type.
* `name` - (Required, String) task name.
* `src_topic_id` - (Required, String) src topic id.
* `task_type` - (Required, Int) task type.
* `dst_resources` - (Optional, List) data transform des resources.
* `enable_flag` - (Optional, Int) task enable flag.

The `dst_resources` object supports the following:

* `alias` - (Required, String) alias.
* `topic_id` - (Required, String) dst topic id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls data_transform can be imported using the id, e.g.

```
terraform import tencentcloud_cls_data_transform.data_transform data_transform_id
```

