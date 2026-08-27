Provides a resource to create a CLS splunk deliver.

Example Usage

```hcl
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf-example"
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf-example"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
}

resource "tencentcloud_cls_splunk_deliver" "example" {
  topic_id = tencentcloud_cls_topic.example.id
  name     = "tf-example"

  net_info {
    host     = "10.0.0.1"
    port     = 8088
    token    = "your-splunk-token"
    net_type = 2
    is_ssl   = true
  }

  metadata_info {
    format = "json"
    meta_fields = [
      "__SOURCE__",
      "__FILENAME__",
      "__TIMESTAMP__",
    ]
    enable_tag = true
  }
}
```

Import

cls splunk deliver can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_splunk_deliver.example task-xxx#topic-yyy
```