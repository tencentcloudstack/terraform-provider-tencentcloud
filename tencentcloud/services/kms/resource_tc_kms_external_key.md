Provide a resource to create a KMS external key.

Example Usage

Create a basic instance.

```hcl
resource "tencentcloud_kms_external_key" "example" {
  alias       = "tf-example-kms-externalkey"
  description = "example of kms external key"

  tags = {
    createdBy = "terraform"
  }
}
```

Create kms instance with HSM

```hcl
resource "tencentcloud_kms_external_key" "example" {
  alias          = "tf-example-kms-externalkey"
  description    = "example of kms external key"
  hsm_cluster_id = "cls-hsm-mwpd9cjm"
  
  tags = {
    createdBy = "terraform"
  }
}
```

Specify the encryption algorithm and public key.

```hcl
resource "tencentcloud_kms_external_key" "example" {
  alias               = "tf-example-kms-externalkey"
  description         = "example of kms external key"
  wrapping_algorithm  = "RSAES_PKCS1_V1_5"
  key_material_base64 = "your_public_key_base64_encoded"
  is_enabled          = true

  tags = {
    createdBy = "terraform"
  }
}
```

Disable the external kms key.

```hcl
resource "tencentcloud_kms_external_key" "example" {
  alias               = "tf-example-kms-externalkey"
  description         = "example of kms external key"
  wrapping_algorithm  = "RSAES_PKCS1_V1_5"
  key_material_base64 = "your_public_key_base64_encoded"
  is_enabled          = false

  tags = {
    createdBy = "terraform"
  }
}
```

Import

KMS external keys can be imported using the id, e.g.

```
$ terraform import tencentcloud_kms_external_key.example 25068921-2101-11f0-bf1f-5254000328e1
```
