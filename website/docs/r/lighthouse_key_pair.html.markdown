---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_key_pair"
sidebar_current: "docs-tencentcloud-resource-lighthouse_key_pair"
description: |-
  Provides a resource to create a lighthouse key_pair
---

# tencentcloud_lighthouse_key_pair

Provides a resource to create a lighthouse key_pair

## Example Usage

```hcl
resource "tencentcloud_lighthouse_key_pair" "key_pair" {
  key_name = "key_name_test"
}
```

## Argument Reference

The following arguments are supported:

* `key_name` - (Required, String, ForceNew) Key pair name, which can contain up to 25 digits, letters, and underscores.
* `public_key` - (Optional, String, ForceNew) Public key content of the key pair, which is in the OpenSSH RSA format.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Creation time. Expressed according to the ISO8601 standard, and using UTC time. Format: YYYY-MM-DDThh:mm:ssZ.
* `private_key` - Key to private key.


## Import

lighthouse key_pair can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_key_pair.key_pair key_pair_id
```

