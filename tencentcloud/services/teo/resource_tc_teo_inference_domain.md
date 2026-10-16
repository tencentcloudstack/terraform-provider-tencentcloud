Provides a resource to create a TEO inference domain.

~> **NOTE:** Before creating an inference domain, please verify the domain ownership of the site.

Example Usage

```hcl
resource "tencentcloud_teo_inference_domain" "example" {
  zone_id    = "zone-39quuimqg8r6"
  service_id = "sid-xxxxxx"
  domain     = "www.demo.com"
  auth_switch = "On"
}
```

Import

TEO inference domain can be imported using the id, e.g.

```
terraform import tencentcloud_teo_inference_domain.example zone-39quuimqg8r6#sid-xxxxxx#www.demo.com
```