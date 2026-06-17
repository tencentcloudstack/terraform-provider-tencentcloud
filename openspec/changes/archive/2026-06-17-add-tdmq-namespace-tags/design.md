## Context

The `tencentcloud_tdmq_namespace` resource manages TDMQ (Pulsar) namespaces via the `CreateEnvironment`, `DescribeEnvironments`, `ModifyEnvironmentAttributes`, and `DeleteEnvironments` APIs. Currently, the resource supports `environ_name`, `msg_ttl`, `cluster_id`, `remark`, and `retention_policy` parameters, but does not support the `Tags` parameter that was added to the `CreateEnvironment` API.

The TDMQ `CreateEnvironment` API accepts a `Tags` field (type `[]*Tag` where `Tag` has `TagKey` and `TagValue` string fields). The `DescribeEnvironments` API returns tags in the `Environment.Tags` field. However, the `ModifyEnvironmentAttributes` API does **not** support updating tags, meaning tags can only be set during creation and are immutable afterward.

## Goals / Non-Goals

**Goals:**
- Add `Tags` parameter to `tencentcloud_tdmq_namespace` resource schema as Optional + ForceNew
- Pass Tags to the `CreateEnvironment` API during resource creation
- Read Tags from the `DescribeEnvironments` API response during resource read
- Add Tags to the immutable args list in the update method to provide a clear error message if users attempt to change tags
- Update documentation and add unit tests for the new parameter

**Non-Goals:**
- Supporting tag updates via ModifyEnvironmentAttributes (the API does not support it)
- Changing existing schema fields or behavior
- Adding Tags support to other TDMQ resources (out of scope for this change)

## Decisions

### Decision 1: Tags schema type — TypeList with tag_key/tag_value sub-fields

**Choice**: Use `schema.TypeList` with nested `tag_key` and `tag_value` string fields.

**Rationale**: The `CreateEnvironment` API's `Tags` field is `[]*Tag` with explicit `TagKey`/`TagValue` struct fields (not a simple map). This matches the pattern used by other resources in the codebase that deal with `[]*Tag` style API parameters (e.g., `cdwdoris_instance`). The `TypeMap` approach used by `tdmq_instance` and `tdmq_professional_cluster` works for simpler tag APIs but doesn't align with the structured Tag type in this API.

### Decision 2: Tags as ForceNew parameter

**Choice**: Mark Tags as `ForceNew: true` since `ModifyEnvironmentAttributes` does not support Tags.

**Rationale**: Since the update API cannot modify tags, any change to tags would require recreating the resource. ForceNew ensures Terraform handles this correctly. Additionally, adding "tags" to the `immutableArgs` list in the update method provides a clear error message when users attempt to change tags.

### Decision 3: Service layer function signature update

**Choice**: Add a `tags []*tdmq.Tag` parameter to the `CreateTdmqNamespace` function.

**Rationale**: This is the minimal change needed to pass tags through to the API request. The alternative would be to refactor the function to accept a generic map, but that would be a larger change with no additional benefit.

## Risks / Trade-offs

- [ForceNew triggers resource recreation on tag changes] → Mitigation: This is the correct behavior since the API doesn't support tag updates. The immutableArgs check in the update method provides an early error message.
- [Tags returned by DescribeEnvironments may be nil] → Mitigation: Check for nil before flattening Tags into the schema, consistent with the nil-checking pattern used for other fields in the read method.
