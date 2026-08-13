---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_secret_version_v2"
sidebar_current: "docs-tencentcloud-resource-ssm_secret_version_v2"
description: |-
  Provides a resource to create a SSM secret version using the terraform-plugin-framework implementation.
---

# tencentcloud_ssm_secret_version_v2

Provides a resource to create a SSM secret version using the terraform-plugin-framework implementation.

This is the framework-based counterpart of `tencentcloud_ssm_secret_version`. It additionally supports a write-only attribute `secret_string_wo` whose value is never persisted in Terraform plan or state.

-> **Note:** Exactly one of `secret_binary`, `secret_string` or `secret_string_wo` must be set. A maximum of 10 versions can exist under one credential. Only credentials in the enabled or disabled state can have new versions added.

## Example Usage

### Text-type credential, plaintext

```hcl
resource "tencentcloud_ssm_secret" "example" {
  secret_name             = "tf-example"
  description             = "desc."
  recovery_window_in_days = 0
  is_enabled              = true
}

resource "tencentcloud_ssm_secret_version_v2" "v1" {
  secret_name   = tencentcloud_ssm_secret.example.secret_name
  version_id    = "v1"
  secret_string = "this is secret string"
}
```

### Write-only credential (value not persisted in state)

```hcl
resource "tencentcloud_ssm_secret_version_v2" "wo" {
  secret_name      = tencentcloud_ssm_secret.example.secret_name
  version_id       = "wo"
  secret_string_wo = "this value is never persisted in plan or state"
}
```

### Binary credential, encoded using base64

```hcl
resource "tencentcloud_ssm_secret_version_v2" "v2" {
  secret_name   = tencentcloud_ssm_secret.example.secret_name
  version_id    = "v2"
  secret_binary = "MTIzMTIzMTIzMTIzMTIzQQ=="
}
```

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required, String) Specifies the name of the credential to which the new version is to be added.
* `version_id` - (Required, String) Specifies the version ID for the newly added version. It can be up to 64 bytes in length and must consist of a combination of letters, numbers, and the characters `-`, `_`, or `.`, starting with a letter or a number.
* `secret_binary` - (Optional, String) Binary credential information, encoded using Base64. Conflicts with secret_string and secret_string_wo.
* `secret_string_wo` - (Optional, String) Write-only text-based credential information. Its value is sent to the API but never persisted in Terraform plan or state. Conflicts with secret_binary and secret_string.
* `secret_string` - (Optional, String) Text-based credential information in plaintext (Base64 encoding is not required). Conflicts with secret_binary and secret_string_wo.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource identifier in the form `<secret_name>#<version_id>`.


## Import

SSM secret version (v2) can be imported using `secretName#versionId`, e.g.
```
$ terraform import tencentcloud_ssm_secret_version_v2.v1 test#v1
```

