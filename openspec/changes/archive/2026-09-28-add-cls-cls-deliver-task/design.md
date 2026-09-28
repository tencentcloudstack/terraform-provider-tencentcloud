## Context

The Terraform Provider for TencentCloud manages cloud resources via `tencentcloud-sdk-go`. The CLS (Cloud Log Service) product supports cross-account/cross-region log delivery tasks but no corresponding Terraform resource exists. This change adds a new RESOURCE_KIND_GENERAL resource `tencentcloud_cls_cls_deliver_task` implementing full CRUD using four cloud APIs from the `cls/v20201016` SDK package (already vendored).

The resource will live in `tencentcloud/services/cls/resource_tc_cls_cls_deliver_task.go` and follow the established CRUD skeleton of `tencentcloud_igtm_strategy` (nested block schema, single-id with `FILED_SP` separator is NOT needed here because the resource has a single `task_id`, retry with `tccommon.WriteRetryTimeout`/`ReadRetryTimeout`, immutable-args handling for ForceNew fields).

## Goals / Non-Goals

**Goals:**
- Provide full lifecycle management (CRUD + import) for CLS deliver tasks via Terraform.
- Expose all CRUD-relevant API parameters as schema fields with correct optionality (Required vs Optional) matching the cloud API.
- Surface read-only computed fields (status, progress, create_time, update_time, uin) for observability.
- Provide a service-layer helper to query a single deliver task by id (used by Read/Update/Delete after Create).
- Include gomonkey-based unit tests covering Create/Read/Update/Delete.

**Non-Goals:**
- Not a data source (RESOURCE_KIND_DATASOURCE). This change only adds the manage resource.
- Not exposing `filters`/`offset`/`limit` of `DescribeCLSDeliverTasks` as user-facing schema; the Read path uses a fixed `taskId` filter internally.
- Not supporting pagination for the describe call in the service layer (a single task lookup returns `Infos[0]`).

## Decisions

### Decision 1: Resource id uses single `task_id` (not a composite id)

**Rationale**: The `CreateCLSDeliverTask` response returns a single `TaskId`, and `ModifyCLSDeliverTask`/`DeleteCLSDeliverTask` only require `TaskId`. Unlike `igtm_strategy` (which needs `instanceId`+`strategyId`), this resource needs no companion id. So `d.SetId(taskId)` with a plain string is correct; no `FILED_SP` join needed.

**Alternative considered**: Making `task_id` a ForceNew field and joining it with something — rejected because no second key exists in the API.

### Decision 2: Read path queries `DescribeCLSDeliverTasks` with a `taskId` filter

**Rationale**: There is no single-task Describe API. The list API `DescribeCLSDeliverTasks` accepts a `Filters` array where `Key=taskId` filters by task id. The service helper `DescribeClsClsDeliverTaskById` builds this filter, takes `Infos[0]`, and returns `*cls.CLSDeliverTaskInfo`. Set `Limit=100` (API max) and a fixed offset 0.

**Alternative**: Calling DescribeCLSDeliverTasks without filter and searching client-side — rejected (inefficient, filter is supported).

### Decision 3: Schema flattening of nested objects into Terraform blocks

**Rationale**: `SourceTopicConfig`, `TargetTopicConfig`, and `DeliverRule` are single nested objects (MaxItems: 1 blocks). `SourceTopicConfig.Topics` is a list of `SourceTopicInfo` (each with `TopicId`). Mirrors the `igtm_strategy` nested block style.

### Decision 4: `compliance` and `has_services_log` use `TypeInt` (uint64 in SDK)

**Rationale**: SDK fields `Compliance`, `HasServicesLog`, `Enable`, `DataScope`, `TopicFilterType`, `AccountType` are all `*uint64`. Terraform schema uses `TypeInt` and converts via `helper.IntUint64()`.

### Decision 5: `task_id` is Optional + Computed (set after create, used in import), NOT ForceNew

**Rationale**: `task_id` is the resource's own id; it is set after create. Import uses it. It is not user-settable on create. Per provider conventions for a single-id resource, `task_id` is `Computed: true` (and the actual `d.Id()`).

### Decision 6: Mutable args for Update detection

Update detects changes on: `task_name`, `source_topic_config`, `target_topic_config`, `deliver_rule`, `enable`, `has_services_log`. `compliance` is create-only (not in ModifyCLSDeliverTask) and treated as immutable.

### Decision 7: API field verification against vendor struct definitions

Confirmed against `vendor/.../cls/v20201016/models.go`:

`CreateCLSDeliverTaskRequest` (models.go:2493): `TaskName *string`, `SourceTopicConfig *SourceTopicConfig`, `TargetTopicConfig *TargetTopicConfig`, `DeliverRule *DeliverRule`, `Compliance *uint64`, `HasServicesLog *uint64`. Response `CreateCLSDeliverTaskResponseParams` (2540): `TaskId *string`.

`DescribeCLSDeliverTasksRequest` (10608): `Filters []*Filter`, `Offset *uint64`, `Limit *uint64` (max 100). Response (10643): `Infos []*CLSDeliverTaskInfo`, `Total *uint64`.

`ModifyCLSDeliverTaskRequest` (18958): `TaskId *string`, `TaskName *string`, `SourceTopicConfig *SourceTopicConfig`, `TargetTopicConfig *TargetTopicConfig`, `DeliverRule *DeliverRule`, `Enable *uint64`, `HasServicesLog *uint64`. NOTE: `Compliance` is NOT in Modify — create-only.

`DeleteCLSDeliverTaskRequest` (7591): `TaskId *string`. Response empty (only `RequestId`).

Nested structs:
- `SourceTopicConfig` (25935): `TopicFilterType *uint64`, `LogsetId *string`, `Topics []*SourceTopicInfo`.
- `SourceTopicInfo` (25946): `TopicId *string`.
- `TargetTopicConfig` (26097): `AccountType *uint64`, `Region *string`, `LogsetId *string`, `TopicId *string`, `RoleArn *string`, `ExternalId *string`.
- `DeliverRule` (9950): `DataScope *uint64`.
- `CLSDeliverTaskInfo` (722): `TaskId *string`, `TaskName *string`, `Uin *uint64`, `SourceTopicConfig *SourceTopicConfig`, `TargetTopicConfig *TargetTopicConfig`, `DeliverRule *DeliverRule`, `Compliance *uint64`, `Status *uint64`, `Enable *uint64`, `Progress *uint64`, `HasServicesLog *uint64`, `CreateTime *uint64`, `UpdateTime *uint64`.
- `Filter` (17202): `Key *string`, `Values []*string`.

## Risks / Trade-offs

- [Risk] `DescribeCLSDeliverTasks` returns a list; a missing task returns an empty `Infos` array → Mitigation: service helper returns nil; Read logs `[CRUD] cls_cls_deliver_task id=<id>` then `d.SetId("")`.
- [Risk] `compliance` is create-only but appears in the response — user edits to it after create would be silently ignored → Mitigation: treat `compliance` as immutable (not in mutableArgs); on update, do not send it.
- [Risk] `Enable` field has dual meaning (create-time not present, modify-time toggles run/pause) → Mitigation: `enable` is Optional+Computed, only sent in ModifyCLSDeliverTask when `d.HasChange("enable")`.

## Migration Plan

N/A — purely additive new resource. No existing state to migrate. Rollback = delete the new files and revert provider.go/provider.md additions.