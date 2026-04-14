---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_alias_domain"
sidebar_current: "docs-tencentcloud-resource-teo_alias_domain"
description: |-
  Provides a resource to manage TEO alias domain
---

# tencentcloud_teo_alias_domain

Provides a resource to manage TEO alias domain

## Example Usage

```hcl
resource "tencentcloud_teo_zone" "example" {
  zone_name = "example.com"
  type      = "partial"
  area      = "overseas"
  plan_id   = "edgeone-2kfv1h391n6w"
}

resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = tencentcloud_teo_zone.example.id
  alias_name  = "alias.example.com"
  target_name = "example.com"
}
```

### Example Usage with Paused Status

```hcl
resource "tencentcloud_teo_zone" "example" {
  zone_name = "example.com"
  type      = "partial"
  area      = "overseas"
  plan_id   = "edgeone-2kfv1h391n6w"
}

resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = tencentcloud_teo_zone.example.id
  alias_name  = "alias.example.com"
  target_name = "example.com"
  paused      = true
}
```

### format, e.g.:

```hcl
bash
terraform import tencentcloud_teo_alias_domain.example zone-id #alias.example.com
```

## Argument Reference

The following arguments are supported:

* `alias_name` - (Required, String, ForceNew) Alias domain name.
* `target_name` - (Required, String) Target domain name.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `paused` - (Optional, Bool) Indicates whether the alias domain is disabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

