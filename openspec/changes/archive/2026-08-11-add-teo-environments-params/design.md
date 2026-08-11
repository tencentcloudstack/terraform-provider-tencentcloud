## Context

The `tencentcloud_teo_environments` data source (`tencentcloud/services/teo/data_source_tc_teo_environments.go`) reads environment data from the cloud API `DescribeEnvironments` (package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`).

Current state:
- The data source schema exposes `zone_id` (Required) and `env_infos` (Computed list) plus `result_output_file`.
- `env_infos.current_config_group_version_infos` currently exposes `version_id`, `version_number`, `group_id`, `group_type`, `description`, `status`, `create_time`.
- The service layer method `TeoService.DescribeTeoEnvironmentsByFilter(ctx, param)` currently returns only `[]*teov20220901.EnvInfo` (i.e. `response.Response.EnvInfos`). It does NOT return the top-level `TotalCount` field present on `DescribeEnvironmentsResponseParams`.

Cloud API confirmation (from vendored SDK `models.go`):
- `DescribeEnvironmentsResponseParams.TotalCount` is `*uint64` → maps to schema field `total_count` (int).
- `ConfigGroupVersionInfo.SourceVersion` is `*string`, nested under `EnvInfo.CurrentConfigGroupVersionInfos` → maps to schema field `source_version` (string) under `env_infos.current_config_group_version_infos`.

Both fields already exist in the vendored SDK, so no vendor/dependency update is required.

## Goals / Non-Goals

**Goals:**
- Surface `total_count` (from `response.TotalCount`) as a new top-level computed field on the data source.
- Surface `source_version` (from `response.EnvInfos[].CurrentConfigGroupVersionInfos[].SourceVersion`) as a new nested computed field under `env_infos.current_config_group_version_infos`.
- Keep the change fully backward compatible: existing schemas, configurations, and state remain valid.
- Provide unit tests (gomonkey mock approach) and updated documentation.

**Non-Goals:**
- Do not change the `zone_id` input parameter or the data source query behavior.
- Do not introduce pagination parameters (`limit`/`offset`) into the schema; pagination handling stays internal if any exists.
- Do not modify any other teo resources or data sources.
- Do not update the vendored SDK dependency.

## Decisions

### Decision 1: Modify service layer to also return TotalCount
**Choice**: Change `TeoService.DescribeTeoEnvironmentsByFilter` return signature so that `TotalCount` can be surfaced, rather than reading `TotalCount` inside the data source by re-calling the API.

**Rationale**: The service layer is the single place that calls `DescribeEnvironments`. The data source already depends on this method. To expose `TotalCount`, the cleanest approach is to have the service layer return it alongside the env list, so the data source can set it from the same API response.

**Implementation detail**: Rather than change the existing method's return type (which could ripple through callers), we will adjust the method to return the total count as an additional return value. We will verify there is only one caller (the data source Read function) — confirmed via grep. The method will return `(ret []*teov20220901.EnvInfo, totalCount *uint64, errRet error)`. When `response.Response.TotalCount` is non-nil it is returned; otherwise `nil`.

**Alternative considered**: Re-implement the API call directly in the data source Read to access `TotalCount`. Rejected because it duplicates the API call logic and breaks the established service-layer abstraction pattern used throughout this provider.

### Decision 2: Field types and placement
**Choice**:
- `total_count`: top-level computed field, `schema.TypeInt`, mapped from `response.TotalCount` (`*uint64`). Conversion: dereference the `*uint64` to `int` when non-nil.
- `source_version`: nested computed field under `env_infos.current_config_group_version_infos`, `schema.TypeString`, mapped from `CurrentConfigGroupVersionInfos[].SourceVersion` (`*string`). Set only when non-nil (consistent with the existing nil-check pattern in the Read function).

**Rationale**: Matches the cloud API types and the existing code style where each field is guarded by a nil check before `d.Set`/map assignment.

### Decision 3: Backward compatibility
**Choice**: Both new fields are `Computed` (output-only). No existing field is removed or renamed. No `Required`/`Optional` input field is added.

**Rationale**: Computed additions never break existing Terraform configurations or state files.

## Risks / Trade-offs

- **[Risk] Service layer signature change affects other callers** → Mitigation: grep confirmed the only caller is `dataSourceTencentCloudTeoEnvironmentsRead`. The change is localized.
- **[Risk] `TotalCount` may be nil/zero in some API responses** → Mitigation: guard with a nil check before setting `total_count`; if nil, leave the field unset (zero value), consistent with the provider's established nil-handling pattern.
- **[Risk] `SourceVersion` may be empty for some config group versions** → Mitigation: guard with a nil check; only set when the API returns a non-nil value.
- **[Trade-off]** We modify an existing service-layer method rather than add a new one, to avoid code duplication. This is acceptable given the single caller.
