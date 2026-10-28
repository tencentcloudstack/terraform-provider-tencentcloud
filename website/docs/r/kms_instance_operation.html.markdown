---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_instance_operation"
sidebar_current: "docs-tencentcloud-resource-kms_instance_operation"
description: |-
  Provides a resource to create a kms kms_instance operation
---

# tencentcloud_kms_instance_operation

Provides a resource to create a kms kms_instance operation

## Example Usage

```hcl
resource "tencentcloud_kms_instance_operation" "example" {
  instance_charge_type = "PREPAID"
  product_version      = "professional"
  client_token         = "tf-token-xxx"

  instance_charge_prepaid {
    period     = 1
    renew_flag = 1
  }

  tags {
    tag_key   = "env"
    tag_value = "prod"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_charge_type` - (Required, String, ForceNew) Instance charge type. Valid value: `PREPAID`.
* `client_token` - (Optional, String, ForceNew) Reentrant client id for idempotency.
* `instance_charge_prepaid` - (Optional, List, ForceNew) Prepaid instance charge info.
* `product_version` - (Optional, String, ForceNew) Product version. Valid value: `professional` (professional edition).
* `tags` - (Optional, List, ForceNew) Instance tags.

The `instance_charge_prepaid` object supports the following:

* `period` - (Required, Int) Purchase duration in months. Valid range: `1` to `100`.
* `renew_flag` - (Required, Int) Auto-renewal flag. `0`: not set, `1`: auto-renew, `2`: no auto-renew.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `bill_id` - Billing frozen serial no.
* `deal_name` - Billing order no.
* `instance_id` - Instance id.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateKmsInstance](https://www.tencentcloud.com/search/?wd=CreateKmsInstance&type=document&lang=en)

