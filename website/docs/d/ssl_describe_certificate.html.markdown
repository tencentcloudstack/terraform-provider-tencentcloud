---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_certificate"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_certificate"
description: |-
  Use this data source to query detailed information of ssl describe_certificate
---

# tencentcloud_ssl_describe_certificate

Use this data source to query detailed information of ssl describe_certificate

## Example Usage

```hcl
data "tencentcloud_ssl_describe_certificate" "describe_certificate" {
  certificate_id = "8cj4g8h8"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String) Certificate ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - result list.
  * `alias` - Remark name.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `c_a_common_names` - All general names of the CA certificateNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `c_a_encrypt_algorithms` - All encryption methods of CA certificateNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `c_a_end_times` - CA certificate all maturity timeNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `cert_begin_time` - Certificate takes effect time.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `cert_end_time` - The certificate is invalid time.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `certificate_extra` - Certificate extension information.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `company_type` - Type of company. Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `domain_number` - Certificate can be configured in the number of domain names.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `origin_certificate_id` - Original certificate ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `renew_order` - New order certificate ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `replaced_by` - Re -issue the original ID of the certificate.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `replaced_for` - Re -issue a new ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `s_m_cert` - Is it a national secret certificateNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `certificate_type` - Certificate type: CA = CA certificate, SVR = server certificate.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `deployable` - Whether it can be deployed.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `domain` - domain name.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `dv_auth_detail` - DV certification information.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `dv_auth_domain` - DV authentication value domain name.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `dv_auth_key_sub_domain` - DV certification sub -domain name.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `dv_auth_key` - DV certification key.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `dv_auth_path` - DV authentication value path.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `dv_auth_value` - DV certification value.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `dv_auths` - DV certification information.Note: This field may return NULL, indicating that the valid value cannot be obtained.
      * `dv_auth_domain` - DV authentication value domain name.Note: This field may return NULL, indicating that the valid value cannot be obtained.
      * `dv_auth_key` - DV certification key.Note: This field may return NULL, indicating that the valid value cannot be obtained.
      * `dv_auth_path` - DV authentication value path.Note: This field may return NULL, indicating that the valid value cannot be obtained.
      * `dv_auth_sub_domain` - DV certification sub -domain name,Note: This field may return NULL, indicating that the valid value cannot be obtained.
      * `dv_auth_value` - DV certification value.Note: This field may return NULL, indicating that the valid value cannot be obtained.
      * `dv_auth_verify_type` - DV certification type.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `dv_revoke_auth_detail` - DV certificate revoking verification valueNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `dv_auth_domain` - DV authentication value domain name.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `dv_auth_key` - DV certification key.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `dv_auth_path` - DV authentication value path.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `dv_auth_sub_domain` - DV certification sub -domain name,Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `dv_auth_value` - DV certification value.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `dv_auth_verify_type` - DV certification type.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `from` - Certificate source: Trustasia,uploadNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `insert_time` - application time.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `is_dv` - Whether it is the DV version.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `is_vip` - Whether it is a VIP customer.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `is_vulnerability` - Whether the vulnerability scanning function is enabled.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `is_wildcard` - Whether it is a pan -domain certificate certificate.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `order_id` - Order ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `owner_uin` - Account UIN.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `package_type_name` - Certificate type name.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `package_type` - Types of Certificate Package: 1 = Geotrust DV SSL CA -G3, 2 = Trustasia TLS RSA CA, 3 = SecureSite Enhanced Enterprise Edition (EV Pro), 4 = SecureSite enhanced (EV), 5 = SecureSite Enterprise Professional Edition (OVPro), 6 = SecureSite Enterprise (OV), 7 = SecureSite Enterprise (OV) compatriots, 8 = Geotrust enhanced type (EV), 9 = Geotrust Enterprise (OV), 10 = Geotrust Enterprise (OV) pass,11 = Trustasia Domain Multi -domain SSL certificate, 12 = Trustasia domain model (DV) passing, 13 = Trustasia Enterprise Passing Character (OV) SSL certificate (D3), 14 = Trustasia Enterprise (OV) SSL certificate (D3), 15= Trustasia Enterprise Multi -domain name (OV) SSL certificate (D3), 16 = Trustasia enhanced (EV) SSL certificate (D3), 17 = Trustasia enhanced multi -domain name (EV) SSL certificate (D3), 18 = GlobalSign enterprise type enterprise type(OV) SSL certificate, 19 = GlobalSign Enterprise Type -type STL Certificate, 20 = GlobalSign enhanced (EV) SSL certificate, 21 = Trustasia Enterprise Tongzhi Multi -domain name (OV) SSL certificate (D3), 22 = GlobalSignignMulti -domain name (OV) SSL certificate, 23 = GlobalSign Enterprise Type -type multi -domain name (OV) SSL certificate, 24 = GlobalSign enhanced multi -domain name (EV) SSL certificate.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `product_zh_name` - Certificate issuer name.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `project_id` - Project ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `renew_able` - Whether you can issue a certificate.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `status_msg` - status information.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `status_name` - status description.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `status` - = Submitted information, to be uploaded to confirmation letter, 9 = Certificate is revoked, 10 = revoked, 11 = Re -issuance, 12 = Upload and revoke the confirmation letter.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `subject_alt_name` - The certificate contains multiple domain names (containing the main domain name).Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `submitted_data` - Submitted information information.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `admin_email` - Administrator mailbox address.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `admin_first_name` - Administrator name.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `admin_last_name` - The surname of the administrator.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `admin_phone_num` - Administrator phone number.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `admin_position` - Administrator position.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `certificate_domain` - Domain information.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `contact_email` - Contact mailbox address,Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `contact_first_name` - Contact name.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `contact_last_name` - Contact surname.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `contact_number` - Contact phone number.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `contact_position` - Contact position.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `csr_content` - CSR content.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `csr_type` - CSR type, (online = online CSR, PARSE = paste CSR).Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `domain_list` - DNS information.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `key_password` - Private key password.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `organization_address` - address.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `organization_city` - city.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `organization_country` - nation.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `organization_division` - department.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `organization_name` - Enterprise or unit name.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `organization_region` - Province.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `phone_area_code` - Local region code.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `phone_number` - Landline number.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `postal_code` - Postal code.Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `verify_type` - Verification type.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `validity_period` - Validity period: unit (month).Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `verify_type` - Verification type: DNS_AUTO = Automatic DNS verification, DNS = manual DNS verification, file = file verification, email = email verification.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `vulnerability_report` - Vulnerability scanning evaluation report.Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `vulnerability_status` - Vulnerability scanning status.Note: This field may return NULL, indicating that the valid value cannot be obtained.


