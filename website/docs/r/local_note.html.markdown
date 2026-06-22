---
subcategory: "Provider Meta"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_local_note"
sidebar_current: "docs-tencentcloud-resource-local_note"
description: |-
  Provides a local-only resource that stores a free-form note inside the
Terraform state. This reference resource does not call any cloud API; it
exists to demonstrate the framework resource lifecycle (Create / Read /
Update / Delete) and serves as a template for new framework resources.
---

# tencentcloud_local_note

Provides a local-only resource that stores a free-form note inside the
Terraform state. This reference resource does not call any cloud API; it
exists to demonstrate the framework resource lifecycle (Create / Read /
Update / Delete) and serves as a template for new framework resources.

## Example Usage

```hcl
resource "tencentcloud_local_note" "example" {
  content = "hello, framework"
}
```

## Argument Reference

The following arguments are supported:

* `title` - (Required, String) Human-readable title of the note.
* `content` - (Optional, String) Free-form content of the note. Defaults to empty string when not specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Auto-generated, immutable identifier of the note.
* `last_updated` - RFC3339 timestamp of the last successful Create/Update.


## Import

A local note can be imported using its id (the SHA-256 of `content`):

```
terraform import tencentcloud_local_note.example 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
```

