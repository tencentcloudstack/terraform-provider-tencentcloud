---
subcategory: "tdcpg"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdcpg_cluster"
sidebar_current: "docs-tencentcloud-resource-tdcpg_cluster"
description: |-
  Provides a resource to create a tdcpg cluster
---

# tencentcloud_tdcpg_cluster

Provides a resource to create a tdcpg cluster

## Example Usage

```hcl
resource "tencentcloud_tdcpg_cluster" "cluster" {
  zone                 = "ap-guangzhou-3"
  master_user_password = ""
  cpu                  = 1
  memory               = 1
  vpc_id               = "vpc_id"
  subnet_id            = "subnet_id"
  pay_mode             = "POSTPAID_BY_HOUR"
  cluster_name         = "cluster_name"
  db_version           = "10.17"
  instance_count       = 1
  period               = 1
  project_id           = 0
}
```

## Argument Reference

The following arguments are supported:

* `cpu` - (Required, Int) cpu cores.
* `master_user_password` - (Required, String) user password.
* `memory` - (Required, Int) memory size.
* `pay_mode` - (Required, String) pay mode, the value is either PREPAID or POSTPAID_BY_HOUR.
* `subnet_id` - (Required, String) subnet id.
* `vpc_id` - (Required, String) vpc id.
* `zone` - (Required, String) available zone.
* `cluster_name` - (Optional, String) cluster name.
* `db_version` - (Optional, String) community version number, default to 10.17.
* `instance_count` - (Optional, Int) instance count.
* `period` - (Optional, Int) purchase time, required when PayMode is PREPAID, the value range is 1~60, default to 1.
* `project_id` - (Optional, Int) project id, default to 0, means default project.
* `storage` - (Optional, Int) max storage, the unit is GB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdcpg cluster can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdcpg_cluster.cluster cluster_id
```

