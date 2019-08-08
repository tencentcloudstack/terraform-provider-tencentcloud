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
  clb_id                 = "lb-p7olt9e5"
  source_listener_id     = "lbl-jc1dx6ju#lb-p7olt9e5"
  target_listener_id     = "lbl-asj1hzuo#lb-p7olt9e5"
  rewrite_source_rule_id = "loc-ft8fmngv#lbl-jc1dx6ju#lb-p7olt9e5"
  rewrite_target_rule_id = "loc-4xxr2cy7#lbl-asj1hzuo#lb-p7olt9e5"
  result_output_file     = "mytestpath"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required)  ID of the CLB to be queried.
* `rewrite_source_rule_id` - (Required) Rule ID of source listener to be queried.
* `source_listener_id` - (Required) Id of source listener to be queried.
* `result_output_file` - (Optional) Used to save results.
* `rewrite_target_rule_id` - (Optional) Rule ID of target listener to be queried.
* `target_listener_id` - (Optional) Id of source listener to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `redirection_list` - A list of cloud load redirection configurations. Each element contains the following attributes:
  * `clb_id` -  ID of the CLB.
  * `rewrite_source_rule_id` - Rule IDd of source listener.
  * `rewrite_target_rule_id` - Rule ID of target listener.
  * `source_listener_id` - Id of source listener.
  * `target_listener_id` - Id of source listener.


