## Context

The TencentCloud Terraform Provider already manages several DLC (Data Lake Compute) resources (`tencentcloud_dlc_data_engine`, `tencentcloud_dlc_work_group`, `tencentcloud_dlc_store_location_config`, etc.). However, the **meta database** (`MetaDatabase`) — the top-level catalog object that groups tables and governs data lifecycle policies — has no Terraform resource. Users must create and delete meta databases manually in the console.

The DLC v20210125 SDK exposes three relevant APIs:
- `CreateMetaDatabase` — accepts `DatasourceConnectionName`, `MetaDatabaseInfo` (DatabaseName + Comment), `GovernPolicy` (RuleType + GovernEngine), and `SmartPolicy` (deeply nested BaseInfo + Policy). Returns `BatchId` and `TaskIdSet`, indicating an **async** operation.
- `DescribeDatabase` — accepts `DatabaseName` + `DatasourceConnectionName`, returns `DatabaseResponseInfo` (which includes computed fields like CreateTime, ModifiedTime, Location, UserAlias, UserSubUin, GovernPolicy, DatabaseId, CatalogName, CatalogType, IsInformationSchema, Properties).
- `DeleteMetaDatabase` — accepts `DatabaseName` + `DatasourceConnectionName`, also returns `BatchId` and `TaskIdSet` (async).

There is **no Update / Modify / Alter** API for `MetaDatabase` in the SDK (confirmed by grep of `client.go`). This means the resource is CRD-only.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_dlc_meta_database` that supports create, read, delete, and import.
- Map all `CreateMetaDatabase` input parameters to a Terraform schema, including the deeply nested `SmartPolicy` → `Policy` → `Resources` / `Written` / `Lifecycle` / `Index` / `ChangeTable` / `TableExpiration` blocks.
- Surface all computed fields returned by `DescribeDatabase` so the Terraform state reflects the real cloud state.
- Handle async semantics: after `CreateMetaDatabase` and `DeleteMetaDatabase`, poll `DescribeDatabase` until the resource is confirmed present / gone.
- Enforce immutability: since no Update API exists, any change to a non-ID top-level field must return an error so Terraform triggers a destroy + recreate.

**Non-Goals:**
- Implementing an in-place Update path (the cloud API does not support it).
- Managing DLC tables, data engines, or other DLC resources (out of scope).
- Managing tags on the meta database (the Create API does not accept tags).

## Decisions

### Decision 1: CRD-only with `immutableArgs` enforcement

**Rationale:** The DLC SDK has no `ModifyMetaDatabase` / `UpdateMetaDatabase` / `AlterMetaDatabase` API. The only valid lifecycle operations are create, read, and delete. Rather than omitting the `Update` function entirely (which would cause Terraform to error on any plan with changes), we implement `Update` but immediately check all non-ID top-level fields against an `immutableArgs` list. If any of them changed, we return `fmt.Errorf("argument `%s` cannot be changed", v)`, which Terraform surfaces to the user, telling them to recreate the resource.

**Alternative considered:** Setting `ForceNew: true` on every field. Rejected because it would cause Terraform to silently destroy and recreate the resource on any field change, which is surprising for deeply nested policy blocks. The `immutableArgs` pattern gives the user an explicit error message and matches the established pattern used by `tencentcloud_waf_clb_domain` and other CRD-only resources.

**Implementation:** The `database_name` field is `Required` and `ForceNew: true` (it is part of the resource ID). All other top-level schema fields (`datasource_connection_name`, `comment`, `govern_policy`, `smart_policy`) are `Optional` without `ForceNew`; their immutability is enforced by the `immutableArgs` check in `Update`.

### Decision 2: Composite resource ID

**Rationale:** `DescribeDatabase` and `DeleteMetaDatabase` both require `DatabaseName` and `DatasourceConnectionName`. The `DatasourceConnectionName` is optional (defaults to `DataLakeCatalog`), so the ID must carry it when the user provides it.

**Format:** `datasource_connection_name` + `tccommon.FILED_SP` + `database_name` when `datasource_connection_name` is non-empty; bare `database_name` otherwise.

**Import:** The `.md` doc must document that import requires the composite ID when `datasource_connection_name` is used.

### Decision 3: Async polling via `DescribeDatabase`

**Rationale:** `CreateMetaDatabase` returns `BatchId` and `TaskIdSet`, confirming it is asynchronous. The resource may not be immediately queryable after the create call returns. Similarly, `DeleteMetaDatabase` returns the same fields, meaning the resource may linger briefly before being removed.

**Implementation:** After `CreateMetaDatabase` succeeds, we set `d.SetId(...)` and then call `DescribeDatabase` in a retry loop (using `resource.Retry` with `tccommon.ReadRetryTimeout`) until the database is found or the timeout expires. After `DeleteMetaDatabase` succeeds, we call `DescribeDatabase` in a retry loop until it returns not-found (resource gone) or the timeout expires.

**Note on retry placement:** Per project rules, the `d.SetId(...)` call is placed **outside** the retry block (after the retry returns success), and the retry block only contains the API call itself. The polling loop after create/delete is a separate concern from the retry-wrapped API call.

### Decision 4: Schema structure for nested SmartPolicy

**Rationale:** The `SmartPolicy` parameter is deeply nested: `BaseInfo` (flat), `Policy` → `Inherit`, `Resources` (list), `Written` → `AdvancePolicy` → `SortOrders` (list), `Lifecycle`, `Index`, `ChangeTable`, `TableExpiration`. We model each level as a `schema.Resource` with nested `schema.TypeList` where the cloud API uses arrays.

**Flatten approach:** Per project rules (rule 13), we do NOT wrap the entire resource in a single top-level `smart_policy_set` wrapper. Instead, `smart_policy` is a single-element `schema.TypeList` (or a block) whose element is a `schema.Resource` with `base_info`, `policy` sub-blocks. The `Resources` field inside `Policy` is a `schema.TypeList` because the cloud API accepts an array. `SortOrders` inside `AdvancePolicy` is similarly a `schema.TypeList`.

### Decision 5: Read uses `DatabaseResponseInfo` not `DatabaseInfo`

**Rationale:** `DescribeDatabase` returns `response.Response.DatabaseInfo` of type `*DatabaseResponseInfo` (not `*DatabaseInfo`). `DatabaseResponseInfo` includes additional computed fields: `CreateTime`, `ModifiedTime`, `UserAlias`, `UserSubUin`, `GovernPolicy`, `DatabaseId`, `CatalogName`, `CatalogType`, `IsInformationSchema`, `Properties`. These are surfaced as computed schema fields so the Terraform state is complete.

### Decision 6: Unit test with gomonkey mock

**Rationale:** Per project rules (rule for new resources), the test file must use gomonkey mock to stub the cloud API calls rather than the terraform test suite. This allows business-logic unit testing without real cloud credentials.

## Risks / Trade-offs

- **[Risk] Async create may exceed default retry timeout** → Mitigation: use `tccommon.ReadRetryTimeout` for polling and add a `Timeouts` block to the schema so users can override the default.
- **[Risk] Deeply nested schema increases complexity** → Mitigation: follow the established pattern from `tencentcloud_igtm_strategy` (multi-level `schema.Resource` nesting). Keep build/flatten helpers focused per block.
- **[Risk] User changes an immutable field and gets confused** → Mitigation: the error message explicitly names the field and states it cannot be changed; the `.md` doc notes that the resource is CRD-only.
- **[Trade-off] No in-place update** → The resource must be destroyed and recreated when any non-ID field changes. This is inherent to the cloud API design and cannot be avoided.
- **[Risk] `datasource_connection_name` defaults to `DataLakeCatalog`** → Mitigation: when the user omits it, the ID is just `database_name`; when present, the composite ID preserves it so read/delete can reconstruct the request correctly.
