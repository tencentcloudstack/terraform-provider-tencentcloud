Provides a resource to create a teo acceleration_domain

Example Usage

```hcl
resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
    zone_id     = "zone-2o0i41pv2h8c"
    domain_name = "aaa.makn.cn"

    origin_info {
        origin      = "150.109.8.1"
        origin_type = "IP_DOMAIN"
    }
}
```

Import

teo acceleration_domain can be imported using the id, e.g.

```
terraform import tencentcloud_teo_acceleration_domain.acceleration_domain acceleration_domain_id
```