## 1. Service Layer

- [x] 1.1 Add `CreateService(ctx context.Context, payMode, clientToken string) (response *ssm.CreateServiceResponse, errRet error)` method to `SsmService` in `tencentcloud/services/ssm/service_tencentcloud_ssm.go`
  - Construct `ssm.NewCreateServiceRequest()`; set `request.PayMode = helper.String(payMode)`, `request.ClientToken = helper.String(clientToken)`
  - Wrap the call in `resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {...})`
  - Inside retry: `ratelimit.Check(request.GetAction())`, call `me.client.UseSsmClient().CreateServiceWithContext(ctx, request)`
  - On error: return `tccommon.RetryError(e)`
  - On nil response/response.Response: return `resource.NonRetryableError(fmt.Errorf("CreateService failed, Response is nil"))`
  - On success: `log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())`
  - Defer error logging: `log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())`
  - Return `*ssm.CreateServiceResponse`
- [x] 1.2 Add `GetServiceStatus(ctx context.Context) (response *ssm.GetServiceStatusResponseParams, errRet error)` method to `SsmService` in `tencentcloud/services/ssm/service_tencentcloud_ssm.go`
  - Construct `ssm.NewGetServiceStatusRequest()` (empty)
  - Wrap the call in `resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {...})`
  - Inside retry: `ratelimit.Check(request.GetAction())`, call `me.client.UseSsmClient().GetServiceStatusWithContext(ctx, request)`
  - On error: return `tccommon.RetryError(e)`
  - On nil response/response.Response: return `resource.NonRetryableError(fmt.Errorf("GetServiceStatus failed, Response is nil"))`
  - On success: `log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", ...)`; set `response = result.Response`
  - Return `*ssm.GetServiceStatusResponseParams`

## 2. Operation Resource Implementation

- [x] 2.1 Create `tencentcloud/services/ssm/resource_tc_ssm_service_operation.go` following `resource_tc_dnspod_modify_record_group_operation.go` / `resource_tc_sts_assume_role_operation.go`
  - Package `ssm`, imports: `context`, `fmt`, `log`, `github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource`, `github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema`, `tccommon`, `helper`
  - Factory `ResourceTencentCloudSsmServiceOperation() *schema.Resource` returns a `*schema.Resource` with `Create`/`Read`/`Delete` handlers (no `Update`)
  - Schema: `pay_mode` (Required, ForceNew, TypeString); `client_token` is NOT a schema field — generated internally via `helper.BuildToken()`
  - `Create`:
    - Read `pay_mode` from resource data; generate `clientToken := helper.BuildToken()`
    - `ssmService := SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}`
    - `createResp, err := ssmService.CreateService(ctx, payMode, clientToken)`; on err log and return
    - Log `createResp.Response.ResourceId`/`OrderId` (nil-safe) via `log.Printf("[DEBUG]...")`
    - Set `d.SetId(helper.BuildToken())` BEFORE polling (async resource — retain id even if polling fails)
    - Poll `GetServiceStatus` inside a `resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {...})`:
      - `status, e := ssmService.GetServiceStatus(ctx)`; on e return `tccommon.RetryError(e)`
      - if `status == nil || status.ServiceEnabled == nil` return `resource.NonRetryableError(...)`
      - if `*status.ServiceEnabled == true` → log enabled status, return nil (done)
      - else `log.Printf("[DEBUG] ...keep polling")`; return `resource.RetryableError(...)`
    - On poll-retry exhausted (reqErr != nil): log and return reqErr
    - Call `resourceTencentCloudSsmServiceOperationRead(d, meta)` and return
  - `Read`: no-op (return nil)
  - `Delete`: no-op (return nil)
  - Use `ssm_service_operation` (snake_case) in all log/error messages per convention

## 3. Resource Registration

- [x] 3.1 Edit `tencentcloud/provider.go`
  - Add `"tencentcloud_ssm_service_operation": ssm.ResourceTencentCloudSsmServiceOperation()` to `ResourcesMap`
- [x] 3.2 Edit `tencentcloud/framework/registry.go`
  - Remove `ssm.NewSsmServiceAction` from the `actionFactories` slice (replaced by SDKv2 operation resource)

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/ssm/resource_tc_ssm_service_operation.md` following `resource_tc_dnspod_modify_record_group_operation.md` / `resource_tc_sts_assume_role_operation.md` format
  - One-line description mentioning SSM (Secrets Manager)
  - Example Usage block: `resource "tencentcloud_ssm_service_operation" "example" { pay_mode = "PostPay" }`
  - No Import section, no manual Argument/Attribute Reference sections
- [ ] 4.2 Run `make doc` (in finalize phase) to generate website documentation

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/ssm/resource_tc_ssm_service_operation_test.go` using gomonkey mocks (no acceptance suite)
  - Mock meta implementing `tccommon.ProviderMeta` returning a `*connectivity.TencentCloudClient`
  - Mock `UseSsmClient` (method-level via `ApplyMethodReturn`), `CreateServiceWithContext`, `GetServiceStatusWithContext`
  - `TestSsmServiceOperation_Create_Success`: CreateService ok + GetServiceStatus returns ServiceEnabled==true; assert no error, id set, PayMode=="PostPay", ClientToken non-empty
  - `TestSsmServiceOperation_Create_PollNotYetEnabled`: first GetServiceStatus ServiceEnabled==false, second true; assert no error, at least two polls
  - `TestSsmServiceOperation_Create_APIError`: CreateServiceWithContext returns error; assert error contains message
  - `TestSsmServiceOperation_Read`: assert no error
  - `TestSsmServiceOperation_Delete`: assert no error
  - `TestSsmServiceOperation_Schema`: assert Create/Read/Delete set, Update nil, pay_mode Required/ForceNew, client_token absent

## 6. Code Quality (Finalize Phase)

- [ ] 6.1 Run `gofmt` to format all modified Go files (via tfpacer-finalize skill)
- [x] 6.2 Verify all error returns are properly handled; nil-safe access on every response field before logging/reading
- [ ] 6.3 Update changelog file via `tfpacer-finalize` skill (resource name `tencentcloud_ssm_service_operation`)