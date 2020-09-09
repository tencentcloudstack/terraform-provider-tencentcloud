---
subcategory: "CVM"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_key_pair"
sidebar_current: "docs-tencentcloud-resource-key_pair"
description: |-
  Provides a key pair resource.
---

# tencentcloud_key_pair

Provides a key pair resource.

## Example Usage

```hcl
resource "tencentcloud_key_pair" "foo" {
  key_name   = "terraform_test"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}
```

## Argument Reference

The following arguments are supported:

* `key_name` - (Required) The key pair's name. It is the only in one TencentCloud account.
* `public_key` - (Required, ForceNew) You can import an existing public key and using TencentCloud key pair to manage it.
* `project_id` - (Optional, ForceNew) Specifys to which project the key pair belongs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Key pair can be imported using the id, e.g.

```
$ terraform import tencentcloud_key_pair.foo skey-17634f05
```

