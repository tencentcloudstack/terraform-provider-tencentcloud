## 1. Schema Definition

- [x] 1.1 Add the `key` field to the `tencentcloud_waf_cc_session` resource schema in `tencentcloud/services/waf/resource_tc_waf_cc_session.go`
  - Type: `schema.TypeString`
  - Optional: true (not Required, not Computed, not ForceNew)
  - Description: "Precise-match session key, configured when Category is precise matching."

## 2. Create Operation

- [x] 2.1 In `resourceTencentCloudWafCcSessionCreate`, read `key` from schema via `d.GetOk("key")` and set `request.Key` using `helper.String(v.(string))` when present, following the same pattern as the other string fields (e.g. `source`, `category`).

## 3. Read Operation

- [x] 3.1 In `resourceTencentCloudWafCcSessionRead`, after fetching `ccSession`, add a nil-guarded `d.Set("key", ccSession.Key)` block (only set when `ccSession.Key != nil`), matching the existing pattern for other response fields.

## 4. Update Operation

- [x] 4.1 In `resourceTencentCloudWafCcSessionUpdate`, confirm `key` is NOT added to the `immutableArgs` slice (it must remain mutable).
- [x] 4.2 In `resourceTencentCloudWafCcSessionUpdate`, read `key` from schema via `d.GetOk("key")` and set `request.Key` when present, following the same pattern used in the create operation.

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/waf/resource_tc_waf_cc_session.md` example to include the new `key` parameter in the Example Usage block.
- [x] 5.2 Run `make doc` (during finalize phase) to regenerate the `website/docs/` documentation from the updated `.md` file.

## 6. Testing

- [x] 6.1 Update `tencentcloud/services/waf/resource_tc_waf_cc_session_test.go` to add `key` to the test HCL configs and add `resource.TestCheckResourceAttr` assertions for the `key` parameter in both the basic and update test steps.

## 7. Verification (deferred to finalize phase)

- [x] 7.1 Run `gofmt` on modified Go files (during finalize phase via tfpacer-finalize skill).
- [x] 7.2 Verify the provider registration in `tencentcloud/provider.go` does not require changes (no new resource, only a new schema field on an existing resource).
