## Context

CLS exposes two account-level APIs in the vendored SDK `cls/v20201016` (both verified present — no SDK upgrade needed):

- `OpenClsService()` — activates CLS for the account. **Request has no input parameters**; response only contains `RequestId`.
- `GetClsService()` — returns `Status *int64` (0: service opened, 1: service not opened). **Request has no input parameters**.

This is a **one-time operation resource**, following the style of `tencentcloud_dlc_update_row_filter_operation`: it has `Create`, `Read`, `Delete` (no `Update`), performs the action on Create, and treats Delete as a state-only no-op.

## Goals / Non-Goals

**Goals:**
- Trigger CLS activation via `OpenClsService` on create, with retry.
- Auto-generate the resource ID (per business rule), using the CLS-package convention `helper.BuildToken()`.
- Surface the activation status via `GetClsService` as a computed attribute, with nil-safe access and retry.
- Ship resource example `.md`, website docs, unit test, and provider registration.

**Non-Goals:**
- No de-activation on delete (the API has no such capability); Delete is a no-op.
- No required user-facing input arguments (both APIs take none).
- No Timeouts block (synchronous calls, no async status polling needed).

## Decisions

### Resource name and file layout
- Resource: `tencentcloud_cls_open_service_operation`
- Files under `tencentcloud/services/cls/`:
  - `resource_tc_cls_open_service_operation.go`
  - `resource_tc_cls_open_service_operation.md`
  - `resource_tc_cls_open_service_operation_test.go`
- Website doc: `website/docs/r/cls_open_service_operation.html.markdown`
- Constructor `ResourceTencentCloudClsOpenServiceOperation()`, registered in `provider.go`.

### Schema
- Both `OpenClsService` and `GetClsService` requests have no input parameters, so there are no required/optional arguments.
- `status` (TypeInt, Computed) ← from `GetClsService.Status`, description "Account service status. 0: opened, 1: not opened.".

### CRUD mapping (operation-type, following `tencentcloud_dlc_update_row_filter_operation`)
- **Create**: call `OpenClsService` inside `resource.Retry(WriteRetryTimeout, ...)`; on success `d.SetId(helper.BuildToken())`; then call Read.
- **Read**: call `GetClsService` inside `resource.Retry(ReadRetryTimeout, ...)`; nil-safe set `status` from the response.
- **Delete**: no-op, return nil (state removal only).

### Client method
- Use `UseClsV20201016Client()` (the versioned CLS client used by recent CLS resources) with the `WithContext` variants and a `ctx` derived from `tccommon.NewResourceLifeCycleHandleFuncContext`.

## Risks / Trade-offs

- [Activation is account-global and irreversible] → Mitigation: document clearly that Delete only removes the resource from state and does not deactivate CLS.
- [`OpenClsService` is idempotent-ish: re-running on an already-open account] → Mitigation: rely on the API's own behavior; create remains safe to re-apply. If it errors on already-open, the retry helper surfaces the error to the user.
- [No user inputs / computed-only resource] → acceptable for a one-time activation operation; the computed `status` gives a readable signal.

## Migration Plan

Purely additive; no migration. New resource registered alongside existing CLS resources. Rollback = revert the additive files and provider registration.
