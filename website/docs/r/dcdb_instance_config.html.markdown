---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_instance_config"
sidebar_current: "docs-tencentcloud-resource-dcdb_instance_config"
description: |-
  Provides a resource to create a dcdb instance_config
---

# tencentcloud_dcdb_instance_config

Provides a resource to create a dcdb instance_config

## Example Usage

```hcl
resource "tencentcloud_dcdb_instance_config" "instance_config" {
  instance_id        = local.dcdb_id
  rs_access_strategy = 0
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `rs_access_strategy` - (Required, Int) RS nearest access mode, 0-no policy, 1-nearest access.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dcdb instance_config can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_instance_config.instance_config instance_config_id
```

