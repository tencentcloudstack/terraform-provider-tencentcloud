---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_upload_revoke_letter_operation"
sidebar_current: "docs-tencentcloud-resource-ssl_upload_revoke_letter_operation"
description: |-
  Provides a resource to create a ssl upload_revoke_letter
---

# tencentcloud_ssl_upload_revoke_letter_operation

Provides a resource to create a ssl upload_revoke_letter

## Example Usage

```hcl
resource "tencentcloud_ssl_upload_revoke_letter_operation" "upload_revoke_letter" {
  certificate_id = "8xRYdDlc"
  revoke_letter  = filebase64("./c.pdf")
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String, ForceNew) Certificate ID.
* `revoke_letter` - (Required, String, ForceNew) The format of the base64-encoded certificate confirmation letter file should be jpg, jpeg, png, or pdf, and the size should be between 1kb and 1.4M.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ssl upload_revoke_letter can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_upload_revoke_letter_operation.upload_revoke_letter upload_revoke_letter_id
```

