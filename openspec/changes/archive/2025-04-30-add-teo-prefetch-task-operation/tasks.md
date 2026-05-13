## 1. Resource Schema & CRUD Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_prefetch_task_operation.go` with schema definition including: zone_id (Required, ForceNew), targets (Required, ForceNew, TypeList of TypeString), mode (Optional, ForceNew), headers (Optional, ForceNew, TypeList of TypeMap with name/value), prefetch_media_segments (Optional, ForceNew), job_id (Computed), tasks (Computed, TypeList with nested schema for task results)
- [x] 1.2 Implement Create method: call CreatePrefetchTask API with retry (WriteRetryTimeout), extract JobId, then poll DescribePrefetchTasks using job-id filter (6*ReadRetryTimeout) until status is not "processing"; handle success/failed/timeout statuses; set resource ID as zoneId:jobId
- [x] 1.3 Implement empty Read method (RESOURCE_KIND_OPERATION pattern)
- [x] 1.4 Implement empty Delete method (RESOURCE_KIND_OPERATION pattern)

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_prefetch_task` resource in `tencentcloud/provider.go`
- [x] 2.2 Add resource entry in `tencentcloud/provider.md`

## 3. Unit Tests

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_prefetch_task_operation_test.go` with gomonkey mock tests for: successful prefetch task creation, failed task scenario

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_prefetch_task_operation.md` with description, example usage (following gendoc README format)
