## 1. Service Layer

- [x] 1.1 Add `DescribePostgresqlReadonlyInstanceV2ById` method to `service_tencentcloud_postgresql.go` that calls `DescribeDBInstanceAttribute` with retry and returns `(*postgresql.DBInstance, has bool, err error)`
- [x] 1.2 Add `CreatePostgresqlReadonlyInstanceV2` method to `service_tencentcloud_postgresql.go` that wraps `CreateReadOnlyDBInstance` call (request passed in, returns response)
- [x] 1.3 Add `IsolatePostgresqlReadonlyInstanceV2` method to `service_tencentcloud_postgresql.go` that wraps `IsolateDBInstances` with `DBInstanceIdSet` from instance ID

## 2. Resource Schema Definition

- [x] 2.1 Create `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance_v2.go` with `ResourceTencentCloudPostgresqlReadonlyInstanceV2()` function defining schema for all CreateReadOnlyDBInstance parameters (zone, master_db_instance_id, spec_code, storage, instance_count, period as Required+ForceNew; vpc_id, subnet_id, instance_charge_type, auto_voucher, voucher_ids, auto_renew_flag, project_id, activity_id, read_only_group_id, tag_list, security_group_ids, need_support_ipv6, name, db_version, dedicated_cluster_id, deletion_protection, tags as Optional; deal_names, bill_id, db_instance_id_set, billing_parameters, db_instance_id as Computed)
- [x] 2.2 Add `Timeouts` block (Create/Delete) to the schema for async operations
- [x] 2.3 Add `Importer` with `schema.ImportStatePassthrough` to support `terraform import`

## 3. Create Operation

- [x] 3.1 Implement `resourceTencentCloudPostgresqlReadonlyInstanceV2Create` - build `CreateReadOnlyDBInstanceRequest` from schema fields
- [x] 3.2 Wrap API call in `resource.Retry(tccommon.WriteRetryTimeout)` with `tccommon.RetryError()` error wrapping
- [x] 3.3 After retry, check `response == nil` / `Response == nil` / `DBInstanceIdSet` empty → log logId+d.Id() and return `NonRetryableError`
- [x] 3.4 Set `d.SetId()` to `DBInstanceIdSet[0]` outside retry block, after error handling
- [x] 3.5 Poll `DescribeDBInstanceAttribute` until instance status is `running` (async creation)
- [x] 3.6 Set computed fields: `deal_names`, `bill_id`, `db_instance_id_set`, `billing_parameters` from create response
- [x] 3.7 Handle `tags` field via tag service if provided

## 4. Read Operation

- [x] 4.1 Implement `resourceTencentCloudPostgresqlReadonlyInstanceV2Read` - call `DescribeDBInstanceAttribute` via service layer with `resource.Retry(tccommon.ReadRetryTimeout)`
- [x] 4.2 If `response == nil` or `DBInstance == nil`, print `log.Printf("[CRUD] postgresql_readonly_instance_v2 id=%s", d.Id())` then `d.SetId("")`
- [x] 4.3 Backfill all schema fields from `DBInstance` structure with nil-checks before each `d.Set()` (zone, vpc_id, subnet_id, db_version, storage, name, project_id, need_support_ipv6, deletion_protection, master_db_instance_id, etc.)
- [x] 4.4 Set `db_instance_id` = d.Id()

## 5. Update Operation

- [x] 5.1 Implement `resourceTencentCloudPostgresqlReadonlyInstanceV2Update` with immutableArgs check - add all non-ForceNew top-level fields to `immutableArgs` array
- [x] 5.2 If any immutable field has changed, return error via `helper.ImmutableArgsChek`
- [x] 5.3 Call Read at the end to refresh state

## 6. Delete Operation

- [x] 6.1 Implement `resourceTencentCloudPostgresqlReadonlyInstanceV2Delete` - call `IsolateDBInstances` via service layer with `resource.Retry(tccommon.WriteRetryTimeout)` and `tccommon.RetryError()` wrapping
- [x] 6.2 Poll `DescribeDBInstanceAttribute` until instance status is `isolated` (async isolation)

## 7. Provider Registration

- [x] 7.1 Add `tencentcloud_postgresql_readonly_instance_v2` resource registration in `tencentcloud/provider.go`
- [x] 7.2 Add `tencentcloud_postgresql_readonly_instance_v2` entry in `tencentcloud/provider.md`

## 8. Documentation

- [x] 8.1 Create `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance_v2.md` with one-line description (mentioning PostgreSQL product), Example Usage (using jsonencode for json string fields if needed), and Import section (mentioning instance ID)

## 9. Unit Tests

- [x] 9.1 Create `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance_v2_test.go` using gomonkey to mock cloud API calls
- [x] 9.2 Add test case for Create success path (mock CreateReadOnlyDBInstance returns valid response, verify d.Id() is set)
- [x] 9.3 Add test case for Create empty response (mock returns empty DBInstanceIdSet, verify NonRetryableError)
- [x] 9.4 Add test case for Read success path (mock DescribeDBInstanceAttribute returns running instance, verify fields set)
- [x] 9.5 Add test case for Read empty response (mock returns nil, verify d.SetId("") called)
- [x] 9.6 Add test case for Delete success path (mock IsolateDBInstances success, verify isolation polling)

## 10. Verification

- [x] 10.1 Verify all new Go files compile correctly (check imports, function signatures, type usage against vendor SDK)
- [x] 10.2 Verify CRUD code parameters match cloud API: Create params exist in CreateReadOnlyDBInstance, Read params exist in DescribeDBInstanceAttribute, Delete params exist in IsolateDBInstances
- [x] 10.3 Verify all function error returns are checked (use `_ = func()` for functions that cannot fail)
