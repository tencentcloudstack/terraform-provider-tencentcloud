---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_ckafka_consumer"
sidebar_current: "docs-tencentcloud-resource-cls_ckafka_consumer"
description: |-
  Provides a resource to create a CLS ckafka consumer
---

# tencentcloud_cls_ckafka_consumer

Provides a resource to create a CLS ckafka consumer

## Example Usage

### Create a Ckafka delivery task

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
    enable_tag = true
    meta_fields = [
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

### Create a cross-account Ckafka delivery task

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
    instance_id = "ckafka-5k7vmdzo"
    topic_name  = "tf-test"
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

## Argument Reference

The following arguments are supported:

* `topic_id` - (Required, String, ForceNew) topic id.
* `advanced_config` - (Optional, List) advanced configuration.
* `ckafka` - (Optional, List) ckafka info.
* `compression` - (Optional, Int) compression method. 0 for NONE, 2 for SNAPPY, 3 for LZ4.
* `content` - (Optional, List) metadata information.
* `dsl_filter` - (Optional, String) log pre-filtering - pre-filter the raw data written to ckafka.
* `external_id` - (Optional, String) external id for cross-account delivery.
* `need_content` - (Optional, Bool) whether to deliver the metadata information of the log.
* `role_arn` - (Optional, String) role access description name for cross-account delivery.

The `advanced_config` object supports the following:

* `partition_fields` - (Optional, Set) list of fields used to calculate hash. Supports up to 5 fields.
* `partition_hash_status` - (Optional, Bool) whether to deliver messages with the same hash value of fields to the same ckafka partition. Default value: `false`.

The `ckafka` object supports the following:

* `instance_id` - (Required, String) instance id.
* `topic_name` - (Required, String) topic name of ckafka.
* `instance_name` - (Optional, String) instance name.If the delivery task is created using a Role ARN, this field can be left empty.
* `topic_id` - (Optional, String) topic id of ckafka. If the delivery task is created using a Role ARN, this field can be left empty.
* `vip` - (Optional, String) vip.If the delivery task is created using a Role ARN, this field can be left empty.
* `vport` - (Optional, String) vport.If the delivery task is created using a Role ARN, this field can be left empty.

The `content` object supports the following:

* `enable_tag` - (Required, Bool) whether to deliver the TAG info.
* `meta_fields` - (Required, Set) metadata info list.
* `auto_convert_number` - (Optional, Bool) whether to automatically convert numeric strings (such as `123`) in the first-level values of the JSON structure to numeric types (int / float). Default value: `false`.
* `json_type` - (Optional, Int) JSON format for delivery. Valid values:
- `0`: consistent with the original log, no escaping.
- `1`: escaped.
* `tag_json_not_tiled` - (Optional, Bool) whether to tiling tag json.
* `timestamp_accuracy` - (Optional, Int) delivery timestamp precision,1 for second, 2 for millisecond.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateConsumer](https://www.tencentcloud.com/document/api/614/44253)
* Read: [DescribeConsumer](https://www.tencentcloud.com/document/api/614/44251)
* Update: [ModifyConsumer](https://www.tencentcloud.com/document/api/614/44250)
* Delete: [DeleteConsumer](https://www.tencentcloud.com/document/api/614/44252)

## Import

cls ckafka_consumer can be imported using the id, e.g.

```
terraform import tencentcloud_cls_ckafka_consumer.example 7e34a3a7-635e-4da8-9005-88106c1fde69
```

