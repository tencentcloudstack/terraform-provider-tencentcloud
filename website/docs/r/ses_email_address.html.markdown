---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_email_address"
sidebar_current: "docs-tencentcloud-resource-ses_email_address"
description: |-
  Provides a resource to create a ses email_address
---

# tencentcloud_ses_email_address

Provides a resource to create a ses email_address

## Example Usage

```hcl
resource "tencentcloud_ses_email_address" "email_address" {
  email_address     = "aaa@iac-tf.cloud"
  email_sender_name = "aaa"
}
```

## Argument Reference

The following arguments are supported:

* `email_address` - (Required, String, ForceNew) Your sender address. (You can create up to 10 sender addresses for each domain.).
* `email_sender_name` - (Optional, String, ForceNew) Sender name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ses email_address can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_email_address.email_address aaa@iac-tf.cloud
```

