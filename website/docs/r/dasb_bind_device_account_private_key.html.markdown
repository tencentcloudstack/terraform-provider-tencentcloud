---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_bind_device_account_private_key"
sidebar_current: "docs-tencentcloud-resource-dasb_bind_device_account_private_key"
description: |-
  Provides a resource to create a dasb bind_device_account_private_key
---

# tencentcloud_dasb_bind_device_account_private_key

Provides a resource to create a dasb bind_device_account_private_key

## Example Usage

```hcl
resource "tencentcloud_dasb_bind_device_account_private_key" "example" {
  device_account_id    = 16
  private_key          = "MIICXAIBAAKBgQCqGKukO1De7zhZj6+H0qtjTkVxwTCpvKe4eCZ0FPqri0cb2JZfXJ/DgYSF6vUpwmJG8wVQZKjeGcjDOL5UlsuusFncCzWBQ7RKNUSesmQRMSGkVb1/3j+skZ6UtW+5u09lHNsj6tQ51s1SPrCBkedbNf0Tp0GbMJDyR4e9T04ZZwIDAQABAoGAFijko56+qGyN8M0RVyaRAXz++xTqHBLh"
  private_key_password = "TerraformPassword"
}
```

## Argument Reference

The following arguments are supported:

* `device_account_id` - (Required, Int, ForceNew) Host account ID.
* `private_key` - (Required, String, ForceNew) Host account private key, the latest length is 128 bytes, the maximum length is 8192 bytes.
* `private_key_password` - (Optional, String, ForceNew) Host account private key password, maximum length 256 bytes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



