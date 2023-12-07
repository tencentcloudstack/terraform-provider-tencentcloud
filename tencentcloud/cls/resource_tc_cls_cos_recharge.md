Provides a resource to create a cls cos_recharge

~> **NOTE:** This resource can not be deleted if you run `terraform destroy`.

Example Usage

```hcl
resource "tencentcloud_cls_cos_recharge" "cos_recharge" {
  bucket        = "cos-lock-1308919341"
  bucket_region = "ap-guangzhou"
  log_type      = "minimalist_log"
  logset_id     = "dd426d1a-95bc-4bca-b8c2-baa169261812"
  name          = "cos_recharge_for_test"
  prefix        = "test"
  topic_id      = "7e34a3a7-635e-4da8-9005-88106c1fde69"

  extract_rule_info {
    backtracking            = 0
    is_gbk                  = 0
    json_standard           = 0
    keys                    = []
    metadata_type           = 0
    un_match_up_load_switch = false

    filter_key_regex {
      key   = "__CONTENT__"
      regex = "dasd"
    }
  }
}
```

Import

cls cos_recharge can be imported using the id, e.g.

```
terraform import tencentcloud_cls_cos_recharge.cos_recharge topic_id#cos_recharge_id
```