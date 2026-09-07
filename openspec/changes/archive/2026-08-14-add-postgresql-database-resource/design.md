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
- Follow the established provider patterns: service-layer encapsulation of SDK calls, retry logic via `resource.Retry(tccommon.WriteRetryTimeout, ...)`, `ratelimit.Check`, `tccommon.RetryError`, and `tccommon.FILED_SP` separator for composite IDs.

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
- `encoding` (Optional, Computed, ForceNew) — maps to `request.Encoding`.
- `collate` (Optional, Computed, ForceNew) — maps to `request.Collate`.
- `ctype` (Optional, Computed, ForceNew) — maps to `request.Ctype`.

**Rationale**: The cloud API only supports modifying `DatabaseOwner` via `ModifyDatabaseOwner`. The `encoding`, `collate`, and `ctype` fields are set at creation time and cannot be modified afterward, so they are marked `ForceNew`. The `db_instance_id` and `database_name` form the composite ID and are also `ForceNew`. The `encoding`, `collate`, and `ctype` fields are additionally marked `Computed` because the backend applies defaults when they are omitted (e.g. `collate`/`ctype` default to `C`); without `Computed`, Terraform would diff the server-defaulted value against `null` and force an unnecessary destroy-and-recreate replacement.

### 3. Read strategy
**Decision**: Encapsulate the read in a `PostgresqlService` method `DescribePostgresqlDatabaseById`, which calls `DescribeDatabases` with `DBInstanceId` and a `database-name` filter, iterates the returned `Databases` array to find the exact match by name, and returns the matched `Database` (or `nil` when not found). The Read function then sets schema fields from the returned `Database` struct.

**Rationale**: `DescribeDatabases` supports filter `database-name` (string, fuzzy match). Since we know the exact name, we filter and then match precisely from the result set. Wrapping the SDK call in the service layer applies `ratelimit.Check` and keeps the query logic reusable (e.g. for a future data source).

### 4. Update strategy
**Decision**: In the Update function, check if `database_owner` has changed. If so, call the `PostgresqlService` method `ModifyPostgresqlDatabaseOwner`, which invokes `ModifyDatabaseOwner`. For the immutable fields (`encoding`, `collate`, `ctype`), they are `ForceNew` so Terraform handles recreation automatically. Use the `mutableArgs` / immutable args check pattern: since only `database_owner` is mutable, any change to `encoding`/`collate`/`ctype` triggers `ForceNew` recreation at the Terraform SDK level, so no explicit immutable check is needed in the update method — but we include the pattern for safety.

**Rationale**: The cloud API `ModifyDatabaseOwner` only accepts `DBInstanceId`, `DatabaseName`, and `DatabaseOwner`. No other fields can be updated.

### 5. Delete strategy
**Decision**: In the Delete function, call the `PostgresqlService` method `DeletePostgresqlDatabaseById`, which invokes `DeleteDatabase` with `DBInstanceId` and `DatabaseName` and wraps the call with `ratelimit.Check` and `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

### 6. Service-layer encapsulation
**Decision**: Encapsulate all SDK calls in the `PostgresqlService` (in `service_tencentcloud_postgresql.go`) as four reusable methods: `CreatePostgresqlDatabase`, `DescribePostgresqlDatabaseById`, `ModifyPostgresqlDatabaseOwner`, and `DeletePostgresqlDatabaseById`. The service accesses the postgres client via `me.client.UsePostgresqlClient()`. Write methods apply `ratelimit.Check(request.GetAction())` and `resource.Retry(tccommon.WriteRetryTimeout, ...)` (returning `tccommon.RetryError` on SDK failure and `NonRetryableError` on a nil response); the read method applies `ratelimit.Check(request.GetAction())`.

**Rationale**: This follows the established provider pattern — the `sqlserver` package routes all DB CRUD through `SqlserverService` methods (`resource_tc_sqlserver_db.go`), and the `postgresql` package already routes account reads/deletes through `PostgresqlService` methods. `UsePostgresqlClient()` returns the postgres client (`*postgresql.Client`) for the `postgres/v20170312` SDK package (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312`), which is already vendored. Keeping SDK calls in the service layer centralizes rate-limit and retry logic and keeps it reusable.

### 7. File location
**Decision**: Place the resource in `tencentcloud/services/postgresql/resource_tc_postgresql_database.go` following the naming convention `resource_tc_<service>_<name>.go`.

**Rationale**: Existing postgres resources are in the `postgresql` service directory. The resource is named `tencentcloud_postgresql_database` (using the `postgresql` prefix consistent with the other resources in the `postgresql` service directory).

### 8. Testing
**Decision**: Use gomonkey mock for unit tests (no TF_ACC test suite). Mock the `PostgresqlService` methods (`CreatePostgresqlDatabase`, `DescribePostgresqlDatabaseById`, `ModifyPostgresqlDatabaseOwner`, `DeletePostgresqlDatabaseById`) to test business logic in CRUD functions.

**Rationale**: Per project rules, new terraform resources should use mock-based unit tests rather than the terraform test suite. Mocking the service layer (rather than the raw SDK client) matches the resource's actual call path and keeps the tests decoupled from the SDK.

## Risks / Trade-offs

- **[Risk] DescribeDatabases uses fuzzy match for `database-name` filter** → Mitigation: After calling the API, iterate the returned `Databases` array and match the exact `DatabaseName` to ensure we read the correct database, not a fuzzy-matched one.
- **[Risk] CreateDatabase returns no ID** → Mitigation: Use composite ID from input parameters (`db_instance_id#database_name`). Verify the database was actually created by calling Read after Create.
- **[Risk] Race condition: database may take time to appear after creation** → Mitigation: After `CreateDatabase` succeeds, call the Read function to verify the database exists and populate state.
- **[Risk] `encoding`/`collate`/`ctype` changes require recreation** → Mitigation: Mark these as `ForceNew` in the schema. This is expected behavior since the cloud API does not support modifying these fields post-creation.
