## 1. Schema Definition

- [x] 1.1 Add `instance_concurrency_config` (TypeList, Optional) to resource schema in `resource_tc_scf_function.go`, with sub-fields: `dynamic_enabled` (TypeString, Optional), `max_concurrency` (TypeInt, Optional), `instance_isolation_enabled` (TypeString, Optional), `type` (TypeString, Optional), `mix_node_config` (TypeList, Optional with `node_spec` TypeString and `num` TypeInt sub-fields), `session_config` (TypeList, Optional with `session_source`, `session_name`, `maximum_concurrency_session_per_instance`, `maximum_ttl_in_seconds`, `maximum_idle_time_in_seconds` sub-fields)

## 2. Service Layer (`service_tencentcloud_scf.go`)

- [x] 2.1 Add `instanceConcurrencyConfig *scf.InstanceConcurrencyConfig` field to `scfFunctionInfo` struct
- [x] 2.2 In `CreateFunction`, set `request.InstanceConcurrencyConfig = info.instanceConcurrencyConfig`
- [x] 2.3 In `ModifyFunctionConfig`, set `request.InstanceConcurrencyConfig = info.instanceConcurrencyConfig`

## 3. CRUD Logic (`resource_tc_scf_function.go`)

- [x] 3.1 In `resourceTencentCloudScfFunctionCreate`: parse `instance_concurrency_config` from ResourceData and populate `scfFunctionInfo.instanceConcurrencyConfig`
- [x] 3.2 In `resourceTencentCloudScfFunctionRead`: after `GetFunction` response, if `resp.InstanceConcurrencyConfig` is not nil, parse and set `instance_concurrency_config` in state (including all sub-fields with nil checks)
- [x] 3.3 In `resourceTencentCloudScfFunctionUpdate`: add `HasChange("instance_concurrency_config")` check and populate `scfFunctionInfo.instanceConcurrencyConfig` for `UpdateFunctionConfiguration` call

## 4. Documentation

- [x] 4.1 Create/update `resource_tc_scf_function.md` with `instance_concurrency_config` example usage in the Example Usage section

## 5. Unit Tests

- [x] 5.1 Add unit test cases in `resource_tc_scf_function_test.go` for `instance_concurrency_config` using gomonkey mock to verify Create, Read, and Update logic

## 6. Registration

- [x] 6.1 Verify resource is already registered in `tencentcloud/provider.go` (no changes expected for existing resource)

## 7. Verification

- [x] 7.1 Run `gofmt` on all changed Go files to ensure formatting (deferred to tfpacer-finalize skill)
- [x] 7.2 Run `make doc` to generate website documentation (deferred to tfpacer-finalize skill)