## 1. Service Layer

- [x] 1.1 Create `tencentcloud/services/ga2/service_tencentcloud_ga2.go` with GA2 service struct and helper methods for calling `DescribeAccelerateAreas` (with pagination and retry logic)

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go` with schema definition including `global_accelerator_id`, `accelerator_areas`, `accelerate_area_set`, `task_id`, and Timeouts block
- [x] 2.2 Implement Create function: call `CreateAccelerateAreas`, then poll `DescribeAccelerateAreas` until areas appear
- [x] 2.3 Implement Read function: call `DescribeAccelerateAreas` with pagination, flatten results into state
- [x] 2.4 Implement Update function: call `ModifyAccelerateAreas` with updated areas (including `AcceleratorAreaId`), then poll until changes are reflected
- [x] 2.5 Implement Delete function: read current area IDs via `DescribeAccelerateAreas`, call `DeleteAccelerateAreas`, then poll until areas are removed

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_ga2_accelerate_area` resource in `tencentcloud/provider.go`
- [x] 3.2 Add `tencentcloud_ga2_accelerate_area` entry in `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.md` with example usage and import section

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area_test.go` with gomonkey-based unit tests covering Create, Read, Update, and Delete flows
- [x] 5.2 Run unit tests with `go test -gcflags=all=-l` to verify they pass
