## Context

The Terraform Provider for TencentCloud currently supports managing PostgreSQL DB instances (e.g. `tencentcloud_postgresql_instance`), but does not provide a resource to manage individual databases *within* a DB instance. Users must manually create databases after provisioning an instance, which breaks the infrastructure-as-code workflow.

The TencentCloud postgres SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312`, already vendored) exposes four APIs that enable full CRUD management of a database:

- `CreateDatabase` — creates a database; returns only RequestId (no ID).
- `DescribeDatabases` — lists databases for a DB instance, with filter support.
- `ModifyDatabaseOwner` — changes the database owner.
- `DeleteDatabase` — deletes a database.

## Goals / Non-Goals

**Goals:**
- Provide `tencentcloud_postgresql_database` resource (RESOURCE_KIND_GENERAL) with full CRUD lifecycle.
- Allow users to create, read, update (owner), and delete a PostgreSQL database within a DB instance.
- Support import via composite ID.
- Follow the established provider patterns: retry logic via `tccommon.ReadRetryTimeout` / `tccommon.WriteRetryTimeout`, `tccommon.RetryError`, `tccommon.FILED_SP` separator for composite IDs.

**Non-Goals:**
- Not managing database-level privileges or connections (out of scope — the cloud API does not provide create/update/delete for those).
- Not exposing the `filters`, `items`, or `databases` fields from `DescribeDatabases` as schema fields (those are query-only fields for the Describe API, not resource attributes).

## Decisions

### 1. Composite ID strategy
**Decision**: Use `db_instance_id` + `FILED_SP` + `database_name` as the resource ID.

**Rationale**: `CreateDatabase` does not return any ID. The database is uniquely identified by the combination of `DBInstanceId` + `DatabaseName`. This follows the existing provider convention (e.g. `tencentcloud_igtm_strategy` uses `instanceId#strategyId`).

### 2. Schema design
**Decision**: Define the following schema fields:
- `db_instance_id` (Required, ForceNew) — maps to `request.DBInstanceId`.
- `database_name` (Required, ForceNew) — maps to `request.DatabaseName`.
- `database_owner` (Required) — maps to `request.DatabaseOwner`; updatable via `ModifyDatabaseOwner`.
- `encoding` (Optional, ForceNew) — maps to `request.Encoding`.
- `collate` (Optional, ForceNew) — maps to `request.Collate`.
- `ctype` (Optional, ForceNew) — maps to `request.Ctype`.

**Rationale**: The cloud API only supports modifying `DatabaseOwner` via `ModifyDatabaseOwner`. The `encoding`, `collate`, and `ctype` fields are set at creation time and cannot be modified afterward, so they are marked `ForceNew`. The `db_instance_id` and `database_name` form the composite ID and are also `ForceNew`.

### 3. Read strategy
**Decision**: In the Read function, call `DescribeDatabases` with `DBInstanceId` and a `database-name` filter to fetch the specific database. Iterate through the `Databases` array in the response to find the matching entry by name. Set schema fields from the `Database` struct.

**Rationale**: `DescribeDatabases` supports filter `database-name` (string, fuzzy match). Since we know the exact name, we filter and then match precisely from the result set. This follows the project pattern of using `tccommon.ReadRetryTimeout` for read retries.

### 4. Update strategy
**Decision**: In the Update function, check if `database_owner` has changed. If so, call `ModifyDatabaseOwner`. For the immutable fields (`encoding`, `collate`, `ctype`), they are `ForceNew` so Terraform handles recreation automatically. Use the `mutableArgs` / immutable args check pattern: since only `database_owner` is mutable, any change to `encoding`/`collate`/`ctype` triggers `ForceNew` recreation at the Terraform SDK level, so no explicit immutable check is needed in the update method — but we include the pattern for safety.

**Rationale**: The cloud API `ModifyDatabaseOwner` only accepts `DBInstanceId`, `DatabaseName`, and `DatabaseOwner`. No other fields can be updated.

### 5. Delete strategy
**Decision**: In the Delete function, call `DeleteDatabase` with `DBInstanceId` and `DatabaseName`. Use `tccommon.WriteRetryTimeout` for retries.

### 6. API client
**Decision**: Use `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient()` to access the postgres client.

**Rationale**: `UsePostgresqlClient()` returns the postgres client (`*postgre.Client`) for the `postgres/v20170312` SDK package (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312`), which is already vendored.

### 7. File location
**Decision**: Place the resource in `tencentcloud/services/postgresql/resource_tc_postgresql_database.go` following the naming convention `resource_tc_<service>_<name>.go`.

**Rationale**: Existing postgres resources are in the `postgresql` service directory. The resource is named `tencentcloud_postgresql_database` (using the `postgresql` prefix consistent with the other resources in the `postgresql` service directory).

### 8. Testing
**Decision**: Use gomonkey mock for unit tests (no TF_ACC test suite). Mock the cloud API client methods to test business logic in CRUD functions.

**Rationale**: Per project rules, new terraform resources should use mock-based unit tests rather than the terraform test suite.

## Risks / Trade-offs

- **[Risk] DescribeDatabases uses fuzzy match for `database-name` filter** → Mitigation: After calling the API, iterate the returned `Databases` array and match the exact `DatabaseName` to ensure we read the correct database, not a fuzzy-matched one.
- **[Risk] CreateDatabase returns no ID** → Mitigation: Use composite ID from input parameters (`db_instance_id#database_name`). Verify the database was actually created by calling Read after Create.
- **[Risk] Race condition: database may take time to appear after creation** → Mitigation: After `CreateDatabase` succeeds, call the Read function which uses `tccommon.ReadRetryTimeout` retry logic to wait for the database to appear.
- **[Risk] `encoding`/`collate`/`ctype` changes require recreation** → Mitigation: Mark these as `ForceNew` in the schema. This is expected behavior since the cloud API does not support modifying these fields post-creation.
