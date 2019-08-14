---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_redirection"
sidebar_current: "docs-tencentcloud-resource-clb_redirection"
description: |-
  Provides a resource to create a CLB redirection.
---

# tencentcloud_clb_redirection

Provides a resource to create a CLB redirection.

## Example Usage

```hcl
resource "tencentcloud_clb_redirection" "foo" {
  clb_id                  = "lb-p7olt9e5"
  source_listener_id      = "lbl-jc1dx6ju"
  target_listener_id      = "lbl-asj1hzuo"
  source_listener_rule_id = "loc-ft8fmngv"
  target_listener_rule_id = "loc-4xxr2cy7"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, ForceNew) Id of CLB instance.
* `source_listener_id` - (Required, ForceNew) Id of source listener.
* `source_listener_rule_id` - (Required, ForceNew) Rule ID of source listener.
* `target_listener_id` - (Required, ForceNew) Id of source listener.
* `target_listener_rule_id` - (Required, ForceNew) Rule ID of target listener.


## Import

CLB redirection can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_redirection.foo loc-ft8fmngv#loc-4xxr2cy7#lbl-jc1dx6ju#lbl-asj1hzuo#lb-p7olt9e5
```

