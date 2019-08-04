---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_server_attachments"
sidebar_current: "docs-tencentcloud-datasource-clb_server_attachments"
description: |-
  Use this data source to query detailed information of CLB
---

# tencentcloud_clb_server_attachments

Use this data source to query detailed information of CLB

## Example Usage

```hcl
data "tencentcloud_clb" "clblab" {
  listener_id   = "lbl-hh141sn9#lb-k2zjp9lv"
  clb_id        = "lb-k2zjp9lv"
  location_id   = "loc-4xxr2cy7"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required)  ID of the CLB.
* `listener_id` - (Required)  ID of the CLB listener.
* `location_id` - (Required)  ID of the CLB listener rule. 
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `attachment_list` - A list of cloud load redirection configurations. Each element contains the following attributes:
  * `clb_id` - Id of the cloud load balancer. 
  * `listener_id` - Id of the cloud load balance listener. 
  * `location_id` - Id of the cloud load balance listener rule. 
  * `protocol_type` - Type of protocol within the listener, and available values include 'TCP', 'UDP', 'HTTP', 'HTTPS' and 'TCP_SSL'. 
  * `targets` - Backend infos.
    * `instance_id` - Id of the backend server.
    * `port` - Port of the backend server.
    * `weight` - Weight of the backend server.


