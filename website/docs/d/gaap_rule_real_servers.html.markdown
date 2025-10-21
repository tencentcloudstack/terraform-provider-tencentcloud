---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_rule_real_servers"
sidebar_current: "docs-tencentcloud-datasource-gaap_rule_real_servers"
description: |-
  Use this data source to query detailed information of gaap rule real servers
---

# tencentcloud_gaap_rule_real_servers

Use this data source to query detailed information of gaap rule real servers

## Example Usage

```hcl
data "tencentcloud_gaap_rule_real_servers" "rule_real_servers" {
  rule_id = "rule-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `rule_id` - (Required, String) Rule Id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bind_real_server_set` - Bind Real Server info.
  * `down_ip_list` - When the real server is a domain name, the domain name is resolved to one or more IPs, and this field represents the list of abnormal IPs. When the status is abnormal, but the field is empty, it indicates that the domain name resolution is abnormal.
  * `real_server_failover_role` - The primary and secondary roles of the real server:master represents primary, slave represents secondary, and this parameter must be in the active and standby mode of the real server when the listener is turned on.
  * `real_server_id` - Real Server Id.
  * `real_server_ip` - Real Server Ip or domain.
  * `real_server_port` - Real Server PortNote: This field may return null, indicating that a valid value cannot be obtained.
  * `real_server_status` - RealServerStatus: 0 indicates normal;1 indicates an exception.When the health check status is not enabled, it is always normal.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `real_server_weight` - Real Server Weight.
* `real_server_set` - Real Server Set.
  * `in_ban_blacklist` - Is it on the banned blacklist? 0 indicates not on the blacklist, and 1 indicates on the blacklist.
  * `project_id` - Project Id.
  * `real_server_id` - Real Server Id.
  * `real_server_ip` - Real Server IP or domain.
  * `real_server_name` - Real Server Name.


