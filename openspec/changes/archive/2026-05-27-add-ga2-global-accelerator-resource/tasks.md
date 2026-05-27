## 1. Service Layer

- [x] 1.1 Add `DescribeGa2GlobalAcceleratorById` method to `tencentcloud/services/ga2/service_tencentcloud_ga2.go` that queries a single global accelerator instance by ID using `DescribeGlobalAccelerators` API with `global-accelerator-id` filter

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator.go` with schema definition including all input fields (`name`, `instance_charge_type`, `description`, `cross_border_type`, `cross_border_promise_flag`, `tags`) and computed fields (`create_time`, `state`, `status`, `ddos_id`, `cname`), with Timeouts block and Import support
- [x] 2.2 Implement `resourceTencentCloudGa2GlobalAcceleratorCreate` function that calls `CreateGlobalAccelerator` API, validates response, sets resource ID, and waits for async task completion via `WaitForGa2TaskFinish`
- [x] 2.3 Implement `resourceTencentCloudGa2GlobalAcceleratorRead` function that calls `DescribeGa2GlobalAcceleratorById`, handles not-found case, and sets all attributes from response
- [x] 2.4 Implement `resourceTencentCloudGa2GlobalAcceleratorUpdate` function that calls `ModifyGlobalAccelerator` API with changed fields (`name`, `description`, `cross_border_type`, `cross_border_promise_flag`) and waits for async task completion
- [x] 2.5 Implement `resourceTencentCloudGa2GlobalAcceleratorDelete` function that calls `DeleteGlobalAccelerator` API and waits for async task completion

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_ga2_global_accelerator` resource in `tencentcloud/provider.go` resource map
- [x] 3.2 Add `tencentcloud_ga2_global_accelerator` entry in `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator.md` with Example Usage and Import sections

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_test.go` with gomonkey-based unit tests covering Create, Read, Update, and Delete operations
