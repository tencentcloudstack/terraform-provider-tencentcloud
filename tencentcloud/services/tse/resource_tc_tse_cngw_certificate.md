Provides a resource to create a tse cngw_certificate

Example Usage

```hcl

resource "tencentcloud_tse_cngw_certificate" "cngw_certificate" {
  gateway_id   = "gateway-ddbb709b"
  bind_domains = ["example1.com"]
  cert_id      = "vYSQkJ3K"
  name         = "xxx1"
}

```

Import

tse cngw_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_certificate.cngw_certificate gatewayId#Id
```