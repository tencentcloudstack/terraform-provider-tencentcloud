---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_cos_shipper"
sidebar_current: "docs-tencentcloud-resource-cls_cos_shipper"
description: |-
  Provides a resource to create a cls cos shipper.
---

# tencentcloud_cls_cos_shipper

Provides a resource to create a cls cos shipper.

## Example Usage

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "example" {
  bucket = "private-bucket-${local.app_id}"
  acl    = "private"
}

resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf-example"
  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf-example"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cls_cos_shipper" "example" {
  bucket       = tencentcloud_cos_bucket.example.id
  topic_id     = tencentcloud_cls_topic.example.id
  interval     = 300
  max_size     = 200
  partition    = "/%Y/%m/%d/%H/"
  prefix       = "ap-guangzhou-fffsasad-1649734752"
  shipper_name = "ap-guangzhou-fffsasad-1649734752"

  compress {
    format = "lzop"
  }

  content {
    format = "json"

    json {
      enable_tag = true
      meta_fields = [
        "__FILENAME__",
        "__SOURCE__",
        "__TIMESTAMP__",
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) Destination bucket in the shipping rule to be created.
* `prefix` - (Required, String) Prefix of the shipping directory in the shipping rule to be created.
* `shipper_name` - (Required, String) Shipping rule name.
* `topic_id` - (Required, String) ID of the log topic to which the shipping rule to be created belongs.
* `compress` - (Optional, List) Compression configuration of shipped log.
* `content` - (Optional, List) Format configuration of shipped log content.
* `end_time` - (Optional, Int) End time for data shipping, which cannot be set to a future time. If you do not specify this parameter, it indicates continuous data shipping.
* `filename_mode` - (Optional, Int) Naming a shipping file. Valid values: 0 (by random number); 1 (by shipping time). Default value: 0.
* `filter_rules` - (Optional, List) Filter rules for shipped logs. Only logs matching the rules can be shipped. All rules are in the AND relationship, and up to five rules can be added. If the array is empty, no filtering will be performed, and all logs will be shipped.
* `interval` - (Optional, Int) Shipping time interval in seconds. Default value: 300. Value range: 300~900.
* `max_size` - (Optional, Int) Maximum size of a file to be shipped, in MB. Default value: 256. Value range: 100~256.
* `partition` - (Optional, String) Partition rule of shipped log, which can be represented in strftime time format.
* `start_time` - (Optional, Int) Start time for data shipping, which cannot be earlier than the lifecycle start time of the log topic. If you do not specify this parameter, it will be set to the time when you create the data shipping task.
* `storage_type` - (Optional, String) COS bucket storage type. support: STANDARD_IA, ARCHIVE, DEEP_ARCHIVE, STANDARD, MAZ_STANDARD, MAZ_STANDARD_IA, INTELLIGENT_TIERING.

The `compress` object supports the following:

* `format` - (Required, String) Compression format. Valid values: gzip, lzop, none (no compression).

The `content` object supports the following:

* `format` - (Required, String) Content format. Valid values: json, csv.
* `csv` - (Optional, List) CSV format content description.Note: this field may return null, indicating that no valid values can be obtained.
* `json` - (Optional, List) JSON format content description.Note: this field may return null, indicating that no valid values can be obtained.

The `csv` object of `content` supports the following:

* `delimiter` - (Required, String) Field delimiter.
* `escape_char` - (Required, String) Field delimiter.
* `keys` - (Required, Set) Names of keys.Note: this field may return null, indicating that no valid values can be obtained.
* `non_existing_field` - (Required, String) Content used to populate non-existing fields.
* `print_key` - (Required, Bool) Whether to print key on the first row of the CSV file.

The `filter_rules` object supports the following:

* `key` - (Required, String) Filter rule key.
* `regex` - (Required, String) Filter rule.
* `value` - (Required, String) Filter rule value.

The `json` object of `content` supports the following:

* `enable_tag` - (Required, Bool) Enablement flag.
* `meta_fields` - (Required, Set) Metadata information list
Note: this field may return null, indicating that no valid values can be obtained..

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls cos shipper can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_cos_shipper.example 5d1b7b2a-c163-4c48-bb01-9ee00584d761
```

