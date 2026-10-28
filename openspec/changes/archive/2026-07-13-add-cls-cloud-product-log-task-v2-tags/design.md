## Context

The `tencentcloud_cls_cloud_product_log_task_v2` resource manages CLS cloud product log collection tasks. Its CRUD lifecycle is implemented in `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go`, backed by these cloud APIs from the `cls/v20201016` SDK package:

- Create: `CreateCloudProductLogCollection`
- Read: `DescribeCloudProductLogTasks` (via service method `DescribeClsCloudProductLogTaskById`), plus `DescribeClsLogset` (`DescribeLogsets`) and `DescribeClsTopicById` (`DescribeTopics`) to fetch logset/topic names.
- Update: `ModifyCloudProductLogCollection`
- Delete: `DeleteCloudProductLogCollection`, optionally `DeleteTopic` and `DeleteLogset` when `force_delete` is true.

The resource already exposes `instance_id`, `assumer_name`, `log_type`, `cloud_product_region`, `cls_region`, `logset_name`, `topic_name`, `extend`, `logset_id`, `topic_id`, `force_delete`, `is_delete_topic`, and `is_delete_logset`. The `extend` field is the only currently-mutable (non-ForceNew) parameter updated via `ModifyCloudProductLogCollection`.

The vendored SDK confirms that both `CreateCloudProductLogCollectionRequest` and `ModifyCloudProductLogCollectionRequest` already include a `Tags []*Tag` field (where `Tag` has `Key` and `Value` string pointers). No SDK upgrade is required.

The `CloudProductLogTaskInfo` response (from `DescribeCloudProductLogTasks`) exposes `TopicTags []*Tag` and `LogsetTags []*Tag` separately. The `LogsetInfo` and `TopicInfo` structs also carry `Tags []*Tag`. These are the available read-back sources for tags.

## Goals / Non-Goals

**Goals:**
- Add an optional `tags` (TypeMap of string) parameter to the resource schema.
- Pass `tags` to the create API so tags are bound at creation time.
- Pass `tags` to the modify API on update so tags can be changed in-place (no recreation).
- Read tags back into state so `terraform plan` stays clean after apply.
- Add a tags usage example to the resource documentation.
- Add unit test cases (gomonkey mocks) covering tags in create and update flows.

**Non-Goals:**
- Changing any existing schema field (no ForceNew/Computed changes to existing parameters).
- Adding tags support to the legacy `tencentcloud_cls_cloud_product_log_task` resource (non-v2).
- Introducing a separate tag-management resource or data source.
- Modifying the service-layer method signatures for create/modify beyond what is needed (tags will be set inline in the resource CRUD functions, consistent with how `extend` is handled today).

## Decisions

### Decision 1: `tags` schema type — TypeMap vs TypeList

**Choice**: Use `schema.TypeMap` with `Elem: &schema.Schema{Type: schema.TypeString}`.

**Rationale**: Terraform idiomatic tag representation for simple key-value string pairs is a map. The cloud API `Tag` struct has `Key`/`Value` string fields, which maps cleanly to a TypeMap. A TypeList of nested objects would add unnecessary complexity for the user. This also matches the pattern used by the reference proposal (`add-clb-target-group-missing-params`) for simple tags.

**Alternative considered**: `schema.TypeList` of nested `Key`/`Value` blocks — rejected as overly verbose for flat string tags.

### Decision 2: `tags` is mutable (not ForceNew)

**Choice**: `tags` is `Optional` without `ForceNew`, and is included in the `mutableArgs` check in the update function.

**Rationale**: The `ModifyCloudProductLogCollection` API explicitly accepts a `Tags` field, confirming the cloud API supports in-place tag modification. Marking it ForceNew would cause needless resource recreation.

**Alternative considered**: ForceNew — rejected because the modify API supports it.

### Decision 3: Read-back source for tags

**Choice**: Read tags from the `DescribeCloudProductLogTasks` response (`CloudProductLogTaskInfo`). The response exposes `TopicTags` and `LogsetTags`. Since the create/modify API binds the same `Tags` to both the logset and topic, we read from `TopicTags` (falling back to `LogsetTags` if `TopicTags` is empty) and flatten into the `tags` map.

**Rationale**: This avoids extra API calls; `DescribeClsCloudProductLogTaskById` is already invoked during Read. The existing Read logic already calls `DescribeClsLogset`/`DescribeClsTopicById` for names, which also expose `Tags`, but using the task info response keeps the tags read consistent with how the task was configured (the create/modify API binds tags to both logset and topic).

**Alternative considered**: Read `Tags` from `LogsetInfo.Tags` / `TopicInfo.Tags` via the existing `DescribeClsLogset` / `DescribeClsTopicById` calls. This is viable too, but adds dependency on those calls succeeding (they are only invoked when `logset_id`/`topic_id` are present). Using the task info response is more direct.

### Decision 4: Expansion/flattening helpers

**Choice**: Convert the Terraform `tags` map to `[]*clsv20201016.Tag` inline in the create/update functions (small loop), and flatten `[]*clsv20201016.Tag` back to a map in the read function. No separate helper file needed.

**Rationale**: The conversion is trivial (a few lines) and keeps the change localized, consistent with the existing code style where parameter mapping is inline.

### Decision 5: Update flow integration

**Choice**: Add `"tags"` to the existing `mutableArgs` slice in `resourceTencentCloudClsCloudProductLogTaskV2Update`. When `tags` has changed, include `request.Tags` in the `ModifyCloudProductLogCollection` request. Keep the existing `extend` mutable arg behavior unchanged.

**Rationale**: The update function already has a `mutableArgs` + `needChange` pattern for `extend`. Reusing it for `tags` is the minimal, consistent change.

## Risks / Trade-offs

- **[Tags read-back mismatch]** The create/modify API binds tags to both logset and topic, but `TopicTags` and `LogsetTags` in the response could theoretically diverge if modified out-of-band. → **Mitigation**: Read from `TopicTags` (primary) and fall back to `LogsetTags`; document that out-of-band tag changes may surface as drift, which is standard Terraform behavior.
- **[Empty tags vs nil tags]** Passing an empty (non-nil) `[]*Tag` vs `nil` on modify when the user removes all tags. → **Mitigation**: When the `tags` map is empty/unset, still call modify with an empty slice so the API clears the tags, matching the user's intent to remove tags. Guard against sending tags only when needed to avoid spurious updates.
- **[Tag limit]** The API documents a maximum of 10 tag key-value pairs. → **Mitigation**: Rely on API-side validation; do not add client-side count validation (the API is the source of truth and limits may change).
- **[Backward compatibility]** Adding an optional field is safe for existing state/configs. → **Mitigation**: The field is `Optional` with no `Required` and no `Default`; existing resources are unaffected.

## Migration Plan

This is a purely additive, backward-compatible change. No migration is required:

1. Users upgrade the provider version.
2. Existing configurations continue to work unchanged (the new `tags` parameter is optional).
3. Users can opt-in to `tags` by adding the parameter to existing or new resources.
4. State is populated with tags on the next `terraform refresh`/`apply`.

**Rollback**: Reverting the provider version removes the `tags` parameter from the schema. Terraform will show the field as unknown/removed in state, but existing resources are not affected because the cloud-side tags persist independently.

## Open Questions

- None blocking. The SDK fields required for create, modify, and read all exist in the vendored SDK and have been verified.
