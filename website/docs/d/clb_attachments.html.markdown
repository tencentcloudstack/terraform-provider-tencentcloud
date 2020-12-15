---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_attachments"
sidebar_current: "docs-tencentcloud-datasource-clb_attachments"
description: |-
  Use this data source to query detailed information of CLB attachments
---

# tencentcloud_clb_attachments

Use this data source to query detailed information of CLB attachments

## Example Usage

```hcl
data "tencentcloud_clb_attachments" "clblab" {
  listener_id = "lbl-hh141sn9"
  clb_id      = "lb-k2zjp9lv"
  rule_id     = "loc-4xxr2cy7"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required) ID of the CLB to be queried.
* `listener_id` - (Required) ID of the CLB listener to be queried.
* `result_output_file` - (Optional) Used to save results.
* `rule_id` - (Optional) ID of the CLB listener rule. If the protocol of listener is `HTTP`/`HTTPS`, this para is required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `attachment_list` - A list of cloud load balancer attachment configurations. Each element contains the following attributes:
  * `clb_id` - ID of the CLB.
  * `listener_id` - ID of the CLB listener.
  * `protocol_type` - Type of protocol within the listener, and available values include `TCP`, `UDP`, `HTTP`, `HTTPS` and `TCP_SSL`. NOTES: `TCP_SSL` is testing internally, please apply if you need to use.
  * `rule_id` - ID of the CLB listener rule.
  * `targets` - Information of the backends to be attached.
    * `instance_id` - Id of the backend server.
    * `port` - Port of the backend server.
    * `weight` - Forwarding weight of the backend service, the range of [0, 100], defaults to `10`.


