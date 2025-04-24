Provide a resource to create a KMS key.

Example Usage

Create and enable a instance.

```hcl
resource "tencentcloud_kms_key" "example" {
  alias                         = "tf-example-kms-key"
  description                   = "example of kms key"
  key_rotation_enabled          = false
  is_enabled                    = true
  pending_delete_window_in_days = 7

  tags = {
    createdBy = "Terraform"
  }
}
```

Create kms instance with HSM

```hcl
resource "tencentcloud_kms_key" "example" {
  alias                         = "tf-example-kms-key"
  description                   = "example of kms key"
  key_rotation_enabled          = false
  is_enabled                    = true
  pending_delete_window_in_days = 7
  hsm_cluster_id                = "cls-hsm-mwpd9cjm"

  tags = {
    createdBy = "Terraform"
  }
}
```

Specify the Key Usage as an asymmetry method.

```hcl
resource "tencentcloud_kms_key" "example" {
  alias       = "tf-example-kms-key"
  description = "example of kms key"
  key_usage   = "ASYMMETRIC_DECRYPT_RSA_2048"
  is_enabled  = false
}
```

Disable the kms key instance.

```hcl
resource "tencentcloud_kms_key" "example" {
  alias                = "tf-example-kms-key"
  description          = "example of kms key"
  key_rotation_enabled = false
  is_enabled           = false

  tags = {
    createdBy = "Terraform"
  }
}
```

Import

KMS keys can be imported using the id, e.g.

```
$ terraform import tencentcloud_kms_key.example 287e8f40-7cbb-11eb-9a3a-5254004f7f94
```