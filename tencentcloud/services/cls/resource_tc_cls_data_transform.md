Provides a resource to create a CLS data transform

Example Usage

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

Import

CLS data transform can be imported using the id, e.g.

```
terraform import tencentcloud_cls_data_transform.example 7b4bcb05-9154-4cdc-a479-f6b5743846e5
```