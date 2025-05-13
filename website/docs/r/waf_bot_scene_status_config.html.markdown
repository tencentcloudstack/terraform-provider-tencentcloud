---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_bot_scene_status_config"
sidebar_current: "docs-tencentcloud-resource-waf_bot_scene_status_config"
description: |-
  Provides a resource to create a WAF bot scene status config
---

# tencentcloud_waf_bot_scene_status_config

Provides a resource to create a WAF bot scene status config

## Example Usage

```hcl
resource "tencentcloud_waf_bot_scene_status_config" "example" {
  domain   = "example.com"
  scene_id = "3024324123"
  status   = true
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain.
* `scene_id` - (Required, String, ForceNew) Scene ID.
* `status` - (Required, Bool) Bot status. true - enable; false - disable.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `priority` - Priority.
* `scene_name` - Scene name.
* `type` - Scene type, default: Default scenario, custom: Non default scenario.


## Import

WAF bot scene status config can be imported using the id, e.g.

```
terraform import tencentcloud_waf_bot_scene_status_config.example example.com#3024324123
```

