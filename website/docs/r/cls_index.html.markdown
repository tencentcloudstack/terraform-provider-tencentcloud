---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_index"
sidebar_current: "docs-tencentcloud-resource-cls_index"
description: |-
  Provides a resource to create a cls index.
---

# tencentcloud_cls_index

Provides a resource to create a cls index.

## Example Usage

```hcl
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf_example"
  tags = {
    "demo" = "test"
  }
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf_example"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 30
  storage_type         = "hot"
  describes            = "Test Demo."
  hot_period           = 10
  tags = {
    "test" = "test",
  }
}

locals {
  tokenizer_value = "@&?|#()='\",;:<>[]{}"
}

resource "tencentcloud_cls_index" "example" {
  topic_id = tencentcloud_cls_topic.example.id

  rule {
    full_text {
      case_sensitive = true
      tokenizer      = local.tokenizer_value
      contain_z_h    = true
    }

    key_value {
      case_sensitive = true
      key_values {
        key = "hello"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = local.tokenizer_value
          type        = "text"
        }
      }

      key_values {
        key = "world"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = local.tokenizer_value
          type        = "text"
        }
      }
    }

    tag {
      case_sensitive = true
      key_values {
        key = "terraform"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = local.tokenizer_value
          type        = "text"
        }
      }
    }

    dynamic_index {
      status = true
    }
  }
  status                  = true
  include_internal_fields = true
  metadata_flag           = 1
}
```

## Argument Reference

The following arguments are supported:

* `topic_id` - (Required, String) Log topic ID.
* `include_internal_fields` - (Optional, Bool) Internal field marker of full-text index. Default value: false. Valid value: false: excluding internal fields; true: including internal fields.
* `metadata_flag` - (Optional, Int) Metadata flag. Default value: 0. Valid value: 0: full-text index (including the metadata field with key-value index enabled); 1: full-text index (including all metadata fields); 2: full-text index (excluding metadata fields)..
* `rule` - (Optional, List) Index rule.
* `status` - (Optional, Bool) Whether to take effect. Default value: true.

The `dynamic_index` object of `rule` supports the following:

* `status` - (Required, Bool) index automatic configuration switch.

The `full_text` object of `rule` supports the following:

* `case_sensitive` - (Required, Bool) Case sensitivity.
* `contain_z_h` - (Required, Bool) Whether Chinese characters are contained.
* `tokenizer` - (Required, String) Full-Text index delimiter. Each character in the string represents a delimiter.

The `key_value` object of `rule` supports the following:

* `case_sensitive` - (Required, Bool) Case sensitivity.
* `key_values` - (Optional, List) Key-Value pair information of the index to be created. Up to 100 key-value pairs can be configured.

The `key_values` object of `key_value` supports the following:

* `key` - (Required, String) When a key value or metafield index needs to be configured for a field, the metafield Key does not need to be prefixed with __TAG__. and is consistent with the one when logs are uploaded. __TAG__. will be prefixed automatically for display in the console..
* `value` - (Optional, List) Field index description information.

The `key_values` object of `tag` supports the following:

* `key` - (Required, String) When a key value or metafield index needs to be configured for a field, the metafield Key does not need to be prefixed with __TAG__. and is consistent with the one when logs are uploaded. __TAG__. will be prefixed automatically for display in the console..
* `value` - (Optional, List) Field index description information.

The `rule` object supports the following:

* `dynamic_index` - (Optional, List) The key value index is automatically configured. If it is empty, it means that the function is not enabled.
* `full_text` - (Optional, List) Full-Text index configuration.
* `key_value` - (Optional, List) Key-Value index configuration.
* `tag` - (Optional, List) Metafield index configuration.

The `tag` object of `rule` supports the following:

* `case_sensitive` - (Required, Bool) Case sensitivity.
* `key_values` - (Optional, List) Key-Value pair information of the index to be created. Up to 100 key-value pairs can be configured.

The `value` object of `key_values` supports the following:

* `type` - (Required, String) Field type. Valid values: long, text, double.
* `contain_z_h` - (Optional, Bool) Whether Chinese characters are contained.
* `sql_flag` - (Optional, Bool) Whether the analysis feature is enabled for the field.
* `tokenizer` - (Optional, String) Field delimiter, which is meaningful only if the field type is text. Each character in the entered string represents a delimiter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls cos index can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_index.example 0937e56f-4008-49d2-ad2d-69c52a9f11cc
```

