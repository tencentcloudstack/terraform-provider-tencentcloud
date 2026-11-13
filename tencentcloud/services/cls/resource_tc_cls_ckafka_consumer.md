Provides a resource to create a CLS ckafka consumer

Example Usage

Create a Ckafka delivery task

```hcl
resource "tencentcloud_cls_ckafka_consumer" "example" {
  compression  = 0
  need_content = true
  topic_id     = "7e34a3a7-635e-4da8-9005-88106c1fde69"

  ckafka {
    instance_id   = "ckafka-qzoeaqx8"
    instance_name = "ckafka-instance"
    topic_id      = "topic-c6tm4kpm"
    topic_name    = "name"
    vip           = "172.16.112.23"
    vport         = "9092"
  }

  content {
    enable_tag         = true
    meta_fields        = [
      "__FILENAME__",
      "__HOSTNAME__",
      "__PKGID__",
      "__SOURCE__",
      "__TIMESTAMP__",
    ]
    tag_json_not_tiled = true
    timestamp_accuracy = 2
    json_type          = 1
  }
}
```

Create a cross-account Ckafka delivery task

```hcl
resource "tencentcloud_cls_ckafka_consumer" "example" {
  compression  = 2
  dsl_filter   = "log_keep(op_str_eq(v(\"type\"), \"kibana_stats\"))"
  external_id  = "hello123"
  need_content = true
  role_arn     = "qcs::cam::uin/100022770164:roleName/uinA_writeCLS_to_CKafka"
  topic_id     = "0c68cdff-cb95-4deb-876a-46d805a7a1fc"
  advanced_config {
    partition_fields      = ["test"]
    partition_hash_status = true
  }
  ckafka {
    instance_id   = "ckafka-5k7vmdzo"
    topic_name    = "tf-test"
  }
  content {
    auto_convert_number = true
    enable_tag          = true
    json_type           = 0
    meta_fields         = ["__FILENAME__", "__HOSTNAME__", "__SOURCE__", "__TIMESTAMP__"]
    tag_json_not_tiled  = true
    timestamp_accuracy  = 2
  }
}

```

Import

cls ckafka_consumer can be imported using the id, e.g.

```
terraform import tencentcloud_cls_ckafka_consumer.example 7e34a3a7-635e-4da8-9005-88106c1fde69
```