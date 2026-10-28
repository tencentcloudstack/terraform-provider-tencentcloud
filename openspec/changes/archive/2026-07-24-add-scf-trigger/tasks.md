## 1. Service layer helper

- [x] 1.1 Add `DescribeScfTriggerById(ctx, functionName, namespace, triggerName)` method to `tencentcloud/services/scf/service_tencentcloud_scf.go` that calls `ListTriggers` filtered by `FunctionName` + `Namespace` + `Filters(TriggerName=...)`, paginates internally, and returns the first matching `*scf.TriggerInfo` (or nil if not found).

## 2. Resource schema definition

- [x] 2.1 Create `tencentcloud/services/scf/resource_tc_scf_trigger.go` with `ResourceTencentCloudScfTrigger()` defining the schema: `function_name` (Required, ForceNew), `trigger_name` (Required, ForceNew), `type` (Required, ForceNew), `namespace` (Optional, ForceNew, Default "default"), `trigger_desc` (Optional), `qualifier` (Optional), `enable` (Optional, TypeString), `custom_argument` (Optional), `description` (Optional), and computed fields `available_status`, `add_time`, `mod_time`. Include `Importer` with `schema.ImportStatePassthrough`.

## 3. CRUD implementation

- [x] 3.1 Implement `resourceTencentCloudScfTriggerCreate`: build a `CreateTriggerRequest` from schema, call `CreateTrigger` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)` wrapping errors with `tccommon.RetryError()`, verify the response/`TriggerInfo` is non-nil (return `NonRetryableError` if empty), then set the composite id `function_name#namespace#trigger_name` and call Read.
- [x] 3.2 Implement `resourceTencentCloudScfTriggerRead`: parse the 3-part composite id, call `service.DescribeScfTriggerById`; if nil, log `[CRUD] scf_trigger id=<id>` then `d.SetId("")`; otherwise set fields from `TriggerInfo` with nil checks (convert `Enable` int64 1/0 → "OPEN"/"CLOSE"), and set computed `available_status`/`add_time`/`mod_time` only when non-nil.
- [x] 3.3 Implement `resourceTencentCloudScfTriggerUpdate`: parse the 3-part id, build an `UpdateTriggerRequest` with identity fields, set mutable fields (`enable`, `qualifier`, `trigger_desc`, `description`, `custom_argument`) when `d.HasChange`, call `UpdateTrigger` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`, then call Read.
- [x] 3.4 Implement `resourceTencentCloudScfTriggerDelete`: parse the 3-part id, build a `DeleteTriggerRequest` with `FunctionName`, `TriggerName`, `Type`, `Namespace`, `Qualifier`, and `TriggerDesc` (when applicable) from state, call `DeleteTrigger` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

## 4. Provider registration

- [x] 4.1 Register `tencentcloud_scf_trigger` in the resources map of `tencentcloud/provider.go` (refer to the `tencentcloud_igtm_strategy` registration pattern).
- [x] 4.2 Register `tencentcloud_scf_trigger` in `tencentcloud/provider.md`.

## 5. Resource documentation

- [x] 5.1 Create `tencentcloud/services/scf/resource_tc_scf_trigger.md` following the gendoc/README.md conventions: one-line description mentioning SCF, Example Usage block (use `jsonencode()` where JSON strings are needed), and Import section explaining the composite id `function_name#namespace#trigger_name`. Do NOT add `Argument Reference` / `Attribute Reference` sections (auto-generated).

## 6. Unit tests

- [x] 6.1 Create `tencentcloud/services/scf/resource_tc_scf_trigger_test.go` using gomonkey mocks for the SCF client methods (not the Terraform acceptance framework), covering create/read/update/delete business logic, and run with `go test -gcflags=all=-l` to verify the unit tests pass.

## 7. Code correctness verification

- [x] 7.1 Verify all CRUD code paths use only parameters that exist in the corresponding SCF API (Create params in `CreateTriggerRequest`, Update params in `UpdateTriggerRequest`, Delete params in `DeleteTriggerRequest`); confirm against the vendored `scf/v20180416/models.go`.
- [x] 7.2 Verify all functions returning an error have it checked (use `_ = func()` for functions that cannot fail) to avoid unused-variable compile errors.
