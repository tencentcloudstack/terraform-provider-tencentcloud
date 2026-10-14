---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_ssh_key_pair_secret"
sidebar_current: "docs-tencentcloud-resource-ssm_ssh_key_pair_secret"
description: |-
  Provides a resource to create a ssm ssh key pair secret
---

# tencentcloud_ssm_ssh_key_pair_secret

Provides a resource to create a ssm ssh key pair secret

## Example Usage

```hcl
resource "tencentcloud_kms_key" "example" {
  alias                = "tf-example-kms-key"
  description          = "example of kms key"
  key_rotation_enabled = false
  is_enabled           = true

  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ssm_ssh_key_pair_secret" "example" {
  secret_name   = "tf-example"
  project_id    = 0
  description   = "desc."
  kms_key_id    = tencentcloud_kms_key.example.id
  ssh_key_name  = "tf_example_ssh"
  status        = "Enabled"
  clean_ssh_key = true

  tags = {
    createdBy = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, Int) ID of the project to which the created SSH key belongs.
* `secret_name` - (Required, String, ForceNew) Secret name, which must be unique in the same region. It can contain 128 bytes of letters, digits, hyphens and underscores and must begin with a letter or digit.
* `clean_ssh_key` - (Optional, Bool) Specifies whether to delete the SSH key from both the secret and the SSH key list in the CVM console. This field is only take effect when delete SSH key secrets. Valid values: `True`: deletes SSH key from both the secret and SSH key list in the CVM console. Note that the deletion will fail if the SSH key is already bound to a CVM instance.`False`: only deletes the SSH key information in the secret.
* `description` - (Optional, String) Description, such as what it is used for. It contains up to 2,048 bytes.
* `kms_key_id` - (Optional, String) Specifies a KMS CMK to encrypt the secret.If this parameter is left empty, the CMK created by Secrets Manager by default will be used for encryption.You can also specify a custom KMS CMK created in the same region for encryption.
* `ssh_key_name` - (Optional, String) Name of the SSH key pair, which only contains digits, letters and underscores and must start with a digit or letter. The maximum length is 25 characters.
* `status` - (Optional, String) Enable or Disable Secret. Valid values is `Enabled` or `Disabled`. Default is `Enabled`.
* `tags` - (Optional, Map) Tags of secret.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Credential creation time in UNIX timestamp format.
* `resource_id` - The resource ID associated with the secret.
* `secret_type` - `0`: user-defined secret. `1`: Tencent Cloud services secret. `2`: SSH key secret. `3`: Tencent Cloud API key secret. Note: this field may return `null`, indicating that no valid values can be obtained.
* `ssh_key_id` - The key pair ID is the unique identifier of the key pair in the cloud server.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSSHKeyPairSecret](https://cloud.tencent.com/document/api/1140/60331) - This API is used to create a secret that hosts SSH keys.
* Create: [DescribeSecret](https://cloud.tencent.com/document/api/1140/40526) - This API is used to obtain the detailed attribute information of a Secret.
* Create: [DisableSecret](https://cloud.tencent.com/document/api/1140/40525) - This API is used to disable a Secret and will change its status to `Disabled`. The plaintext of a disabled Secret cannot be obtained through APIs.
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Read: [DescribeSecret](https://cloud.tencent.com/document/api/1140/40526) - This API is used to obtain the detailed attribute information of a Secret.
* Read: [GetSSHKeyPairValue](https://cloud.tencent.com/document/api/1140/60330) - This API is used to obtain the plaintext value of the SSH key secret.
* Read: [ListSecrets](https://cloud.tencent.com/document/api/1140/40519) - This API is used to obtain the detailed information list of all Secrets. You can specify the filter fields and sorting order as needed.
* Update: [DescribeSecret](https://cloud.tencent.com/document/api/1140/40526) - This API is used to obtain the detailed attribute information of a Secret.
* Update: [DisableSecret](https://cloud.tencent.com/document/api/1140/40525) - This API is used to disable a Secret and will change its status to `Disabled`. The plaintext of a disabled Secret cannot be obtained through APIs.
* Update: [EnableSecret](https://cloud.tencent.com/document/api/1140/40524) - This API is used to enable a Secret and will change its status to `Enabled`. You can call the `GetSecretValue` API to obtain the plaintext of this Secret.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Update: [UpdateDescription](https://cloud.tencent.com/document/api/1140/40516) - This API is used to update the description of a Secret. This API can only update Secrets in `Enabled` or `Disabled` status.
* Delete: [DeleteSecret](https://cloud.tencent.com/document/api/1140/40528) - This API is used to delete a Secret. You can set whether to delete the Secret immediately or on schedule using the `RecoveryWindowInDays` parameter.
* Delete: [DisableSecret](https://cloud.tencent.com/document/api/1140/40525) - This API is used to disable a Secret and will change its status to `Disabled`. The plaintext of a disabled Secret cannot be obtained through APIs.

## Import

ssm ssh_key_pair_secret can be imported using the id, e.g.

```
terraform import tencentcloud_ssm_ssh_key_pair_secret.ssh_key_pair_secret ssh_key_pair_secret_name
```

