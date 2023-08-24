---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_play_domain_cert_attachment"
sidebar_current: "docs-tencentcloud-resource-css_play_domain_cert_attachment"
description: |-
  Provides a resource to create a css play_domain_cert_attachment. This resource is used for binding the play domain and specified certification together.
---

# tencentcloud_css_play_domain_cert_attachment

Provides a resource to create a css play_domain_cert_attachment. This resource is used for binding the play domain and specified certification together.

## Example Usage

```hcl
data "tencentcloud_ssl_certificates" "foo" {
  name = "your_ssl_cert"
}

resource "tencentcloud_css_play_domain_cert_attachment" "play_domain_cert_attachment" {
  cloud_cert_id = data.tencentcloud_ssl_certificates.foo.certificates.0.id
  domain_info {
    domain_name = "your_domain_name"
    status      = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_info` - (Required, List, ForceNew) The playback domains to bind and whether to enable HTTPS for them. If `CloudCertId` is unspecified, and a domain is already bound with a certificate, this API will only update the HTTPS configuration of the domain.
* `cloud_cert_id` - (Optional, String, ForceNew) Tencent cloud ssl certificate Id. Refer to `tencentcloud_ssl_certificate` to create or obtain the resource ID.

The `domain_info` object supports the following:

* `domain_name` - (Required, String) domain name.
* `status` - (Required, Int) Whether to enable the https rule for the domain name. 1: enable, 0: disabled, -1: remain unchanged.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cert_expire_time` - certificate expiration time.
* `cert_id` - certificate ID.
* `cert_type` - certificate type. 0: Self-owned certificate, 1: Tencent Cloud ssl managed certificate.
* `certificate_alias` - certificate remarks. Synonymous with CertName.
* `update_time` - The time when the rule was last updated.


## Import

css play_domain_cert_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_css_play_domain_cert_attachment.play_domain_cert_attachment domainName#cloudCertId
```

