---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_cos_recharge"
sidebar_current: "docs-tencentcloud-resource-cls_cos_recharge"
description: |-
  Provides a resource to create a cls cos_recharge
---

# tencentcloud_cls_cos_recharge

Provides a resource to create a cls cos_recharge

~> **NOTE:** This resource can not be deleted if you run `terraform destroy`.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `bucket_region` - (Required, String) cos bucket region.
* `bucket` - (Required, String) cos bucket.
* `log_type` - (Required, String) log type.
* `logset_id` - (Required, String) logset id.
* `name` - (Required, String) recharge name.
* `prefix` - (Required, String) cos file prefix.
* `topic_id` - (Required, String, ForceNew) topic id.
* `compress` - (Optional, String) supported gzip, lzop, snappy.
* `extract_rule_info` - (Optional, List) extract rule info.

The `extract_rule_info` object supports the following:

* `address` - (Optional, String) syslog address.
* `backtracking` - (Optional, Int) backtracking data volume in incremental acquisition mode.
* `begin_regex` - (Optional, String) begin line regex.
* `delimiter` - (Optional, String) log delimiter.
* `filter_key_regex` - (Optional, List) rules that need to filter logs.
* `is_gbk` - (Optional, Int) gbk encoding.
* `json_standard` - (Optional, Int) is standard json.
* `keys` - (Optional, Set) key list.
* `log_regex` - (Optional, String) log regex.
* `meta_tags` - (Optional, List) metadata tag list.
* `metadata_type` - (Optional, Int) metadata type.
* `parse_protocol` - (Optional, String) parse protocol.
* `path_regex` - (Optional, String) metadata path regex.
* `protocol` - (Optional, String) syslog protocol.
* `time_format` - (Optional, String) time format.
* `time_key` - (Optional, String) time key.
* `un_match_log_key` - (Optional, String) parsing failure log key.
* `un_match_up_load_switch` - (Optional, Bool) whether to upload the parsing failure log.

The `filter_key_regex` object supports the following:

* `key` - (Required, String) need filter log key.
* `regex` - (Required, String) need filter log regex.

The `meta_tags` object supports the following:

* `key` - (Optional, String) metadata key.
* `value` - (Optional, String) metadata value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls cos_recharge can be imported using the id, e.g.

```
terraform import tencentcloud_cls_cos_recharge.cos_recharge topic_id#cos_recharge_id
```

