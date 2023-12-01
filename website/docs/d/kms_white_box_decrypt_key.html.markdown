---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_white_box_decrypt_key"
sidebar_current: "docs-tencentcloud-datasource-kms_white_box_decrypt_key"
description: |-
  Use this data source to query detailed information of kms white_box_decrypt_key
---

# tencentcloud_kms_white_box_decrypt_key

Use this data source to query detailed information of kms white_box_decrypt_key

## Example Usage

```hcl
data "tencentcloud_kms_white_box_decrypt_key" "example" {
  key_id = "244dab8c-6dad-11ea-80c6-5254006d0810"
}
```

## Argument Reference

The following arguments are supported:

* `key_id` - (Required, String) Globally unique identifier for the white box key.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `decrypt_key` - White box decryption key, base64 encoded.


