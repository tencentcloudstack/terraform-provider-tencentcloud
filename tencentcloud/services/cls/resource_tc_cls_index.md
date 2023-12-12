Provides a resource to create a cls index.

Example Usage

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

Import

cls cos index can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_index.index 0937e56f-4008-49d2-ad2d-69c52a9f11cc
```