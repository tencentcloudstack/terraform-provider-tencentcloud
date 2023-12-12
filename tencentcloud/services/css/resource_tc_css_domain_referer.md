Provides a resource to create a css domain_referer

Example Usage

```hcl
resource "tencentcloud_css_domain_referer" "domain_referer" {
  allow_empty = 1
  domain_name = "test122.jingxhu.top"
  enable      = 0
  rules       = "example.com"
  type        = 1
}
```

Import

css domain_referer can be imported using the id, e.g.

```
terraform import tencentcloud_css_domain_referer.domain_referer domainName
```