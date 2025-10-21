---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_black_list_delete"
sidebar_current: "docs-tencentcloud-resource-ses_black_list_delete"
description: |-
  Provides a resource to create a ses black_list
---

# tencentcloud_ses_black_list_delete

Provides a resource to create a ses black_list

~> **NOTE:** Used to remove email addresses from blacklists.

## Example Usage

```hcl
resource "tencentcloud_ses_black_list_delete" "black_list" {
  email_address = "terraform-tf@gmail.com"
}
```

## Argument Reference

The following arguments are supported:

* `email_address` - (Required, String, ForceNew) Email addresses to be unblocklisted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



