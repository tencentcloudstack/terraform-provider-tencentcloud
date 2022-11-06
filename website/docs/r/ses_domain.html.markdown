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
}
```

## Argument Reference

The following arguments are supported:

* `email_identity` - (Required, String, ForceNew) Your sender domain. You are advised to use a third-level domain, for example, mail.qcloud.com.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ses domain can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_domain.domain iac.cloud
```

