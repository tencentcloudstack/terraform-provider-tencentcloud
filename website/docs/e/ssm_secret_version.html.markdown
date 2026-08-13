---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_secret_version"
sidebar_current: "docs-tencentcloud-ephemeral_resource-ssm_secret_version"
description: |-
  Provides an ephemeral resource to read a SSM secret version value.
---

# tencentcloud_ssm_secret_version

Provides an ephemeral resource to read a SSM secret version value.

-> **Note:** A maximum of 10 versions can be supported under one credential. Only new versions can be added to credentials in the enabled and disabled states.

## Example Usage

### Text type credential information plaintext

```hcl
ephemeral "tencentcloud_ssm_secret_version" "v1" {
  secret_name = "my-secret"
  version_id  = "v1"
}
```

### Binary credential information, encoded using base64

```hcl
ephemeral "tencentcloud_ssm_secret_version" "v2" {
  secret_name = "my-secret"
  version_id  = "v2"
}
```

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required, String) Specifies the name of the credential to which the version to be read belongs.
* `version_id` - (Required, String) Specifies the version ID of the secret version to read.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `secret_binary` - The binary credential information of the secret version, encoded using Base64.
* `secret_string` - The text-based credential information of the secret version.


