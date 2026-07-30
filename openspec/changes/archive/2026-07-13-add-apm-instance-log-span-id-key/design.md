## Context

The `tencentcloud_apm_instance` resource (in `tencentcloud/services/apm/resource_tc_apm_instance.go`) manages an APM (应用性能监控) business system. The resource already exposes the `log_trace_id_key` parameter, which maps to the CLS index key for `traceId` and is valid when `log_index_type = 1` (key-value index).

The companion API field `LogSpanIdKey` (the CLS index key for `spanId`) is:
- Present in the `ModifyApmInstanceRequest` struct (update API).
- Present in the `ApmInstanceDetail` struct returned by `DescribeApmInstances` (read API).
- NOT present in `CreateApmInstanceRequest` (create API).

This matches the existing pattern for nearly every config parameter on this resource: `CreateApmInstance` only accepts a small subset (`Name`, `Description`, `TraceDuration`, `Tags`, `SpanDailyCounters`, `PayMode`, `Free`); all other parameters (including `LogTraceIdKey`) are applied through a post-create `ModifyApmInstance` call performed inside the Create function, and through `ModifyApmInstance` in the Update function.

## Goals / Non-Goals

**Goals:**
- Add the `log_span_id_key` Optional string parameter to the `tencentcloud_apm_instance` resource schema, mirroring `log_trace_id_key`.
- Persist the parameter via the post-create `ModifyApmInstance` call and the Update `ModifyApmInstance` call.
- Read the parameter back from `DescribeApmInstances` (`ApmInstanceDetail.LogSpanIdKey`) into Terraform state.
- Add a unit test (mock-based, gomonkey) covering the new parameter.
- Update the resource `.md` documentation.

**Non-Goals:**
- No changes to the `tencentcloud_apm_instances` data source (separate capability, already exposes `LogSpanIdKey` as a Computed field).
- No changes to `CreateApmInstance` (the API does not accept `LogSpanIdKey`; the post-create Modify call is the supported mechanism).
- No new external dependencies.

## Decisions

### Decision 1: Schema field definition mirrors `log_trace_id_key`
- Name: `log_span_id_key`
- Type: `schema.TypeString`
- Optional: true (no `Computed`, matching `log_trace_id_key`)
- Description: "Index key of spanId. It is valid when the CLS index type is key-value index."

**Rationale**: `log_trace_id_key` is the direct sibling field and uses the exact same shape. Keeping the two consistent avoids surprising users.

### Decision 2: Apply via post-create `ModifyApmInstance`, not `CreateApmInstance`
`CreateApmInstanceRequest` does not have a `LogSpanIdKey` field, so the value cannot be passed at creation time. The existing resource Create function already performs a second `ModifyApmInstance` call (the `configRequest`) to set all the config-only parameters after the instance is created. We add `configRequest.LogSpanIdKey` there, following the same `d.GetOk("log_span_id_key")` pattern used for `log_trace_id_key`.

**Alternative considered**: Making the field `Computed`-only until first update. Rejected — it would prevent users from setting the value at create time, which the API supports through the post-create Modify call that every other config parameter already uses.

### Decision 3: Update flow uses the same `d.GetOk` pattern
In `resourceTencentCloudApmInstanceUpdate`, the resource rebuilds the entire `ModifyApmInstanceRequest` from `d.GetOk`/`d.GetOkExists` for every field (there is no per-field `d.HasChange` gating). We add `request.LogSpanIdKey` using `d.GetOk("log_span_id_key")`, identical to how `log_trace_id_key` is handled. This keeps the diff behavior consistent with the rest of the resource.

### Decision 4: Read flow nil-guards the field
In `resourceTencentCloudApmInstanceRead`, we add a nil check for `instance.LogSpanIdKey` before calling `d.Set("log_span_id_key", ...)`, matching the pattern used for every other field in the read function (per project rule: do not call `setXX()` when the response field is nil).

### Decision 5: Mock-based unit test (gomonkey)
Per project rules, new terraform resources use mock-based unit tests (gomonkey) rather than the Terraform test suite. `tencentcloud_apm_instance` is an existing resource, but the test guidance for newly-added parameters on existing files is to keep using the existing test suite style. We will add a test case to `resource_tc_apm_instance_test.go` following the existing test patterns in that file, asserting the `log_span_id_key` attribute round-trips through Create/Read/Update.

## Risks / Trade-offs

- [Risk] API may ignore `LogSpanIdKey` when `log_index_type` is not key-value index (0). → Mitigation: This is a documented API precondition shared with `log_trace_id_key`; the provider does not need to enforce it. Users who set the field without key-value indexing will get the API's default behavior, same as the sibling field.
- [Risk] Drift if the value is set out-of-band. → Mitigation: The Read flow reads `ApmInstanceDetail.LogSpanIdKey` into state, so external changes are detected on refresh.
- [Trade-off] The Create function sends `LogSpanIdKey` in the post-create Modify call even if the user did not set it. → This is consistent with the existing behavior for all other config fields on this resource (`d.GetOk` returns false and the field is left nil when unset, so it is simply not sent). No change to existing behavior.

## Migration Plan

No migration required. The field is Optional and defaults to unset; existing configurations are unaffected. Rollback is achieved by reverting the code change — no state migration is involved because the field is not Computed (unset state is equivalent to absent).
