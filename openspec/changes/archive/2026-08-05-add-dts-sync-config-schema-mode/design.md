## Context

The `tencentcloud_dts_sync_config` resource (`tencentcloud/services/dts/resource_tc_dts_sync_config.go`) wraps the DTS `ConfigureSyncJob` (write) and `DescribeSyncJobs` (read) APIs from the vendored SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206`. The resource exposes a nested `objects` block containing a `databases` list, where each database entry currently maps fields like `db_name`, `new_db_name`, `db_mode`, `schema_name`, `new_schema_name`, `table_mode`, etc. to the SDK's `Database` struct.

The SDK's `Database` struct already defines a `SchemaMode *string` field (`json:"SchemaMode"`), described as the "schema selection mode" with enum values `All` (all objects under the current object) and `Partial` (some objects), used by PostgreSQL and SQL Server sync links. Neither the Terraform schema, the Update (`ConfigureSyncJob`) build path, nor the Read (`DescribeSyncJobs`) flatten path currently touch this field, so users cannot set or read it through Terraform.

## Goals / Non-Goals

**Goals:**
- Expose `schema_mode` as an `Optional` string field inside the `objects.databases[]` block so users can configure the per-database schema selection mode.
- Wire the new field into the Update path (`ConfigureSyncJob`) so the configured value is sent to the cloud API.
- Wire the new field into the Read path (`DescribeSyncJobs` via `DescribeDtsSyncConfigById`) so the value is read back into Terraform state.
- Maintain full backwards compatibility: no changes to existing fields, no breaking schema changes.

**Non-Goals:**
- No new top-level resource fields; the change is scoped strictly to the `objects.databases[]` block.
- No changes to the `DBItem` struct (`DescribeModifiedDatabases`-style APIs) — only the `Database` struct used by `ConfigureSyncJob`/`DescribeSyncJobs`.
- No validation of enum values (`All`/`Partial`) beyond what the cloud API enforces; the field passes through the string as-is.
- No changes to retry logic, state refresh, or the service-layer method signatures.

## Decisions

1. **Field placement: nested under `objects.databases[]`.**
   - Rationale: `SchemaMode` is a property of a single `Database` object in the SDK (`Objects.Databases[].SchemaMode`), mirroring the existing `db_mode`/`table_mode`/`view_mode` siblings. Placing it in the same block keeps the Terraform schema isomorphic with the cloud API object graph.
   - Alternative considered: a top-level field — rejected because it would not map to a per-database value.

2. **Schema attributes: `Optional`, `TypeString`, no `Computed`.**
   - Rationale: Follows the pattern of sibling mode fields (`table_mode`, `view_mode`, `function_mode`, etc.) in the same block, which are all `Optional` string fields. Marking it `Computed` is unnecessary because the cloud API returns the value only when set, and the Read path already guards with a nil check.
   - Alternative considered: `Optional + Computed` — rejected to match the established sibling-field convention in this resource and avoid terraform diff churn.

3. **Write path: set `database.SchemaMode` from state in the `objects` build block.**
   - Rationale: The Update function already iterates `databasesMap` and constructs a `dts.Database` per entry; adding a single `if v, ok := databasesMap["schema_mode"]; ok { database.SchemaMode = helper.String(v.(string)) }` line is consistent with the surrounding code for `db_mode`, `schema_name`, etc.

4. **Read path: set `databasesMap["schema_mode"]` only when `databases.SchemaMode != nil`.**
   - Rationale: Per project rules, `setXX()` must be guarded by a nil check on the response field. The Read function already follows this pattern for every sibling field (e.g., `if databases.DbMode != nil { databasesMap["db_mode"] = databases.DbMode }`).

5. **No `_extension.go` file.**
   - Rationale: The change is a minimal field addition handled entirely within the existing resource file; no external/multi-resource shared logic is introduced.

## Risks / Trade-offs

- [Risk] Users set a `schema_mode` value the cloud API rejects for non-PG/SQLServer links. → Mitigation: The cloud API returns an error surfaced through the existing `RetryError` wrapping in the Update path; no provider-side validation is added to stay consistent with sibling mode fields.
- [Risk] Reading back a `nil` `SchemaMode` could clear user-set state. → Mitigation: The Read path only writes the key when the response field is non-nil, matching the existing pattern; absent values leave state untouched.
- [Trade-off] Not marking `Computed` means a value the API auto-populates won't be read into state unless the API returns it. This is acceptable and consistent with `db_mode`/`table_mode`/`view_mode` siblings.
