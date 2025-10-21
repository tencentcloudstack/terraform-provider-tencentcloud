---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_eip"
sidebar_current: "docs-tencentcloud-datasource-dayu_eip"
description: |-
  Use this data source to query dayu eip rules
---

# tencentcloud_dayu_eip

Use this data source to query dayu eip rules

## Example Usage

```hcl
data "tencentcloud_dayu_eip" "test" {
  resource_id = "bgpip-000004xg"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String) Id of the resource.
* `bind_status` - (Optional, List: [`String`]) The binding state of the instance, value range [BINDING, BIND, UNBINDING, UNBIND], default is [BINDING, BIND, UNBINDING, UNBIND].
* `limit` - (Optional, Int) The number of pages, default is `10`.
* `offset` - (Optional, Int) The page start offset, default is `0`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of layer 4 rules. Each element contains the following attributes:


