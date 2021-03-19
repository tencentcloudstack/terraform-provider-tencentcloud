---
subcategory: "SSM"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_secret_version"
sidebar_current: "docs-tencentcloud-resource-ssm_secret_version"
description: |-
  Provide a resource to create a SSM secret version.
---

# tencentcloud_ssm_secret_version

Provide a resource to create a SSM secret version.

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

resource "tencentcloud_ssm_secret_version" "v2" {
  secret_name   = tencentcloud_ssm_secret.foo.secret_name
  version_id    = "v2"
  secret_binary = "MTIzMTIzMTIzMTIzMTIzQQ=="
}
```

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required, ForceNew) Name of secret which cannot be repeated in the same region. The maximum length is 128 bytes. The name can only contain English letters, numbers, underscore and hyphen '-'. The first character must be a letter or number.
* `version_id` - (Required, ForceNew) Version of secret. The maximum length is 64 bytes. The version_id can only contain English letters, numbers, underscore and hyphen '-'. The first character must be a letter or number.
* `secret_binary` - (Optional) The base64-encoded binary secret. secret_binary and secret_string must be set only one, and the maximum support is 4096 bytes. When secret status is `Disabled`, this field will not update anymore.
* `secret_string` - (Optional) The string text of secret. secret_binary and secret_string must be set only one, and the maximum support is 4096 bytes. When secret status is `Disabled`, this field will not update anymore.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SSM secret version can be imported using the secretName#versionId, e.g.
```
$ terraform import tencentcloud_ssm_secret_version.v2 test#v2
```

