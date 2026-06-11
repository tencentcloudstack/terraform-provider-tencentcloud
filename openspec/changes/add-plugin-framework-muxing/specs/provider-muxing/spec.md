## ADDED Requirements

### Requirement: Provider binary SHALL serve a tf5muxserver combining SDKv2 and framework providers

The provider's `main` package SHALL build a single Protocol v5 server using `tf5muxserver.NewMuxServer` that fronts both the SDKv2 provider (`tencentcloud.Provider().GRPCProvider`) and the framework provider (`providerserver.NewProtocol5(tencentcloud.New(primary))`). The published provider address `registry.terraform.io/tencentcloudstack/tencentcloud` MUST NOT change.

#### Scenario: Single binary serves both stacks
- **WHEN** `terraform-provider-tencentcloud` is launched (with or without `-debuggable`)
- **THEN** it serves exactly one gRPC endpoint via `tf5server.Serve`
- **AND** that endpoint is the muxer combining the two providers in this order: SDKv2 first, framework second

#### Scenario: Backward compatibility for existing users
- **WHEN** an existing user upgrades from a pre-mux version of the provider
- **THEN** their existing `terraform plan` and `terraform apply` against any existing resource succeed without any state migration
- **AND** the provider address in their state file remains `registry.terraform.io/tencentcloudstack/tencentcloud`

### Requirement: Mux SHALL reject duplicate resource and data source type names across stacks

The build and CI process SHALL fail before release if the same Terraform resource type name (e.g. `tencentcloud_foo_bar`) is registered in BOTH the SDKv2 provider (`Provider().ResourcesMap`) and the framework provider (`FrameworkProvider.Resources()`). The same SHALL apply to data sources.

#### Scenario: Duplicate registration fails CI
- **WHEN** a developer registers `tencentcloud_foo_bar` in both SDKv2 `ResourcesMap` and framework `Resources()`
- **AND** CI runs the mux validation check
- **THEN** the check exits non-zero with a diagnostic identifying the duplicated type name
- **AND** the PR cannot be merged until the duplicate is removed from one stack

#### Scenario: Mux startup never panics on collisions in production
- **WHEN** the released binary starts up
- **THEN** `tf5muxserver.NewMuxServer` returns no error
- **AND** because CI already enforced uniqueness, no runtime collision is possible

### Requirement: Configure ordering between stacks SHALL be deterministic and tested

The mux setup SHALL place the SDKv2 provider before the framework provider in the providers slice, ensuring SDKv2 `ConfigureContextFunc` runs first and populates the shared client pointer before the framework provider's `Configure` reads it. A unit test SHALL verify this ordering and the pointer hand-off.

#### Scenario: SDKv2 configures before framework
- **WHEN** the unit test invokes the muxer's `ConfigureProvider` RPC with valid credentials
- **THEN** the SDKv2 stack's `ConfigureContextFunc` runs and populates the shared client
- **AND** the framework stack's `Configure` then reads a non-nil pointer and returns no error diagnostics

### Requirement: Mux MUST remain on Protocol v5 to preserve Terraform 0.13+ compatibility

The mux implementation SHALL use `tf5muxserver` and `providerserver.NewProtocol5`. Migration to Protocol v6 (`tf6muxserver` + `tf5to6server`) is explicitly out of scope and SHALL be tracked as a separate, future change.

#### Scenario: Terraform 0.13 client connects successfully
- **WHEN** a Terraform 0.13.x client downloads and initializes the muxed provider
- **THEN** plugin handshake succeeds at Protocol v5
- **AND** existing resources continue to plan and apply
