## Context

The `tencentcloud_dlc_data_engine` resource (in `tencentcloud/services/dlc/resource_tc_dlc_data_engine.go`) manages a DLC (Data Lake Compute) virtual cluster through its full lifecycle. It currently exposes ~25 schema fields covering engine type, cluster sizing, billing, scheduling, and session templates, but provides no way to bind tags at creation.

The cloud SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125`) already supports tags:
- `CreateDataEngineRequest.Tags` — `[]*TagInfo` (create-only input). `TagInfo` has `TagKey *string` and `TagValue *string`.
- `DataEngineInfo.TagList` — `[]*TagInfo` (returned by `DescribeDataEngine` and `DescribeDataEngines`). The resource `Read` uses `DescribeDataEngine` (singular), whose `Response.DataEngine` is a `*DataEngineInfo`, so `TagList` is available there too.
- `UpdateDataEngineRequest` has **no** `Tags` field — tags cannot be modified after creation.

The existing resource already follows the pattern of treating create-only fields (e.g. `engine_type`, `cluster_type`, `pay_mode`) as immutable: they are listed in `immutableArgs` in `Update`, and several use `ForceNew` / are not wired into the update request. We will follow the same pattern for tags.

## Goals / Non-Goals

**Goals:**
- Expose a `tags` block on `tencentcloud_dlc_data_engine` so users can bind tags at creation time via IaC.
- Keep state in sync by reading `TagList` back during `Read`.
- Stay fully backward compatible (optional field, no change to existing configurations/state).

**Non-Goals:**
- Modifying tags on an existing engine in place — the DLC `UpdateDataEngine` API does not support it, so this is intentionally excluded.
- Adding tags to the DLC data sources (`tencentcloud_dlc_data_engine_network`, session parameters) — out of scope for this change.

## Decisions

### Decision 1: Model `tags` as a TypeList block with `tag_key` / `tag_value` sub-fields

**Choice:** A `schema.TypeList` named `tags`, with an `Elem` `schema.Resource` containing `tag_key` (string, required) and `tag_value` (string, optional).

**Rationale:** This mirrors the cloud API `TagInfo` struct (`TagKey`/`TagValue`) 1:1, matches the naming convention already used in this resource (snake_case schema names mapping to PascalCase SDK fields), and is consistent with how other TencentCloud provider resources represent tag lists. It avoids the special-cased `key=value` map style because DLC tags are a list of `TagInfo` pairs, which maps cleanly to a list block.

**Alternative considered:** A `schema.TypeMap` keyed by tag key. Rejected — the API is a list of `{TagKey, TagValue}` structs (allows repeated keys in the raw API), and a list block preserves ordering and is the established pattern in this codebase.

### Decision 2: Mark `tags` as `ForceNew: true`

**Choice:** The `tags` schema field uses `ForceNew: true`.

**Rationale:** `UpdateDataEngineRequest` has no `Tags` field, so the API cannot update tags. Without `ForceNew`, a change to `tags` would be silently ignored by `Update` (the value would only be sent on the next create). Marking it `ForceNew` makes the behavior explicit: changing tags recreates the engine, matching how the existing code treats other create-only fields (`engine_type`, `cluster_type`, `mode`, etc. are in `immutableArgs`). We will also add `"tags"` to the `immutableArgs` slice in `Update` for consistency with the existing create-only arguments — although `ForceNew` means Terraform plans a recreate before `Update` runs, listing it keeps the guard uniform.

### Decision 3: Populate `request.Tags` from the block in `Create`, and flatten `dataEngine.TagList` in `Read`

**Choice:**
- In `resourceTencentCloudDlcDataEngineCreate`, iterate the `tags` block with `helper.InterfacesInterfaces`/`d.Get("tags")`, build `[]*dlc.TagInfo`, and assign to `request.Tags`.
- In `resourceTencentCloudDlcDataEngineRead`, after obtaining `dataEngine` (a `*DataEngineInfo`), flatten `dataEngine.TagList` into a `[]map[string]interface{}` (with `tag_key`/`tag_value` only when non-nil) and `d.Set("tags", ...)`.

**Rationale:** This follows the exact marshalling/unmarshalling pattern already used in this resource for `data_engine_config_pairs` and `session_resource_template.running_time_parameters` (both are TypeList blocks of `DataEngineConfigPair`). Nil checks on response fields follow the project rule that `setXX()` is only called when the response field is non-nil.

## Risks / Trade-offs

- **[Recreation on tag change]** → Mitigation: `ForceNew` is the only correct behavior given the API limitation; it is documented in the field `Description` and in the resource `.md` so users understand changing tags recreates the engine.
- **[State drift for engines created outside Terraform]** → Mitigation: `Read` flattens `TagList` whenever present and leaves `tags` unset/empty when the API returns no tags (nil), so imported or pre-existing engines refresh cleanly.
- **[Existing buggy `config_value` assignment in `data_engine_config_pairs`]** → Out of scope; we will not alter unrelated existing code. The new `tags` marshalling will correctly set both `TagKey` and `TagValue`.
