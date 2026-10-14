---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_free_certificate"
sidebar_current: "docs-tencentcloud-resource-ssl_free_certificate"
description: |-
  Provide a resource to create a Free Certificate.
---

# tencentcloud_ssl_free_certificate

Provide a resource to create a Free Certificate.

~> **NOTE:** Once certificat created, it cannot be removed within 1 hours.

## Example Usage

### only support type 2. 2=TrustAsia TLS RSA CA.

```hcl
resource "tencentcloud_ssl_free_certificate" "example" {
  dv_auth_method    = "DNS_AUTO"
  domain            = "example.com"
  package_type      = "2"
  contact_email     = "test@example.com"
  contact_phone     = "18352458901"
  validity_period   = 12
  csr_encrypt_algo  = "RSA"
  csr_key_parameter = "2048"
  csr_key_password  = "csr_pwd"
  alias             = "example_free_cert"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Specify domain name.
* `dv_auth_method` - (Required, String) Specify DV authorize method. Available values: `DNS_AUTO` - automatic DNS auth, `DNS` - manual DNS auth, `FILE` - auth by file.
* `alias` - (Optional, String) Specify alias for remark.
* `contact_email` - (Optional, String) Email address.
* `contact_phone` - (Optional, String) Phone number.
* `csr_encrypt_algo` - (Optional, String) Specify CSR encrypt algorithm, only support `RSA` for now.
* `csr_key_parameter` - (Optional, String) Specify CSR key parameter, only support `"2048"` for now.
* `csr_key_password` - (Optional, String) Specify CSR key password.
* `old_certificate_id` - (Optional, String, ForceNew) Specify old certificate ID, used for re-apply.
* `package_type` - (Optional, String) Type of package. Only support `"83"` (TrustAsia TLS RSA CA).
* `project_id` - (Optional, Int) ID of projects which this certification belong to.
* `validity_period` - (Optional, String) Specify validity period in month, only support `"3"` months for now.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cert_begin_time` - Certificate begin time.
* `cert_end_time` - Certificate end time.
* `certificate_private_key` - Certificate private key.
* `certificate_public_key` - Certificate public key.
* `deployable` - Indicates whether the certificate deployable.
* `dv_auths` - DV certification information.
  * `dv_auth_key` - DV authentication key.
  * `dv_auth_value` - DV authentication value.
  * `dv_auth_verify_type` - DV authentication type.
* `insert_time` - Certificate insert time.
* `product_zh_name` - Product zh name.
* `renewable` - Indicates whether the certificate renewable.
* `status_msg` - Certificate status message.
* `status_name` - Certificate status name.
* `status` - Certificate status. 0 = Approving, 1 = Approved, 2 = Approve failed, 3 = expired, 4 = DNS record added, 5 = OV/EV Certificate and confirm letter needed, 6 = Order canceling, 7 = Order canceled, 8 = Submitted and confirm letter needed, 9 = Revoking, 10 = Revoked, 11 = re-applying, 12 = Revoke and confirm letter needed, 13 = Free SSL and confirm letter needed.
* `vulnerability_status` - Vulnerability status.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ApplyCertificate](https://cloud.tencent.com/document/api/400/41678) - This API is used to apply for a free certificate.
* Read: [DescribeCertificateDetail](https://cloud.tencent.com/document/api/400/41673) - This API is used to get certificate details.
* Update: [ModifyCertificateAlias](https://cloud.tencent.com/document/api/400/41669) - This API is used to modify a certificate alias by passing in the certificate ID and new alias.
* Update: [ModifyCertificateProject](https://cloud.tencent.com/document/api/400/41668) - This API is used to modify the projects of multiple certificates.
* Delete: [DeleteCertificate](https://cloud.tencent.com/document/api/400/41675) - This API is used to delete a certificate.

## Import

FreeCertificate instance can be imported, e.g.
```
$ terraform import tencentcloud_ssl_free_certificate.test free_certificate-id
```

