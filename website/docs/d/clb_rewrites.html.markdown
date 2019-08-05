---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_rewrites"
sidebar_current: "docs-tencentcloud-datasource-clb_rewrites"
description: |-
  Use this data source to query detailed information of CLB rewrite
---

# tencentcloud_clb_rewrites

Use this data source to query detailed information of CLB rewrite

## Example Usage

```hcl
data "tencentcloud_clb_rewrites" "clblab" {
  clb_id                = "lb-p7olt9e5"
  source_listener_id    = "lbl-jc1dx6ju#lb-p7olt9e5"
  target_listener_id    = "lbl-asj1hzuo#lb-p7olt9e5"
  rewrite_source_loc_id = "loc-ft8fmngv#lbl-jc1dx6ju#lb-p7olt9e5"
  rewrite_target_loc_id = "loc-4xxr2cy7#lbl-asj1hzuo#lb-p7olt9e5"
  result_output_file    = "mytestpath"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required)  ID of the CLB to be queried.
* `rewrite_source_loc_id` - (Required) Id of rule id of source listener.
* `source_listener_id` - (Required) Id of source listener.
* `result_output_file` - (Optional) Used to save results.
* `rewrite_target_loc_id` - (Optional) Id of rule id of target listener.
* `target_listener_id` - (Optional) Id of source listener.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rewrite_list` - A list of cloud load redirection configurations. Each element contains the following attributes:
  * `clb_id` -  ID of the CLB to be queried.
  * `rewrite_source_loc_id` - Id of rule id of source listener.
  * `rewrite_target_loc_id` - Id of rule id of target listener.
  * `source_listener_id` - Id of source listener.
  * `target_listener_id` - Id of source listener.


