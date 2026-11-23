Provides a resource to create a teo application proxy

Example Usage

```hcl
resource "tencentcloud_teo_application_proxy" "example" {
    accelerate_type      = 0
    plat_type            = "domain"
    proxy_name           = "test"
    proxy_type           = "instance"
    security_type        = 1
    session_persist_time = 0
    status               = "online"
    zone_id              = "zone-2o0l8g7zisgt"

    ipv6 {
        switch = "off"
    }
}
```

Import

teo application_proxy can be imported using the zoneId#proxyId, e.g.
```
terraform import tencentcloud_teo_application_proxy.example zone-2983wizgxqvm#proxy-6972528a-373a-11ed-afca-52540044a456
```
