## Context

TEO (EdgeOne) exposes a `CreatePlanForZone` API in the vendored SDK `teo/v20220901` (verified present — no SDK upgrade needed):

- `CreatePlanForZone(request)` — purchases a plan for a zone that has not yet bound a plan. Request inputs: `ZoneId *string`, `PlanType *string`. Response outputs: `ResourceNames []*string`, `DealNames []*string`, `RequestId *string`.
- The API is synchronous; no async task/polling is needed.

This is a **one-time operation resource** (RESOURCE_KIND_OPERATION), following the style of `tencentcloud_teo_identify_zone_operation`: it has `Create`, `Read`, `Delete` (no `Update`), performs the action on Create, and treats Read/Delete as no-ops. Unlike `tencentcloud_teo_plan` (which uses the separate `CreatePlan` API and is a full CRUD resource), this resource targets a zone that has no bound plan and only performs the purchase action.

## Goals / Non-Goals

**Goals:**
- Trigger plan purchase for a zone via `CreatePlanForZone` on create, with retry and nil-safe response handling.
- Auto-generate the resource ID (per business rule), using `helper.BuildToken()`.
- Surface the response's `ResourceNames` and `DealNames` as computed attributes.
- Ship resource example `.md`, website docs, unit test, and provider registration.

**Non-Goals:**
- No revert/cancel of a purchased plan on delete (the API has no such capability); Delete is a no-op.
- No Read API to query purchase status; Read is a no-op that returns nil.
- No Timeouts block (synchronous call, no async status polling needed).

## Decisions

### Resource name and file layout
- Resource: `tencentcloud_teo_plan_for_zone`
- Files under `tencentcloud/services/teo/`:
  - `resource_tc_teo_plan_for_zone_operation.go`
  - `resource_tc_teo_plan_for_zone_operation.md`
  - `resource_tc_teo_plan_for_zone_operation_test.go`
- Website doc: `website/docs/r/teo_plan_for_zone.html.markdown`
- Constructor `ResourceTencentCloudTeoPlanForZone()`, registered in `provider.go`.

### Schema
- Required (ForceNew):
  - `zone_id` (TypeString): Zone ID.
  - `plan_type` (TypeString): Plan type to purchase. Valid values: `sta`, `sta_with_bot`, `sta_cm`, `sta_cm_with_bot`, `sta_global`, `sta_global_with_bot`, `ent`, `ent_with_bot`, `ent_cm`, `ent_cm_with_bot`, `ent_global`, `ent_global_with_bot`.
- Computed:
  - `resource_names` (TypeList, element TypeString): List of purchased resource names returned by the API.
  - `deal_names` (TypeList, element TypeString): List of purchased order/deal names returned by the API.

### CRUD mapping (operation-type, following `tencentcloud_teo_identify_zone_operation`)
- **Create**: build `CreatePlanForZoneRequest` from `zone_id` and `plan_type`; call the service-layer `TeoPlanForZone` helper inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`; on success `d.SetId(helper.BuildToken())`; then nil-safe set computed `resource_names` and `deal_names`; finally call Read.
- **Read**: no-op, return nil (no query API for the purchase result).
- **Delete**: no-op, return nil (state removal only).

### Service-layer method
- Add `TeoPlanForZone(zoneId, planType string) (resourceNames, dealNames []*string, errRet error)` to `service_tencentcloud_teo.go`, following the `TeoIdentifyZone` pattern: build request, ratelimit check, call `me.client.UseTeoV20220901Client().CreatePlanForZone(request)`, nil-safe access to `response.Response`.

### Retry placement
- Per the business rules, the retry wraps only the API call. The `d.SetId()` and `d.Set(...)` calls happen after the retry block (outside it), in the success path after the retry error handling.

## Risks / Trade-offs

- [Plan purchase is irreversible and account/billing-impacting] → Mitigation: document clearly that Delete only removes the resource from state and does not cancel the plan. The `plan_type` is ForceNew, so changing it forces recreation (a new purchase).
- [No Read API to verify the purchase] → Mitigation: Read is a no-op; the computed `resource_names` and `deal_names` are set once at create time from the API response and persisted in state.
- [Re-running create on a zone that already has a plan] → Mitigation: rely on the API's own behavior; `CreatePlanForZone` will surface an error (e.g., `InvalidParameter.ZoneHasBeenBound`) which the retry helper surfaces to the user.
