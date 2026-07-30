## ADDED Requirements

### Requirement: Framework resources and data sources SHALL follow a unified file naming convention

Each framework resource SHALL live at `tencentcloud/services/<service>/resource_tc_<service>_<name>_framework.go`, and each framework data source at `tencentcloud/services/<service>/data_source_tc_<service>_<name>_framework.go`. The `_framework` suffix is mandatory to distinguish from SDKv2 implementations in the same directory. Acceptance tests SHALL use `..._framework_test.go`.

The Terraform-facing type name SHALL remain `tencentcloud_<service>_<name>` (no suffix), so that user HCL is identical regardless of stack.

#### Scenario: New framework resource is named correctly
- **WHEN** a developer adds a framework resource `tencentcloud_foo_bar`
- **THEN** the file is placed at `tencentcloud/services/foo/resource_tc_foo_bar_framework.go`
- **AND** the test file is `resource_tc_foo_bar_framework_test.go`
- **AND** the resource's `Metadata` returns `TypeName = "tencentcloud_foo_bar"`

### Requirement: Each service package SHALL expose framework registration entry points

Each service package that contributes framework artifacts SHALL provide a `framework.go` file that exports two functions:

```
func FrameworkResources() []func() resource.Resource
func FrameworkDataSources() []func() datasource.DataSource
```

These functions SHALL be the only place where framework resources are listed for that service. If a service has no framework artifacts, the file MAY be omitted or return empty slices.

#### Scenario: Provider aggregates per-service registrations
- **WHEN** `FrameworkProvider.Resources` is called
- **THEN** it iterates known service packages and concatenates each service's `FrameworkResources()` slice
- **AND** the resulting slice contains every framework resource exactly once

#### Scenario: Service with no framework artifacts compiles cleanly
- **WHEN** a service has only SDKv2 resources
- **THEN** the absence of `framework.go` SHALL NOT cause build failures
- **AND** the provider aggregator skips that service silently

### Requirement: Framework resources SHALL receive client through ProviderData type assertion

Each framework resource and data source SHALL implement `Configure(ctx, req, resp)` by performing a type assertion on `req.ProviderData` to `*fwtransport.ProviderData`. On success, the resource stores the embedded `*connectivity.TencentCloudClient`. On failure (nil or wrong type), it MUST append an unexpected-provider-data diagnostic and return without panic.

#### Scenario: Resource Configure handles ProviderData correctly
- **WHEN** the framework runtime calls `Configure` with a valid `*fwtransport.ProviderData`
- **THEN** the resource's internal client field is set
- **AND** no diagnostics are appended

#### Scenario: Resource Configure handles nil ProviderData defensively
- **WHEN** the framework runtime calls `Configure` with `req.ProviderData == nil` (early lifecycle phase)
- **THEN** the resource MUST return early without diagnostics
- **AND** Terraform's normal lifecycle re-invokes Configure later

#### Scenario: Resource Configure rejects wrong type
- **WHEN** `req.ProviderData` is non-nil but not `*fwtransport.ProviderData`
- **THEN** the resource appends an `Error` diagnostic with summary "Unexpected Provider Data Type"
- **AND** does not dereference the pointer

### Requirement: Framework resources SHALL use the shared fwhelper utilities for retries, timeouts, and error translation

Framework resources SHALL NOT re-implement retry, timeout parsing, or SDK error-to-diagnostic translation. They MUST use the helpers in `tencentcloud/internal/fwhelper`:

- `fwhelper.RetryFramework(ctx, timeout, fn)` for eventual-consistency retries.
- `terraform-plugin-framework-timeouts` (re-exported via fwhelper) for Create/Read/Update/Delete timeout blocks.
- `fwhelper.WrapSDKError(err)` to convert `*sdkErrors.TencentCloudSDKError` into structured `diag.Diagnostic` values.
- `fwhelper.StringValueOrNull` / `Int64ValueOrNull` etc. for `*T` ↔ `types.T` conversions in Read paths.

#### Scenario: Resource retries on eventual consistency without bespoke code
- **WHEN** an API returns `ResourceNotFound` immediately after creation
- **THEN** the resource calls `fwhelper.RetryFramework(ctx, timeouts.Read, ...)` to retry
- **AND** does not inline its own retry loop

#### Scenario: SDK errors are surfaced as structured diagnostics
- **WHEN** a TencentCloud SDK call returns `*sdkErrors.TencentCloudSDKError` with a request id
- **THEN** the resource appends the diagnostic produced by `fwhelper.WrapSDKError(err)`
- **AND** the diagnostic's detail includes the request id and SDK error code

### Requirement: Framework resources SHALL ship with acceptance tests using ProtoV5ProviderFactories

Acceptance tests for framework resources SHALL configure `resource.TestCase` with `ProtoV5ProviderFactories` (NOT `Providers` or `ProviderFactories`), pointing at the muxed provider so the test exercises the same surface as production.

#### Scenario: Acceptance test runs against muxed provider
- **WHEN** a framework resource's acceptance test runs with `TF_ACC=1`
- **THEN** the test harness initializes a `tf5muxserver` with both stacks
- **AND** the test plan/apply behaves identically to production

### Requirement: Framework resources SHALL have website documentation generated alongside SDKv2 resources

Every framework resource and data source SHALL have a corresponding markdown file under `website/docs/r/` or `website/docs/d/`, generated via the project's documentation pipeline (`make doc`). The documentation generator SHALL be capable of introspecting framework schemas; if the existing generator only handles SDKv2, this change SHALL extend or replace it.

#### Scenario: make doc covers framework resources
- **WHEN** a developer runs `make doc` after adding a framework resource
- **THEN** the corresponding `website/docs/r/<service>_<name>.html.markdown` is created or updated
- **AND** the file documents all framework schema attributes including their type, optionality, and description
