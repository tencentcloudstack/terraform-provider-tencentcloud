---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zone_ddos_policy"
sidebar_current: "docs-tencentcloud-datasource-teo_zone_ddos_policy"
description: |-
  Use this data source to query detailed information of teo zoneDDoSPolicy
---

# tencentcloud_teo_zone_ddos_policy

Use this data source to query detailed information of teo zoneDDoSPolicy

## Example Usage

```hcl
data "tencentcloud_teo_zone_ddos_policy" "zoneDDoSPolicy" {
  zone_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Site ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domains` - All subdomain info. Note: This field may return null, indicating that no valid value can be obtained.
  * `accelerate_type` - Acceleration function switch. Valid values:- `on`: Enable.- `off`: Disable.
  * `host` - Subdomain.
  * `security_type` - Security function switch. Valid values:- `on`: Enable.- `off`: Disable.
  * `status` - Status of the subdomain. Valid values:- `init`: waiting to config NS.- `offline`: need to enable site accelerating.- `process`: processing the config deployment.- `online`: normal status. Note: This field may return null, indicating that no valid value can be obtained.
* `shield_areas` - Shielded areas of the zone.
  * `application` - DDoS layer 7 application.
    * `accelerate_type` - Acceleration function switch. Valid values:- `on`: Enable.- `off`: Disable.
    * `host` - Subdomain.
    * `security_type` - Security function switch. Valid values:- `on`: Enable.- `off`: Disable.
    * `status` - Status of the subdomain. Valid values:- `init`: waiting to config NS.- `offline`: need to enable site accelerating.- `process`: processing the config deployment.- `online`: normal status. Note: This field may return null, indicating that no valid value can be obtained.
  * `entity_name` - When `Type` is `domain`, this field is `ZoneName`. When `Type` is `application`, this field is `ProxyName`. Note: This field may return null, indicating that no valid value can be obtained.
  * `entity` - When `Type` is `domain`, this field is `ZoneId`. When `Type` is `application`, this field is `ProxyId`. Note: This field may return null, indicating that no valid value can be obtained.
  * `policy_id` - Policy ID.
  * `tcp_num` - TCP forwarding rule number of layer 4 application.
  * `type` - Valid values: `domain`, `application`.
  * `udp_num` - UDP forwarding rule number of layer 4 application.
  * `zone_id` - Site ID.


