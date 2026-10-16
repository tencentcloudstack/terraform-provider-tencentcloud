---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_enable_unit_rule"
sidebar_current: "docs-tencentcloud-resource-tsf_enable_unit_rule"
description: |-
  Provides a resource to create a tsf enable_unit_rule
---

# tencentcloud_tsf_enable_unit_rule

Provides a resource to create a tsf enable_unit_rule

## Example Usage

```hcl
resource "tencentcloud_tsf_enable_unit_rule" "enable_unit_rule" {
  rule_id = "unit-rl-is9m4nxz"
  switch  = "enabled"
}
```

## Argument Reference

The following arguments are supported:

* `rule_id` - (Required, String) api ID.
* `switch` - (Required, String) switch, on: `enabled`, off: `disabled`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DisableUnitRule](https://www.tencentcloud.com/search/?wd=DisableUnitRule&type=document&lang=en)
* Create: [EnableUnitRule](https://www.tencentcloud.com/search/?wd=EnableUnitRule&type=document&lang=en)
* Read: [DescribeUnitRule](https://www.tencentcloud.com/search/?wd=DescribeUnitRule&type=document&lang=en)
* Update: [DisableUnitRule](https://www.tencentcloud.com/search/?wd=DisableUnitRule&type=document&lang=en)
* Update: [EnableUnitRule](https://www.tencentcloud.com/search/?wd=EnableUnitRule&type=document&lang=en)

## Import

tsf enable_unit_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_enable_unit_rule.enable_unit_rule enable_unit_rule_id
```

