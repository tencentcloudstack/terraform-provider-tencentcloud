---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_secrets"
sidebar_current: "docs-tencentcloud-datasource-ssm_secrets"
description: |-
  Use this data source to query detailed information of SSM secret
---

# tencentcloud_ssm_secrets

Use this data source to query detailed information of SSM secret

## Example Usage

```hcl
data "tencentcloud_ssm_secrets" "example" {
  secret_name = tencentcloud_ssm_secret.example.secret_name
  state       = 1
}

resource "tencentcloud_ssm_secret" "example" {
  secret_name = "tf_example"
  description = "desc."

  tags = {
    createdBy = "terraform"
  }
}
```

### OR you can filter by tags

```hcl
data "tencentcloud_ssm_secrets" "example" {
  secret_name = tencentcloud_ssm_secret.example.secret_name
  state       = 1

  tags = {
    createdBy = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `order_type` - (Optional, Int) The order to sort the create time of secret. `0` - desc, `1` - asc. Default value is `0`.
* `product_name` - (Optional, String) This parameter only takes effect when the SecretType parameter value is 1. When the SecretType value is 1, if the Product Name value is empty, it means to query all types of cloud product credentials. If the Product Name value is MySQL, it means to query MySQL database credentials. If the Product Name value is Tdsql mysql, it means to query Tdsql (MySQL version) credentials.
* `result_output_file` - (Optional, String) Used to save results.
* `secret_name` - (Optional, String) Secret name used to filter result.
* `secret_type` - (Optional, Int) 0- represents user-defined credentials, defaults to 0. 1- represents the user's cloud product credentials. 2- represents SSH key pair credentials. 3- represents cloud API key pair credentials.
* `state` - (Optional, Int) Filter by state of secret. `0` - all secrets are queried, `1` - only Enabled secrets are queried, `2` - only Disabled secrets are queried, `3` - only PendingDelete secrets are queried.
* `tags` - (Optional, Map) Tags to filter secret.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `secret_list` - A list of SSM secrets.
  * `associated_instance_ids` - When the credential type is SSH key pair credential, this field is valid and is used to represent the CVM instance ID associated with the SSH key pair.
  * `create_time` - Create time of secret.
  * `create_uin` - Uin of Creator.
  * `delete_time` - Delete time of CMK.
  * `description` - Description of secret.
  * `kms_key_id` - KMS keyId used to encrypt secret.
  * `kms_key_type` - KMS CMK type used to encrypt credentials, DEFAULT represents the default key created by SecretsManager, and CUSTOMER represents the user specified key.
  * `next_rotation_time` - Next rotation start time, uinx timestamp.
  * `product_name` - Cloud product name, only effective when SecretType is 1, which means the credential type is cloud product credential.
  * `project_id` - When the credential type is SSH key pair credential, this field is valid and represents the item ID to which the SSH key pair belongs.
  * `resource_id` - The cloud product instance ID number corresponding to the cloud product credentials.
  * `resource_name` - When the credential type is SSH key pair credential, this field is valid and is used to represent the name of the SSH key pair credential.
  * `rotation_begin_time` - The user specified rotation start time.
  * `rotation_frequency` - The frequency of rotation, in days, takes effect when rotation is on.
  * `rotation_status` - 1: - Turn on the rotation; 0- No rotation Note: This field may return null, indicating that a valid value cannot be obtained.
  * `secret_name` - Name of secret.
  * `secret_type` - 0- User defined credentials; 1- Cloud product credentials; 2- SSH key pair credentials; 3- Cloud API key pair credentials.
  * `status` - Status of secret.
  * `target_uin` - When the credential type is a cloud API key pair credential, this field is valid and is used to represent the user UIN to which the cloud API key pair belongs.


