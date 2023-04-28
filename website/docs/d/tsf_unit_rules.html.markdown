---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_unit_rules"
sidebar_current: "docs-tencentcloud-datasource-tsf_unit_rules"
description: |-
  Use this data source to query detailed information of tsf unit_rules
---

# tencentcloud_tsf_unit_rules

Use this data source to query detailed information of tsf unit_rules

## Example Usage

```hcl
data "tencentcloud_tsf_unit_rules" "unit_rules" {
  gateway_instance_id = "gw-ins-lvdypq5k"
  status              = "disabled"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_instance_id` - (Required, String) gateway instance id.
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, String) Enabled state, disabled: unpublished, enabled: published.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Pagination list information.
  * `content` - record entity list.
    * `created_time` - created time.
    * `description` - Rule description.
    * `gateway_instance_id` - Gateway Entity ID.
    * `id` - rule ID.
    * `name` - rule name.
    * `status` - Use status: enabled/disabled.
    * `unit_rule_item_list` - list of rule items.
      * `description` - Rule description.
      * `dest_namespace_id` - Destination Namespace ID.
      * `dest_namespace_name` - destination namespace name.
      * `id` - rule item ID.
      * `name` - rule item name.
      * `priority` - Rule order, the smaller the higher the priority: the default is 0.
      * `relationship` - Logical relationship: AND/OR.
      * `unit_rule_id` - Unitization rule ID.
      * `unit_rule_tag_list` - List of rule labels.
        * `id` - rule ID.
        * `tag_field` - tag name.
        * `tag_operator` - Operator: IN/NOT_IN/EQUAL/NOT_EQUAL/REGEX.
        * `tag_type` - Tag Type: U(User Tag).
        * `tag_value` - tag value.
        * `unit_rule_item_id` - Unitization rule item ID.
    * `updated_time` - Updated time.
  * `total_count` - total number of records.


