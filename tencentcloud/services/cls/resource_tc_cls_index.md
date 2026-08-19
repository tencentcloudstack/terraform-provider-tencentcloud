Provides a resource to create a cls index.

Example Usage

```hcl
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf_example"
  tags        = {
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
  tags                 = {
    "test" = "test",
  }
}

locals {
  tokenizer_value = <<-EOT
    @&?|#()='\",;:<>[]{}/ \n\t\r\\
  EOT
}

resource "tencentcloud_cls_index" "example" {
  topic_id = "abc97756-e620-47a4-aa2b-08561e79f086"

  rule {
    full_text {
      case_sensitive = true
      tokenizer      = local.tokenizer_value
      contain_z_h    = true
    }

    key_value {
      case_sensitive = true
      key_values {
        key = "key1"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = local.tokenizer_value
          type        = "text"
          alias       = "alias1"
        }
      }

      key_values {
        key = "key2"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = local.tokenizer_value
          type        = "json"
          alias       = "alias2"
          child_node {
            key = "key3"
            value {
              contain_z_h = true
              sql_flag    = true
              tokenizer   = local.tokenizer_value
              type        = "json"
              alias       = "alias3"
              child_node {
                key = "key4"
                value {
                  contain_z_h = true
                  sql_flag    = true
                  tokenizer   = local.tokenizer_value
                  type        = "text"
                  alias       = "alias4"
                }
              }
              child_node {
                key = "key5"
                value {
                  contain_z_h = true
                  sql_flag    = true
                  tokenizer   = local.tokenizer_value
                  type        = "text"
                  alias       = "name5"
                }
              }
            }
          }
          child_node {
            key = "key6"
            value {
              contain_z_h = true
              sql_flag    = true
              tokenizer   = local.tokenizer_value
              type        = "text"
              alias       = "name6"
            }
          }
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

Import

cls cos index can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_index.example 0937e56f-4008-49d2-ad2d-69c52a9f11cc
```