---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_data_transform"
sidebar_current: "docs-tencentcloud-resource-cls_data_transform"
description: |-
  Provides a resource to create a CLS data transform
---

# tencentcloud_cls_data_transform

Provides a resource to create a CLS data transform

## Example Usage

```hcl
resource "tencentcloud_cls_logset" "logset_src" {
  logset_name = "tf-example-src"
  tags = {
    createdBy = "terraform"
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
    createdBy = "terraform"
  }
}

resource "tencentcloud_cls_logset" "logset_dst" {
  logset_name = "tf-example-dst"
  tags = {
    createdBy = "terraform"
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
    createdBy = "terraform"
  }
}

resource "tencentcloud_cls_data_transform" "example" {
  func_type    = 1
  src_topic_id = tencentcloud_cls_topic.topic_src.id
  name         = "tf-example"
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

* `etl_content` - (Required, String) Data transform content. If `func_type` is `2`, must use `log_auto_output`.
* `func_type` - (Required, Int) Task type. `1`: Specify the theme; `2`: Dynamic creation.
* `name` - (Required, String) Task name.
* `src_topic_id` - (Required, String) Source topic ID.
* `task_type` - (Required, Int) Task type. `1`: Use random data from the source log theme for processing preview; `2`: Use user-defined test data for processing preview; `3`: Create real machining tasks.
* `dst_resources` - (Optional, List) Data transform des resources. If `func_type` is `1`, this parameter is required. If `func_type` is `2`, this parameter does not need to be filled in.
* `enable_flag` - (Optional, Int) Task enable flag. `1`: enable, `2`: disable, Default is `1`.

The `dst_resources` object supports the following:

* `alias` - (Required, String) Alias.
* `topic_id` - (Required, String) Dst topic ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CLS data transform can be imported using the id, e.g.

```
terraform import tencentcloud_cls_data_transform.example 7b4bcb05-9154-4cdc-a479-f6b5743846e5
```

