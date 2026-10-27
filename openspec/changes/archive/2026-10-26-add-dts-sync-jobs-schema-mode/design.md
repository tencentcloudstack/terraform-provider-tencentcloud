## Context

The `tencentcloud_dts_sync_jobs` data source (`RESOURCE_KIND_DATASOURCE`) reads sync job details from the DTS cloud API `DescribeSyncJobs`. It flattens each returned `SyncJobInfo` into the `list` computed block. Inside `list[].objects[].databases[]`, the data source currently exposes several fields from the cloud API `Database` struct (e.g. `db_name`, `new_db_name`, `db_mode`, `schema_name`, `new_schema_name`, `table_mode`, `view_mode`, `function_mode`, `procedure_mode`, etc.) but does **not** expose `SchemaMode`. The cloud API `Database.SchemaMode` field already exists in the vendored SDK, so this is a pure provider-side additive change.

### Cloud API reference (vendored, `dts/v20211206/models.go`)

The full response path for the new field is:
`DescribeSyncJobsResponse.Response.JobList []*SyncJobInfo` → `SyncJobInfo.Objects *Objects` → `Objects.Databases []*Database` → `Database.SchemaMode *string`.

Key struct excerpts (from `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206/models.go`):

```go
// line 4236
type DescribeSyncJobsResponseParams struct {
	TotalCount *uint64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
	JobList []*SyncJobInfo `json:"JobList,omitnil,omitempty" name:"JobList"`
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

// line 4247
type DescribeSyncJobsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeSyncJobsResponseParams `json:"Response"`
}

// line 8401 (inside SyncJobInfo)
	Objects *Objects `json:"Objects,omitnil,omitempty" name:"Objects"`

// line 6475
type Objects struct {
	Mode *string `json:"Mode,omitnil,omitempty" name:"Mode"`
	Databases []*Database `json:"Databases,omitnil,omitempty" name:"Databases"`
	AdvancedObjects []*string `json:"AdvancedObjects,omitnil,omitempty" name:"AdvancedObjects"`
	OnlineDDL *OnlineDDL `json:"OnlineDDL,omitnil,omitempty" name:"OnlineDDL"`
	DatabasesOpFilter []*DBOpFilter `json:"DatabasesOpFilter,omitnil,omitempty" name:"DatabasesOpFilter"`
}

// line 1842
type Database struct {
	DbName *string `json:"DbName,omitnil,omitempty" name:"DbName"`
	NewDbName *string `json:"NewDbName,omitnil,omitempty" name:"NewDbName"`
	DbMode *string `json:"DbMode,omitnil,omitempty" name:"DbMode"`
	SchemaName *string `json:"SchemaName,omitnil,omitempty" name:"SchemaName"`
	NewSchemaName *string `json:"NewSchemaName,omitnil,omitempty" name:"NewSchemaName"`
	// schema选择模式，pg和sqlserver需要使用。枚举值：All（当前对象下的所有对象）、Partial（部分对象）
	SchemaMode *string `json:"SchemaMode,omitnil,omitempty" name:"SchemaMode"`
	TableMode *string `json:"TableMode,omitnil,omitempty" name:"TableMode"`
	Tables []*Table `json:"Tables,omitnil,omitempty" name:"Tables"`
	ViewMode *string `json:"ViewMode,omitnil,omitempty" name:"ViewMode"`
	Views []*View `json:"Views,omitnil,omitempty" name:"Views"`
	FunctionMode *string `json:"FunctionMode,omitnil,omitempty" name:"FunctionMode"`
	Functions []*string `json:"Functions,omitnil,omitempty" name:"Functions"`
	ProcedureMode *string `json:"ProcedureMode,omitnil,omitempty" name:"ProcedureMode"`
	Procedures []*string `json:"Procedures,omitnil,omitempty" name:"Procedures"`
	TriggerMode *string `json:"TriggerMode,omitnil,omitempty" name:"TriggerMode"`
	Triggers []*string `json:"Triggers,omitnil,omitempty" name:"Triggers"`
	EventMode *string `json:"EventMode,omitnil,omitempty" name:"EventMode"`
	Events []*string `json:"Events,omitnil,omitempty" name:"Events"`
}
```

`SchemaMode` is a `*string` with `omitnil`, so it may be nil in the API response.

### Existing schema structure (target resource)

File: `tencentcloud/services/dts/data_source_tc_dts_sync_jobs.go`

The `list` block → `objects` block → `databases` block currently defines these fields (TypeList of `schema.Resource`):
- `db_name` (TypeString, Computed)
- `new_db_name` (TypeString, Computed)
- `db_mode` (TypeString, Computed)
- `schema_name` (TypeString, Computed)
- `new_schema_name` (TypeString, Computed)
- `table_mode` (TypeString, Computed)
- `tables` (TypeList, Computed)
- `view_mode` (TypeString, Computed)
- `views` (TypeList, Computed)
- `function_mode` (TypeString, Computed)
- `functions` (TypeSet, Computed)
- `procedure_mode` (TypeString, Computed)
- `procedures` (TypeSet, Computed)

There is **no** `schema_mode` field currently. The new field will be inserted among the sibling string fields (logical placement near `new_schema_name` / `table_mode`).

### Existing Read handler flattening (target resource)

In the Read handler, the `databases[]` flattening loop builds `databasesMap` and sets each field with a nil-guard (`if databases.Xxx != nil { databasesMap["xxx"] = databases.Xxx }`). The new `schema_mode` field follows the same nil-guard pattern:

```go
if databases.SchemaMode != nil {
	databasesMap["schema_mode"] = databases.SchemaMode
}
```

## Goals / Non-Goals

**Goals:**
- Expose the cloud API `Database.SchemaMode` value through the `tencentcloud_dts_sync_jobs` data source as a computed `schema_mode` field inside `list.objects.databases[]`.
- Maintain full backwards compatibility (additive computed field only).

**Non-Goals:**
- Do not modify the `tencentcloud_dts_sync_config` resource (a separate resource already has its own `schema_mode` capability — `dts-sync-config-schema-mode`).
- Do not change the `DescribeSyncJobs` request parameters or add new query filters.
- Do not modify any other nested block structure beyond adding the single `schema_mode` field.

## Decisions

### Decision 1: Add `schema_mode` as a Computed-only field inside the existing `databases[]` nested schema

**Rationale**: This is a `RESOURCE_KIND_DATASOURCE`, so all `list.*` fields are computed (read-only). The field maps directly to the cloud API `Database.SchemaMode` response field. No request-side change is needed because `DescribeSyncJobs` already returns `SchemaMode` in each `Database` entry.

**Alternatives considered**: None — this is the only location consistent with the cloud API response shape and the existing flattening structure.

### Decision 2: Follow the existing nil-guard pattern in the Read flattening logic

**Rationale**: The existing code sets each `databasesMap` entry only when the source pointer is non-nil (`if databases.SchemaName != nil { databasesMap["schema_name"] = databases.SchemaName }`). The new field uses the identical pattern to avoid writing empty/zero values into state and to stay consistent with the codebase convention for datasource flattening.

### Decision 3: Use mock-based (gomonkey) unit test rather than Terraform acceptance test

**Rationale**: Per project conventions, newly added datasource tests use gomonkey to mock the cloud API and test only business logic, avoiding the need for real cloud credentials.

## Risks / Trade-offs

- **[Risk] Field returned as nil for non-PG/SQLServer links** → Mitigation: The nil-guard in the flattening logic skips setting the field when nil, so links that do not populate `SchemaMode` simply leave the field unset in state. This is consistent with how all sibling fields behave.
- **[Risk] Existing state compatibility** → Mitigation: Adding a computed field is backwards compatible; Terraform will populate it on the next refresh without requiring configuration changes.
