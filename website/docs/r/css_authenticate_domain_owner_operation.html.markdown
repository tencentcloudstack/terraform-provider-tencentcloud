---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_authenticate_domain_owner_operation"
sidebar_current: "docs-tencentcloud-resource-css_authenticate_domain_owner_operation"
description: |-
  Provides a resource to verify the domain ownership by specified way when DomainNeedVerifyOwner failed in domain creation.
---

# tencentcloud_css_authenticate_domain_owner_operation

Provides a resource to verify the domain ownership by specified way when DomainNeedVerifyOwner failed in domain creation.

## Example Usage

dnsCheck way:

```hcl
resource "tencentcloud_css_authenticate_domain_owner_operation" "dnsCheck" {
  domain_name = "your_domain_name"
  verify_type = "dnsCheck"
}
```

fileCheck way:

```hcl
resource "tencentcloud_css_authenticate_domain_owner_operation" "fileCheck" {
  domain_name = "your_domain_name"
  verify_type = "fileCheck"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String, ForceNew) The domain name to verify.
* `verify_type` - (Optional, String, ForceNew) Authentication type. Possible values:`dnsCheck`: Immediately verify whether the resolution record of the configured dns is consistent with the content to be verified, and save the record if successful.`fileCheck`: Immediately verify whether the web file is consistent with the content to be verified, and save the record if successful.`dbCheck`: Check if authentication has been successful.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



