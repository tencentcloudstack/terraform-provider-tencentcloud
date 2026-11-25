---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_custom_group_references"
sidebar_current: "docs-tencentcloud-data_source-teo_custom_group_references"
description: |-
  Use this data source to query the reference information of a specified TEO (EdgeOne) custom group.
---

# tencentcloud_teo_custom_group_references

Use this data source to query the reference information of a specified TEO (EdgeOne) custom group.

## Example Usage

```hcl
data "tencentcloud_teo_custom_group_references" "example" {
  zone_id  = "zone-3fkff38fyw8s"
  group_id = "cg-2q1pkwfyagzz"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) Custom group ID.
* `zone_id` - (Required, String) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `entity_id` - Entity ID. Its meaning depends on EntityType: zone ID for `WebSec.ZonePolicy`, domain for `WebSec.HostPolicy`, template ID for `WebSec.Template`, API resource ID for `WebSec.APIResource`, JavaScript injection rule ID for `WebSec.JSInjectionRule`, instance ID for `DDoS.L4Proxy` and `DDoS.L3Transit`.
  * `entity_name` - Entity name. Its meaning depends on EntityType: template name for `WebSec.Template`, API resource name for `WebSec.APIResource`, JavaScript injection rule name for `WebSec.JSInjectionRule`; empty for `WebSec.ZonePolicy`, `WebSec.HostPolicy`, `DDoS.L4Proxy` and `DDoS.L3Transit`.
  * `entity_type` - Entity type. Valid values: `WebSec.ZonePolicy` (zone-level security policy), `WebSec.HostPolicy` (host-level security policy), `WebSec.Template` (policy template), `WebSec.APIResource` (API resource), `WebSec.JSInjectionRule` (JavaScript injection rule), `DDoS.L4Proxy` (layer-4 proxy DDoS protection), `DDoS.L3Transit` (layer-3 transit DDoS protection).
  * `sub_entity_id` - Sub entity ID, i.e. the rule ID of the corresponding SubEntityType; empty for `DDoS.L4Proxy.IpAccessControl` and `DDoS.L3Transit.IpAccessControl`.
  * `sub_entity_name` - Sub entity name, i.e. the rule name of the corresponding SubEntityType; empty for `DDoS.L4Proxy.IpAccessControl` and `DDoS.L3Transit.IpAccessControl`.
  * `sub_entity_type` - Sub entity type. Valid values: `WebSec.ExceptionRule` (protection exception rule), `WebSec.BasicAccessRule` (basic access control), `WebSec.PreciseMatchRule` (precise match rule), `WebSec.RateLimitRule` (rate limit rule), `WebSec.BotCustomRule` (custom bot rule), `DDoS.L4Proxy.IpAccessControl` (layer-4 proxy IP allowlist/blocklist), `DDoS.L3Transit.IpAccessControl` (layer-3 transit IP allowlist/blocklist).
  * `zone_id` - Site ID.
* `references` - List of configurations that reference the custom group.


