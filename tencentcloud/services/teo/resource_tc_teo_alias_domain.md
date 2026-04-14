Provides a resource to manage TEO alias domain

Example Usage

```hcl
resource "tencentcloud_teo_zone" "example" {
  zone_name = "example.com"
  type        = "partial"
  area        = "overseas"
  plan_id     = "edgeone-2kfv1h391n6w"
}

resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = tencentcloud_teo_zone.example.id
  alias_name  = "alias.example.com"
  target_name = "example.com"
}
```

Example Usage with Paused Status

```hcl
resource "tencentcloud_teo_zone" "example" {
  zone_name = "example.com"
  type        = "partial"
  area        = "overseas"
  plan_id     = "edgeone-2kfv1h391n6w"
}

resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = tencentcloud_teo_zone.example.id
  alias_name  = "alias.example.com"
  target_name = "example.com"
  paused      = true
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Site ID.
* `alias_name` - (Required, ForceNew) Alias domain name.
* `target_name` - (Required) Target domain name.
* `paused` - (Optional, Computed) Indicates whether the alias domain is disabled. Default is `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID in the format of `zone_id#alias_name`.
* `paused` - Whether the alias domain is disabled.

## Import

Alias domain can be imported using the `zone_id#alias_name` format, e.g.:

```bash
terraform import tencentcloud_teo_alias_domain.example zone-id#alias.example.com
```
