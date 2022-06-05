---
subcategory: "CLS"
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
resource "tencentcloud_cls_index" "index" {
  topic_id = "0937e56f-4008-49d2-ad2d-69c52a9f11cc"

  rule {
    full_text {
      case_sensitive = true
      tokenizer      = "@&?|#()='\",;:<>[]{}/ \n\t\r\\"
      contain_z_h    = true
    }

    key_value {
      case_sensitive = true
      key_values {
        key = "hello"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = "@&?|#()='\",;:<>[]{}/ \n\t\r\\"
          type        = "text"
        }
      }

      key_values {
        key = "world"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = "@&?|#()='\",;:<>[]{}/ \n\t\r\\"
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
          tokenizer   = "@&?|#()='\",;:<>[]{}/ \n\t\r\\"
          type        = "text"
        }
      }
    }
  }
  status                  = true
  include_internal_fields = true
  metadata_flag           = 1
}
```

## Argument Reference

The following arguments are supported:

* `topic_id` - (Required) Log topic ID.
* `include_internal_fields` - (Optional) Internal field marker of full-text index. Default value: false. Valid value: false: excluding internal fields; true: including internal fields.
* `metadata_flag` - (Optional) Metadata flag. Default value: 0. Valid value: 0: full-text index (including the metadata field with key-value index enabled); 1: full-text index (including all metadata fields); 2: full-text index (excluding metadata fields)..
* `rule` - (Optional) Index rule.
* `status` - (Optional) Whether to take effect. Default value: true.

The `full_text` object supports the following:

* `case_sensitive` - (Required) Case sensitivity.
* `contain_z_h` - (Required) Whether Chinese characters are contained.
* `tokenizer` - (Required) Full-Text index delimiter. Each character in the string represents a delimiter.

The `key_value` object supports the following:

* `case_sensitive` - (Required) Case sensitivity.
* `key_values` - (Optional) Key-Value pair information of the index to be created. Up to 100 key-value pairs can be configured.

The `key_values` object supports the following:

* `key` - (Required) When a key value or metafield index needs to be configured for a field, the metafield Key does not need to be prefixed with __TAG__. and is consistent with the one when logs are uploaded. __TAG__. will be prefixed automatically for display in the console..
* `value` - (Optional) Field index description information.

The `rule` object supports the following:

* `full_text` - (Optional) Full-Text index configuration.
* `key_value` - (Optional) Key-Value index configuration.
* `tag` - (Optional) Metafield index configuration.

The `tag` object supports the following:

* `case_sensitive` - (Required) Case sensitivity.
* `key_values` - (Optional) Key-Value pair information of the index to be created. Up to 100 key-value pairs can be configured.

The `value` object supports the following:

* `type` - (Required) Field type. Valid values: long, text, double.
* `contain_z_h` - (Optional) Whether Chinese characters are contained.
* `sql_flag` - (Optional) Whether the analysis feature is enabled for the field.
* `tokenizer` - (Optional) Field delimiter, which is meaningful only if the field type is text. Each character in the entered string represents a delimiter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls cos index can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_index.index 0937e56f-4008-49d2-ad2d-69c52a9f11cc
```

