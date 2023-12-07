Use this data source to query detailed information of tse gateway_certificates

Example Usage

```hcl

data "tencentcloud_tse_gateway_certificates" "gateway_certificates" {
  gateway_id = "gateway-ddbb709b"
  filters {
    key = "BindDomain"
    value = "example.com"
  }
}
```