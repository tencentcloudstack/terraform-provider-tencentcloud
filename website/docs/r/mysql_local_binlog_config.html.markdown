---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_local_binlog_config"
sidebar_current: "docs-tencentcloud-resource-mysql_local_binlog_config"
description: |-
  Provides a resource to create a mysql local_binlog_config
---

# tencentcloud_mysql_local_binlog_config

Provides a resource to create a mysql local_binlog_config

## Example Usage

```hcl
resource "tencentcloud_mysql_local_binlog_config" "local_binlog_config" {
  instance_id = "cdb-fitq5t9h"
  save_hours  = 140
  max_usage   = 50
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID in the format of cdb-c1nl9rpv. It is the same as the instance ID displayed in the TencentDB console.
* `max_usage` - (Required, Int) Space utilization of local binlog. Value range: [30,50].
* `save_hours` - (Required, Int) Retention period of local binlog. Valid range: 72-168 hours. When there is disaster recovery instance, the valid range will be 120-168 hours.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql local_binlog_config can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_local_binlog_config.local_binlog_config instance_id
```

