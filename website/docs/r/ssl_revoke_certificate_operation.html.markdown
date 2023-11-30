---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_revoke_certificate_operation"
sidebar_current: "docs-tencentcloud-resource-ssl_revoke_certificate_operation"
description: |-
  Provides a resource to create a ssl revoke_certificate
---

# tencentcloud_ssl_revoke_certificate_operation

Provides a resource to create a ssl revoke_certificate

## Example Usage

```hcl
resource "tencentcloud_ssl_revoke_certificate_operation" "revoke_certificate" {
  certificate_id = "7zUGkVab"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String, ForceNew) Certificate ID.
* `reason` - (Optional, String, ForceNew) Reasons for revoking certificate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ssl revoke_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_revoke_certificate_operation.revoke_certificate revoke_certificate_id
```

