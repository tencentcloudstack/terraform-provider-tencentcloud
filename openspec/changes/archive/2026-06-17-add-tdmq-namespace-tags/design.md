## Context

The `tencentcloud_tdmq_namespace` resource manages TDMQ (Pulsar) namespaces via the `CreateEnvironment`, `DescribeEnvironments`, `ModifyEnvironmentAttributes`, and `DeleteEnvironments` APIs. Currently, the resource supports `environ_name`, `msg_ttl`, `cluster_id`, `remark`, and `retention_policy` parameters, but does not support the `Tags` parameter that was added to the `CreateEnvironment` API.

The TDMQ `CreateEnvironment` API accepts a `Tags` field (type `[]*Tag` where `Tag` has `TagKey` and `TagValue` string fields). The `DescribeEnvironments` API returns tags in the `Environment.Tags` field. Although the `ModifyEnvironmentAttributes` API does **not** support updating tags directly, tags can be managed via the generic tag service APIs (`TagResources`/`UnTagResources`) through the `ModifyTags` helper in the tag service layer.

## Goals / Non-Goals

**Goals:**
- Add `Tags` parameter to `tencentcloud_tdmq_namespace` resource schema as Optional (mutable)
- Pass Tags to the `CreateEnvironment` API during resource creation
- Read Tags from the `DescribeEnvironments` API response during resource read
- In the Update method, detect tag changes and use `svctag.ModifyTags` (which calls `TagResources`/`UnTagResources`) to apply tag additions, modifications, and deletions
- Update documentation and add unit tests for the new parameter

**Non-Goals:**
- Changing existing schema fields or behavior
- Adding Tags support to other TDMQ resources (out of scope for this change)

## Decisions

### Decision 1: Tags schema type — TypeList with tag_key/tag_value sub-fields

**Choice**: Use `schema.TypeList` with nested `tag_key` and `tag_value` string fields.

**Rationale**: The `CreateEnvironment` API's `Tags` field is `[]*Tag` with explicit `TagKey`/`TagValue` struct fields (not a simple map). This matches the pattern used by other resources in the codebase that deal with `[]*Tag` style API parameters (e.g., `cdwdoris_instance`). The `TypeMap` approach used by `tdmq_instance` and `tdmq_professional_cluster` works for simpler tag APIs but doesn't align with the structured Tag type in this API.

### Decision 2: Tags as mutable parameter with update via TagResources/UnTagResources

**Choice**: Tags are Optional and mutable. Tag updates are handled via the tag service's `ModifyTags` function which internally calls `TagResources`/`UnTagResources` APIs.

**Rationale**: The generic tag service APIs support modifying tags on any resource. In the Update method, the TypeList tags are converted to `map[string]interface{}` format, then `svctag.DiffTags` computes the replacements and deletions, and `tagService.ModifyTags` applies them using the resource name `tdmq/environment/{region}/{environId}`.

### Decision 3: Service layer function signature update

**Choice**: Add a `tags []*tdmq.Tag` parameter to the `CreateTdmqNamespace` function.

**Rationale**: This is the minimal change needed to pass tags through to the API request. The alternative would be to refactor the function to accept a generic map, but that would be a larger change with no additional benefit.

## Risks / Trade-offs

- [Tags returned by DescribeEnvironments may be nil] → Mitigation: Check for nil before flattening Tags into the schema, consistent with the nil-checking pattern used for other fields in the read method.
- [TypeList to map conversion in Update] → Mitigation: The conversion logic iterates over the TypeList items and builds a `map[string]interface{}` keyed by `tag_key` with `tag_value` as the value, which is the format expected by `svctag.DiffTags`.
