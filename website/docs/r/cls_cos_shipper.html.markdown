---
subcategory: "CLS"
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
resource "tencentcloud_cls_cos_shipper" "shipper" {
  bucket       = "preset-scf-bucket-1308919341"
  interval     = 300
  max_size     = 200
  partition    = "/%Y/%m/%d/%H/"
  prefix       = "ap-guangzhou-fffsasad-1649734752"
  shipper_name = "ap-guangzhou-fffsasad-1649734752"
  topic_id     = "4d07fba0-b93e-4e0b-9a7f-d58542560bbb"

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

* `bucket` - (Required) Destination bucket in the shipping rule to be created.
* `prefix` - (Required) Prefix of the shipping directory in the shipping rule to be created.
* `shipper_name` - (Required) Shipping rule name.
* `topic_id` - (Required) ID of the log topic to which the shipping rule to be created belongs.
* `compress` - (Optional) Compression configuration of shipped log.
* `content` - (Optional) Format configuration of shipped log content.
* `filter_rules` - (Optional) Filter rules for shipped logs. Only logs matching the rules can be shipped. All rules are in the AND relationship, and up to five rules can be added. If the array is empty, no filtering will be performed, and all logs will be shipped.
* `interval` - (Optional) Shipping time interval in seconds. Default value: 300. Value range: 300~900.
* `max_size` - (Optional) Maximum size of a file to be shipped, in MB. Default value: 256. Value range: 100~256.
* `partition` - (Optional) Partition rule of shipped log, which can be represented in strftime time format.

The `compress` object supports the following:

* `format` - (Required) Compression format. Valid values: gzip, lzop, none (no compression).

The `content` object supports the following:

* `format` - (Required) Content format. Valid values: json, csv.
* `csv` - (Optional) CSV format content description.Note: this field may return null, indicating that no valid values can be obtained.
* `json` - (Optional) JSON format content description.Note: this field may return null, indicating that no valid values can be obtained.

The `csv` object supports the following:

* `delimiter` - (Required) Field delimiter.
* `escape_char` - (Required) Field delimiter.
* `keys` - (Required) Names of keys.Note: this field may return null, indicating that no valid values can be obtained.
* `non_existing_field` - (Required) Content used to populate non-existing fields.
* `print_key` - (Required) Whether to print key on the first row of the CSV file.

The `filter_rules` object supports the following:

* `key` - (Required) Filter rule key.
* `regex` - (Required) Filter rule.
* `value` - (Required) Filter rule value.

The `json` object supports the following:

* `enable_tag` - (Required) Enablement flag.
* `meta_fields` - (Required) Metadata information list
Note: this field may return null, indicating that no valid values can be obtained..

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls cos shipper can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_cos_shipper.shipper 5d1b7b2a-c163-4c48-bb01-9ee00584d761
```

