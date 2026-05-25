## 1. Service Layer Setup

- [x] 1.1 Create `tencentcloud/services/ga2/` directory and `service_tencentcloud_ga2.go` with GA2 client initialization and API wrapper functions (CreateAccelerateAreas, DescribeAccelerateAreas, ModifyAccelerateAreas, DeleteAccelerateAreas)

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go` with schema definition including `global_accelerator_id`, `accelerator_areas` (input), `accelerate_area_set` (computed output), and Timeouts block
- [x] 2.2 Implement Create function: call CreateAccelerateAreas, then poll DescribeAccelerateAreas until areas appear
- [x] 2.3 Implement Read function: call DescribeAccelerateAreas with pagination, populate `accelerate_area_set`
- [x] 2.4 Implement Update function: call ModifyAccelerateAreas, then poll DescribeAccelerateAreas until changes are reflected
- [x] 2.5 Implement Delete function: read current area IDs, call DeleteAccelerateAreas, then poll DescribeAccelerateAreas until areas are gone

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_ga2_accelerate_area` in `tencentcloud/provider.go` resource map
- [x] 3.2 Add `tencentcloud_ga2_accelerate_area` entry in `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.md` with Example Usage and Import sections

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area_test.go` with gomonkey-based unit tests covering Create, Read, Update, and Delete operations
- [x] 5.2 Run unit tests with `go test -gcflags=all=-l` to verify they pass
