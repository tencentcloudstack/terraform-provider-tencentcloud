---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_secret_version"
sidebar_current: "docs-tencentcloud-resource-ssm_secret_version"
description: |-
  Provide a resource to create a SSM secret version.
---

# tencentcloud_ssm_secret_version

Provide a resource to create a SSM secret version.

-> **Note:** A maximum of 10 versions can be supported under one credential. Only new versions can be added to credentials in the enabled and disabled states.

## Example Usage

### Text type credential information plaintext

```hcl
resource "tencentcloud_ssm_secret" "example" {
  secret_name             = "tf-example"
  description             = "desc."
  recovery_window_in_days = 0
  is_enabled              = true

  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ssm_secret_version" "v1" {
  secret_name   = tencentcloud_ssm_secret.example.secret_name
  version_id    = "v1"
  secret_string = "this is secret string"
}
```

### Binary credential information, encoded using base64

```hcl
resource "tencentcloud_ssm_secret_version" "v2" {
  secret_name   = tencentcloud_ssm_secret.example.secret_name
  version_id    = "v2"
  secret_binary = "MTIzMTIzMTIzMTIzMTIzQQ=="
}
```

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required, String, ForceNew) Specifies the name of the credential to which the new version is to be added.
* `version_id` - (Required, String, ForceNew) Specifies the version ID for the newly added version. It can be up to 64 bytes in length and must consist of a combination of letters, numbers, and the characters `-`, `_`, or `.`, starting with a letter or a number.
* `secret_binary` - (Optional, String) Binary credential information, encoded using Base64. You must set exactly one of SecretBinary or SecretString.
* `secret_string` - (Optional, String) Text-based credential information in plaintext (Base64 encoding is not required). You must set exactly one of SecretBinary or SecretString.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SSM secret version can be imported using the secretName#versionId, e.g.
```
$ terraform import tencentcloud_ssm_secret_version.v1 test#v1
```

