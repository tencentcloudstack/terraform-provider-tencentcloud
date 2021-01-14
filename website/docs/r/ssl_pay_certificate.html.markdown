---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_pay_certificate"
sidebar_current: "docs-tencentcloud-resource-ssl_pay_certificate"
description: |-
  Provide a resource to create a payment SSL.
---

# tencentcloud_ssl_pay_certificate

Provide a resource to create a payment SSL.

~> **NOTE:** Provides the creation of a paid certificate, including the submission of certificate information and order functions;
currently, it does not support re-issuing certificates, revoking certificates, and deleting certificates; the certificate remarks
and belonging items can be updated. The Destroy operation will only cancel the certificate order, and will not delete the
certificate and refund the fee. If you need a refund, you need to check the current certificate status in the console
as `Review Cancel`, and then you can click `Request a refund` to refund the fee.

## Example Usage

```hcl
resource "tencentcloud_ssl_pay_certificate" "ssl" {
  product_id = 33
  domain_num = 1
  alias      = "test-ssl"
  project_id = 0
  information {
    csr_type              = "online"
    certificate_domain    = "www.domain.com"
    organization_name     = "test"
    organization_division = "test"
    organization_address  = "test"
    organization_country  = "CN"
    organization_city     = "test"
    organization_region   = "test"
    postal_code           = "0755"
    phone_area_code       = "0755"
    phone_number          = "12345678901"
    verify_type           = "DNS"
    admin_first_name      = "test"
    admin_last_name       = "test"
    admin_phone_num       = "12345678901"
    admin_email           = "test@tencent.com"
    admin_position        = "dev"
    contact_first_name    = "test"
    contact_last_name     = "test"
    contact_email         = "test@tencent.com"
    contact_number        = "12345678901"
    contact_position      = "dev"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_num` - (Required, ForceNew) Number of domain names included in the certificate.
* `information` - (Required, ForceNew) Certificate information.
* `product_id` - (Required, ForceNew) Certificate commodity ID. Valid value ranges: (3~42). `3` means SecureSite Enhanced Enterprise Edition (EV Pro), `4` means SecureSite Enhanced (EV), `5` means SecureSite Enterprise Professional Edition (OV Pro), `6` means SecureSite Enterprise (OV), `7` means SecureSite Enterprise Type (OV) wildcard, `8` means Geotrust enhanced (EV), `9` means Geotrust enterprise (OV), `10` means Geotrust enterprise (OV) wildcard, `11` means TrustAsia domain type multi-domain SSL certificate, `12` means TrustAsia domain type ( DV) wildcard, `13` means TrustAsia enterprise wildcard (OV) SSL certificate (D3), `14` means TrustAsia enterprise (OV) SSL certificate (D3), `15` means TrustAsia enterprise multi-domain (OV) SSL certificate (D3), `16` means TrustAsia Enhanced (EV) SSL Certificate (D3), `17` means TrustAsia Enhanced Multiple Domain (EV) SSL Certificate (D3), `18` means GlobalSign Enterprise (OV) SSL Certificate, `19` means GlobalSign Enterprise Wildcard (OV) SSL Certificate, `20` means GlobalSign Enhanced (EV) SSL Certificate, `21` means TrustAsia Enterprise Wildcard Multiple Domain (OV) SSL Certificate (D3), `22` means GlobalSign Enterprise Multiple Domain (OV) SSL Certificate, `23` means GlobalSign Enterprise Multiple Wildcard Domain name (OV) SSL certificate, `24` means GlobalSign enhanced multi-domain (EV) SSL certificate, `25` means Wotrus domain type certificate, `26` means Wotrus domain type multi-domain certificate, `27` means Wotrus domain type wildcard certificate, `28` means Wotrus enterprise type certificate, `29` means Wotrus enterprise multi-domain certificate, `30` means Wotrus enterprise wildcard certificate, `31` means Wotrus enhanced certificate, `32` means Wotrus enhanced multi-domain certificate, `33` means DNSPod national secret domain name certificate, `34` means DNSPod national secret domain name certificate Multi-domain certificate, `35` means DNSPod national secret domain name wildcard certificate, `37` means DNSPod national secret enterprise certificate, `38` means DNSPod national secret enterprise multi-domain certificate, `39` means DNSPod national secret enterprise wildcard certificate, `40` means DNSPod national secret increase Strong certificate, `41` means DNSPod national secret enhanced multi-domain certificate, `42` means TrustAsia domain-type wildcard multi-domain certificate.
* `alias` - (Optional) Remark name.
* `project_id` - (Optional) The ID of project.
* `time_span` - (Optional) Certificate period, currently only supports 1 year certificate purchase.

The `information` object supports the following:

* `admin_email` - (Required, ForceNew) The administrator's email address.
* `admin_first_name` - (Required, ForceNew) The first name of the administrator.
* `admin_last_name` - (Required, ForceNew) The last name of the administrator.
* `admin_phone_num` - (Required, ForceNew) Manager mobile phone number.
* `admin_position` - (Required, ForceNew) Manager position.
* `certificate_domain` - (Required, ForceNew) Domain name for binding certificate.
* `contact_email` - (Required, ForceNew) Contact email address.
* `contact_first_name` - (Required, ForceNew) Contact first name.
* `contact_last_name` - (Required, ForceNew) Contact last name.
* `contact_number` - (Required, ForceNew) Contact phone number.
* `contact_position` - (Required, ForceNew) Contact position.
* `organization_address` - (Required, ForceNew) Company address.
* `organization_city` - (Required, ForceNew) Company city.
* `organization_country` - (Required, ForceNew) Country name, such as China: CN.
* `organization_division` - (Required, ForceNew) Department name.
* `organization_name` - (Required, ForceNew) Company name.
* `organization_region` - (Required, ForceNew) The province where the company is located.
* `phone_area_code` - (Required, ForceNew) Company landline area code.
* `phone_number` - (Required, ForceNew) Company landline number.
* `postal_code` - (Required, ForceNew) Company postal code.
* `verify_type` - (Required, ForceNew) Certificate verification method. Valid values: `DNS_AUTO`, `DNS`, `FILE`. `DNS_AUTO` means automatic DNS verification, this verification type is only supported for domain names resolved by Tencent Cloud and the resolution status is normal, `DNS` means manual DNS verification, `FILE` means file verification.
* `csr_content` - (Optional, ForceNew) CSR content uploaded.
* `csr_type` - (Optional, ForceNew) CSR generation method. Valid values: `online`, `parse`. `online` means online generation, `parse` means manual upload.
* `domain_list` - (Optional, ForceNew) Array of uploaded domain names, multi-domain certificates can be uploaded.
* `key_password` - (Optional, ForceNew) Private key password.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `certificate_id` - Returned certificate ID.
* `order_id` - Order ID returned.
* `status` - SSL certificate status.


## Import

payment SSL instance can be imported, e.g.

```
$ terraform import tencentcloud_ssl_pay_certificate.ssl iPQNn61x#33#1#1
```

