---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_commit_certificate_information"
sidebar_current: "docs-tencentcloud-resource-ssl_commit_certificate_information"
description: |-
  Provides a resource to create a ssl commit_certificate_information
---

# tencentcloud_ssl_commit_certificate_information

Provides a resource to create a ssl commit_certificate_information

## Example Usage

```hcl
resource "tencentcloud_ssl_pay_certificate" "example" {
  product_id       = 33
  domain_num       = 1
  alias            = "example-ssl-update"
  project_id       = 0
  wait_commit_flag = true
  information {
    csr_type              = "online"
    certificate_domain    = "www.domain.com"
    organization_name     = "test-update"
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
resource "tencentcloud_ssl_commit_certificate_information" "example" {
  product_id     = 33
  certificate_id = tencentcloud_ssl_pay_certificate.example.certificate_id
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String, ForceNew) Certificate Id.
* `product_id` - (Required, Int, ForceNew) Certificate commodity ID. Valid value ranges: (3~42). `3` means SecureSite enhanced Enterprise Edition (EV Pro), `4` means SecureSite enhanced (EV), `5` means SecureSite Enterprise Professional Edition (OV Pro), `6` means SecureSite Enterprise (OV), `7` means SecureSite Enterprise Type (OV) wildcard, `8` means Geotrust enhanced (EV), `9` means Geotrust enterprise (OV), `10` means Geotrust enterprise (OV) wildcard, `11` means TrustAsia domain type multi-domain SSL certificate, `12` means TrustAsia domain type ( DV) wildcard, `13` means TrustAsia enterprise wildcard (OV) SSL certificate (D3), `14` means TrustAsia enterprise (OV) SSL certificate (D3), `15` means TrustAsia enterprise multi-domain (OV) SSL certificate (D3), `16` means TrustAsia Enhanced (EV) SSL Certificate (D3), `17` means TrustAsia Enhanced Multiple Domain (EV) SSL Certificate (D3), `18` means GlobalSign Enterprise (OV) SSL Certificate, `19` means GlobalSign Enterprise Wildcard (OV) SSL Certificate, `20` means GlobalSign Enhanced (EV) SSL Certificate, `21` means TrustAsia Enterprise Wildcard Multiple Domain (OV) SSL Certificate (D3), `22` means GlobalSign Enterprise Multiple Domain (OV) SSL Certificate, `23` means GlobalSign Enterprise Multiple Wildcard Domain name (OV) SSL certificate, `24` means GlobalSign enhanced multi-domain (EV) SSL certificate, `25` means Wotrus domain type certificate, `26` means Wotrus domain type multi-domain certificate, `27` means Wotrus domain type wildcard certificate, `28` means Wotrus enterprise type certificate, `29` means Wotrus enterprise multi-domain certificate, `30` means Wotrus enterprise wildcard certificate, `31` means Wotrus enhanced certificate, `32` means Wotrus enhanced multi-domain certificate, `33` means WoTrus National Secret Domain name Certificate, `34` means WoTrus National Secret Domain name Certificate (multiple domain names), `35` WoTrus National Secret Domain name Certificate (wildcard), `37` means WoTrus State Secret Enterprise Certificate, `38` means WoTrus State Secret Enterprise Certificate (multiple domain names), `39` means WoTrus State Secret Enterprise Certificate (wildcard), `40` means WoTrus National secret enhanced certificate, `41` means WoTrus National Secret enhanced Certificate (multiple domain names), `42` means TrustAsia- Domain name Certificate (wildcard multiple domain names), `43` means DNSPod Enterprise (OV) SSL Certificate, `44` means DNSPod- Enterprise (OV) wildcard SSL certificate, `45` means DNSPod Enterprise (OV) Multi-domain name SSL Certificate, `46` means DNSPod enhanced (EV) SSL certificate, `47` means DNSPod enhanced (EV) multi-domain name SSL certificate, `48` means DNSPod Domain name Type (DV) SSL Certificate, `49` means DNSPod Domain name Type (DV) wildcard SSL certificate, `50` means DNSPod domain name type (DV) multi-domain name SSL certificate, `51` means DNSPod (State Secret) Enterprise (OV) SSL certificate, `52` DNSPod (National Secret) Enterprise (OV) wildcard SSL certificate, `53` means DNSPod (National Secret) Enterprise (OV) multi-domain SSL certificate, `54` means DNSPod (National Secret) Domain Name (DV) SSL certificate, `55` means DNSPod (National Secret) Domain Name Type (DV) wildcard SSL certificate, `56` means DNSPod (National Secret) Domain Name Type (DV) multi-domain SSL certificate.
* `confirm_letter` - (Optional, String, ForceNew) The base64-encoded certificate confirmation file should be in jpg, jpeg, png, pdf, and the size should be between 1kb and 1.4M. Note: it only works when product_id is set to 8, 9 or 10.
* `verify_type` - (Optional, String, ForceNew) Domain name verification method.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



