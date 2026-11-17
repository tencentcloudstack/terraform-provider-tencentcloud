## 1. Vendor Preparation

- [x] 1.1 Bump the vendored SSM SDK so `CreateService` is available (it does not exist in `ssm v1.3.164`)
  - `go.mod`: `tencentcloud-sdk-go/tencentcloud/ssm` `v1.3.164` → `v1.3.199`, `tencentcloud-sdk-go/tencentcloud/common` `v1.3.194` → `v1.3.199`
  - `go.sum` + `vendor/modules.txt` updated accordingly
  - `vendor/.../ssm/v20190923/{client,models,errors}.go` now contain `CreateServiceWithContext` / `CreateServiceRequest` / `CreateServiceResponse`

## 2. Service Layer

- [x] 2.1 Add `CreateService(ctx context.Context, payMode, clientToken string) (response *ssm.CreateServiceResponse, errRet error)` method to `SsmService` in `tencentcloud/services/ssm/service_tencentcloud_ssm.go`
  - Construct `ssm.NewCreateServiceRequest()`; set `request.PayMode = helper.String(payMode)`, `request.ClientToken = helper.String(clientToken)`
  - Wrap the call in `resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {...})`
  - Inside retry: `ratelimit.Check(request.GetAction())`, call `me.client.UseSsmClient().CreateServiceWithContext(ctx, request)`
  - On error: return `tccommon.RetryError(e)`
  - On success: `log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())`
  - Defer error logging: `log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())`
  - On nil response/response.Response: return `resource.NonRetryableError(fmt.Errorf("CreateService failed, Response is nil"))`
  - Return `*ssm.CreateServiceResponse`
- [x] 2.2 Add `GetServiceStatus(ctx context.Context) (response *ssm.GetServiceStatusResponseParams, errRet error)` method to `SsmService` in `tencentcloud/services/ssm/service_tencentcloud_ssm.go`
  - Construct `ssm.NewGetServiceStatusRequest()` (empty)
  - Wrap the call in `resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {...})`
  - Inside retry: `ratelimit.Check(request.GetAction())`, call `me.client.UseSsmClient().GetServiceStatusWithContext(ctx, request)`
  - On error: return `tccommon.RetryError(e)`
  - On nil response/response.Response: return `resource.NonRetryableError(fmt.Errorf("GetServiceStatus failed, Response is nil"))`
  - On success: `log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", ...)`; set `response = result.Response`
  - Return `*ssm.GetServiceStatusResponseParams`

## 3. Action Implementation

- [x] 3.1 Create `tencentcloud/services/ssm/action_tc_ssm_service.go` following `action_tc_dlc_initialize_tc_lake.go` / `action_tc_antiddos_unblock_resources.go`
  - Package `ssm`, imports: `context`, `fmt`, `log`, `github.com/hashicorp/terraform-plugin-framework/action`, `github.com/hashicorp/terraform-plugin-framework/action/schema`, `github.com/hashicorp/terraform-plugin-framework/types`, `github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource`, `tccommon`, `fw`, `helper`
  - Compile-time assertion: `var _ action.ActionWithConfigure = &SsmServiceAction{}`
  - Factory `NewSsmServiceAction() action.Action` returns `&SsmServiceAction{}`
  - Struct `SsmServiceAction` embedding `fw.ActionWithConfigure`
  - Model `SsmServiceActionModel` with only `PayMode types.String` (tfsdk `pay_mode`); `client_token` is NOT part of the model/schema
  - `Metadata`: set `resp.TypeName = "tencentcloud_ssm_service"`
  - `Schema`: `Description` mentioning SSM service opening; `Attributes` map with only `pay_mode` (Required StringAttribute)
  - `Invoke`:
    - Decode config into model, check diagnostics
    - Validate `pay_mode` (null/unknown/empty → `resp.Diagnostics.AddError("Missing pay_mode", ...)`, return)
    - Guard `a.Client() == nil` → AddError "Provider not configured", return
    - Generate `clientToken := helper.BuildToken()` (not read from config)
    - `service := NewSsmService(a.Client())`
    - `createResp, err := service.CreateService(ctx, payMode, clientToken)`; on err `resp.Diagnostics.AddError("opening SSM service", err.Error())`, return
    - Log `createResp.Response.ResourceId`/`OrderId` (nil-safe) via `log.Printf("[DEBUG]...")`
    - Poll `GetServiceStatus` inside a `resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {...})`:
      - `status, e := service.GetServiceStatus(ctx)`; on e return `tccommon.RetryError(e)`
      - if `status == nil || status.ServiceEnabled == nil` return `resource.NonRetryableError(fmt.Errorf("ssm_service status ServiceEnabled is nil"))`
      - if `*status.ServiceEnabled == true` → log enabled status, return nil (done)
      - else `log.Printf("[DEBUG] ssm_service still not enabled, keep polling")`; return `resource.RetryableError(fmt.Errorf("ssm_service not enabled yet"))`
    - On poll-retry exhausted (reqErr != nil): `resp.Diagnostics.AddError("polling SSM service status", reqErr.Error())`, return
    - On success: return (no state/output written)

## 4. Framework Registration & Index

- [x] 4.1 Edit `tencentcloud/framework/registry.go`
  - Add `ssm.NewSsmServiceAction` to the `actionFactories` slice
  - Do NOT add a new `ssm` import (already present); do NOT edit `tencentcloud/provider.go`
- [x] 4.2 Edit `tencentcloud/provider.md`
  - Add `tencentcloud_ssm_service` under the SSM product node's `Action` sub-section

## 5. Documentation

- [x] 5.1 Create `tencentcloud/services/ssm/action_tc_ssm_service.md` following `action_tc_dlc_initialize_tc_lake.md` / `action_tc_antiddos_unblock_resources.md` format
  - One-line description mentioning SSM
  - NOTE about Terraform 1.14+ support
  - Example Usage block: `action "tencentcloud_ssm_service" "example" { config { pay_mode = "PostPay" } }` (no `client_token`)
  - No Import section, no manual Argument/Attribute Reference sections
- [ ] 5.2 Run `make doc` (in finalize phase) to generate `website/docs/actions/ssm_service.html.markdown` and the `tencentcloud.erb` action index entry

## 6. Unit Tests

- [x] 6.1 Create `tencentcloud/services/ssm/action_tc_ssm_service_test.go` following `action_tc_antiddos_unblock_resources_test.go` (gomonkey mocks, no acceptance suite)
  - Helper `ssmServiceActionSchema(t)` to build the action schema
  - Helper `newSsmServiceActionInvokeRequest(t, payMode)` building an `action.InvokeRequest` carrying the single `pay_mode` String tftype (empty string → null, for the missing-input test)
  - Helper `setSsmServiceActionClient(a, client)` injecting a mock client via `Configure` (`sharedmeta.ProviderMeta{Client: client}`)
  - `TestSsmServiceAction_Invoke_Success`: mock `UseSsmClient` + `CreateServiceWithContext` returning a successful response and `GetServiceStatusWithContext` returning `ServiceEnabled == true`; assert no diagnostics, `PayMode == "PostPay"`, generated `ClientToken` non-empty
  - `TestSsmServiceAction_Invoke_PollNotYetEnabled`: first `GetServiceStatusWithContext` returns `ServiceEnabled == false`, second returns `ServiceEnabled == true`; assert no diagnostics and ≥2 polls
  - `TestSsmServiceAction_Invoke_MissingPayMode`: null `pay_mode`; assert "Missing pay_mode" diagnostic and no API call
  - `TestSsmServiceAction_Invoke_CreateServiceError`: `CreateServiceWithContext` returns error; assert "opening SSM service" diagnostic
  - `TestSsmServiceAction_Invoke_PollError`: `GetServiceStatusWithContext` returns a non-retryable error; assert "polling SSM service status" diagnostic
  - `TestSsmServiceAction_Invoke_ClientNotConfigured`: invoke without Configure; assert "Provider not configured" diagnostic
  - `TestSsmServiceAction_MetadataAndSchema`: assert `TypeName == "tencentcloud_ssm_service"`, attributes map contains exactly `pay_mode` (Required) and no `client_token`, no output attributes, no diagnostics

## 7. Code Quality (Finalize Phase)

- [ ] 7.1 Run `gofmt` to format all modified Go files (via tfpacer-finalize skill)
- [x] 7.2 Verify all error returns are properly handled (use `_ =` for functions that cannot fail); nil-safe access on every response field before logging/reading
- [x] 7.3 Keep `.changelog/4611.txt` aligned with the action implementation: `release-note:new-resource` listing `tencentcloud_ssm_service` (the repo convention for framework actions, see `.changelog/4528.txt`/`.changelog/4529.txt`/`.changelog/4534.txt`)
- [ ] 7.4 Re-run `make doc` / changelog verification in the finalize phase
