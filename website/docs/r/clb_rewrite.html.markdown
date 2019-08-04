---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_rewrite"
sidebar_current: "docs-tencentcloud-resource-clb_rewrite"
description: |-
  Provide a resource to create a CLB instance.
---

# tencentcloud_clb_rewrite

Provide a resource to create a CLB instance.

## Example Usage

```hcl
resource "tencentcloud_clb_rewrite" "rewrite" {
  clb_id                = "lb-p7olt9e5"
  source_listener_id    = "lbl-jc1dx6ju#lb-p7olt9e5"
  target_listener_id    = "lbl-asj1hzuo#lb-p7olt9e5"
  rewrite_source_loc_id = "loc-ft8fmngv#lbl-jc1dx6ju#lb-p7olt9e5"
  rewrite_target_loc_id = "loc-4xxr2cy7#lbl-asj1hzuo#lb-p7olt9e5"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, ForceNew) Id of CLB instance.
* `rewrite_source_loc_id` - (Required, ForceNew) Id of rule id of source listener. 
* `rewrite_target_loc_id` - (Required, ForceNew) Id of rule id of target listener. 
* `source_listener_id` - (Required, ForceNew) Id of source listener. 
* `target_listener_id` - (Required, ForceNew) Id of source listener. 


## Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_rewrite.rewrite loc-ft8fmngv#loc-4xxr2cy7#lbl-jc1dx6ju#lbl-asj1hzuo#lb-p7olt9e5
```

