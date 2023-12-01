---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip_public_address_adjust"
sidebar_current: "docs-tencentcloud-resource-eip_public_address_adjust"
description: |-
  Provides a resource to create a eip public_address_adjust
---

# tencentcloud_eip_public_address_adjust

Provides a resource to create a eip public_address_adjust

## Example Usage

```hcl
resource "tencentcloud_eip_public_address_adjust" "public_address_adjust" {
  instance_id = "ins-cr2rfq78"
  address_id  = "eip-erft45fu"
}
```

## Argument Reference

The following arguments are supported:

* `address_id` - (Optional, String, ForceNew) A unique ID that identifies an EIP instance. The unique ID of EIP is in the form:`eip-erft45fu`.
* `instance_id` - (Optional, String, ForceNew) A unique ID that identifies the CVM instance. The unique ID of CVM is in the form:`ins-osckfnm7`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



