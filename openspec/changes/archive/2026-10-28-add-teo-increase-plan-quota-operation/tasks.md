## 1. Framework Action Implementation

- [x] 1.1 Rewrite `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.go` as a terraform-plugin-framework `action.Action` implementing `action.Action` and `action.ActionWithConfigure`: schema (plan_id, quota_type, quota_number as Required; deal_name as Computed), `Metadata`, `Schema`, `Configure` (type-assert `*sharedmeta.ProviderMeta`), and `Invoke` calling `IncreasePlanQuotaWithContext` API with `tccommon.WriteRetryTimeout` retry
- [x] 1.2 Add `SetClientForTest` helper method to support unit testing

## 2. Unit Tests

- [x] 2.1 Rewrite `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation_test.go` with unit tests using gomonkey to mock the TEO client (`IncreasePlanQuotaWithContext`), covering successful invocation, API error, nil response, nil deal_name, missing plan_id, missing quota_type, schema validation, metadata, and Configure (nil/wrong-type) scenarios

## 3. Framework Registration

- [x] 3.1 Register the new action `teo.NewTeoIncreasePlanQuotaAction` in `tencentcloud/framework/registry.go` (`actionFactories`)
- [x] 3.2 List the action `tencentcloud_teo_increase_plan_quota` under the TEO `Action` section in `tencentcloud/provider.md`
- [x] 3.3 Remove the old SDKv2 resource registration from `tencentcloud/provider.go` (ResourcesMap) and `tencentcloud/provider.md` (Resource section)

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/action_tc_teo_increase_plan_quota.md` with description, Example Usage (action block), following the framework action doc format
- [x] 4.2 Remove the old `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.md` (no longer an SDKv2 resource)

## 5. Validation

- [ ] 5.1 Verify the code compiles successfully (via build check in later steps)
- [ ] 5.2 Verify `make doc` generates the action documentation under `website/docs/a/`