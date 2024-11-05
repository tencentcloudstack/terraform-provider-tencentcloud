Provides a resource to create a tcss image registry

Example Usage

```hcl
resource "tencentcloud_tcss_image_registry" "example" {
  name             = "terraform"
  username         = "root"
  password         = "Password@demo"
  url              = "https://example.com"
  registry_type    = "harbor"
  net_type         = "public"
  registry_version = "V1"
  registry_region  = "default"
  need_scan        = true
  conn_detect_config {
    quuid = "backend"
    uuid  = "backend"
  }
}
```
