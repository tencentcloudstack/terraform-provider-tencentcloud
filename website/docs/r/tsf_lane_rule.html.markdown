---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_lane_rule"
sidebar_current: "docs-tencentcloud-resource-tsf_lane_rule"
description: |-
  Provides a resource to create a tsf lane_rule
---

# tencentcloud_tsf_lane_rule

Provides a resource to create a tsf lane_rule

## Example Usage

```hcl
resource "tencentcloud_tsf_lane_rule" "lane_rule" {
  rule_name = "terraform-rule-name"
  remark    = "terraform-test"
  rule_tag_list {
    tag_name     = "xxx"
    tag_operator = "EQUAL"
    tag_value    = "222"
  }
  rule_tag_relationship = "RELEATION_AND"
  lane_id               = "lane-abw5oo5a"
  enable                = false
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required, Bool) open state, true/false, default: false.
* `lane_id` - (Required, String) lane ID.
* `remark` - (Required, String) Lane rule notes.
* `rule_name` - (Required, String) lane rule name.
* `rule_tag_list` - (Required, List) list of swimlane rule labels.
* `rule_tag_relationship` - (Required, String) lane rule label relationship.
* `program_id_list` - (Optional, Set: [`String`]) Program id list.

The `rule_tag_list` object supports the following:

* `tag_name` - (Required, String) label name.
* `tag_operator` - (Required, String) label operator.
* `tag_value` - (Required, String) tag value.
* `create_time` - (Optional, Int) creation time.
* `lane_rule_id` - (Optional, String) lane rule ID.
* `tag_id` - (Optional, String) label ID.
* `update_time` - (Optional, Int) update time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - creation time.
* `priority` - Priority.
* `rule_id` - Rule id.
* `update_time` - update time.


