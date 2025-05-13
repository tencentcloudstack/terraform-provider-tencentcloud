---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_bot_status_config"
sidebar_current: "docs-tencentcloud-resource-waf_bot_status_config"
description: |-
  Provides a resource to create a WAF bot status config
---

# tencentcloud_waf_bot_status_config

Provides a resource to create a WAF bot status config

## Example Usage

```hcl
resource "tencentcloud_waf_bot_status_config" "example" {
  domain = "example.com"
  status = "1"
}
```

### Or

```hcl
resource "tencentcloud_waf_bot_status_config" "example" {
  domain      = "example.com"
  status      = "0"
  instance_id = "waf_2kxtlbky11bbcr4b"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain.
* `status` - (Required, String) Bot status. 1 - enable; 0 - disable.
* `instance_id` - (Optional, String) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `current_global_scene` - The currently enabled scenario with a global matching range and the highest priority.
  * `priority` - Priority.
  * `scene_id` - Scene ID.
  * `scene_name` - Scene name.
  * `update_time` - Update time.
* `custom_rule_nums` - Total number of custom rules, excluding BOT whitelist.
* `scene_count` - Scene total count.
* `valid_scene_count` - Number of effective scenarios.


