## 1. Resource Schema Definition

- [x] 1.1 Create `tencentcloud/services/cls/resource_tc_cls_cls_deliver_task.go` with `ResourceTencentCloudClsClsDeliverTask()` returning `*schema.Resource` with Importer (`schema.ImportStatePassthrough`)
- [x] 1.2 Define top-level schema fields: `task_name` (Required string), `compliance` (Required int), `has_services_log` (Optional int), `enable` (Optional int), `task_id`/`uin`/`status`/`progress`/`create_time`/`update_time`/`total` (Computed)
- [x] 1.3 Define `source_topic_config` block (Required, MaxItems 1) with `topic_filter_type` (Required int), `logset_id` (Required string), `topics` (Optional list of blocks with `topic_id` Required string)
- [x] 1.4 Define `target_topic_config` block (Required, MaxItems 1) with `account_type` (Required int), `region` (Required string), `logset_id` (Required string), `topic_id` (Required string), `role_arn` (Optional string), `external_id` (Optional string)
- [x] 1.5 Define `deliver_rule` block (Required, MaxItems 1) with `data_scope` (Required int)

## 2. Service Layer

- [x] 2.1 Add `DescribeClsClsDeliverTaskById(ctx, taskId) (*cls.CLSDeliverTaskInfo, error)` to `tencentcloud/services/cls/service_tencentcloud_cls.go` calling `DescribeCLSDeliverTasks` with a `taskId` filter and `Limit=100`, returning `Infos[0]` (nil if empty)

## 3. CRUD Implementation

- [x] 3.1 Implement `resourceTencentCloudClsClsDeliverTaskCreate`: build `CreateCLSDeliverTaskRequest` (TaskName, SourceTopicConfig, TargetTopicConfig, DeliverRule, Compliance, HasServicesLog), call via `resource.Retry(tccommon.WriteRetryTimeout)`, check Response/TaskId non-nil with `[CRUD] cls_cls_deliver_task id=<id>` log, `d.SetId(taskId)`, then Read
- [x] 3.2 Implement `resourceTencentCloudClsClsDeliverTaskRead`: call service `DescribeClsClsDeliverTaskById`, nil → log `[CRUD] cls_cls_deliver_task id=<id>` then `d.SetId("")`; else set fields only when non-nil
- [x] 3.3 Implement `resourceTencentCloudClsClsDeliverTaskUpdate`: detect changes in mutableArgs (`task_name`, `source_topic_config`, `target_topic_config`, `deliver_rule`, `enable`, `has_services_log`); build `ModifyCLSDeliverTaskRequest` (TaskId + changed fields; NOT Compliance), call via `resource.Retry(tccommon.WriteRetryTimeout)`, then Read
- [x] 3.4 Implement `resourceTencentCloudClsClsDeliverTaskDelete`: build `DeleteCLSDeliverTaskRequest` with TaskId from `d.Id()`, call via `resource.Retry(tccommon.WriteRetryTimeout)`

## 4. Provider Registration

- [x] 4.1 Register `tencentcloud_cls_cls_deliver_task` → `cls.ResourceTencentCloudClsClsDeliverTask()` in `tencentcloud/provider.go` resource map (near other cls resources)
- [x] 4.2 Add `tencentcloud_cls_cls_deliver_task` entry to `tencentcloud/provider.md`

## 5. Documentation

- [x] 5.1 Create `tencentcloud/services/cls/resource_tc_cls_cls_deliver_task.md` with one-line description ("Provides a resource to create CLS deliver task"), Example Usage (using nested blocks), and Import section (import by task_id)

## 6. Unit Tests

- [x] 6.1 Create `tencentcloud/services/cls/resource_tc_cls_cls_deliver_task_test.go` using gomonkey mocks; add TestClsDeliverTask_Create_Basic, TestClsDeliverTask_Read_Basic, TestClsDeliverTask_Read_EmptyResult, TestClsDeliverTask_Update_Basic, TestClsDeliverTask_Delete_Basic covering business logic without real API calls

## 7. Finalization (tfpacer-finalize skill)

- [ ] 7.1 Run `gofmt` on changed Go files via tfpacer-finalize skill
- [ ] 7.2 Generate website docs via `make doc` (tfpacer-finalize skill)
- [ ] 7.3 Add changelog entry under `.changelog/` (tfpacer-finalize skill)