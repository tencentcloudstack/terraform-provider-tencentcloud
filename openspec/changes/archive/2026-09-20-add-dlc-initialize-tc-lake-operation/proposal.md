## Why

The DLC (Data Lake Compute) cloud API exposes an `InitializeTCLake` operation that activates (opens) the TCLake service for the current account. This is a one-time, idempotent initialization step that practitioners currently cannot trigger through Terraform. Adding a framework-based action resource lets users provision TCLake initialization as part of their IaC workflow alongside other DLC resources.

## What Changes

- Add a new Terraform Plugin Framework **action** resource `tencentcloud_dlc_initialize_tc_lake` implemented via the `action.Action` interface (not the SDKv2 operation resource pattern).
- Implement the `Invoke` handler to call the DLC `InitializeTCLake` API (`InitializeTCLakeWithContext`) with an empty request (the API takes no input parameters).
- Capture the API response fields `instance_id` (`response.Response.InstanceId`) and `is_success` (`response.Response.IsSuccess`) and log them via `log.Printf` for observability. The Terraform Plugin Framework action schema does NOT support Computed/output attributes (`InvokeResponse` only carries `Diagnostics` and `SendProgress`), so no output attributes are exposed in the action schema. This matches the reference implementations `NewTeoConfirmOriginAclUpdate` and `NewBdrcRunCopyPairTasks`, which invoke an API and surface only success/error diagnostics.
- Register the action factory `dlc.NewDlcInitializeTCLake` in `tencentcloud/framework/registry.go`'s `actionFactories` slice. Do NOT register it in the SDKv2 `provider.go` `ResourcesMap`.
- Add a `InitializeTCLake` method to the DLC service layer (`service_tencentcloud_dlc.go`) wrapping the API call with `tccommon.WriteRetryTimeout` retry and `tccommon.RetryError` error wrapping. The method returns the `*InitializeTCLakeResponse` so the Invoke handler can log the `InstanceId` and `IsSuccess` fields (but they are not persisted as Terraform state).
- Provide a markdown documentation file `action_tc_dlc_initialize_tc_lake.md` with a one-line description mentioning DLC, an Example Usage section using the framework action block syntax, and a NOTE about Terraform 1.14+ support.
- Add unit tests in `action_tc_dlc_initialize_tc_lake_test.go` using gomonkey mocks (no Terraform acceptance test suite), covering successful invoke, API error, and schema/metadata validation.

## Capabilities

### New Capabilities
- `dlc-initialize-tc-lake-action`: Adds a framework action `tencentcloud_dlc_initialize_tc_lake` that activates TCLake via the DLC `InitializeTCLake` API. This is a one-time operation; no cloud-side state is persisted and no output attributes are exposed (the API response `InstanceId`/`IsSuccess` are logged for observability only).

### Modified Capabilities

None. There is no existing spec for a DLC TCLake initialization resource, so this is a pure addition.

## Impact

### Affected Code
- `tencentcloud/services/dlc/action_tc_dlc_initialize_tc_lake.go` - New action implementation (factory, model, schema, Invoke)
- `tencentcloud/services/dlc/service_tencentcloud_dlc.go` - New `InitializeTCLake` service method
- `tencentcloud/services/dlc/action_tc_dlc_initialize_tc_lake_test.go` - New unit tests
- `tencentcloud/services/dlc/action_tc_dlc_initialize_tc_lake.md` - New action documentation
- `tencentcloud/framework/registry.go` - Register `dlc.NewDlcInitializeTCLake` in `actionFactories` and add `dlc` import

### API References
- **DLC API** (`dlc/v20210125`):
  - `InitializeTCLake`: activates/open TCLake. Takes no request parameters. Returns `response.Response.InstanceId` (实例Id, `*string`) and `response.Response.IsSuccess` (是否成功, `*bool`).

### Backward Compatibility
- **Fully Compatible**: This is a pure addition of a new action resource. No existing resources, data sources, or schemas are modified. No state migration is required.