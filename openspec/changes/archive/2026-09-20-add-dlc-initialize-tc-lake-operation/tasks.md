## 1. Service Layer

- [x] 1.1 Add `InitializeTCLake(ctx context.Context) (*dlc.InitializeTCLakeResponse, error)` method to `DlcService` in `tencentcloud/services/dlc/service_tencentcloud_dlc.go`
  - Construct `dlc.NewInitializeTCLakeRequest()` (empty request, no parameters)
  - Wrap the call in `resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {...})`
  - Inside retry: `ratelimit.Check(request.GetAction())`, call `me.client.UseDlcClient().InitializeTCLakeWithContext(ctx, request)`
  - On error: return `tccommon.RetryError(e)`
  - On success: `log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())`
  - Defer error logging: `log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", ...)`
  - Return the `*dlc.InitializeTCLakeResponse` to the caller

## 2. Action Implementation

- [x] 2.1 Create `tencentcloud/services/dlc/action_tc_dlc_initialize_tc_lake.go` following `action_tc_teo_confirm_origin_acl_update.go` / `action_tc_bdrc_run_copy_pair_tasks.go`
  - Package `dlc`, imports: `context`, `fmt`, `github.com/hashicorp/terraform-plugin-framework/action`, `github.com/hashicorp/terraform-plugin-framework/action/schema`, `github.com/hashicorp/terraform-provider-tencentcloud/tencentcloud/framework/fw`
  - Compile-time assertion: `var _ action.ActionWithConfigure = &DlcInitializeTCLake{}`
  - Factory `NewDlcInitializeTCLake() action.Action` returns `&DlcInitializeTCLake{}`
  - Struct `DlcInitializeTCLake` embedding `fw.ActionWithConfigure`
  - Empty model `DlcInitializeTCLakeModel` (kept for symmetry with references)
  - `Metadata`: set `resp.TypeName = "tencentcloud_dlc_initialize_tc_lake"`
  - `Schema`: `Description` mentioning DLC TCLake initialization, empty `Attributes: map[string]schema.Attribute{}`
  - `Invoke`:
    - Decode config into model (empty), check diagnostics
    - Guard `a.Client() == nil` → add "Provider not configured" diagnostic, return
    - Call `service := NewDlcService(a.Client()); response, err := service.InitializeTCLake(ctx)`
    - On err: `resp.Diagnostics.AddError(fmt.Sprintf("initializing DLC TCLake"), err.Error())`, return
    - On success: return (no state/output written)

## 3. Framework Registration

- [x] 3.1 Edit `tencentcloud/framework/registry.go`
  - Add import `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dlc"`
  - Add `dlc.NewDlcInitializeTCLake` to the `actionFactories` slice
  - Do NOT add anything to `tencentcloud/provider.go` `ResourcesMap`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/dlc/action_tc_dlc_initialize_tc_lake.md` following `action_tc_bdrc_run_copy_pair_tasks.md` format
  - One-line description mentioning DLC
  - NOTE about Terraform 1.14+ support
  - Example Usage block: `action "tencentcloud_dlc_initialize_tc_lake" "example" { config {} }`
  - No Import section, no manual Argument/Attribute Reference sections
- [ ] 4.2 Run `make doc` (in finalize phase) to generate website documentation

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/dlc/action_tc_dlc_initialize_tc_lake_test.go` following `action_tc_bdrc_run_copy_pair_tasks_test.go` (gomonkey mocks, no acceptance suite)
  - Helper to build an `action.InvokeRequest` with an empty config object (no attributes)
  - Helper `setDlcInitializeTCLakeClient` to inject a mock client via `Configure`
  - Test `TestDlcInitializeTCLake_Invoke_Success`: mock `UseDlcClient` + `InitializeTCLakeWithContext` returning a successful response; assert no diagnostics
  - Test `TestDlcInitializeTCLake_Invoke_APIError`: mock `InitializeTCLakeWithContext` returning an error; assert error diagnostic contains the API error
  - Test `TestDlcInitializeTCLake_Invoke_ClientNotConfigured`: invoke without Configure; assert "Provider not configured" diagnostic
  - Test `TestDlcInitializeTCLake_MetadataAndSchema`: assert `TypeName == "tencentcloud_dlc_initialize_tc_lake"`, attributes map is empty, no diagnostics

## 6. Code Quality (Finalize Phase)

- [ ] 6.1 Run `gofmt` to format all modified Go files
- [x] 6.2 Verify all error returns are properly handled (use `_ =` for functions that cannot fail)
- [ ] 6.3 Create changelog file via `tfpacer-finalize` skill