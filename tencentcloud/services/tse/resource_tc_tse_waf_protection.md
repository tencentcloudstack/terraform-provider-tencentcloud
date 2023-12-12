Provides a resource to create a tse waf_protection

Example Usage

```hcl
resource "tencentcloud_tse_waf_protection" "waf_protection" {
  gateway_id = "gateway-ed63e957"
  type       = "Route"
  list       = ["7324a769-9d87-48ce-a904-48c3defc4abd"]
  operate    = "open"
}
```