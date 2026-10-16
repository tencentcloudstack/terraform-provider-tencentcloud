Provides a resource to create a TEO inference domain.

~> **NOTE:** Before creating an inference domain, please ensure that the ownership of the domain has been verified and the domain has been bound to the site.

~> **NOTE:** The inference domain must point to an existing inference service under the same site.

Example Usage

```hcl
resource "tencentcloud_teo_inference_domain" "example" {
  zone_id        = "zone-3fkff38fyw8s"
  domain_name    = "inference.demo.com"
  service_id     = "is-xxxxxxxx"
  origin_protocol = "FOLLOW"
  http_origin_port  = 80
  https_origin_port = 443
}
```

Import

TEO inference domain can be imported using the zoneId#domainName, e.g.

```
terraform import tencentcloud_teo_inference_domain.example zone-3fkff38fyw8s#inference.demo.com
```