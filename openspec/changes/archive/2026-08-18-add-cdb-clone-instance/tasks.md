## 1. Schema & Resource Definition

- [x] 1.1 Create `tencentcloud/services/cdb/resource_tc_cdb_clone_instance.go` with the resource schema definition including all top-level fields (`instance_id`, `specified_rollback_time`, `specified_backup_id`, `uniq_vpc_id`, `uniq_subnet_id`, `memory`, `volume`, `instance_name`, `security_group`, `resource_tags`, `cpu`, `protect_mode`, `deploy_mode`, `slave_zone`, `backup_zone`, `device_type`, `instance_nodes`, `deploy_group_id`, `dry_run`, `cage_id`, `project_id`, `pay_type`, `period`, `cluster_topology`, `src_region`, `specified_sub_backup_id`, `master_zone`, `zone`, `fourth_zone`, `async_request_id`) with correct types, Required/Optional/ForceNew/Computed flags per the spec
- [x] 1.2 Add nested block schemas for `resource_tags` (key, value) and `cluster_topology` (per `ClusterTopology` struct fields) with `MaxItems: 1` for `cluster_topology`
- [x] 1.3 Add `Timeouts` block to the resource schema: Create (default 60 min), Update (default 60 min), Delete (default 20 min)
- [x] 1.4 Add `Importer` with `schema.ImportStatePassthrough` to the resource schema
- [x] 1.5 Register the `resourceTencentCloudCdbCloneInstanceCreate`, `resourceTencentCloudCdbCloneInstanceRead`, `resourceTencentCloudCdbCloneInstanceUpdate`, `resourceTencentCloudCdbCloneInstanceDelete` functions and `ResourceTencentCloudCdbCloneInstance()` factory

## 2. CRUD Implementation

- [x] 2.1 Implement `resourceTencentCloudCdbCloneInstanceCreate`: call `CreateCloneInstance` with retry (`tccommon.WriteRetryTimeout`); validate response and `AsyncRequestId` non-empty (return `NonRetryableError` if empty); set `async_request_id`; poll `DescribeAsyncRequestInfo` until `SUCCESS`; query `DescribeCloneList` with source `InstanceId` to extract `DstInstanceId`; set resource ID to `DstInstanceId`; call Read
- [x] 2.2 Implement `resourceTencentCloudCdbCloneInstanceRead`: call `DescribeDBInstances` with retry (`tccommon.ReadRetryTimeout`), `InstanceIds=[d.Id()]`, `QueryClusterInfo=true`; if items empty, log `[CRUD] cdb_clone_instance id=%s` then `d.SetId("")`; set schema fields only when response fields are non-nil; do not overwrite create-only/source fields
- [x] 2.3 Implement `resourceTencentCloudCdbCloneInstanceUpdate`: check `immutableArgs` for non-updatable changed fields (return error if found); build `UpgradeDBInstanceRequest` with `InstanceId=d.Id()` and changed updatable fields (`memory`, `volume`, `cpu`, `protect_mode`, `deploy_mode`, `slave_zone`, `backup_zone`, `device_type`, `cluster_topology`, `fourth_zone`); call with retry; validate response and `AsyncRequestId`; poll `DescribeAsyncRequestInfo` until `SUCCESS`; call Read
- [x] 2.4 Implement `resourceTencentCloudCdbCloneInstanceDelete`: call `IsolateDBInstance` with retry; poll `DescribeDBInstanceById` until `MYSQL_STATUS_ISOLATED` (return nil if instance gone); call `OfflineIsolatedInstances` with `InstanceIds=[d.Id()]`; poll `DescribeIsolatedDBInstanceById` until instance disappears

## 3. Service Layer

- [x] 3.1 Add `CreateCloneInstance` helper method in `tencentcloud/services/cdb/service_tencentcloud_mysql.go` to encapsulate the `CreateCloneInstance` API call with parameter mapping (including `ResourceTags` []*TagInfo, `SecurityGroup` []*string, `ClusterTopology` struct)
- [x] 3.2 Add `UpgradeCdbCloneInstance` helper method in the CDB service layer to encapsulate the `UpgradeDBInstance` API call with parameter mapping
- [x] 3.3 Reuse existing `DescribeDBInstanceById`, `DescribeIsolatedDBInstanceById`, `DescribeAsyncRequestInfo`, `DescribeMysqlCloneListByFilter`, `OfflineIsolatedInstances`, `IsolateDBInstance` helpers (verify they exist and are exported in the MysqlService)

## 4. Provider Registration

- [x] 4.1 Register `tencentcloud_cdb_clone_instance` resource in `tencentcloud/provider.go` with the factory function `cdb.ResourceTencentCloudCdbCloneInstance()`
- [x] 4.2 Update `tencentcloud/provider.md` to include the `tencentcloud_cdb_clone_instance` resource entry

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/cdb/resource_tc_cdb_clone_instance_test.go` with gomonkey-based unit tests (mock the CDB SDK client methods) for Create (mock `CreateCloneInstance`, `DescribeAsyncRequestInfo`, `DescribeCloneList`, `DescribeDBInstances`), Read, Update (mock `UpgradeDBInstance`, `DescribeAsyncRequestInfo`, `DescribeDBInstances`), and Delete (mock `IsolateDBInstance`, `DescribeDBInstanceById`, `OfflineIsolatedInstances`, `DescribeIsolatedDBInstanceById`)
- [x] 5.2 Verify the generated code is compilable in the current environment (do NOT run `go build`/`go vet`/`go test`)

## 6. Documentation

- [x] 6.1 Create `tencentcloud/services/cdb/resource_tc_cdb_clone_instance.md` following gendoc format: one-line description mentioning CDB, Example Usage (clone with rollback time, clone with backup id), and Import section (RESOURCE_KIND_GENERAL supports import)
- [x] 6.2 Ensure `make doc` generates the corresponding `website/docs/r/` documentation (executed in the finalize phase via tfpacer-finalize skill)
