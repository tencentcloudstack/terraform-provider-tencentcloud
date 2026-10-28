## ADDED Requirements

### Requirement: Framework Provider SHALL declare a schema mirroring SDKv2 provider fields

The framework provider SHALL declare provider-block attributes whose names and semantics mirror the SDKv2 provider, so that a single `provider "tencentcloud" {}` block in user HCL is accepted by the muxed server without `unknown attribute` errors.

The schema MUST cover at minimum: `secret_id`, `secret_key`, `security_token`, `region`, `protocol`, `domain`, `shared_credentials_dir`, `profile`, `cam_role_name`, and the nested `assume_role` block. All attributes MUST be declared `Optional` (sensitive ones marked `Sensitive`).

#### Scenario: User configures provider once and both stacks accept it
- **WHEN** a user writes a `provider "tencentcloud" { secret_id = "..." region = "ap-guangzhou" }` block and runs `terraform plan` against the muxed binary
- **THEN** Terraform validates the block against the merged schema without raising "unsupported argument" diagnostics
- **AND** both SDKv2 and framework provider servers receive the configuration via mux

#### Scenario: Framework schema stays in sync with SDKv2
- **WHEN** a developer adds a new top-level provider attribute to SDKv2
- **THEN** they MUST add the equivalent attribute to the framework provider schema in the same change
- **AND** CI's mux schema check fails if the two schemas diverge on attribute names

### Requirement: Framework Provider Configure SHALL reuse the SDKv2-built client

The framework provider's `Configure` method SHALL NOT independently parse credentials, environment variables, or shared credentials files. Instead it SHALL obtain the `*connectivity.TencentCloudClient` constructed by the SDKv2 provider and inject it as `ProviderData` for resources, data sources, and ephemeral resources.

A process-level shared pointer (e.g. `atomic.Pointer[connectivity.TencentCloudClient]`) SHALL be the handover channel: SDKv2's `ConfigureContextFunc` writes it; framework's `Configure` reads it.

#### Scenario: Framework resource receives the same client as SDKv2
- **WHEN** the muxer calls `Configure` on both stacks with the same provider block
- **AND** the framework provider's `Configure` runs after SDKv2's
- **THEN** `resp.ResourceData` and `resp.DataSourceData` MUST be set to a `*ProviderData` whose `Client` field points to the exact same `*connectivity.TencentCloudClient` instance returned by `tencentcloud.Provider().Meta()`

#### Scenario: Framework Configure handles SDKv2 not-yet-configured case
- **WHEN** framework `Configure` runs and the shared client pointer is `nil`
- **THEN** it MUST append an `Error` diagnostic with summary "TencentCloud provider not configured" and a detail explaining that the SDKv2 stack must configure first
- **AND** it MUST NOT panic or set partially populated `ProviderData`

#### Scenario: Credentials are parsed exactly once
- **WHEN** the muxed provider is configured
- **THEN** environment-variable resolution, shared-credentials-file lookup, `assume_role` STS calls, and `cam_role_name` metadata calls MUST occur exactly once (in SDKv2's `ConfigureContextFunc`)
- **AND** no equivalent parsing logic SHALL exist in the framework `Configure` path

### Requirement: Framework Provider SHALL aggregate resources and data sources via per-service registration functions

The provider SHALL collect framework-side resources and data sources from per-service aggregation functions: each service package under `tencentcloud/services/<service>/` exposes `FrameworkResources()` and `FrameworkDataSources()`. The top-level `FrameworkProvider.Resources()` and `DataSources()` methods SHALL concatenate these lists.

#### Scenario: New framework resource is wired into the provider
- **WHEN** a developer adds `resource_tc_foo_bar_framework.go` defining `NewFooBarResource()`
- **AND** appends it to `services/foo/framework.go`'s `FrameworkResources()` slice
- **THEN** running `terraform plan` against the built binary lists `tencentcloud_foo_bar` as a known resource type
- **AND** no other registration step (e.g. global `init()`) is required

### Requirement: Framework Provider SHALL expose extension points for ephemeral resources, functions, and list resources

`FrameworkProvider` SHALL implement `EphemeralResources`, `Functions`, and `ListResources` to return slices that aggregate per-service contributions, even if the slices are empty in the initial release. This SHALL keep the door open for future capabilities without further provider-level refactors.

#### Scenario: Adding a provider-defined function does not require provider plumbing changes
- **WHEN** a future change adds the first provider-defined function in some service package
- **THEN** the developer only edits that service's `framework.go` to register it
- **AND** the top-level `FrameworkProvider.Functions` already iterates and exposes it without modification
