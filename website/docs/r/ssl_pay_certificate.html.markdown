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

* `domain_num` - (Required, Int, ForceNew) Number of domain names included in the certificate.
* `information` - (Required, List, ForceNew) Certificate information.
* `product_id` - (Required, Int, ForceNew) Certificate commodity ID. Valid value ranges: (3~42). `3` means SecureSite enhanced Enterprise Edition (EV Pro), `4` means SecureSite enhanced (EV), `5` means SecureSite Enterprise Professional Edition (OV Pro), `6` means SecureSite Enterprise (OV), `7` means SecureSite Enterprise Type (OV) wildcard, `8` means Geotrust enhanced (EV), `9` means Geotrust enterprise (OV), `10` means Geotrust enterprise (OV) wildcard, `11` means TrustAsia domain type multi-domain SSL certificate, `12` means TrustAsia domain type ( DV) wildcard, `13` means TrustAsia enterprise wildcard (OV) SSL certificate (D3), `14` means TrustAsia enterprise (OV) SSL certificate (D3), `15` means TrustAsia enterprise multi-domain (OV) SSL certificate (D3), `16` means TrustAsia Enhanced (EV) SSL Certificate (D3), `17` means TrustAsia Enhanced Multiple Domain (EV) SSL Certificate (D3), `18` means GlobalSign Enterprise (OV) SSL Certificate, `19` means GlobalSign Enterprise Wildcard (OV) SSL Certificate, `20` means GlobalSign Enhanced (EV) SSL Certificate, `21` means TrustAsia Enterprise Wildcard Multiple Domain (OV) SSL Certificate (D3), `22` means GlobalSign Enterprise Multiple Domain (OV) SSL Certificate, `23` means GlobalSign Enterprise Multiple Wildcard Domain name (OV) SSL certificate, `24` means GlobalSign enhanced multi-domain (EV) SSL certificate, `25` means Wotrus domain type certificate, `26` means Wotrus domain type multi-domain certificate, `27` means Wotrus domain type wildcard certificate, `28` means Wotrus enterprise type certificate, `29` means Wotrus enterprise multi-domain certificate, `30` means Wotrus enterprise wildcard certificate, `31` means Wotrus enhanced certificate, `32` means Wotrus enhanced multi-domain certificate, `33` means WoTrus National Secret Domain name Certificate, `34` means WoTrus National Secret Domain name Certificate (multiple domain names), `35` WoTrus National Secret Domain name Certificate (wildcard), `37` means WoTrus State Secret Enterprise Certificate, `38` means WoTrus State Secret Enterprise Certificate (multiple domain names), `39` means WoTrus State Secret Enterprise Certificate (wildcard), `40` means WoTrus National secret enhanced certificate, `41` means WoTrus National Secret enhanced Certificate (multiple domain names), `42` means TrustAsia- Domain name Certificate (wildcard multiple domain names), `43` means DNSPod Enterprise (OV) SSL Certificate, `44` means DNSPod- Enterprise (OV) wildcard SSL certificate, `45` means DNSPod Enterprise (OV) Multi-domain name SSL Certificate, `46` means DNSPod enhanced (EV) SSL certificate, `47` means DNSPod enhanced (EV) multi-domain name SSL certificate, `48` means DNSPod Domain name Type (DV) SSL Certificate, `49` means DNSPod Domain name Type (DV) wildcard SSL certificate, `50` means DNSPod domain name type (DV) multi-domain name SSL certificate, `51` means DNSPod (State Secret) Enterprise (OV) SSL certificate, `52` DNSPod (National Secret) Enterprise (OV) wildcard SSL certificate, `53` means DNSPod (National Secret) Enterprise (OV) multi-domain SSL certificate, `54` means DNSPod (National Secret) Domain Name (DV) SSL certificate, `55` means DNSPod (National Secret) Domain Name Type (DV) wildcard SSL certificate, `56` means DNSPod (National Secret) Domain Name Type (DV) multi-domain SSL certificate.
* `alias` - (Optional, String) Remark name.
* `confirm_letter` - (Optional, String) The base64-encoded certificate confirmation file should be in jpg, jpeg, png, pdf, and the size should be between 1kb and 1.4M. Note: it only works when product_id is set to 8, 9 or 10.
* `project_id` - (Optional, Int) The ID of project.
* `time_span` - (Optional, Int) Certificate period, currently only supports 1 year certificate purchase.

The `information` object supports the following:

* `admin_email` - (Required, String, ForceNew) The administrator's email address.
* `admin_first_name` - (Required, String, ForceNew) The first name of the administrator.
* `admin_last_name` - (Required, String, ForceNew) The last name of the administrator.
* `admin_phone_num` - (Required, String, ForceNew) Manager mobile phone number.
* `admin_position` - (Required, String, ForceNew) Manager position.
* `certificate_domain` - (Required, String, ForceNew) Domain name for binding certificate.
* `contact_email` - (Required, String, ForceNew) Contact email address.
* `contact_first_name` - (Required, String, ForceNew) Contact first name.
* `contact_last_name` - (Required, String, ForceNew) Contact last name.
* `contact_number` - (Required, String, ForceNew) Contact phone number.
* `contact_position` - (Required, String, ForceNew) Contact position.
* `organization_address` - (Required, String, ForceNew) Company address.
* `organization_city` - (Required, String, ForceNew) Company city.
* `organization_country` - (Required, String, ForceNew) Country name, such as China: CN.
* `organization_division` - (Required, String, ForceNew) Department name.
* `organization_name` - (Required, String, ForceNew) Company name.
* `organization_region` - (Required, String, ForceNew) The province where the company is located.
* `phone_area_code` - (Required, String, ForceNew) Company landline area code.
* `phone_number` - (Required, String, ForceNew) Company landline number.
* `postal_code` - (Required, String, ForceNew) Company postal code.
* `verify_type` - (Required, String, ForceNew) Certificate verification method. Valid values: `DNS_AUTO`, `DNS`, `FILE`. `DNS_AUTO` means automatic DNS verification, this verification type is only supported for domain names resolved by Tencent Cloud and the resolution status is normal, `DNS` means manual DNS verification, `FILE` means file verification.
* `csr_content` - (Optional, String, ForceNew) CSR content uploaded.
* `csr_type` - (Optional, String, ForceNew) CSR generation method. Valid values: `online`, `parse`. `online` means online generation, `parse` means manual upload.
* `domain_list` - (Optional, Set, ForceNew) Array of uploaded domain names, multi-domain certificates can be uploaded.
* `key_password` - (Optional, String, ForceNew) Private key password.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `certificate_id` - Returned certificate ID.
* `dv_auths` - DV certification information.
  * `dv_auth_key` - DV authentication key.
  * `dv_auth_value` - DV authentication value.
  * `dv_auth_verify_type` - DV authentication type.
* `order_id` - Order ID returned.
* `status` - SSL certificate status.


## Import

payment SSL instance can be imported, e.g.

```
$ terraform import tencentcloud_ssl_pay_certificate.ssl iPQNn61x#33#1#1
```

