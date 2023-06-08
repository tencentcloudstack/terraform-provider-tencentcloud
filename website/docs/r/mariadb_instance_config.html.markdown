---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_instance_config"
sidebar_current: "docs-tencentcloud-resource-mariadb_instance_config"
description: |-
  Provides a resource to create a mariadb instance_config
---

# tencentcloud_mariadb_instance_config

Provides a resource to create a mariadb instance_config

## Example Usage

```hcl
resource "tencentcloud_mariadb_instance_config" "test" {
  instance_id        = "tdsql-9vqvls95"
  vpc_id             = "vpc-ii1jfbhl"
  subnet_id          = "subnet-3ku415by"
  rs_access_strategy = 1
  extranet_access    = 0
  vip                = "127.0.0.1"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) instance id.
* `extranet_access` - (Optional, Int) External network status, 0-closed; 1- Opening; Default not enabled.
* `rs_access_strategy` - (Optional, Int) RS proximity mode, 0- no strategy, 1- access to the nearest available zone.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mariadb instance_config can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_instance_config.test id
```

