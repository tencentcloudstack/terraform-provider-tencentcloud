---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_redirections"
sidebar_current: "docs-tencentcloud-datasource-clb_redirections"
description: |-
  Use this data source to query detailed information of CLB redirections
---

# tencentcloud_clb_redirections

Use this data source to query detailed information of CLB redirections

## Example Usage

```hcl
data "tencentcloud_clb_redirections" "foo" {
  clb_id             = "lb-p7olt9e5"
  source_listener_id = "lbl-jc1dx6ju"
  target_listener_id = "lbl-asj1hzuo"
  source_rule_id     = "loc-ft8fmngv"
  target_rule_id     = "loc-4xxr2cy7"
  result_output_file = "mytestpath"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required) Id of the CLB to be queried.
* `source_listener_id` - (Required) Id of source listener to be queried.
* `source_rule_id` - (Required) Rule id of source listener to be queried.
* `result_output_file` - (Optional) Used to save results.
* `target_listener_id` - (Optional) Id of target listener to be queried.
* `target_rule_id` - (Optional) Rule id of target listener to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `redirection_list` - A list of cloud load balancer redirection configurations. Each element contains the following attributes:
  * `clb_id` - Id of the CLB.
  * `source_listener_id` - Id of source listener.
  * `source_rule_id` - Rule id of source listener.
  * `target_listener_id` - Id of target listener.
  * `target_rule_id` - Rule id of target listener.


