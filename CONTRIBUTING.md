# Contributing to terraform-provider-tencentcloud

Thank you for taking the time to contribute! This document captures the
project conventions you should follow when adding or modifying provider
resources and data sources.

## Provider Architecture: SDKv2 + Framework

The provider serves both [`terraform-plugin-sdk/v2`][sdkv2] and
[`terraform-plugin-framework`][framework] resources via a single
[`tf5muxserver`][mux] in `main.go`. The two stacks share the same
`*connectivity.TencentCloudClient`, so credentials and SDK clients are
configured exactly once.

[sdkv2]: https://github.com/hashicorp/terraform-plugin-sdk
[framework]: https://github.com/hashicorp/terraform-plugin-framework
[mux]: https://github.com/hashicorp/terraform-plugin-mux

### When to use which stack

| Situation | Use |
|---|---|
| New resource or data source | **framework** (preferred) |
| Modifying an existing SDKv2 resource (small change) | SDKv2 (do not migrate) |
| Migrating an existing SDKv2 resource to framework | Separate change with full acceptance test, **one resource per PR** |

The default for **all new development** is the framework stack, mirroring
the approach taken by [`terraform-provider-google`][google] and
[`terraform-provider-aws`][aws].

[google]: https://github.com/hashicorp/terraform-provider-google
[aws]: https://github.com/hashicorp/terraform-provider-aws

### Hard rules (enforced by CI)

1. **No type-name collisions across stacks.** A given Terraform resource
   type name (e.g. `tencentcloud_foo_bar`) must be registered in *exactly
   one* of: `Provider().ResourcesMap` (SDKv2) or
   `FrameworkProvider.Resources()` (framework). The same applies to data
   sources. CI runs `make check-mux` to enforce this.
2. **Backward compatibility is non-negotiable.** Do not change the schema,
   ID format, or state shape of any released resource without a state
   migration and a separate change document.
3. **No new credential parsing in framework provider.** Credentials are
   parsed by the SDKv2 provider's `providerConfigure` function and shared
   via `sharedmeta.SetSharedMeta`. The framework provider must only
   *read* the shared client, never re-parse environment variables or
   shared credentials files.
4. **Vendor sync.** Any change to `go.mod` must be followed by
   `go mod vendor` in the same commit; CI compares `vendor/modules.txt`.

## Adding a Framework Resource

framework resources/data sources/functions/ephemerals/lists/actions are
organized under `tencentcloud/framework/` in a **product-by-type** layout.
There is **no** `tencentcloud/services/<service>/framework.go` middle layer —
the registry imports product subpackages directly.

```text
tencentcloud/framework/<product>/<type>/
├── <resource_name>_resource.go               # implementation
├── <resource_name>_resource_test.go          # unit/acceptance test
└── ...
```

Concrete examples already in the repo:

```text
tencentcloud/framework/meta/resources/local_note_resource.go      # package metaresources
tencentcloud/framework/meta/datasources/provider_runtime_data_source.go
tencentcloud/framework/meta/functions/parse_resource_id_function.go
tencentcloud/framework/meta/ephemerals/temp_credential_ephemeral_resource.go
tencentcloud/framework/cvm/actions/reboot_instance_action.go      # package cvmactions
```

Steps:

1. Pick the product directory. If the resource is clearly tied to a
   TencentCloud service (CVM, VPC, CBS, ...), use
   `tencentcloud/framework/<product>/<type>/`. Only resources that are
   **cross-product or not bound to any specific service** belong under
   `tencentcloud/framework/meta/<type>/`.
2. The package name is `<product><type-plural>`, e.g. `cvmresources`,
   `metafunctions`, `cvmactions`. This prefix avoids cross-product
   collisions when `registry.go` imports several `<x>resources` subpackages.
3. Implement the resource using the `terraform-plugin-framework` interfaces.
   Use `tencentcloud/framework/internal/helper` helpers for retries, error
   translation, type conversions, and the `timeouts` block.
4. In your `Configure` method, type-assert
   `req.ProviderData.(*sharedmeta.ProviderMeta)` and store the
   `Client` field. Be defensive against `nil` provider data.
5. Wire the factory into `tencentcloud/framework/registry.go` by adding an
   `import` for your product subpackage and appending the factory in the
   matching aggregator (`frameworkResources` / `frameworkDataSources` /
   `frameworkFunctions` / `frameworkEphemeralResources` /
   `frameworkListResources` / `frameworkActions`). No change to
   `tencentcloud/framework/provider.go` is needed.
6. Use `tcfwacctest.AccProtoV5ProviderFactories` (alias for
   `tencentcloud/framework/acctest`, NOT `AccProviders`) in acceptance
   tests so the test exercises the muxed binary. The framework-only
   factory now lives under `tencentcloud/framework/acctest/`; the
   shared `AccPreCheck` / test helpers remain in `tencentcloud/acctest/`,
   so framework tests typically import both packages with aliases
   `tcacctest` (PreCheck/test_util) and `tcfwacctest` (factories).
7. Add documentation to `website/docs/r/<service>_<name>.html.markdown`
   (or `d/...` for data sources). The current `gendoc` generator does not
   yet understand framework schemas — handwrite the doc until the
   generator is upgraded (tracked separately).

> Note: the framework `Action` interface's execution method is named
> **`Invoke`** (not `Run`). Only the framework `Function` interface uses
> `Run`. The schema-side method for `Function` is `Definition` (returning
> `function.Definition{Parameters, Return}`), not `Schema`.

## Useful Make targets

| Target | Purpose |
|---|---|
| `make build` | Compile the provider |
| `make fmt` | Format code with `gofmt` |
| `make lint` | golangci-lint + tfproviderlint |
| `make test` | Unit tests, 30s timeout |
| `make check-mux` | Verify SDKv2 + framework mux compatibility (no panics, no duplicate type names) |
| `make testacc` | Acceptance tests (requires `TENCENTCLOUD_SECRET_ID`/`KEY`) |
| `make doc` | Regenerate website docs (SDKv2 only — handwrite docs for framework resources for now) |

## PR Checklist

Before opening a PR, please confirm:

- [ ] `make build` succeeds
- [ ] `make fmt` shows no diff
- [ ] `make lint` passes (or pre-existing failures are unrelated)
- [ ] `make check-mux` passes
- [ ] For new framework resources/data sources: their type name is **not**
      already registered in SDKv2 (`tencentcloud/provider.go` ResourcesMap)
- [ ] Acceptance tests added (or skipped with justification)
- [ ] Website docs updated under `website/docs/r/` or `website/docs/d/`
- [ ] If `go.mod` changed: `go mod vendor` was run in the same commit
