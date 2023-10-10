---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_real_servers_status"
sidebar_current: "docs-tencentcloud-datasource-gaap_real_servers_status"
description: |-
  Use this data source to query detailed information of gaap real servers status
---

# tencentcloud_gaap_real_servers_status

Use this data source to query detailed information of gaap real servers status

## Example Usage

```hcl
data "tencentcloud_gaap_real_servers_status" "real_servers_status" {
  real_server_ids = ["rs-qcygnwpd"]
}
```

## Argument Reference

The following arguments are supported:

* `real_server_ids` - (Required, Set: [`String`]) Real Server Ids.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `real_server_status_set` - Real Server Status Set.
  * `bind_status` - Bind Status, 0 indicates unbound, 1 indicates bound by rules or listeners.
  * `group_id` - Bind the group ID of this real server, which is an empty string when not bound.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `proxy_id` - Bind the proxy ID of this real server, which is an empty string when not bound.
  * `real_server_id` - Real Server Id.


