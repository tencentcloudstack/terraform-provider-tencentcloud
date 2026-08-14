## 1. Resource Schema & CRUD Implementation

- [x] 1.1 Create `tencentcloud/services/postgresql/resource_tc_postgresql_database.go` with `ResourceTencentCloudPostgresqlDatabase()` function defining the schema: `db_instance_id` (Required, ForceNew), `database_name` (Required, ForceNew), `database_owner` (Required), `encoding` (Optional, ForceNew), `collate` (Optional, ForceNew), `ctype` (Optional, ForceNew), and Importer support
- [x] 1.2 Implement `resourceTencentCloudPostgresqlDatabaseCreate()` — build `CreateDatabaseRequest` from schema, call `CreateDatabaseWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`, verify response is not nil, set composite ID `db_instance_id#database_name` using `tccommon.FILED_SP`, then call Read
- [x] 1.3 Implement `resourceTencentCloudPostgresqlDatabaseRead()` — split composite ID to get `db_instance_id` and `database_name`, call `DescribeDatabasesWithContext` with `DBInstanceId` and `database-name` filter inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`, iterate `Databases` array to find exact match by name, set schema fields from `Database` struct (skip nil fields), handle not-found by printing log with ID then `d.SetId("")`
- [x] 1.4 Implement `resourceTencentCloudPostgresqlDatabaseUpdate()` — split composite ID, check if `database_owner` changed, call `ModifyDatabaseOwnerWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`, then call Read
- [x] 1.5 Implement `resourceTencentCloudPostgresqlDatabaseDelete()` — split composite ID, call `DeleteDatabaseWithContext` with `DBInstanceId` and `DatabaseName` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`

## 2. Provider Registration

- [x] 2.1 Add `"tencentcloud_postgresql_database": postgresql.ResourceTencentCloudPostgresqlDatabase(),` to the `ResourcesMap` in `tencentcloud/provider.go` (in the PostgreSQL section, after `tencentcloud_postgres_audit_service`)
- [x] 2.2 Add `tencentcloud_postgresql_database` to the PostgreSQL resource list in `tencentcloud/provider.md` (after `tencentcloud_postgres_audit_service`)

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/postgresql/resource_tc_postgresql_database.md` with one-line description referencing TencentDB for PostgreSQL, Example Usage section (HCL), and Import section documenting the composite ID format `db_instance_id#database_name`

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/postgresql/resource_tc_postgresql_database_test.go` using gomonkey mock to mock the postgres API client methods (`CreateDatabaseWithContext`, `DescribeDatabasesWithContext`, `ModifyDatabaseOwnerWithContext`, `DeleteDatabaseWithContext`) and test the CRUD business logic — no terraform test suite, no `go test` execution

## 5. Verification (performed by separate process)

- [x] 5.1 Run `gofmt` formatting on generated Go files (handled by tfpacer-finalize skill)
- [x] 5.2 Run `make doc` to generate website docs (handled by tfpacer-finalize skill)
- [x] 5.3 Verify Go code compiles correctly (handled by separate build process)
