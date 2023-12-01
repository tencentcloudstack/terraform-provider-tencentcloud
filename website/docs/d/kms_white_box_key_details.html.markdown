---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_white_box_key_details"
sidebar_current: "docs-tencentcloud-datasource-kms_white_box_key_details"
description: |-
  Use this data source to query detailed information of kms white_box_key_details
---

# tencentcloud_kms_white_box_key_details

Use this data source to query detailed information of kms white_box_key_details

## Example Usage

```hcl
data "tencentcloud_kms_white_box_key_details" "example" {
  key_status = 0
}
```

## Argument Reference

The following arguments are supported:

* `key_status` - (Optional, Int) Filter condition: status of the key, 0: disabled, 1: enabled.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `key_infos` - List of white box key information.
  * `algorithm` - The type of algorithm used by the key.
  * `alias` - As an alias for a key that is easier to identify and easier to understand, it cannot be empty and is a combination of 1-60 alphanumeric characters - _. The first character must be a letter or number. It cannot be repeated.
  * `create_time` - Key creation time, Unix timestamp.
  * `creator_uin` - Creator.
  * `decrypt_key` - White box decryption key, base64 encoded.
  * `description` - Description of the key.
  * `device_fingerprint_bind` - Is there a device fingerprint bound to the current key?.
  * `encrypt_key` - White box encryption key, base64 encoded.
  * `key_id` - Globally unique identifier for the white box key.
  * `owner_uin` - Creator.
  * `resource_id` - Resource ID, format: creatorUin/$creatorUin/$keyId.
  * `status` - The status of the white box key, the value is: Enabled | Disabled.


