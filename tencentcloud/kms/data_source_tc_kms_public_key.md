Use this data source to query detailed information of kms public_key

Example Usage

```hcl
data "tencentcloud_kms_public_key" "example" {
  key_id = tencentcloud_kms_key.example.id
}

resource "tencentcloud_kms_key" "example" {
  alias                         = "tf-example-kms-key"
  description                   = "example of kms key"
  key_usage                     = "ASYMMETRIC_DECRYPT_RSA_2048"
  is_enabled                    = true
  pending_delete_window_in_days = 7
}
```