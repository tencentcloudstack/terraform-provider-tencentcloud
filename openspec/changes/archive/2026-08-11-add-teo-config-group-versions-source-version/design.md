## Context

The `tencentcloud_teo_config_group_versions` data source (RESOURCE_KIND_DATASOURCE) queries EdgeOne configuration group versions via the `DescribeConfigGroupVersions` API (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`). The data source Read function calls `service.DescribeTeoConfigGroupVersionsByFilter`, which returns `[]*teov20220901.ConfigGroupVersionInfo`. Each element is flattened into the `config_group_version_infos` nested schema block.

Currently the nested schema maps: `version_id`, `version_number`, `group_id`, `group_type`, `description`, `status`, `create_time`. The SDK struct `ConfigGroupVersionInfo` additionally exposes `SourceVersion *string`, which represents the source version ID a version was derived from (format `ver-xxxxxxxx`). This field is not yet surfaced in the Terraform schema.

## Goals / Non-Goals

**Goals:**
- Expose the `SourceVersion` field from the `DescribeConfigGroupVersions` API response as a computed `source_version` attribute within the `config_group_version_infos` block.
- Maintain full backward compatibility with existing schema fields and Terraform configurations.

**Non-Goals:**
- No changes to input parameters (`zone_id`, `group_id`, `filters`) of the data source.
- No changes to the query/filter logic, pagination, or retry behavior.
- No SDK dependency update; the vendor already contains the `SourceVersion` field.

## Decisions

**Decision 1: Add `source_version` as a computed (Optional) field inside the `config_group_version_infos` nested schema.**
- Rationale: The field is an output (response-only) attribute of a datasource, so it must be Computed. Following the existing convention of sibling fields in this block (e.g., `version_number`, `description`) which use `Optional: true, Computed: true` implied via being inside a Computed list, the new field mirrors the existing `version_number` style: `Type: schema.TypeString, Optional: true`.
- Alternatives considered: Making it a top-level field — rejected because `SourceVersion` is per-version element data, not a global query attribute; it belongs inside the per-element nested block alongside `version_id` etc.

**Decision 2: Map the field with a nil-check guard in the Read loop.**
- Rationale: Per project conventions, `setXX()` must only be called when the response field is non-nil. The Read loop already wraps every field assignment in a nil check (e.g., `if configGroupVersionInfos.VersionId != nil`). The new `SourceVersion` follows the same pattern: `if configGroupVersionInfos.SourceVersion != nil { configGroupVersionInfosMap["source_version"] = configGroupVersionInfos.SourceVersion }`.
- Alternatives considered: Unconditional assignment — rejected to avoid nil pointer issues and to match existing code style.

**Decision 3: Update the `.md` documentation file directly; the `website/docs` file is regenerated via `make doc` during finalization.**
- Rationale: Project rules forbid editing `website/` directly; documentation flows through the service-level `.md` file and `make doc`.

## Risks / Trade-offs

- [Risk: API returns `SourceVersion` as nil for versions without a source] → Mitigation: The nil-check guard ensures the field is simply omitted from the map when nil, consistent with all other fields in the block. No error is raised.
- [Risk: Schema change could affect existing state] → Mitigation: The change is purely additive (new computed attribute inside an already-computed nested block). Terraform handles unknown computed attributes gracefully; existing state and configurations are unaffected.
