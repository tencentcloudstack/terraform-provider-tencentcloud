## 1. Service Layer

- [x] 1.1 Add `DescribeGa2ListenerById` method to `tencentcloud/services/ga2/service_tencentcloud_ga2.go` that calls DescribeListeners API with Filters on `listener-id`, supports pagination, and returns `*ga2v20250115.ListenerSet` or nil

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/ga2/resource_tc_ga2_listener.go` with schema definition including all fields (global_accelerator_id, name, port_ranges, description, listener_type, protocol, idle_timeout, get_real_ip_type, client_affinity, client_affinity_time, request_timeout, x_forwarded_for_real_ip, certification_type, cipher_policy_id, server_certificates, client_ca_certificates, listener_id), Timeouts block, and Importer
- [x] 2.2 Implement `resourceTencentCloudGa2ListenerCreate` function calling CreateListener API with retry, polling DescribeTaskResult via WaitForGa2TaskFinish, and setting composite ID
- [x] 2.3 Implement `resourceTencentCloudGa2ListenerRead` function calling DescribeGa2ListenerById and setting all non-nil fields
- [x] 2.4 Implement `resourceTencentCloudGa2ListenerUpdate` function calling ModifyListener API with changed fields, polling DescribeTaskResult
- [x] 2.5 Implement `resourceTencentCloudGa2ListenerDelete` function calling DeleteListener API with retry, polling DescribeTaskResult

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_ga2_listener` in `tencentcloud/provider.go` ResourcesMap
- [x] 3.2 Add `tencentcloud_ga2_listener` entry in `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/ga2/resource_tc_ga2_listener.md` with Example Usage and Import sections

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/ga2/resource_tc_ga2_listener_test.go` with gomonkey-based unit tests covering Create, Read, Update, Delete operations
- [x] 5.2 Run `go test -gcflags=all=-l` to verify unit tests pass
