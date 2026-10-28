## Why

TencentCloud CynosDB (TDSQL-C) supports LibraDB read-only analytics engine instances that can be attached to existing clusters. Currently, there is no Terraform resource to manage the lifecycle (create/read/delete) of these LibraDB instances. Users need a Terraform resource to declaratively attach and detach LibraDB analytics engine instances to/from CynosDB clusters.

## What Changes

- Add a new Terraform resource `tencentcloud_cynosdb_libra_db_instance` of type RESOURCE_KIND_ATTACHMENT to manage the binding of LibraDB read-only analytics engine instances to CynosDB clusters.
- The resource supports:
  - **Create**: Calls `AddLibraDBInstances` API to attach a LibraDB instance to a cluster.
  - **Read**: Calls `DescribeLibraDBInstanceDetail` API to query the instance detail.
  - **Delete**: Calls `IsolateLibraDBCluster` API to isolate (detach) the LibraDB cluster.
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.
- Add documentation in the corresponding `.md` file.

## Capabilities

### New Capabilities
- `cynosdb-libra-db-instance-attachment`: Manage the attachment of LibraDB read-only analytics engine instances to CynosDB clusters, including creating, reading, and deleting (isolating) the instance binding.

### Modified Capabilities
<!-- No existing capabilities are being modified -->

## Impact

- **Code**: New resource file `resource_tc_cynosdb_libra_db_instance_attachment.go`, test file, service layer additions, and documentation.
- **APIs**: Uses CynosDB APIs: `AddLibraDBInstances`, `DescribeLibraDBInstanceDetail`, `IsolateLibraDBCluster` from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107`.
- **Dependencies**: No new vendor dependencies needed (cynosdb SDK already exists in vendor).
- **Provider Registration**: New resource entry in `provider.go` and `provider.md`.
