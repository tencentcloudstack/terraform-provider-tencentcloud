## Context

The TEO (EdgeOne) `IncreasePlanQuota` API allows users to purchase additional quotas for their EdgeOne plans. The API is sychronous (not async), takes three required input parameters (`PlanId`, `QuotaType`, `QuotaNumber`), and returns a `DealName` (order number).

This is a one-time operation resource (RESOURCE_KIND_OPERATION): after calling the API, no state needs to be tracked — there's no resource lifecycle beyond the operation itself.

The vendor directory already contains the TEO v20220901 SDK with `IncreasePlanQuotaRequest`, `IncreasePlanQuotaResponse`, and the corresponding client method.

Existing TEO OPERATION resources (e.g., `resource_tc_teo_check_cname_status_operation.go`, `resource_tc_teo_identify_zone_operation.go`) serve as reference implementations.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_teo_increase_plan_quota` that calls the `IncreasePlanQuota` API
- Support all three input parameters: `plan_id`, `quota_type`, `quota_number`
- Expose the response `deal_name` as a computed output
- Follow the established OPERATION resource pattern (Create does the work, Read/Delete are no-ops)

**Non-Goals:**
- No polling (the API is synchronous, not async)
- No import support (OPERATION resources don't support import)
- No update support (OPERATION resources are one-time only)

## Decisions

### 1. Resource Type: RESOURCE_KIND_OPERATION
**Rationale**: `IncreasePlanQuota` is a one-time purchase operation. There's no resource to manage after the call completes. OPERATION resources have empty Read/Delete and use `helper.BuildToken()` for the ID.

### 2. Use Direct API Call (Not Service Layer)
**Rationale**: For simple OPERATION resources that don't need complex logic or reuse, calling the SDK client directly from the resource file (like `resource_tc_teo_check_cname_status_operation.go`) is simpler and follows the existing pattern. No need for a separate service layer method.

### 3. Schema Design
All three input parameters are `Required` + `ForceNew` since the operation must be idempotent and re-triggered on any change. The output `deal_name` is `Computed` only.

### 4. Retry with WriteRetryTimeout
**Rationale**: Since this is a write operation (not a read), use `tccommon.WriteRetryTimeout` for the retry block. This aligns with other OPERATION resources that perform API calls in Create.

### 5. Error Handling
- Check for nil response after API call, return `NonRetryableError` if response is nil
- Use `tccommon.RetryError()` wrapper for retriable API errors
- Log with appropriate debug/critical levels following existing patterns

## Risks / Trade-offs

- **Risk**: The API may fail due to insufficient balance or invalid quota type → **Mitigation**: Error is surfaced to the user via Terraform's standard error reporting
- **Risk**: Duplicate calls could create duplicate orders → **Mitigation**: This is managed by the cloud API's own idempotency; the Terraform resource is ForceNew on all params so repeated applies only re-trigger on changes
