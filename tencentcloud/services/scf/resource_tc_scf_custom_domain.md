Provides a resource to create a scf custom domain

Example Usage

```hcl
resource "tencentcloud_scf_custom_domain" "scf_custom_domain" {
  domain ="xxxxxx"
  protocol = "HTTP"
  endpoints_config {
    namespace = "default"
    function_name = "xxxxxx"
    qualifier = "$LATEST"
    path_match = "/aa/*"
  }
  waf_config {
    waf_open = "CLOSE"
  }
}
```

Import

scf scf_custom_domain can be imported using the id, e.g.

```
terraform import tencentcloud_scf_custom_domain.scf_custom_domain ${domain}
```
