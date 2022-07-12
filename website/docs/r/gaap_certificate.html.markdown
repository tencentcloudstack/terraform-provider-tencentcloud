---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_certificate"
sidebar_current: "docs-tencentcloud-resource-gaap_certificate"
description: |-
  Provides a resource to create a certificate of GAAP.
---

# tencentcloud_gaap_certificate

Provides a resource to create a certificate of GAAP.

## Example Usage

```hcl
resource "tencentcloud_gaap_certificate" "foo" {
  type    = "BASIC"
  content = "test:tx2KGdo3zJg/."
  name    = "test_certificate"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String, ForceNew) Content of the certificate, and URL encoding. When the certificate is basic authentication, use the `user:xxx password:xxx` format, where the password is encrypted with `htpasswd` or `openssl`; When the certificate is `CA` or `SSL`, the format is `pem`.
* `type` - (Required, String, ForceNew) Type of the certificate. Valid value: `BASIC`, `CLIENT`, `SERVER`, `REALSERVER` and `PROXY`. `BASIC` means basic certificate; `CLIENT` means client CA certificate; `SERVER` means server SSL certificate; `REALSERVER` means realserver CA certificate; `PROXY` means proxy SSL certificate.
* `key` - (Optional, String, ForceNew) Key of the `SSL` certificate.
* `name` - (Optional, String) Name of the certificate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `begin_time` - Beginning time of the certificate.
* `create_time` - Creation time of the certificate.
* `end_time` - Ending time of the certificate.
* `issuer_cn` - Issuer name of the certificate.
* `subject_cn` - Subject name of the certificate.


## Import

GAAP certificate can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_certificate.foo cert-d5y6ei3b
```

