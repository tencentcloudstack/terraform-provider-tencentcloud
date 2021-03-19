---
subcategory: "SSM"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_secret"
sidebar_current: "docs-tencentcloud-resource-ssm_secret"
description: |-
  Provide a resource to create a SSM secret.
---

# tencentcloud_ssm_secret

Provide a resource to create a SSM secret.

## Example Usage

```hcl
resource "tencentcloud_ssm_secret" "foo" {
  secret_name             = "test"
  description             = "test secret"
  recovery_window_in_days = 0
  is_enabled              = true

  init_secret {
    version_id    = "v1"
    secret_string = "123456"
  }

  tags = {
    test-tag = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `init_secret` - (Required) The secret of initial version.
* `secret_name` - (Required, ForceNew) Name of secret which cannot be repeated in the same region. The maximum length is 128 bytes. The name can only contain English letters, numbers, underscore and hyphen '-'. The first character must be a letter or number.
* `description` - (Optional) Description of secret. The maximum is 2048 bytes.
* `is_enabled` - (Optional) Specify whether to enable secret. Default value is `true`.
* `kms_key_id` - (Optional, ForceNew) KMS keyId used to encrypt secret. If it is empty, it means that the CMK created by SSM for you by default is used for encryption. You can also specify the KMS CMK created by yourself in the same region for encryption.
* `recovery_window_in_days` - (Optional) Specify the scheduled deletion date. Default value is `0` that means to delete immediately. 1-30 means the number of days reserved, completely deleted after this date.
* `tags` - (Optional) Tags of secret.

The `init_secret` object supports the following:

* `version_id` - (Required) Version of secret. The maximum length is 64 bytes. The version_id can only contain English letters, numbers, underscore and hyphen '-'. The first character must be a letter or number.
* `secret_binary` - (Optional) The base64-encoded binary secret. secret_binary and secret_string must be set only one, and the maximum support is 4096 bytes. When secret status is `Disabled`, this field will not update anymore.
* `secret_string` - (Optional) The string text of secret. secret_binary and secret_string must be set only one, and the maximum support is 4096 bytes. When secret status is `Disabled`, this field will not update anymore.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Status of secret.


## Import

SSM secret can be imported using the secretName, e.g.
```
$ terraform import tencentcloud_ssm_secret.foo test
```

