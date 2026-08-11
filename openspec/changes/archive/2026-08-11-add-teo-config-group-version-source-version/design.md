## Context

The `tencentcloud_teo_config_group_version_detail` data source reads version detail for an EdgeOne (TEO) config group via the `DescribeConfigGroupVersionDetail` API. The SDK struct `ConfigGroupVersionInfo` (in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`) already defines a `SourceVersion *string` field — the source version ID that the current version was derived from when created. However, the Terraform schema's `config_group_version_info` block omits this field, so the value is never surfaced to users.

**Current state:**
- Data source file: `tencentcloud/services/teo/data_source_tc_teo_config_group_version_detail.go`
- The `config_group_version_info` block schema currently defines: `version_id`, `version_number`, `group_id`, `group_type`, `description`, `status`, `create_time`.
- The Read function maps fields from `respData.ConfigGroupVersionInfo` into a map and sets the `config_group_version_info` list, performing nil-checks before each `set`.

**API behavior analysis:**

| API | SourceVersion in Request | SourceVersion in Response |
|-----|--------------------------|---------------------------|
| `DescribeConfigGroupVersionDetail` | N/A (request takes `ZoneId`, `VersionId`) | Yes (`ConfigGroupVersionInfo.SourceVersion *string`) |

`SourceVersion` is a read-only output field; no request parameter changes are needed.

## Goals / Non-Goals

**Goals:**
- Add a `source_version` computed field to the `config_group_version_info` block of the `tencentcloud_teo_config_group_version_detail` data source.
- Map `source_version` from `respData.ConfigGroupVersionInfo.SourceVersion` in the Read function, with a nil-check consistent with the existing pattern.
- Update the data source `.md` documentation to describe the new field.

**Non-Goals:**
- Adding `SourceVersion` to any other teo resource or data source.
- Changing the API request or adding new API calls.
- Modifying the `config_group_version_info` block's existing fields or their behavior.

## Decisions

### Decision 1: Add `source_version` as a Computed field inside the `config_group_version_info` block

**Rationale:** `SourceVersion` is an output-only field of the `DescribeConfigGroupVersionDetail` response, nested under `ConfigGroupVersionInfo`. The existing schema groups version info fields under a `config_group_version_info` list block; placing `source_version` there keeps the structure consistent with the API response and the existing field mapping. Marking it `Computed: true` (with `Optional: true` to match sibling fields' pattern) reflects that it is populated by the API, not user-supplied.

### Decision 2: Follow the existing nil-check-then-set pattern in Read

**Rationale:** The current Read function checks each field for nil before adding it to the `configGroupVersionInfoMap`. The new `source_version` mapping will follow the same pattern: `if respData.ConfigGroupVersionInfo.SourceVersion != nil { configGroupVersionInfoMap["source_version"] = respData.ConfigGroupVersionInfo.SourceVersion }`. This avoids setting nil pointers into state and matches surrounding code conventions.

### Decision 3: No SDK update required

**Rationale:** The vendored SDK already contains `ConfigGroupVersionInfo.SourceVersion` (models.go line 2034). No `go mod vendor` or dependency bump is needed.

## Risks / Trade-offs

- **[Risk] Field absent in some API responses**: `SourceVersion` may be empty/nil for the initial version of a config group (which has no source).
  - **Mitigation:** The nil-check before set ensures the field is simply omitted from the map when nil, consistent with all other sibling fields. No error is raised.
- **[Risk] Backward compatibility**: Adding a computed field to an existing block.
  - **Mitigation:** Computed-only additive fields are backward compatible; existing configurations and state continue to work unchanged.
