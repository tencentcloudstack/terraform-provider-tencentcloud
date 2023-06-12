---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_unit_namespace"
sidebar_current: "docs-tencentcloud-resource-tsf_unit_namespace"
description: |-
  Provides a resource to create a tsf unit_namespace
---

# tencentcloud_tsf_unit_namespace

Provides a resource to create a tsf unit_namespace

## Example Usage

```hcl
resource "tencentcloud_tsf_unit_namespace" "unit_namespace" {
  gateway_instance_id = "gw-ins-lvdypq5k"
  namespace_id        = "namespace-vwgo38wy"
  namespace_name      = "keep-terraform-cls"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_instance_id` - (Required, String, ForceNew) gateway instance Id.
* `namespace_id` - (Required, String, ForceNew) namespace id.
* `namespace_name` - (Required, String, ForceNew) namespace name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Create time. Note: This field may return null, indicating that no valid value was found.
* `updated_time` - Update time. Note: This field may return null, indicating that no valid value was found.


## Import

tsf unit_namespace can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_unit_namespace.unit_namespace gw-ins-lvdypq5k#namespace-vwgo38wy
```

