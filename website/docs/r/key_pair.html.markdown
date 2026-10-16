---
subcategory: "Cloud Virtual Machine(CVM)"
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
  key_name = "terraform_test"
}

output "private_key" {
  value = tencentcloud_key_pair.foo.private_key
}

output "create_time" {
  value = tencentcloud_key_pair.foo.created_time
}

resource "tencentcloud_key_pair" "foo1" {
  key_name   = "terraform_test"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}
```

## Argument Reference

The following arguments are supported:

* `key_name` - (Required, String) The key pair's name. It is the only in one TencentCloud account.
* `force_stop` - (Optional, Bool) Whether to forcibly shut down a running instance. Default is false. Forcing a shutdown is equivalent to switching off the power button on a physical computer. Forcing a shutdown may result in data loss or file system corruption; therefore, please use this option only when the server cannot be shut down normally.
* `project_id` - (Optional, Int, ForceNew) Specifys to which project the key pair belongs.
* `public_key` - (Optional, String, ForceNew) You can import an existing public key and using TencentCloud key pair to manage it.
* `tags` - (Optional, Map) Tags of the key pair.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Creation time, which follows the `ISO8601` standard and uses `UTC` time in the format of `YYYY-MM-DDThh:mm:ssZ`.
* `private_key` - Content of private key in a key pair. Tencent Cloud do not keep private keys. Please keep it properly.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateKeyPair](https://cloud.tencent.com/document/api/213/15702)
* Create: [ImportKeyPair](https://cloud.tencent.com/document/api/213/15703)
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322)
* Read: [DescribeKeyPairs](https://cloud.tencent.com/document/api/213/15699)
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319)
* Update: [ModifyKeyPairAttribute](https://cloud.tencent.com/document/api/213/15701)
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322)
* Delete: [DeleteKeyPairs](https://cloud.tencent.com/document/api/213/15700)
* Delete: [DescribeKeyPairs](https://cloud.tencent.com/document/api/213/15699)
* Delete: [DisassociateInstancesKeyPairs](https://cloud.tencent.com/document/api/213/15697)

## Import

Key pair can be imported using the id, e.g.

```
$ terraform import tencentcloud_key_pair.foo skey-17634f05
```

