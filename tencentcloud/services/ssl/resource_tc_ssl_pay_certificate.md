Provide a resource to create a payment SSL.

~> **NOTE:** Provides the creation of a paid certificate, including the submission of certificate information and order functions;
currently, it does not support re-issuing certificates, revoking certificates, and deleting certificates; the certificate remarks
and belonging items can be updated. The Destroy operation will only cancel the certificate order, and will not delete the
certificate and refund the fee. If you need a refund, you need to check the current certificate status in the console
as `Review Cancel`, and then you can click `Request a refund` to refund the fee. To update the information of a certificate,
we will automatically roll back your certificate if this certificate is already in the validation stage. This process may take
some time because the CA callback is time-consuming. Please be patient and follow the prompt message. Or, feel free to contact
Tencent Cloud Support.

Example Usage

```hcl
resource "tencentcloud_ssl_pay_certificate" "example" {
    product_id = 33
    domain_num = 1
    alias      = "ssl desc."
    project_id = 0
    information {
        csr_type              = "online"
        certificate_domain    = "www.example.com"
        organization_name     = "Tencent"
        organization_division = "Qcloud"
        organization_address  = "广东省深圳市南山区腾讯大厦1000号"
        organization_country  = "CN"
        organization_city     = "深圳市"
        organization_region   = "广东省"
        postal_code           = "0755"
        phone_area_code       = "0755"
        phone_number          = "86013388"
        verify_type           = "DNS"
        admin_first_name      = "test"
        admin_last_name       = "test"
        admin_phone_num       = "12345678901"
        admin_email           = "test@tencent.com"
        admin_position        = "developer"
        contact_first_name    = "test"
        contact_last_name     = "test"
        contact_email         = "test@tencent.com"
        contact_number        = "12345678901"
        contact_position      = "developer"
    }
}

```

Import

payment SSL instance can be imported, e.g.

```
$ terraform import tencentcloud_ssl_pay_certificate.ssl iPQNn61x#33#1#1
```