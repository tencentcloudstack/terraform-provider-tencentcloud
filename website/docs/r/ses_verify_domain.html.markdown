---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_verify_domain"
sidebar_current: "docs-tencentcloud-resource-ses_verify_domain"
description: |-
  Provides a resource to create a ses verify_domain
---

# tencentcloud_ses_verify_domain

Provides a resource to create a ses verify_domain

~> **NOTE:** Please add the `attributes` information returned by `tencentcloud_ses_domain` to the domain name resolution record through `tencentcloud_dnspod_record`, and then verify it.

## Example Usage

```hcl
resource "tencentcloud_ses_verify_domain" "verify_domain" {
  email_identity = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `email_identity` - (Required, String, ForceNew) Domain name requested for verification.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



