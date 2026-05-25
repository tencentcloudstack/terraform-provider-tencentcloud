## 1. Setup

- [x] 1.1 Create directory structure `tencentcloud/services/ga2/`

## 2. Service Layer

- [x] 2.1 Create `tencentcloud/services/ga2/service_tencentcloud_ga2.go` with GA2 service client initialization and helper functions (DescribeAccelerateAreas with pagination support)

## 3. Resource Implementation

- [x] 3.1 Create `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go` with schema definition (global_accelerator_id as ForceNew, accelerator_areas as Required List, accelerate_area_set as Computed List)
- [x] 3.2 Implement Create function: call CreateAccelerateAreas API with retry, validate non-nil response/TaskId, set resource ID, call Read to poll until areas appear
- [x] 3.3 Implement Read function: call DescribeAccelerateAreas with pagination, set accelerate_area_set in state, handle resource-not-found by clearing ID
- [x] 3.4 Implement Update function: call ModifyAccelerateAreas API with retry, call Read to poll until changes reflected
- [x] 3.5 Implement Delete function: first Read to get AcceleratorAreaIds, then call DeleteAccelerateAreas with retry, poll until areas removed

## 4. Provider Registration

- [x] 4.1 Register `tencentcloud_ga2_accelerate_area` resource in `tencentcloud/provider.go`
- [x] 4.2 Add resource entry in `tencentcloud/provider.md`

## 5. Documentation

- [x] 5.1 Create `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.md` with Example Usage and Import sections

## 6. Unit Tests

- [x] 6.1 Create `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area_test.go` with gomonkey-based unit tests covering Create, Read, Update, Delete flows
- [x] 6.2 Run unit tests with `go test -gcflags=all=-l` to verify tests pass
