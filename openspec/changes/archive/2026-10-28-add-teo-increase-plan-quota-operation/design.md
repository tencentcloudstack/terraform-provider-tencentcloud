## Context

The TEO (EdgeOne) `IncreasePlanQuota` API allows users to purchase additional quotas for their EdgeOne plans. The API is synchronous (not async), takes three required input parameters (`PlanId`, `QuotaType`, `QuotaNumber`), and returns a `DealName` (order number).

This is a one-time operation that fits the terraform-plugin-framework `action.Action` interface: after calling the API, no state needs to be tracked — there's no resource lifecycle beyond the operation itself. Framework actions are invoked via `terraform apply` and do not persist state, which is exactly the semantics of `IncreasePlanQuota`.

The vendor directory already contains the TEO v20220901 SDK with `IncreasePlanQuotaRequest`, `IncreasePlanQuotaResponse`, and the corresponding client method (`IncreasePlanQuota` / `IncreasePlanQuotaWithContext`).

Existing framework references in this provider (e.g., `services/ssm/resource_tc_ssm_secret_version_v2.go`, `services/ssm/ephemeral_tc_ssm_secret_version.go`) serve as reference implementations for the framework pattern, including the `*sharedmeta.ProviderMeta` type-assert in `Configure` and the `sharedmeta` client reuse.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform framework Action `tencentcloud_teo_increase_plan_quota` that calls the `IncreasePlanQuota` API
- Support all three input parameters: `plan_id`, `quota_type`, `quota_number`
- Expose the response `deal_name` as a computed attribute and surface it to the user via a warning diagnostic (framework actions do not write state)
- Follow the established framework action pattern (`action.Action` + `action.ActionWithConfigure`)

**Non-Goals:**
- No polling (the API is synchronous, not async)
- No import support (actions do not support import)
- No update support (actions are one-time invocations)
- No SDKv2 resource registration (actions are registered exclusively via `tencentcloud/framework/registry.go`)

## Decisions

### 1. Resource Type: Framework Action (action.Action)
**Rationale**: `IncreasePlanQuota` is a one-time purchase operation. There's no resource to manage after the call completes. Framework actions are the idiomatic way to model one-time operations in the terraform-plugin-framework: they invoke logic on apply and do not persist state. This replaces the earlier SDKv2 `schema.Resource` (OPERATION) approach with the framework-native action.

### 2. Use Direct SDK Client Call
**Rationale**: The action calls `UseTeoV20220901Client().IncreasePlanQuotaWithContext(ctx, request)` directly. There is no separate service-layer method needed for this single API call, mirroring how SDKv2 OPERATION resources call the SDK client directly.

### 3. Schema Design
All three input parameters are `Required`. The output `deal_name` is `Computed`. Framework actions report results via diagnostics rather than persisted state, so the computed `deal_name` is surfaced to the user through a warning diagnostic containing the order number.

### 4. Retry with WriteRetryTimeout
**Rationale**: Since this is a write operation (not a read), use `tccommon.WriteRetryTimeout` for the retry block with `resource.Retry` and `tccommon.RetryError()` wrapper, consistent with the SDKv2 OPERATION pattern.

### 5. Configure via sharedmeta
**Rationale**: The action implements `action.ActionWithConfigure` and type-asserts `req.ProviderData` to `*sharedmeta.ProviderMeta` to obtain the shared `*connectivity.TencentCloudClient`, exactly like the existing framework resources/ephemerals.

### 6. Registration via framework/registry.go
**Rationale**: Framework actions are aggregated by `frameworkActions()` in `tencentcloud/framework/registry.go` and exposed via `Provider.Actions()`. The action factory `NewTeoIncreasePlanQuotaAction` is appended to `actionFactories`. No edit to the SDKv2 `provider.go` `ResourcesMap` is required.

## Risks / Trade-offs

- **Risk**: The API may fail due to insufficient balance or invalid quota type → **Mitigation**: Error is surfaced to the user via Terraform's standard error diagnostics reporting
- **Risk**: Duplicate calls could create duplicate orders → **Mitigation**: This is managed by the cloud API's own idempotency; framework actions are invoked once per `terraform apply`
- **Trade-off**: Framework actions do not persist state, so the `deal_name` is reported via a warning diagnostic rather than stored in tfstate. Users who need to reference the order number downstream should capture it from the action output.