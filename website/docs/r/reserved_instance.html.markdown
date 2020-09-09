---
subcategory: "CVM"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_reserved_instance"
sidebar_current: "docs-tencentcloud-resource-reserved_instance"
description: |-
  Provides a reserved instance resource.
---

# tencentcloud_reserved_instance

Provides a reserved instance resource.

~> **NOTE:** Reserved instance cannot be deleted and updated. The reserved instance still exist which can be extracted by reserved_instances data source when reserved instance is destroied.

## Example Usage

```hcl
resource "tencentcloud_reserved_instance" "ri" {
  config_id      = "469043dd-28b9-4d89-b557-74f6a8326259"
  instance_count = 2
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) Configuration id of the reserved instance.
* `instance_count` - (Required) Number of reserved instances to be purchased.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `end_time` - Expiry time of the RI.
* `start_time` - Start time of the RI.
* `status` - Status of the RI at the time of purchase.


## Import

Reserved instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_reserved_instance.foo 6cc16e7c-47d7-4fae-9b44-ce5c0f59a920
```

