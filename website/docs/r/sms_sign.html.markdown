---
subcategory: "Short Message Service(SMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sms_sign"
sidebar_current: "docs-tencentcloud-resource-sms_sign"
description: |-
  Provides a resource to create a sms sign
---

# tencentcloud_sms_sign

Provides a resource to create a sms sign

## Example Usage

### Create a sms sign instance

```hcl
resource "tencentcloud_sms_sign" "example" {
  sign_name     = "tf_example_sms_sign"
  sign_type     = 1 # 1：APP,  DocumentType can be chosen（0，1，2，3，4）
  document_type = 4 # Screenshot of application background management (personally developed APP)
  international = 0 # Mainland China SMS
  sign_purpose  = 0 # personal use
  proof_image   = "your_proof_image"
}
```

## Argument Reference

The following arguments are supported:

* `document_type` - (Required, Int) DocumentType is used for enterprise authentication, or website, app authentication, etc. DocumentType: 0, 1, 2, 3, 4, 5, 6, 7, 8.
* `international` - (Required, Int) Whether it is Global SMS: 0: Mainland China SMS; 1: Global SMS.
* `proof_image` - (Required, String) You should Base64-encode the image of the identity certificate corresponding to the signature first, remove the prefix data:image/jpeg;base64, from the resulted string, and then use it as the value of this parameter.
* `sign_name` - (Required, String) Sms sign name, unique.
* `sign_purpose` - (Required, Int) Signature purpose: 0: for personal use; 1: for others.
* `sign_type` - (Required, Int) Sms sign type: 0, 1, 2, 3, 4, 5, 6.
* `commission_image` - (Optional, String) Power of attorney, which should be submitted if SignPurpose is for use by others. You should Base64-encode the image first, remove the prefix data:image/jpeg;base64, from the resulted string, and then use it as the value of this parameter. Note: this field will take effect only when SignPurpose is 1 (for user by others).
* `remark` - (Optional, String) Signature application remarks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



