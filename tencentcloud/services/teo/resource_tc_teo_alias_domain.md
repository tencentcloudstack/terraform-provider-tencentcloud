Provides a resource to create a TEO alias domain

Example Usage

```hcl
resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = "zone-39quuimqg8r6"
  alias_name  = "alias.example.com"
  target_name = "www.example.com"
}
```

Example with paused status

```hcl
resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = "zone-39quuimqg8r6"
  alias_name  = "alias.example.com"
  target_name = "www.example.com"
  paused      = true
}
```

Example with timeouts configuration

```hcl
resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = "zone-39quuimqg8r6"
  alias_name  = "alias.example.com"
  target_name = "www.example.com"

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
```

Import

TEO alias domain can be imported using the id, e.g.

```
terraform import tencentcloud_teo_alias_domain.example zone-39quuimqg8r6#alias.example.com
```

Note: The resource ID format is `zone_id#alias_name`.
