Provides a resource to create a ssl check_certificate_chain

Example Usage

```hcl
resource "tencentcloud_ssl_check_certificate_chain_operation" "check_certificate_chain" {
  certificate_chain = "-----BEGIN CERTIFICATE--·····---END CERTIFICATE-----"
}
```

Import

ssl check_certificate_chain can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_check_certificate_chain_operation.check_certificate_chain check_certificate_chain_id
```