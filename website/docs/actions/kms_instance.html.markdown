---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_instance"
sidebar_current: "docs-tencentcloud-action-kms_instance"
description: |-
  Provides an action to purchase and activate a prepaid KMS (Key Management Service) instance via the CreateKmsInstance API. This is a one-time operation; no cloud-side state is persisted after the action completes.
---

# tencentcloud_kms_instance

Provides an action to purchase and activate a prepaid KMS (Key Management Service) instance via the CreateKmsInstance API. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

## Example Usage

```hcl
action "tencentcloud_kms_instance" "example" {
  config {
    instance_charge_type = "PREPAID"
    product_version      = "professional"

    instance_charge_prepaid {
      period     = 1
      renew_flag = 1
    }

    tags {
      tag_key   = "env"
      tag_value = "prod"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_charge_type` - (Required, String) Instance charge type. Valid value: `PREPAID`.
* `instance_charge_prepaid` - (Optional, Object) Prepaid instance charge info.
* `product_version` - (Optional, String) Product version. Valid value: `professional` (professional edition).
* `tags` - (Optional, List of Object) Instance tags.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## API Reference

This action invokes the following TencentCloud APIs:

* Invoke: [CreateKmsInstance](https://www.tencentcloud.com/search/?wd=CreateKmsInstance&type=document&lang=en)

