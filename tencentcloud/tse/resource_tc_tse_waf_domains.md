Provides a resource to create a tse waf_domains

Example Usage

```hcl
resource "tencentcloud_tse_waf_domains" "waf_domains" {
  domain     = "tse.exmaple.com"
  gateway_id = "gateway-ed63e957"
}
```

Import

tse waf_domains can be imported using the id, e.g.

```
terraform import tencentcloud_tse_waf_domains.waf_domains waf_domains_id
```