---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_white_box_key"
sidebar_current: "docs-tencentcloud-resource-kms_white_box_key"
description: |-
  Provides a resource to create a kms white_box_key
---

# tencentcloud_kms_white_box_key

Provides a resource to create a kms white_box_key

## Example Usage

```hcl
resource "tencentcloud_kms_white_box_key" "example" {
  alias       = "tf_example"
  description = "test desc."
  algorithm   = "SM4"
  status      = "Enabled"
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `algorithm` - (Required, String) All algorithm types for creating keys, supported values: AES_256, SM4.
* `alias` - (Required, String) As an alias for the key to be easier to identify and easier to understand, it cannot be empty and is a combination of 1-60 alphanumeric characters - _. The first character must be a letter or number. Alias are not repeatable.
* `description` - (Optional, String) Description of the key, up to 1024 bytes.
* `status` - (Optional, String) Whether to enable the key. Enabled or Disabled. Default is Enabled.
* `tags` - (Optional, Map) The tags of Key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateWhiteBoxKey](https://cloud.tencent.com/document/api/573/43443)
* Create: [DisableWhiteBoxKey](https://cloud.tencent.com/document/api/573/43437)
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322)
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319)
* Read: [DescribeWhiteBoxKey](https://cloud.tencent.com/document/api/573/43440)
* Update: [DisableWhiteBoxKey](https://cloud.tencent.com/document/api/573/43437)
* Update: [EnableWhiteBoxKey](https://cloud.tencent.com/document/api/573/43435)
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322)
* Delete: [DeleteWhiteBoxKey](https://cloud.tencent.com/document/api/573/43442)

## Import

kms white_box_key can be imported using the id, e.g.

```
terraform import tencentcloud_kms_white_box_key.example 244dab8c-6dad-11ea-80c6-5254006d0810
```

