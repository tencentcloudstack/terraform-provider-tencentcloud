---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_get_parameters_for_import"
sidebar_current: "docs-tencentcloud-datasource-kms_get_parameters_for_import"
description: |-
  Use this data source to query detailed information of kms get_parameters_for_import
---

# tencentcloud_kms_get_parameters_for_import

Use this data source to query detailed information of kms get_parameters_for_import

## Example Usage

```hcl
data "tencentcloud_kms_get_parameters_for_import" "example" {
  key_id             = "786aea8c-4aec-11ee-b601-525400281a45"
  wrapping_algorithm = "RSAES_OAEP_SHA_1"
  wrapping_key_spec  = "RSA_2048"
}
```

## Argument Reference

The following arguments are supported:

* `key_id` - (Required, String) CMK unique identifier.
* `wrapping_algorithm` - (Required, String) Specifies the algorithm for encrypting key material, currently supports RSAES_PKCS1_V1_5, RSAES_OAEP_SHA_1, RSAES_OAEP_SHA_256.
* `wrapping_key_spec` - (Required, String) Specifies the type of encryption key material, currently only supports RSA_2048.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `import_token` - The token required for importing key material is used as the parameter of ImportKeyMaterial.
* `parameters_valid_to` - The validity period of the exported token and public key cannot be imported after this period, and you need to call GetParametersForImport again to obtain it.
* `public_key` - Base64-encoded public key content.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [GetParametersForImport](https://cloud.tencent.com/document/api/573/38497) - This API is used to obtain the parameters of the material to be imported into a CMK. The returned `Token` is used as one of the parameters to execute `ImportKeyMaterial`, and the returned `PublicKey` is used to encrypt the key material. The `Token` and `PublicKey` are valid for 24 hours. If they are expired, you will need to call the API again to get a new `Token` and `PublicKey`.

