## 1. Resource Schema & CRUD Implementation

- [x] 1.1 Create `tencentcloud/services/postgresql/resource_tc_postgresql_database.go` with `ResourceTencentCloudPostgresqlDatabase()` function defining the schema: `db_instance_id` (Required, ForceNew), `database_name` (Required, ForceNew), `database_owner` (Required), `encoding` (Optional, Computed, ForceNew), `collate` (Optional, Computed, ForceNew), `ctype` (Optional, Computed, ForceNew), and Importer support
- [x] 1.2 Implement `resourceTencentCloudPostgresqlDatabaseCreate()` — construct a `PostgresqlService`, call `service.CreatePostgresqlDatabase(...)`, set composite ID `db_instance_id#database_name` using `tccommon.FILED_SP`, then call Read
- [x] 1.3 Implement `resourceTencentCloudPostgresqlDatabaseRead()` — split composite ID to get `db_instance_id` and `database_name`, call `service.DescribePostgresqlDatabaseById(...)`, set schema fields from the returned `Database` struct (skip nil fields), handle not-found by printing log with ID then `d.SetId("")`
- [x] 1.4 Implement `resourceTencentCloudPostgresqlDatabaseUpdate()` — split composite ID, check if `database_owner` changed, call `service.ModifyPostgresqlDatabaseOwner(...)`, then call Read
- [x] 1.5 Implement `resourceTencentCloudPostgresqlDatabaseDelete()` — split composite ID, call `service.DeletePostgresqlDatabaseById(...)`
- [x] 1.6 Implement service-layer methods in `tencentcloud/services/postgresql/service_tencentcloud_postgresql.go`: `CreatePostgresqlDatabase`, `DescribePostgresqlDatabaseById`, `ModifyPostgresqlDatabaseOwner`, `DeletePostgresqlDatabaseById` — each write method wraps the SDK call with `ratelimit.Check(request.GetAction())` and `resource.Retry(tccommon.WriteRetryTimeout, ...)`; the read method applies `ratelimit.Check(request.GetAction())` and returns `nil` when not found

## 2. Provider Registration

- [x] 2.1 Add `"tencentcloud_postgresql_database": postgresql.ResourceTencentCloudPostgresqlDatabase(),` to the `ResourcesMap` in `tencentcloud/provider.go` (in the PostgreSQL section, after `tencentcloud_postgres_audit_service`)
- [x] 2.2 Add `tencentcloud_postgresql_database` to the PostgreSQL resource list in `tencentcloud/provider.md` (after `tencentcloud_postgres_audit_service`)

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/postgresql/resource_tc_postgresql_database.md` with one-line description referencing TencentDB for PostgreSQL, a Note documenting the `database_owner` dependency on an existing account, an Example Usage section (HCL) showing the full dependency chain (vpc → subnet → instance → account → database), and Import section documenting the composite ID format `db_instance_id#database_name`

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/postgresql/resource_tc_postgresql_database_test.go` using gomonkey mock to mock the `PostgresqlService` methods (`CreatePostgresqlDatabase`, `DescribePostgresqlDatabaseById`, `ModifyPostgresqlDatabaseOwner`, `DeletePostgresqlDatabaseById`) and test the CRUD business logic — no terraform test suite, no `go test` execution

## 5. Verification (performed by separate process)

- [x] 5.1 Run `gofmt` formatting on generated Go files (handled by tfpacer-finalize skill)
- [x] 5.2 Run `make doc` to generate website docs (handled by tfpacer-finalize skill)
- [x] 5.3 Verify Go code compiles correctly (handled by separate build process)
