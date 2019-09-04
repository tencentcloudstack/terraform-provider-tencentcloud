---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_certificates"
sidebar_current: "docs-tencentcloud-datasource-ssl_certificates"
description: |-
  Use this data source to query SSL certificate.
---

# tencentcloud_ssl_certificates

Use this data source to query SSL certificate.

## Example Usage

```hcl
resource "tencentcloud_ssl_certificate" "foo" {
  name       = "test-ssl-certificate"
  type       = "CA"
  projece_id = 0
  cert       = "-----BEGIN CERTIFICATE-----\nMIIERzCCAq+gAwIBAgIBAjANBgkqhkiG9w0BAQsFADAoMQ0wCwYDVQQDEwR0ZXN0\nMRcwFQYDVQQKEw50ZXJyYWZvcm0gdGVzdDAeFw0xOTA4MTMwMzE5MzlaFw0yOTA4\nMTAwMzE5MzlaMC4xEzARBgNVBAMTCnNlcnZlciBzc2wxFzAVBgNVBAoTDnRlcnJh\nZm9ybS10ZXN0MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEA1Ryp+DKK\nSNFKZsPtwfR+jzOnQ8YFieIKYgakV688d8YgpolenbmeEPrzT87tunFD7G9f6ALG\nND8rj7npj0AowxhOL/h/v1D9u0UsIaj5i2GWJrqNAhGLaxWiEB/hy5WOiwxDrGei\ngQqJkFM52Ep7G1Yx7PHJmKFGwN9FhIsFi1cNZfVRopZuCe/RMPNusNVZaIi+qcEf\nfsE1cmfmuSlG3Ap0RKOIyR0ajDEzqZn9/0R7VwWCF97qy8TNYk94K/1tq3zyhVzR\nZ83xOSfrTqEfb3so3AU2jyKgYdwr/FZS72VCHS8IslgnqJW4izIXZqgIKmHaRZtM\nN4jUloi6l/6lktt6Lsgh9xECecxziSJtPMaog88aC8HnMqJJ3kScGCL36GYG+Kaw\n5PnDlWXBaeiDe8z/eWK9+Rr2M+rhTNxosAVGfDJyxAXyiX49LQ0v7f9qzwc/0JiD\nbvsUv1cm6OgpoEMP9SXqqBdwGqeKbD2/2jlP48xlYP6l1SoJG3GgZ8dbAgMBAAGj\ndjB0MAwGA1UdEwEB/wQCMAAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0PAQH/\nBAUDAweAADAdBgNVHQ4EFgQULwWKBQNLL9s3cb3tTnyPVg+mpCMwHwYDVR0jBBgw\nFoAUKwfrmq791mY831S6UHARHtgYnlgwDQYJKoZIhvcNAQELBQADggGBAMo5RglS\nAHdPgaicWJvmvjjexjF/42b7Rz4pPfMjYw6uYO8He/f4UZWv5CZLrbEe7MywaK3y\n0OsfH8AhyN29pv2x8g9wbmq7omZIOZ0oCAGduEXs/A/qY/hFaCohdkz/IN8qi6JW\nVXreGli3SrpcHFchSwHTyJEXgkutcGAsOvdsOuVSmplOyrkLHc8uUe8SG4j8kGyg\nEzaszFjHkR7g1dVyDVUedc588mjkQxYeAamJgfkgIhljWKMa2XzkVMcVfQHfNpM1\nn+bu8SmqRt9Wma2bMijKRG/Blm756LoI+skY+WRZmlDnq8zj95TT0vceGP0FUWh5\nhKyiocABmpQs9OK9HMi8vgSWISP+fYgkm/bKtKup2NbZBoO5/VL2vCEPInYzUhBO\njCbLMjNjtM5KriCaR7wDARgHiG0gBEPOEW1PIjZ9UOH+LtIxbNZ4eEIIINLHnBHf\nL+doVeZtS/gJc4G4Adr5HYuaS9ZxJ0W2uy0eQlOHzjyxR6Mf/rpnilJlcQ==\n-----END CERTIFICATE-----"
}

data "tencentcloud_gaap_certificates" "foo" {
  id = "${tencentcloud_ssl_certificate.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) ID of the SSL certificate to be queried.
* `name` - (Optional) Name of the SSL certificate to be queried.
* `result_output_file` - (Optional, ForceNew) Used to save results.
* `type` - (Optional) Type of the SSL certificate to be queried. Available values includes: `CA` and `SVR`.


