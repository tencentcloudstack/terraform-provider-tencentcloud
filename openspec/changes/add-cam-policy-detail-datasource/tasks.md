## 1. Service Layer

- [x] 1.1 Append `DescribeCamPolicyDetailByFilter()` to `service_tencentcloud_cam.go` — wraps `GetPolicy` with retry, accepts `paramMap["PolicyId"]`

## 2. Data Source Implementation

- [x] 2.1 Create `data_source_tc_cam_policy_detail.go` with schema and read handler
- [x] 2.2 Required input: `policy_id`; computed output: `policy_info` block (all response fields)
- [x] 2.3 Read handler: build paramMap, call service with retry, flatten into `policy_info`

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_cam_policy_detail` in `provider.go`

## 4. Documentation & Tests

- [x] 4.1 Create `data_source_tc_cam_policy_detail.md`
- [x] 4.2 Create `data_source_tc_cam_policy_detail_test.go`
