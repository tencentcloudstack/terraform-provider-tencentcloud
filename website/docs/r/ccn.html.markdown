---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn"
sidebar_current: "docs-tencentcloud-resource-ccn"
description: |-
  Provides a resource to create a CCN instance.
---

# tencentcloud_ccn

Provides a resource to create a CCN instance.

## Example Usage

```hcl
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the CCN to be queried, and maximum length does not exceed 60 bytes.
* `description` - (Optional) Description of CCN, and maximum length does not exceed 100 bytes.
* `qos` - (Optional, ForceNew) Service quality of CCN. Valid values: `PT`, `AU`, `AG`. The default is `AU`.
* `tags` - (Optional) Instance tag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of resource.
* `instance_count` - Number of attached instances.
* `state` - States of instance. Valid values: `ISOLATED`(arrears) and `AVAILABLE`.


## Import

Ccn instance can be imported, e.g.

```
$ terraform import tencentcloud_ccn.test ccn-id
```

