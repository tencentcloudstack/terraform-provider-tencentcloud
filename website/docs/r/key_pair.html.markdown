---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_key_pair"
sidebar_current: "docs-tencentcloud-resource-cvm-key-pair"
description: |-
  Provides a TencentCloud key pair resource.
---

# tencentcloud_key_pair

Provides a key pair resource.

## Example Usage

Basic Usage

```hcl
resource "tencentcloud_key_pair" "foo" {
  key_name = "from_terraform_public_key"
  public_key = "ssh-rsa AAAAB3NzaSuperLongString foo@bar"
}
```

## Argument Reference

The following arguments are supported:

* `key_name` - (Force new resource) The key pair's name. It is the only in one TencentCloud account.
* `public_key` - (Force new resource) You can import an existing public key and using TencentCloud key pair to manage it.


## Attributes Reference

* `id` - The id of the key pair, something like `skey-xxxxxxx`, use this for instance creation and resetting.
