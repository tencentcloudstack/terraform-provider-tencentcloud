# Spec: CLS Deliver Task Resource

**Capability**: cls-cls-deliver-task
**Related Change**: add-cls-cls-deliver-task
**Status**: Draft

---

## ADDED Requirements

### Requirement: Resource tencentcloud_cls_cls_deliver_task manages CLS deliver task lifecycle

The provider SHALL expose a `tencentcloud_cls_cls_deliver_task` resource of kind RESOURCE_KIND_GENERAL that creates, reads, updates, and deletes a CLS cross-account/cross-region deliver task via the cloud APIs `CreateCLSDeliverTask`, `DescribeCLSDeliverTasks`, `ModifyCLSDeliverTask`, and `DeleteCLSDeliverTask`.

#### Scenario: Create a deliver task with source and target topic config
- **WHEN** the user applies a configuration specifying `task_name`, a `source_topic_config` block (with `topic_filter_type`, `logset_id`, and a `topics` list of `topic_id`), a `target_topic_config` block (with `account_type`, `region`, `logset_id`, `topic_id`), a `deliver_rule` block (with `data_scope`), and `compliance`
- **THEN** the provider calls `CreateCLSDeliverTask`, sets the resource id to the returned `TaskId`, and stores `task_id`, `status`, `progress`, `create_time`, `update_time`, `uin` as computed fields

#### Scenario: Read an existing deliver task by id
- **WHEN** the provider refreshes state for a resource with a known id
- **THEN** the provider calls `DescribeCLSDeliverTasks` with a `taskId` filter, populates all schema fields from the first `Infos` entry, and if the list is empty logs `[CRUD] cls_cls_deliver_task id=<id>` and clears the id

#### Scenario: Update mutable fields
- **WHEN** the user changes `task_name`, `source_topic_config`, `target_topic_config`, `deliver_rule`, `enable`, or `has_services_log`
- **THEN** the provider calls `ModifyCLSDeliverTask` with `TaskId` and the changed fields

#### Scenario: Delete a deliver task
- **WHEN** the user destroys the resource
- **THEN** the provider calls `DeleteCLSDeliverTask` with `TaskId`

#### Scenario: Import by task id
- **WHEN** the user runs `terraform import tencentcloud_cls_cls_deliver_task.example <task_id>`
- **THEN** the provider uses the task id as the resource id and reads the resource state

### Requirement: Schema fields match cloud API parameters

The resource schema SHALL define fields whose paths (SchemaName) map to the cloud API parameter paths (JsonPath) as follows:

Top-level:
- `task_name` (Required, string) → `request.TaskName`
- `source_topic_config` (Required, list max 1) → `request.SourceTopicConfig`
- `target_topic_config` (Required, list max 1) → `request.TargetTopicConfig`
- `deliver_rule` (Required, list max 1) → `request.DeliverRule`
- `compliance` (Required, int) → `request.Compliance` (create-only)
- `has_services_log` (Optional, int) → `request.HasServicesLog`
- `enable` (Optional, int) → `request.Enable` (update-only)
- `task_id` (Computed, string) → `response.Response.TaskId`
- `uin` (Computed, int) → `response.Response.Infos.Uin`
- `status` (Computed, int) → `response.Response.Infos.Status`
- `progress` (Computed, int) → `response.Response.Infos.Progress`
- `create_time` (Computed, int) → `response.Response.Infos.CreateTime`
- `update_time` (Computed, int) → `response.Response.Infos.UpdateTime`
- `total` (Computed, int) → `response.Response.Total`

`source_topic_config` block:
- `topic_filter_type` (Required, int) → `request.SourceTopicConfig.TopicFilterType`
- `logset_id` (Required, string) → `request.SourceTopicConfig.LogsetId`
- `topics` (Optional, list) → `request.SourceTopicConfig.Topics`, each with `topic_id` (Required, string) → `request.SourceTopicConfig.Topics.TopicId`

`target_topic_config` block:
- `account_type` (Required, int) → `request.TargetTopicConfig.AccountType`
- `region` (Required, string) → `request.TargetTopicConfig.Region`
- `logset_id` (Required, string) → `request.TargetTopicConfig.LogsetId`
- `topic_id` (Required, string) → `request.TargetTopicConfig.TopicId`
- `role_arn` (Optional, string) → `request.TargetTopicConfig.RoleArn`
- `external_id` (Optional, string) → `request.TargetTopicConfig.ExternalId`

`deliver_rule` block:
- `data_scope` (Required, int) → `request.DeliverRule.DataScope`

#### Scenario: Field optionality matches API
- **WHEN** the schema is generated
- **THEN** `task_name`, `source_topic_config`, `target_topic_config`, `deliver_rule`, `compliance` are Required on create; `has_services_log` is Optional; within `source_topic_config`, `topic_filter_type` and `logset_id` are Required while `topics` is Optional with required `topic_id`; within `target_topic_config`, `account_type`, `region`, `logset_id`, `topic_id` are Required while `role_arn` and `external_id` are Optional; within `deliver_rule`, `data_scope` is Required.

### Requirement: CRUD parameters verified against cloud API interfaces

The generated CRUD code MUST only send parameters that exist in the corresponding cloud API request struct:

- Create sends only `TaskName`, `SourceTopicConfig`, `TargetTopicConfig`, `DeliverRule`, `Compliance`, `HasServicesLog` (all exist in `CreateCLSDeliverTaskRequest`).
- Update sends only `TaskId`, `TaskName`, `SourceTopicConfig`, `TargetTopicConfig`, `DeliverRule`, `Enable`, `HasServicesLog` (all exist in `ModifyCLSDeliverTaskRequest`). `Compliance` is NOT sent in update (not present in Modify request).
- Delete sends only `TaskId` (exists in `DeleteCLSDeliverTaskRequest`).

#### Scenario: Compliance is create-only
- **WHEN** the user modifies `compliance` on an existing resource
- **THEN** the provider does NOT send `Compliance` in the Modify call; the field is treated as immutable.

### Requirement: Unit tests use gomonkey mocks

The test file `resource_tc_cls_cls_deliver_task_test.go` SHALL use `github.com/agiledragon/gomonkey/v2` to mock the CLS client methods (`CreateCLSDeliverTaskWithContext`, `DescribeCLSDeliverTasksWithContext`, `ModifyCLSDeliverTaskWithContext`, `DeleteCLSDeliverTaskWithContext`) and test Create/Read/Update/Delete business logic without calling the real cloud API.

#### Scenario: Create test verifies request fields and id
- **WHEN** the Create test runs with a mock client returning `TaskId`
- **THEN** the resource id is set to the task id and captured request fields match the schema input