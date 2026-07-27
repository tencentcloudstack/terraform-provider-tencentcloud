Provides a resource to create a CLS config

Example Usage

```hcl
resource "tencentcloud_cls_config" "example" {
  name       = "tf-example"
  output     = "734f50d1-d621-425c-8768-6f9a5f0412ee"
  path       = "/data/log/**/error.log"
  log_type   = "json_log"
  input_type = "file"
  extract_rule {
    filter_key_regex {
      key   = "key1"
      regex = "value1"
    }

    filter_key_regex {
      key   = "key2"
      regex = "value2"
    }

    is_gbk                  = 0
    json_standard           = 1
    un_match_up_load_switch = true
    un_match_log_key        = "LogParseFailure"
    backtracking            = 0
    metadata_type           = 2
    meta_tags {
      key   = "myKey"
      value = "myValue"
    }

    filter_key_regex {
      key   = "ErrorCode"
      regex = "500"
    }
  }

  exclude_paths {
    type  = "Path"
    value = "/data"
  }

  exclude_paths {
    type  = "File"
    value = "/file"
  }
}
```

Import

CLS config can be imported using the id, e.g.

```
terraform import tencentcloud_cls_config.example 49611ec9-c5f2-4cc9-9e06-15dd7fa43982
```