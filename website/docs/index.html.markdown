---
layout: "tencentcloud"
page_title: "Provider: tencentcloud"
sidebar_current: "docs-tencentcloud-index"
description: |-
  The TencentCloud provider is used to interact with many resources supported by TencentCloud. The provider needs to be configured with the proper credentials before it can be used.
---

# TencentCloud Provider

The TencentCloud provider is used to interact with many resources supported by [TencentCloud](https://intl.cloud.tencent.com).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation on the left to read about the available resources.

-> **Note:** Terraform 0.12.x supported began with provider version 1.9.0 (June 18, 2019).

## Example Usage

```hcl
terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
    }
  }
}

# Configure the TencentCloud Provider
provider "tencentcloud" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
  region     = "ap-guangzhou"
}

# Get availability zones
data "tencentcloud_availability_zones" "default" {}

# Get availability images
data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

# Get availability instance types
data "tencentcloud_instance_types" "default" {
  cpu_core_count = 1
  memory_size    = 1
}

# Create a web server
resource "tencentcloud_instance" "web" {
  instance_name              = "web server"
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  allocate_public_ip         = true
  internet_max_bandwidth_out = 20
  security_groups            = [tencentcloud_security_group.default.id]
  count                      = 1
}

# Create security group
resource "tencentcloud_security_group" "default" {
  name        = "web accessibility"
  description = "make it accessible for both production and stage ports"
}

# Create security group rule allow web and ssh request
resource "tencentcloud_security_group_rule_set" "base" {
  security_group_id = tencentcloud_security_group.default.id

  ingress {
    action      = "ACCEPT"
    cidr_block  = "0.0.0.0/0"
    protocol    = "TCP"
    port        = "80,8080"
    description = "Create security group rule allow web request"
  }

  egress {
    action      = "ACCEPT"
    cidr_block  = "0.0.0.0/0"
    protocol    = "TCP"
    port        = "22"
    description = "Create security group rule allow ssh request"
  }
}
```

## Authentication

The TencentCloud provider offers a flexible means of providing credentials for authentication.
The following methods are supported, in this order, and explained below:

- Static credentials
- Environment variables
- Assume role
- Assume role with SAML
- Assume role with OIDC
- Shared credentials
- Enable pod OIDC
- Cam role name
- MFA certification

### Static credentials

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding an `secret_id` `secret_key` and `region` in-line in the tencentcloud provider block:

Usage:

```hcl
provider "tencentcloud" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
  region     = "ap-guangzhou"
}
```

Use `allowed_account_ids` or `forbidden_account_ids`

```hcl
provider "tencentcloud" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
  region     = "ap-guangzhou"

  allowed_account_ids   = ["100023201586", "100023201349"]
}
```

```hcl
provider "tencentcloud" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
  region     = "ap-guangzhou"

  forbidden_account_ids = ["100023201223"]
}
```

### Environment variables

You can provide your credentials via `TENCENTCLOUD_SECRET_ID` and `TENCENTCLOUD_SECRET_KEY` environment variables,
representing your TencentCloud Secret Id and Secret Key respectively. `TENCENTCLOUD_REGION` is also used, if applicable:

```hcl
provider "tencentcloud" {}
```

Usage:

```shell
$ export TENCENTCLOUD_SECRET_ID="my-secret-id"
$ export TENCENTCLOUD_SECRET_KEY="my-secret-key"
$ export TENCENTCLOUD_REGION="ap-guangzhou"
$ terraform plan
```

### Assume role

If provided with an assume role, Terraform will attempt to assume this role using the supplied credentials. Assume role can be provided by adding an `role_arn`, `session_name`, `session_duration`, `policy`(optional) and `external_id`(optional) in-line in the tencentcloud provider block:

Usage:

```hcl
provider "tencentcloud" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
  region     = "ap-guangzhou"

  assume_role {
    role_arn         = "my-role-arn"
    session_name     = "my-session-name"
    policy           = "my-role-policy"
    session_duration = 3600
  }
}
```

Combining MFA

```hcl
provider "tencentcloud" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
  region     = "ap-guangzhou"

  assume_role {
    role_arn         = "my-role-arn"
    session_name     = "my-session-name"
    policy           = "my-role-policy"
    session_duration = 3600
    serial_number    = "qcs::cam:uin/{my-uin}::mfa/softToken"
    token_code       = "523886"
  }
}
```

The `role_arn`, `session_name`, `session_duration` and `external_id` can also provided via `TENCENTCLOUD_ASSUME_ROLE_ARN`, `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME`, `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION` and `TENCENTCLOUD_ASSUME_ROLE_EXTERNAL_ID` environment variables.

The `serial_number`, `token_code` can also provided via `TENCENTCLOUD_ASSUME_ROLE_SERIAL_NUMBER`, `TENCENTCLOUD_ASSUME_ROLE_TOKEN_CODE` environment variables.

Usage:

```shell
$ export TENCENTCLOUD_SECRET_ID="my-secret-id"
$ export TENCENTCLOUD_SECRET_KEY="my-secret-key"
$ export TENCENTCLOUD_REGION="ap-guangzhou"
$ export TENCENTCLOUD_ASSUME_ROLE_ARN="my-role-arn"
$ export TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME="my-session-name"
$ export TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION=3600

$ export TENCENTCLOUD_ASSUME_ROLE_SERIAL_NUMBER="my-serial-number"
$ export TENCENTCLOUD_ASSUME_ROLE_TOKEN_CODE="my-token-code"
$ terraform plan
```

### Assume role with SAML

If provided with an assume role with SAML, Terraform will attempt to assume this role using the supplied credentials. Assume role can be provided by adding an `role_arn`, `session_name`, `session_duration`, `saml_assertion` and `principal_arn` in-line in the tencentcloud provider block:

-> **Note:** Assume-role-with-SAML is a no-AK auth type, and there is no need setting secret_id and secret_key while using it.

Usage:

```hcl
provider "tencentcloud" {
  assume_role_with_saml {
    role_arn         = "my-role-arn"
    session_name     = "my-session-name"
    session_duration = 3600
    saml_assertion   = "my-saml-assertion"
    principal_arn    = "my-principal-arn"
  }
}
```

The `role_arn`, `session_name`, `session_duration`, `saml_assertion`, `principal_arn` can also provided via `TENCENTCLOUD_ASSUME_ROLE_ARN`, `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME`, `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION`, `TENCENTCLOUD_ASSUME_ROLE_SAML_ASSERTION` and `TENCENTCLOUD_ASSUME_ROLE_PRINCIPAL_ARN` environment variables.

Usage:

```shell
$ export TENCENTCLOUD_ASSUME_ROLE_ARN="my-role-arn"
$ export TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME="my-session-name"
$ export TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION=3600
$ export TENCENTCLOUD_ASSUME_ROLE_SAML_ASSERTION="my-saml-assertion"
$ export TENCENTCLOUD_ASSUME_ROLE_PRINCIPAL_ARN="my-principal-arn"
$ terraform plan
```

### Assume role with OIDC

If provided with an assume role with OIDC, Terraform will attempt to assume this role using the supplied credentials. Assume role can be provided by adding an `role_arn`, `session_name`, `session_duration` and `web_identity_token` or `web_identity_token_file` in-line in the tencentcloud provider block:

-> **Note:** Assume-role-with-OIDC is a no-AK auth type, and there is no need setting secret_id and secret_key while using it.

-> **Note:** If both `web_identity_token` and `web_identity_token_file` are configured, `web_identity_token` will be used preferentially(overriding `web_identity_token_file`).

Content formatting guidelines of `web_identity_token_file`:

The file content must be in JSON format and must contain the key: `web_identity_token`.

```json
{
    "web_identity_token": "eyJ0eXAiOiJKV1QiLCJh......E8T0qyVA7hWM55_g"
}
```

Usage:

Use web_identity_token

```hcl
provider "tencentcloud" {
  assume_role_with_web_identity {
    provider_id        = "OIDC"
    role_arn           = "my-role-arn"
    session_name       = "my-session-name"
    session_duration   = 3600
    web_identity_token = "my-web-identity-token"
  }
}
```

Use web_identity_token_file

```hcl
provider "tencentcloud" {
  assume_role_with_web_identity {
    provider_id             = "OIDC"
    role_arn                = "my-role-arn"
    session_name            = "my-session-name"
    session_duration        = 3600
    web_identity_token_file = "/AbsolutePath/to/your/secrets/web-identity-token-file"
  }
}
```

The `provider_id`, `role_arn`, `session_name`, `session_duration`, `web_identity_token`, `web_identity_token_file` can also provided via `TENCENTCLOUD_ASSUME_ROLE_PROVIDER_ID`, `TENCENTCLOUD_ASSUME_ROLE_ARN`, `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME`, `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION`, `TENCENTCLOUD_ASSUME_ROLE_WEB_IDENTITY_TOKEN` and `TENCENTCLOUD_ASSUME_ROLE_WEB_IDENTITY_TOKEN_FILE` environment variables.

Usage:

```shell
$ export TENCENTCLOUD_SECRET_ID="my-secret-id"
$ export TENCENTCLOUD_SECRET_KEY="my-secret-key"
$ export TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION=3600
$ export TENCENTCLOUD_ASSUME_ROLE_WEB_IDENTITY_TOKEN="my-web-identity-token"
$ export TENCENTCLOUD_ASSUME_ROLE_WEB_IDENTITY_TOKEN_FILE="/AbsolutePath/to/your/secrets/web-identity-token-file"
$ export TENCENTCLOUD_ASSUME_ROLE_PROVIDER_ID="OIDC"
$ terraform plan
```

### Enable pod OIDC

Configure the TencentCloud Provider with TKE OIDC.

-> **Note:** Must ensure CAM OIDC provider and WEBHOOK component are created successfully.

Usage:

```hcl
provider "tencentcloud" {
  enable_pod_oidc = true
}
```

### Cam role name

If provided with a Cam role name, Terraform will just access the metadata URL: `http://metadata.tencentyun.com/latest/meta-data/cam/security-credentials/<cam_role_name>` to obtain the STS credential. The CVM Instance Role also can be set using the `TENCENTCLOUD_CAM_ROLE_NAME` environment variables.

-> **Note:** Cam-role-name is used to grant the role entity the permissions to access services and resources and perform operations in Tencent Cloud. You can associate the CAM role with a CVM instance to call other Tencent Cloud APIs from the instance using the periodically updated temporary Security Token Service (STS) key.

-> **Note:** Cam-role-name is a no-AK auth type, and there is no need setting secret_id and secret_key while using it.

Usage:

```hcl
provider "tencentcloud" {
  cam_role_name = "my-cam-role-name"
}
```

It can also be authenticated together with method Assume role. Authentication process: Perform CAM authentication first, then proceed with Assume role authentication.

Usage:

```hcl
provider "tencentcloud" {
  cam_role_name = "my-cam-role-name"

  assume_role {
    role_arn         = "my-role-arn"
    session_name     = "my-session-name"
    policy           = "my-role-policy"
    session_duration = 3600
    external_id      = "my-external-id"
  }
}
```

### MFA certification

If provided with MFA certification, Terraform will attempt to use the provided credentials for MFA authentication.

Usage:

```hcl
provider "tencentcloud" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
  region     = "ap-guangzhou"

  mfa_certification {
    serial_number    = "qcs::cam:uin/{my-uin}::mfa/softToken"
    token_code       = "523886"
    duration_seconds = 1800
  }
}
```

The `serial_number`, `token_code`, `duration_seconds` can also provided via `TENCENTCLOUD_MFA_CERTIFICATION_SERIAL_NUMBER`, `TENCENTCLOUD_MFA_CERTIFICATION_TOKEN_CODE`, `TENCENTCLOUD_MFA_CERTIFICATION_DURATION_SECONDS` environment variables.

Usage:

```shell
$ export TENCENTCLOUD_MFA_CERTIFICATION_SERIAL_NUMBER="my-serial-number"
$ export TENCENTCLOUD_MFA_CERTIFICATION_TOKEN_CODE="my-token-code"
$ export TENCENTCLOUD_MFA_CERTIFICATION_DURATION_SECONDS=1800
$ terraform plan
```

### CDC cos usage

You can set the cos domain by setting the environment variable `TENCENTCLOUD_COS_DOMAIN`, and configure the cdc scenario as follows:

-> **Note:** Please note that not all cos resources are supported. Please pay attention to the prompts for each resource.

```hcl
locals {
  region = "ap-guangzhou"
  cdc_id = "cluster_xxx"
}

provider "tencentcloud" {
  region     = local.region
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
  cos_domain = "https://${local.cdc_id}.cos-cdc.${local.region}.myqcloud.com/"
}
```

The `cos_domain` can also provided via `TENCENTCLOUD_COS_DOMAIN` environment variables.

Usage:

```shell
$ export TENCENTCLOUD_SECRET_ID="my-secret-id"
$ export TENCENTCLOUD_SECRET_KEY="my-secret-key"
$ export TENCENTCLOUD_REGION="ap-guangzhou"
$ export TENCENTCLOUD_COS_DOMAIN="https://cluster-xxxxxx.cos-cdc.ap-guangzhou.myqcloud.com/"
$ terraform plan
```

### Shared credentials

You can use [Tencent Cloud credentials](https://www.tencentcloud.com/document/product/1013/33464) to specify your credentials. The default location is `$HOME/.tccli` on Linux and macOS, And `"%USERPROFILE%\.tccli"` on Windows. You can optionally specify a different location in the Terraform configuration by providing the `shared_credentials_dir` argument or using the `TENCENTCLOUD_SHARED_CREDENTIALS_DIR` environment variable. This method also supports a `profile` configuration and matching `TENCENTCLOUD_PROFILE` environment variable:

Usage:

On Linux/MacOS

```hcl
provider "tencentcloud" {
  shared_credentials_dir = "/Users/tf_user/.tccli"
  profile                = "default"
}
```

On Windows

```hcl
provider "tencentcloud" {
  shared_credentials_dir = "C:\\Users\\tf_user\\.tccli"
  profile                = "default"
}
```

## Argument Reference

In addition to generic provider arguments (e.g. alias and version), the following arguments are supported in the TencentCloud provider block:

* `secret_id` - (Optional) This is the TencentCloud secret id. It must be provided, but it can also be sourced from the `TENCENTCLOUD_SECRET_ID` environment variable.
* `secret_key` - (Optional) This is the TencentCloud secret key. It must be provided, but it can also be sourced from the `TENCENTCLOUD_SECRET_KEY` environment variable.
* `security_token` - (Optional) TencentCloud security token of temporary access credentials. It can also be sourced from the `TENCENTCLOUD_SECURITY_TOKEN` environment variable. Notice: for supported products, please refer to: [temporary key supported products](https://intl.cloud.tencent.com/document/product/598/10588).
* `region` - (Optional) This is the TencentCloud region. It must be provided, but it can also be sourced from the `TENCENTCLOUD_REGION` environment variables. The default input value is `ap-guangzhou`.
* `shared_credentials_dir` - (Optional) The directory of the shared credentials. It can also be sourced from the `TENCENTCLOUD_SHARED_CREDENTIALS_DIR` environment variable. If not set this defaults to ~/.tccli.
* `profile` - (Optional) The profile name as set in the shared credentials. It can also be sourced from the `TENCENTCLOUD_PROFILE` environment variable. If not set, the default profile created with `tccli configure` will be used.
* `assume_role` - (Optional, Available in 1.33.1+) An `assume_role` block (documented below). If provided, terraform will attempt to assume this role using the supplied credentials. Only one `assume_role` block may be in the configuration.
* `assume_role_with_saml` - (Optional, Available in 1.81.111+) An `assume_role_with_saml` block (documented below). If provided, terraform will attempt to assume this role using the supplied credentials. Only one `assume_role_with_saml` block may be in the configuration.
* `enable_pod_oidc` - (Optional, Available in 1.81.117+) Whether to enable pod oidc.
* `assume_role_with_web_identity` - (Optional, Available in 1.81.111+) An `assume_role_with_web_identity` block (documented below). If provided, terraform will attempt to assume this role using the supplied credentials. Only one `assume_role_with_web_identity` block may be in the configuration.
* `protocol` - (Optional, Available in 1.37.0+) The protocol of the API request. Valid values: `HTTP` and `HTTPS`. Default is `HTTPS`.
* `domain` - (Optional, Available in 1.37.0+) The root domain of the API request, Default is `tencentcloudapi.com`. 
* `cam_role_name` - (Optional, Available in 1.81.117+) The name of the CVM instance CAM role. It can be sourced from the `TENCENTCLOUD_CAM_ROLE_NAME` environment variable. 
* `allowed_account_ids` - (Optional) List of allowed TencentCloud account IDs to prevent you from mistakenly using the wrong one (and potentially end up destroying a live environment). Conflicts with `forbidden_account_ids`, If use `assume_role_with_saml` or `assume_role_with_web_identity`, it is not supported.
* `forbidden_account_ids` - (Optional) List of forbidden TencentCloud account IDs to prevent you from mistakenly using the wrong one (and potentially end up destroying a live environment). Conflicts with `allowed_account_ids`, If use `assume_role_with_saml` or `assume_role_with_web_identity`, it is not supported.

The nested `assume_role` block supports the following:
* `role_arn` - (Required) The ARN of the role to assume. It can also be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARN` environment variable.
* `session_name` - (Required) The session name to use when making the AssumeRole call. It can also be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME` environment variable.
* `session_duration` - (Required) The duration of the session when making the AssumeRole call. Its value ranges from 0 to 43200(seconds), and default is 7200 seconds. It can also be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION` environment variable.
* `policy` - (Optional) A more restrictive policy to apply to the temporary credentials. This gives you a way to further restrict the permissions for the resulting temporary security credentials. You cannot use the passed policy to grant permissions that are in excess of those allowed by the access policy of the role that is being assumed.
* `external_id` - (Optional) External role ID, which can be obtained by clicking the role name in the CAM console. It can contain 2-128 letters, digits, and symbols (=,.@\:/-). Regex: [\\w+=,.@\:/-]*. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_EXTERNAL_ID`.

The nested `assume_role_with_saml` block supports the following:
* `role_arn` - (Required) The ARN of the role to assume. It can also be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARN` environment variable.
* `session_name` - (Required) The session name to use when making the AssumeRole call. It can also be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME` environment variable.
* `session_duration` - (Required) The duration of the session when making the AssumeRole call. Its value ranges from 0 to 43200(seconds), and default is 7200 seconds. It can also be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION` environment variable.
* `saml_assertion` - (Required) SAML assertion information encoded in base64. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SAML_ASSERTION`.
* `principal_arn` - (Required) Player Access Description Name. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_PRINCIPAL_ARN`.

The nested `assume_role_with_web_identity` block supports the following:
* `provider_id` - (Optional) Identity provider name. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_PROVIDER_ID`, Default is OIDC.
* `role_arn` - (Required) The ARN of the role to assume. It can also be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARN` environment variable.
* `session_name` - (Required) The session name to use when making the AssumeRole call. It can also be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME` environment variable.
* `session_duration` - (Required) The duration of the session when making the AssumeRole call. Its value ranges from 0 to 43200(seconds), and default is 7200 seconds. It can also be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION` environment variable.
* `web_identity_token` - (Optional) OIDC token issued by IdP. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_WEB_IDENTITY_TOKEN`. One of `web_identity_token` or `web_identity_token_file` is required.
* `web_identity_token_file` - (Optional) File containing a web identity token from an OpenID Connect (OIDC) or OAuth provider. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_WEB_IDENTITY_TOKEN_FILE`. One of `web_identity_token` or `web_identity_token_file` is required.
