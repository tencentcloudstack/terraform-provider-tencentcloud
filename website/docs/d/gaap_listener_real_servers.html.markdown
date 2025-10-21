---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_listener_real_servers"
sidebar_current: "docs-tencentcloud-datasource-gaap_listener_real_servers"
description: |-
  Use this data source to query detailed information of gaap listener real servers
---

# tencentcloud_gaap_listener_real_servers

Use this data source to query detailed information of gaap listener real servers

## Example Usage

```hcl
data "tencentcloud_gaap_listener_real_servers" "listener_real_servers" {
  listener_id = "listener-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, String) listener ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bind_real_server_set` - Bound real server Information List.
  * `down_i_p_list` - When the real server is a domain name, the domain name is resolved to one or more IPs, and this field represents the list of abnormal IPs. When the status is abnormal, but the field is empty, it indicates that the domain name resolution is abnormal.
  * `real_server_failover_role` - The primary and secondary roles of the real server, &#39;master&#39; represents primary, &#39;slave&#39; represents secondary, and this parameter must be in the active and standby mode of the real server when the listener is turned on.
  * `real_server_i_p` - Real Server IP or domain.
  * `real_server_id` - Real Server Id.
  * `real_server_port` - The port number of the real serverNote: This field may return null, indicating that a valid value cannot be obtained.
  * `real_server_status` - real server health check status, where:0 indicates normal;1 indicates an exception.When the health check status is not enabled, it is always normal.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `real_server_weight` - The weight of this real server.
* `real_server_set` - Real Server Set.
  * `in_ban_blacklist` - Is it on the banned blacklist? 0 indicates not on the blacklist, and 1 indicates on the blacklist.
  * `project_id` - Project Id.
  * `real_server_i_p` - Real Server IP.
  * `real_server_id` - Real Server Id.
  * `real_server_name` - Real Server Name.


