---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_replace_certificate_operation"
sidebar_current: "docs-tencentcloud-resource-ssl_replace_certificate_operation"
description: |-
  Provides a resource to create a ssl replace_certificate
---

# tencentcloud_ssl_replace_certificate_operation

Provides a resource to create a ssl replace_certificate

## Example Usage

```hcl
resource "tencentcloud_ssl_replace_certificate_operation" "replace_certificate" {
  certificate_id = "8L6JsWq2"
  valid_type     = "DNS_AUTO"
  csr_type       = "online"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String, ForceNew) Certificate ID.
* `valid_type` - (Required, String, ForceNew) Verification type: DNS_AUTO = automatic DNS verification (this verification type is only supported for domain names that are resolved by Tencent Cloud and have normal resolution status), DNS = manual DNS verification, FILE = file verification.
* `cert_csr_encrypt_algo` - (Optional, String, ForceNew) CSR encryption method, optional: RSA, ECC, SM2. (Selectable only if CsrType is Online), default is RSA.
* `cert_csr_key_parameter` - (Optional, String, ForceNew) CSR encryption parameter, when CsrEncryptAlgo is RSA, you can choose 2048, 4096, etc., and the default is 2048; when CsrEncryptAlgo is ECC, you can choose prime256v1, secp384r1, etc., and the default is prime256v1;.
* `csr_content` - (Optional, String, ForceNew) CSR Content.
* `csr_key_password` - (Optional, String, ForceNew) KEY Password.
* `csr_type` - (Optional, String, ForceNew) Type, default Original. Available options: Original = original certificate CSR, Upload = manual upload, Online = online generation.
* `reason` - (Optional, String, ForceNew) Reason for reissue.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ssl replace_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_replace_certificate_operation.replace_certificate replace_certificate_id
```

