---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_domain"
sidebar_current: "docs-tencentcloud-resource-ses_domain"
description: |-
  Provides a resource to create a ses domain
---

# tencentcloud_ses_domain

Provides a resource to create a ses domain

## Example Usage

```hcl
resource "tencentcloud_ses_domain" "domain" {
  email_identity = "iac.cloud"
  dkim_option    = 1
  tag_list {
    tag_key   = "env"
    tag_value = "prod"
  }
}
```

## Argument Reference

The following arguments are supported:

* `email_identity` - (Required, String, ForceNew) Your sender domain. You are advised to use a third-level domain, for example, mail.qcloud.com.
* `dkim_option` - (Optional, Int, ForceNew) DKIM key length. 0: 1024-bit, 1: 2048-bit.
* `tag_list` - (Optional, List, ForceNew) Tag list.

The `tag_list` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `attributes` - DNS configuration details.
  * `expected_value` - Values that need to be configured.
  * `send_domain` - Domain name.
  * `type` - Record Type CNAME | A | TXT | MX.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateEmailIdentity](https://cloud.tencent.com/document/api/1288/51048) - This API is used to create a sender domain. Before you can send an email using Tencent Cloud SES, you must create a sender domain as your identity. It can be the domain of your website or mobile app. You must verify the domain to prove that you own it and authorize Tencent Cloud SES to use it to send emails.
* Read: [GetEmailIdentity](https://cloud.tencent.com/document/api/1288/51046) - This API is used to get the configuration details of a sender domain.
* Delete: [DeleteEmailIdentity](https://cloud.tencent.com/document/api/1288/51047) - Delete the sender domain name. Once deleted, that domain name cannot be used anymore to send email.

## Import

ses domain can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_domain.domain iac.cloud
```

