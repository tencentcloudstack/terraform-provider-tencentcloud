## 1. Service Layer

- [x] 1.1 Add `CreateService(ctx context.Context, payMode, clientToken string) (response *ssm.CreateServiceResponse, errRet error)` method to `SsmService` in `tencentcloud/services/ssm/service_tencentcloud_ssm.go`
  - Construct `ssm.NewCreateServiceRequest()`; set `request.PayMode = helper.String(payMode)`, `request.ClientToken = helper.String(clientToken)`
  - Wrap the call in `resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {...})`
  - Inside retry: `ratelimit.Check(request.GetAction())`, call `me.client.UseSsmClient().CreateServiceWithContext(ctx, request)`
  - On error: return `tccommon.RetryError(e)`
  - On success: `log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())`
  - Defer error logging: `log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())`
  - On nil response/response.Response: return `resource.NonRetryableError(fmt.Errorf("CreateService failed, Response is nil"))`
  - Return `*ssm.CreateServiceResponse`
- [x] 1.2 Add `GetServiceStatus(ctx context.Context) (response *ssm.GetServiceStatusResponseParams, errRet error)` method to `SsmService` in `tencentcloud/services/ssm/service_tencentcloud_ssm.go`
  - Construct `ssm.NewGetServiceStatusRequest()` (empty)
  - Wrap the call in `resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {...})`
  - Inside retry: `ratelimit.Check(request.GetAction())`, call `me.client.UseSsmClient().GetServiceStatusWithContext(ctx, request)`
  - On error: return `tccommon.RetryError(e)`
  - On nil response/response.Response: return `resource.NonRetryableError(fmt.Errorf("GetServiceStatus failed, Response is nil"))`
  - On success: `log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", ...)`; set `response = result.Response`
  - Return `*ssm.GetServiceStatusResponseParams`

## 2. Action Implementation

- [x] 2.1 Create `tencentcloud/services/ssm/action_tc_ssm_service.go` following `action_tc_dlc_initialize_tc_lake.go` / `action_tc_antiddos_unblock_resources.go`
  - Package `ssm`, imports: `context`, `fmt`, `github.com/hashicorp/terraform-plugin-framework/action`, `github.com/hashicorp/terraform-plugin-framework/action/schema`, `github.com/hashicorp/terraform-plugin-framework/types`, `log`, sdk `ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"`, `tccommon`, `fw`, `helper`, `connectivity`, `resource`, `ratelimit`
  - Compile-time assertion: `var _ action.ActionWithConfigure = &SsmServiceAction{}`
  - Factory `NewSsmServiceAction() action.Action` returns `&SsmServiceAction{}`
  - Struct `SsmServiceAction` embedding `fw.ActionWithConfigure`
  - Model `SsmServiceActionModel` with `PayMode types.String` / `ClientToken types.String` (tfsdk `pay_mode` / `client_token`)
  - `Metadata`: set `resp.TypeName = "tencentcloud_ssm_service"`
  - `Schema`: `Description` mentioning SSM service opening; `Attributes` map with `pay_mode` (Required StringAttribute) and `client_token` (Required StringAttribute)
  - `Invoke`:
    - Decode config into model, check diagnostics
    - Validate `pay_mode` (null/unknown/empty → `resp.Diagnostics.AddError("Missing pay_mode", ...)`, return)
    - Validate `client_token` (null/unknown/empty → `resp.Diagnostics.AddError("Missing client_token", ...)`, return)
    - Guard `a.Client() == nil` → AddError "Provider not configured", return
    - `service := NewSsmService(a.Client())`
    - `createResp, err := service.CreateService(ctx, payMode, clientToken)`; on err `resp.Diagnostics.AddError(fmt.Sprintf("opening SSM service"), err.Error())`, return
    - Log `createResp.Response.ResourceId`/`OrderId` (nil-safe) via `log.Printf("[DEBUG]...")`
    - Poll `GetServiceStatus` inside a `resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {...})`:
      - `status, e := service.GetServiceStatus(ctx)`; on e return `tccommon.RetryError(e)`
      - if `status.ServiceEnabled == nil` return `resource.NonRetryableError(fmt.Errorf("SSM service status ServiceEnabled is nil"))`
      - if `*status.ServiceEnabled == true` → log enabled status, return nil (done)
      - else `log.Printf("[DEBUG] ssm service still not enabled, keep polling")`; return `resource.RetryableError(fmt.Errorf("SSM service not enabled yet"))`
    - On poll-retry exhausted (reqErr != nil): `resp.Diagnostics.AddError(fmt.Sprintf("polling SSM service status"), reqErr.Error())`, return
    - On success: return (no state/output written)

## 3. Framework Registration

- [x] 3.1 Edit `tencentcloud/framework/registry.go`
  - Add `ssm.NewSsmServiceAction` to the `actionFactories` slice
  - Do NOT add a new `ssm` import (already present); do NOT edit `tencentcloud/provider.go`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/ssm/action_tc_ssm_service.md` following `action_tc_dlc_initialize_tc_lake.md` / `action_tc_antiddos_unblock_resources.md` format
  - One-line description mentioning SSM
  - NOTE about Terraform 1.14+ support
  - Example Usage block: `action "tencentcloud_ssm_service" "example" { config { pay_mode = "PostPay" client_token = "ssm-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" } }`
  - No Import section, no manual Argument/Attribute Reference sections
- [ ] 4.2 Run `make doc` (in finalize phase) to generate website documentation

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/ssm/action_tc_ssm_service_test.go` following `action_tc_antiddos_unblock_resources_test.go` (gomonkey mocks, no acceptance suite)
  - Helper `ssmServiceActionSchema(t)` to build the action schema
  - Helper `newSsmServiceActionInvokeRequest(t, payMode, clientToken)` building an `action.InvokeRequest` carrying the two String tftypes (and a null variant for missing-input tests)
  - Helper `setSsmServiceActionClient(a, client)` injecting a mock client via `Configure` (`sharedmeta.ProviderMeta{Client: client}`)
  - `TestSsmServiceAction_Invoke_Success`: mock `UseSsmClient` + `CreateServiceWithContext` returning a successful response and `GetServiceStatusWithContext` returning `ServiceEnabled == true`; assert no diagnostics
  - `TestSsmServiceAction_Invoke_PollNotYetEnabled`: first `GetServiceStatusWithContext` returns `ServiceEnabled == false`, second returns `ServiceEnabled == true`; assert no diagnostics
  - `TestSsmServiceAction_Invoke_MissingPayMode`: null `pay_mode`; assert "Missing pay_mode" diagnostic
  - `TestSsmServiceAction_Invoke_MissingClientToken`: null `client_token`; assert "Missing client_token" diagnostic
  - `TestSsmServiceAction_Invoke_CreateServiceError`: `CreateServiceWithContext` returns error; assert error diagnostic
  - `TestSsmServiceAction_Invoke_ClientNotConfigured`: invoke without Configure; assert "Provider not configured" diagnostic
  - `TestSsmServiceAction_MetadataAndSchema`: assert `TypeName == "tencentcloud_ssm_service"`, attributes map contains exactly `pay_mode` + `client_token` (both Required), no diagnostics

## 6. Code Quality (Finalize Phase)

- [ ] 6.1 Run `gofmt` to format all modified Go files (via tfpacer-finalize skill)
- [x] 6.2 Verify all error returns are properly handled (use `_ =` for functions that cannot fail); nil-safe access on every response field before logging/reading
- [ ] 6.3 Create changelog file via `tfpacer-finalize` skill