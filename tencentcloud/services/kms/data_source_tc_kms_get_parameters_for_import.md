Use this data source to query detailed information of kms get_parameters_for_import

Example Usage

```hcl
data "tencentcloud_kms_get_parameters_for_import" "example" {
  key_id             = "786aea8c-4aec-11ee-b601-525400281a45"
  wrapping_algorithm = "RSAES_OAEP_SHA_1"
  wrapping_key_spec  = "RSA_2048"
}
```