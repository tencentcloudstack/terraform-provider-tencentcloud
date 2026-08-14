## Why

Terraform Provider for TencentCloud currently lacks the ability to manage PostgreSQL databases within a DB instance. Users need to create, read, update, and delete individual databases on their PostgreSQL instances through Terraform, enabling infrastructure-as-code management of database-level resources. This fills a gap in the postgres service coverage.

## What Changes

- Add new resource `tencentcloud_postgresql_database` (RESOURCE_KIND_GENERAL) that manages the full lifecycle of a PostgreSQL database within a DB instance.
- Implement Create via `CreateDatabase` API, Read via `DescribeDatabases` API, Update via `ModifyDatabaseOwner` API, and Delete via `DeleteDatabase` API.
- Use composite ID `db_instance_id#database_name` since `CreateDatabase` does not return an ID.
- Register the resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.
- Add documentation file `resource_tc_postgresql_database.md`.
- Add unit tests using gomonkey mock (no TF_ACC test suite).

## Capabilities

### New Capabilities
- `postgresql-database-resource`: Resource to manage a PostgreSQL database within a DB instance, supporting create (CreateDatabase), read (DescribeDatabases), update (ModifyDatabaseOwner for owner changes), and delete (DeleteDatabase) operations.

### Modified Capabilities
<!-- None -->

## Impact

- **New files**:
  - `tencentcloud/services/postgresql/resource_tc_postgresql_database.go` — resource schema and CRUD logic
  - `tencentcloud/services/postgresql/resource_tc_postgresql_database_test.go` — unit tests with gomonkey mock
  - `tencentcloud/services/postgresql/resource_tc_postgresql_database.md` — documentation
- **Modified files**:
  - `tencentcloud/provider.go` — register `tencentcloud_postgresql_database` resource
  - `tencentcloud/provider.md` — add resource to provider documentation
- **Cloud APIs**: `CreateDatabase`, `DescribeDatabases`, `ModifyDatabaseOwner`, `DeleteDatabase` from `postgres/v20170312` SDK package (already vendored).
- **Dependencies**: No new external dependencies; uses existing `tencentcloud-sdk-go` postgres package and `tccommon` helpers.
