---
subcategory: "Anti-DDoS(Dayu)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_cc_http_policies"
sidebar_current: "docs-tencentcloud-datasource-dayu_cc_http_policies"
description: |-
  Use this data source to query dayu CC http policies
---

# tencentcloud_dayu_cc_http_policies

Use this data source to query dayu CC http policies

## Example Usage

```hcl
data "tencentcloud_dayu_cc_http_policies" "id_test" {
  resource_type = tencentcloud_dayu_cc_http_policy.test_policy.resource_type
  resource_id   = tencentcloud_dayu_cc_http_policy.test_policy.resource_id
  policy_id     = tencentcloud_dayu_cc_http_policy.test_policy.policy_id
}
data "tencentcloud_dayu_cc_http_policies" "name_test" {
  resource_type = tencentcloud_dayu_cc_http_policy.test_policy.resource_type
  resource_id   = tencentcloud_dayu_cc_http_policy.test_policy.resource_id
  name          = tencentcloud_dayu_cc_http_policy.test_policy.name
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) Id of the resource that the CC http policy works for.
* `resource_type` - (Required) Type of the resource that the CC http policy works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.
* `name` - (Optional) Name of the CC http policy to be queried.
* `policy_id` - (Optional) Id of the CC http policy to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of CC http policies. Each element contains the following attributes:
  * `action` - Action mode.
  * `create_time` - Create time of the CC self-define http policy.
  * `frequency` - Max frequency per minute.
  * `ip_list` - Ip of the CC self-define http policy.
  * `name` - Name of the CC self-define http policy.
  * `policy_id` - Id of the CC self-define http policy.
  * `resource_id` - ID of the resource that the CC self-define http policy works for.
  * `resource_type` - Type of the resource that the CC self-define http policy works for.
  * `smode` - Match mode.
  * `switch` - Indicate the CC self-define http policy takes effect or not.


