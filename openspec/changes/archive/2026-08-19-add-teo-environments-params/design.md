## Context

The `tencentcloud_teo_environments` data source (`tencentcloud/services/teo/data_source_tc_teo_environments.go`) reads environment data from the cloud API `DescribeEnvironments` (package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`).

Current state:
- The data source schema exposes `zone_id` (Required) and `env_infos` (Computed list) plus `result_output_file`.
- `env_infos.current_config_group_version_infos` currently exposes `version_id`, `version_number`, `group_id`, `group_type`, `description`, `status`, `create_time`.
- The service layer method `TeoService.DescribeTeoEnvironmentsByFilter(ctx, param)` returns `[]*teov20220901.EnvInfo` (i.e. `response.Response.EnvInfos`), each of which already carries `CurrentConfigGroupVersionInfos` populated by the cloud API.

Cloud API confirmation (from vendored SDK `models.go`):
- `ConfigGroupVersionInfo.SourceVersion` is `*string`, nested under `EnvInfo.CurrentConfigGroupVersionInfos` → maps to schema field `source_version` (string) under `env_infos.current_config_group_version_infos`.

The field already exists in the vendored SDK, so no vendor/dependency update is required.

## Goals / Non-Goals

**Goals:**
- Surface `source_version` (from `response.EnvInfos[].CurrentConfigGroupVersionInfos[].SourceVersion`) as a new nested computed field under `env_infos.current_config_group_version_infos`.
- Keep the change fully backward compatible: existing schemas, configurations, and state remain valid.
- Provide unit tests and updated documentation.

**Non-Goals:**
- Do not change the `zone_id` input parameter or the data source query behavior.
- Do not introduce pagination parameters (`limit`/`offset`) into the schema; pagination handling stays internal if any exists.
- Do not modify any other teo resources or data sources.
- Do not update the vendored SDK dependency.

## Decisions

### Decision 1: No service layer signature change
**Choice**: Do not modify `TeoService.DescribeTeoEnvironmentsByFilter`'s return signature. The `source_version` field is already reachable via each `EnvInfo.CurrentConfigGroupVersionInfos[].SourceVersion` element returned by the existing method.

**Rationale**: The service layer already returns the full `EnvInfo` list including the nested `CurrentConfigGroupVersionInfos`. The data source Read function iterates over these elements and can populate `source_version` directly without any change to the service method signature.

**Alternative considered**: Re-implement the API call directly in the data source Read to access the nested fields. Rejected because it duplicates the API call logic and breaks the established service-layer abstraction pattern used throughout this provider.

### Decision 2: Field type and placement
**Choice**:
- `source_version`: nested computed field under `env_infos.current_config_group_version_infos`, `schema.TypeString`, mapped from `CurrentConfigGroupVersionInfos[].SourceVersion` (`*string`). Set only when non-nil (consistent with the existing nil-check pattern in the Read function).

**Rationale**: Matches the cloud API type and the existing code style where each field is guarded by a nil check before map assignment.

### Decision 3: Backward compatibility
**Choice**: The new field is `Computed` (output-only). No existing field is removed or renamed. No `Required`/`Optional` input field is added.

**Rationale**: Computed additions never break existing Terraform configurations or state files.

## Risks / Trade-offs

- **[Risk] `SourceVersion` may be empty for some config group versions** → Mitigation: guard with a nil check; only set when the API returns a non-nil value.
- **[Trade-off]** None. No service-layer signature change is required, so there is no ripple effect on other callers.
