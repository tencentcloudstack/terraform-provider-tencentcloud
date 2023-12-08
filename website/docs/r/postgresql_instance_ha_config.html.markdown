---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_instance_ha_config"
sidebar_current: "docs-tencentcloud-resource-postgresql_instance_ha_config"
description: |-
  Provides a resource to set postgresql instance syncMode
---

# tencentcloud_postgresql_instance_ha_config

Provides a resource to set postgresql instance syncMode

## Example Usage

```hcl
resource "tencentcloud_postgresql_instance_ha_config" "example" {
  instance_id              = "postgres-gzg9jb2n"
  sync_mode                = "Semi-sync"
  max_standby_latency      = 10737418240
  max_standby_lag          = 10
  max_sync_standby_latency = 52428800
  max_sync_standby_lag     = 5
}
```



```hcl
resource "tencentcloud_postgresql_instance_ha_config" "example" {
  instance_id         = "postgres-gzg9jb2n"
  sync_mode           = "Async"
  max_standby_latency = 10737418240
  max_standby_lag     = 10
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `max_standby_lag` - (Required, Int) Maximum latency of highly available backup machines. When the delay time of the backup node is less than or equal to this value, and the amount of delay data of the backup node is less than or equal to MaxStandbyLatency, the primary node can be switched. Unit: s; Parameter range: [5, 10].
* `max_standby_latency` - (Required, Int) Maximum latency data volume for highly available backup machines. When the delay data amount of the backup node is less than or equal to this value, and the delay time of the backup node is less than or equal to MaxStandbyLag, it can switch to the main node. Unit: byte; Parameter range: [1073741824, 322122547200].
* `sync_mode` - (Required, String) Master slave synchronization method, Semi-sync: Semi synchronous; Async: Asynchronous. Main instance default value: Semi-sync, Read-only instance default value: Async.
* `max_sync_standby_lag` - (Optional, Int) Maximum delay time for synchronous backup. When the delay time of the standby machine is less than or equal to this value, and the amount of delay data of the standby machine is less than or equal to MaxSyncStandbyLatency, then the standby machine adopts synchronous replication; Otherwise, adopt asynchronous replication. This parameter value is valid for instances where SyncMode is set to Semi sync. When a semi synchronous instance prohibits degradation to asynchronous replication, MaxSyncStandbyLatency and MaxSyncStandbyLag are not set. When semi synchronous instances allow degenerate asynchronous replication, PostgreSQL version 9 instances must have MaxSyncStandbyLatency set and MaxSyncStandbyLag not set, while PostgreSQL version 10 and above instances must have MaxSyncStandbyLatency and MaxSyncStandbyLag set.
* `max_sync_standby_latency` - (Optional, Int) Maximum latency data for synchronous backup. When the amount of data delayed by the backup machine is less than or equal to this value, and the delay time of the backup machine is less than or equal to MaxSyncStandbyLag, then the backup machine adopts synchronous replication; Otherwise, adopt asynchronous replication. This parameter value is valid for instances where SyncMode is set to Semi sync. When a semi synchronous instance prohibits degradation to asynchronous replication, MaxSyncStandbyLatency and MaxSyncStandbyLag are not set. When semi synchronous instances allow degenerate asynchronous replication, PostgreSQL version 9 instances must have MaxSyncStandbyLatency set and MaxSyncStandbyLag not set, while PostgreSQL version 10 and above instances must have MaxSyncStandbyLatency and MaxSyncStandbyLag set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



