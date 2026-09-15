## Context

The BDRC (Backup and Disaster Recovery Center) product provides disaster recovery site pairs (`容灾站点对`) that establish replication relationships between a production site and a recovery site across regions, zones, or clouds. The cloud APIs involved are already vendored under `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330`. There is currently no BDRC service package in `tencentcloud/services/` and no BDRC client registered in `connectivity/client.go`.

The resource follows the RESOURCE_KIND_GENERAL pattern — a standalone resource with full CRUD. The reference implementation style is `tencentcloud_igtm_strategy`.

## Goals / Non-Goals

**Goals:**
- Provide a `tencentcloud_bdrc_disaster_recovery_site_pair` resource that manages the full lifecycle of a BDRC site pair.
- Create the BDRC service package (`tencentcloud/services/bdrc/`) with a service-layer wrapper for the read API.
- Register the BDRC SDK client (`UseBdrcV20260330Client`) in `connectivity/client.go`.
- Register the resource in `provider.go` and `provider.md`.
- Deliver resource documentation and gomonkey-based unit tests.

**Non-Goals:**
- No BDRC data source (`RESOURCE_KIND_DATASOURCE`) in this change.
- No management of copy pairs, protection groups, or VPC mappings — only the site pair resource.
- No acceptance tests (TF_ACC); only gomonkey unit tests are included per project requirements.

## Decisions

### 1. Resource ID — single `SitePairId` (string), not a composite ID

The create API returns a single `SitePairId` (string). The read API (`DescribeDisasterRecoverySitePairs`) accepts `SitePairIds` (a list) plus a required `SitePairType`. Since `SitePairType` can be derived from the `site_pair_product_type` field already stored in state, there is no need to encode multiple values into the Terraform ID.

- **Decision**: Use `SitePairId` directly as `d.Id()`. No `tccommon.FILED_SP` separator needed.
- **Rationale**: Simpler import and ID handling; the read API can be called with just the `SitePairId` plus `SitePairType` from state.

### 2. Immutable fields use `ForceNew`; only `site_pair_name` is updatable

The `ModifySitePairAttribute` API only accepts `SitePairId` and `SitePairName`. All other create-time parameters (`disaster_recovery_type`, `source_region`, `source_zone`, `target_region`, `target_zone`, `source_vpc`, `target_vpc`, `site_pair_product_type`, `copy_type`) are immutable.

- **Decision**: Mark all create parameters except `site_pair_name` as `Required: true, ForceNew: true`. Mark `site_pair_name` and `copy_type` as `Optional`.
  - `site_pair_name` → updatable via `ModifySitePairAttribute`
  - `copy_type` → `Optional, ForceNew: true` (not in the modify API)
- **Rationale**: Matches the SDK contract exactly — only `SitePairName` can be updated after creation.

### 3. Read flow uses `DescribeDisasterRecoverySitePairs` with `SitePairIds` filter

There is no single-resource read API. The list API returns `SitePairSet` (a `[]*SitePair`). The service layer will call `DescribeDisasterRecoverySitePairs` with `SitePairIds` set to the resource ID and `SitePairType` set from state, then return the first matching element.

- **Decision**: Service layer method `DescribeDisasterRecoverySitePairById(ctx, sitePairId, sitePairType)` wraps the list API, sets `Limit` to the API maximum (100), and returns the first `SitePair` whose `SitePairId` matches.
- **Rationale**: Avoids a separate describe-by-id endpoint; uses the list API efficiently with an ID filter.

### 4. All output fields from `SitePair` are computed top-level schema fields

Per the project rule for Describe interfaces (rule 13), list data must be expanded to top-level fields — no nested `xxx_set` wrapper around all fields. The `SitePair` response contains nested arrays/objects (`ProtectedResourceSet`, `ProtectedResourceStatusSet`, `CrossCloudDetails`, `ErrorRecoveryPointObjectiveCopyPairSet`), each of which becomes its own typed schema field at the top level.

- **Decision**: Define nested block schemas for `protected_resource_set`, `protected_resource_status_set`, and `cross_cloud_details` as `TypeList`/`TypeList`/`TypeList` of objects at the top resource schema level. All are `Computed` (read-only).
- **Rationale**: Allows terraform to individually set/read each attribute; complies with the expansion rule.

### 5. Update method checks `immutableArgs` for safety

Since most fields are `ForceNew`, Terraform will not call Update for them. However, per project rule 7 for CRD-style robustness, the update method will still verify that only `site_pair_name` changed; any other change in the `mutableArgs` check triggers `ModifySitePairAttribute`.

- **Decision**: The Update method calls `ModifySitePairAttribute` when `site_pair_name` has changed, passing `SitePairId` from `d.Id()`.
- **Rationale**: Aligns with the only mutable field in the modify API.

### 6. Retry and error handling follows project conventions

- All cloud API calls wrapped in `resource.Retry(tccommon.ReadRetryTimeout/WriteRetryTimeout, ...)` with `tccommon.RetryError(e)` on failure.
- Create checks that `response.Response.SitePairId` is non-nil/non-empty and returns `NonRetryableError` if empty.
- Read checks for empty results: logs `[CRUD]` with `d.Id()` before `d.SetId("")`.
- Delete passes `SitePairIds` as a single-element list.

### 7. SDK client registration in connectivity

Add `bdrcv20260330` import alias, a `bdrcv20260330Conn` field on `TencentCloudClient`, and a `UseBdrcV20260330Client()` method — mirroring the `UseIgtmV20231024Client` pattern.

## Risks / Trade-offs

- **[Read API requires `SitePairType`]** → `SitePairType` is stored as a computed field in state from the create response / read response; the resource always has it available after creation. For import, the user supplies the `SitePairId`; the read method reads `site_pair_type` from `d.Get()` — if state is empty during import, the read will attempt the list API without a type filter or fall back gracefully.
- **[Nested response fields may be nil]** → All `d.Set()` calls are guarded by nil checks before setting, per project rule 8.
- **[`copy_type` is Optional but immutable]** → Marked as `ForceNew` so changes recreate the resource rather than silently failing the update.
