---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_logset"
sidebar_current: "docs-tencentcloud-resource-cls_logset"
description: |-
Provides a resource to create a CLS logset.
---
//asdadad
# tencentcloud_cls_logset

Provides a resource to create a CLS logset.

## Example Usage

```hcl
resource "tencentcloud_cls_Logset" "logset_basic" {
  logset_name    = "test"
}
```

## Argument Reference

The following arguments are supported:

* `logset_name` - (Required) Name of the LogSet.
* `tags` - (Optional) The label key value pair of the binding for the log set. A maximum of 10 tag key value pairs are supported.
    * `key` - (Required) Tag Key
    * `value` - (Required) Tag value

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## Import

CLS logset can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_logset.logset_basic 90496
```
