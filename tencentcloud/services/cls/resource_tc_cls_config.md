Provides a resource to create a cls config

Example Usage

```hcl
resource "tencentcloud_cls_config" "config" {
  name             = "config_hello"
  output           = "4d07fba0-b93e-4e0b-9a7f-d58542560bbb"
  path             = "/var/log/kubernetes"
  log_type         = "json_log"
  extract_rule {
    filter_key_regex {
      key   = "key1"
      regex = "value1"
    }
    filter_key_regex {
      key   = "key2"
      regex = "value2"
    }
    un_match_up_load_switch = true
    un_match_log_key        = "config"
    backtracking            = -1
  }
  exclude_paths {
    type  = "Path"
    value = "/data"
  }
  exclude_paths {
    type  = "File"
    value = "/file"
  }
#  user_define_rule = ""
}
```

Import

cls config can be imported using the id, e.g.

```
terraform import tencentcloud_cls_config.config config_id
```