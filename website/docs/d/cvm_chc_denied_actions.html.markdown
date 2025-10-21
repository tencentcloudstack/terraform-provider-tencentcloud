---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_chc_denied_actions"
sidebar_current: "docs-tencentcloud-datasource-cvm_chc_denied_actions"
description: |-
  Use this data source to query detailed information of cvm chc_denied_actions
---

# tencentcloud_cvm_chc_denied_actions

Use this data source to query detailed information of cvm chc_denied_actions

## Example Usage

```hcl
data "tencentcloud_cvm_chc_denied_actions" "chc_denied_actions" {
  chc_ids = ["chc-xxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `chc_ids` - (Required, Set: [`String`]) CHC host IDs.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `chc_host_denied_action_set` - Actions not allowed for the CHC instance.
  * `chc_id` - CHC instance ID.
  * `deny_actions` - Actions not allowed for the current CHC instance.
  * `state` - CHC instance status.


