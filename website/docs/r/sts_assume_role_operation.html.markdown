---
subcategory: "Security Token Service(STS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sts_assume_role_operation"
sidebar_current: "docs-tencentcloud-resource-sts_assume_role_operation"
description: |-
  Provides a resource to perform an STS (Security Token Service) AssumeRole operation to obtain temporary security credentials.
---

# tencentcloud_sts_assume_role_operation

Provides a resource to perform an STS (Security Token Service) AssumeRole operation to obtain temporary security credentials.

## Example Usage

### Assume role with required parameters

```hcl
resource "tencentcloud_sts_assume_role_operation" "example" {
  role_arn          = "qcs::cam::uin/100000000001:roleName/testRoleName"
  role_session_name = "example-session"
}
```

### Assume role with all parameters

```hcl
resource "tencentcloud_sts_assume_role_operation" "example" {
  role_arn          = "qcs::cam::uin/100000000001:roleName/testRoleName"
  role_session_name = "example-session"
  duration_seconds  = 7200
  policy = jsonencode({
    version = "2.0"
    statement = [
      {
        effect   = "allow"
        action   = ["cos:GetObject"]
        resource = ["*"]
      }
    ]
  })
  external_id     = "external-id-example"
  source_identity = "source-identity-example"
  serial_number   = "qcs::cam:uin/100000000001::mfa/softToken"
  token_code      = "123456"

  tags {
    key   = "env"
    value = "test"
  }

  tags {
    key   = "project"
    value = "demo"
  }
}
```

## Argument Reference

The following arguments are supported:

* `role_arn` - (Required, String, ForceNew) Resource description of the role, which can be obtained by clicking the role name in [Access Management](https://console.cloud.tencent.com/cam/role).
* `role_session_name` - (Required, String, ForceNew) User-defined temporary session name. Length is between 2 and 128, can contain uppercase and lowercase characters, numbers, and special characters: =,.@_-.
* `duration_seconds` - (Optional, Int, ForceNew) Specifies the validity period of the temporary access credential in seconds. Default is 7200 seconds, maximum is 43200 seconds.
* `external_id` - (Optional, String, ForceNew) Role external ID, which can be obtained by clicking the role name in [Access Management](https://console.cloud.tencent.com/cam/role). Length is between 2 and 128.
* `policy` - (Optional, String, ForceNew) Policy description. The policy syntax refers to [CAM Policy Syntax](https://cloud.tencent.com/document/product/598/10603). The policy cannot contain the principal element.
* `serial_number` - (Optional, String, ForceNew) MFA serial number associated with the CAM user making the call. Format: qcs::cam:uin/${ownerUin}::mfa/${mfaType}. mfaType supports softToken.
* `source_identity` - (Optional, String, ForceNew) Caller identity uin.
* `tags` - (Optional, List, ForceNew) Session tag list. A maximum of 50 session tags can be passed, and duplicate tag keys are not supported.
* `token_code` - (Optional, String, ForceNew) MFA authentication code.

The `tags` object supports the following:

* `key` - (Required, String) Tag key, up to 128 characters, case-sensitive.
* `value` - (Required, String) Tag value, up to 256 characters, case-sensitive.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `credentials` - Temporary access credentials.
  * `tmp_secret_id` - Temporary certificate secret ID. Maximum length is 1024 bytes.
  * `tmp_secret_key` - Temporary certificate secret key. Maximum length is 1024 bytes.
  * `token` - Token. The token length is up to 4096 bytes.
* `expiration` - Expiration time of the temporary access credential in ISO8601 format UTC time.
* `expired_time` - Expiration time of the temporary access credential, returned as a Unix timestamp accurate to seconds.


