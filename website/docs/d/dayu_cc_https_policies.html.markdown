---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_cc_https_policies"
sidebar_current: "docs-tencentcloud-datasource-dayu_cc_https_policies"
description: |-
  Use this data source to query dayu CC https policies
---

# tencentcloud_dayu_cc_https_policies

Use this data source to query dayu CC https policies

## Example Usage

```hcl
data "tencentcloud_dayu_cc_https_policies" "name_test" {
  resource_type = tencentcloud_dayu_cc_https_policy.test_policy.resource_type
  resource_id   = tencentcloud_dayu_cc_https_policy.test_policy.resource_id
  name          = tencentcloud_dayu_cc_https_policy.test_policy.name
}
data "tencentcloud_dayu_cc_https_policies" "id_test" {
  resource_type = tencentcloud_dayu_cc_https_policy.test_policy.resource_type
  resource_id   = tencentcloud_dayu_cc_https_policy.test_policy.resource_id
  policy_id     = tencentcloud_dayu_cc_https_policy.test_policy.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) Id of the resource that the CC https policy works for.
* `resource_type` - (Required) Type of the resource that the CC https policy works for, valid values are `bgpip`.
* `name` - (Optional) Name of the CC https policy to be queried.
* `policy_id` - (Optional) Id of the CC https policy to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of CC https policies. Each element contains the following attributes.
  * `create_time` - Create time of the CC self-define https policy.
  * `domain` - Domain that the CC self-define https policy works for.
  * `exe_mode` - Execute mode.
  * `ip_list` - Ip of the CC self-define https policy.
  * `name` - Name of the CC self-define https policy.
  * `policy_id` - Id of the CC self-define https policy.
  * `resource_id` - ID of the resource that the CC self-define https policy works for.
  * `resource_type` - Type of the resource that the CC self-define https policy works for.
  * `rule_id` - Rule id of the domain that the CC self-define https policy works for.
  * `switch` - Indicate the CC self-define https policy takes effect or not.


