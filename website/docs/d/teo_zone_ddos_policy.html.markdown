---
subcategory: "Teo"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zone_ddos_policy"
sidebar_current: "docs-tencentcloud-datasource-teo_zone_ddos_policy"
description: |-
  Use this data source to query zone ddos policy.
---

# tencentcloud_teo_zone_ddos_policy

Use this data source to query zone ddos policy.

## Example Usage

```hcl
data "tencentcloud_teo_zone_ddos_policy" "example" {
  zone_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used for save results.
* `zone_id` - (Optional, String) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `app_id` - App ID.
* `domains` - All subdomain info. Note: This field may return null, indicating that no valid value can be obtained.
  * `accelerate_type` - on: Enable; off: Disable.
  * `host` - Subdomain.
  * `security_type` - on: Enable; off: Disable.
  * `status` - Status of the subdomain. Note: This field may return null, indicating that no valid value can be obtained, init: waiting to config NS; offline: waiting to enable site accelerating; process: config deployment processing; online: normal status.
* `shield_areas` - Shield areas of the zone.
  * `application` - Layer 7 Domain Name Parameters.
    * `accelerate_type` - on: Enable; off: Disable.
    * `host` - Subdomain.
    * `security_type` - on: Enable; off: Disable.
    * `status` - Status of the subdomain. Note: This field may return null, indicating that no valid value can be obtained, init: waiting to config NS; offline: waiting to enable site accelerating; process: config deployment processing; online: normal status.
  * `entity_name` - When `Type` is `domain`, this field is `ZoneName`. When `Type` is `application`, this field is `ProxyName`. Note: This field may return null, indicating that no valid value can be obtained.
  * `entity` - When `Type` is `domain`, this field is `ZoneId`. When `Type` is `application`, this field is `ProxyId`. Note: This field may return null, indicating that no valid value can be obtained.
  * `policy_id` - Policy ID.
  * `share` - Whether the resource is shared.
  * `tcp_num` - TCP forwarding rule number of layer 4 application.
  * `type` - Valid values: `domain`, `application`.
  * `udp_num` - UDP forwarding rule number of layer 4 application.
  * `zone_id` - Site ID.


