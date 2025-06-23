Provides a resource to create a TEO acceleration domain

~> **NOTE:** Before modifying resource content, you need to ensure that the `status` is `online`.

Example Usage

```hcl
resource "tencentcloud_teo_acceleration_domain" "example" {
  zone_id     = "zone-39quuimqg8r6"
  domain_name = "www.demo.com"

  origin_info {
    origin      = "150.109.8.1"
    origin_type = "IP_DOMAIN"
  }

  status            = "online"
  origin_protocol   = "FOLLOW"
  http_origin_port  = 80
  https_origin_port = 443
  ipv6_status       = "follow"
}
```

Import

TEO acceleration domain can be imported using the id, e.g.

```
terraform import tencentcloud_teo_acceleration_domain.example zone-39quuimqg8r6#www.demo.com
```
