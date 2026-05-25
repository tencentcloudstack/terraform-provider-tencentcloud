## 1. Service Layer

- [x] 1.1 Create `tencentcloud/services/ga2/service_tencentcloud_ga2.go` with GA2 client initialization and helper functions for DescribeAccelerateAreas (with pagination), CreateAccelerateAreas, ModifyAccelerateAreas, DeleteAccelerateAreas, and DescribeTaskResult polling logic

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go` with schema definition (global_accelerator_id ForceNew, accelerator_areas input list, accelerate_area_set computed list, task_id computed, Timeouts block) and CRUD functions that call service layer helpers with async task polling after Create/Modify/Delete

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_ga2_accelerate_area` resource in `tencentcloud/provider.go` and add entry in `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.md` with Example Usage and Import sections

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area_test.go` with unit tests using gomonkey to mock cloud API calls, covering Create, Read, Update, and Delete flows
