# tencentcloud/framework

Home of the terraform-plugin-framework Provider entry point and every
framework-side resource / data source.

## Directory layout

Everything related to framework (entry point, registry, tests and the
six business-type implementations) lives under `tencentcloud/framework/`,
organised in a two-level **"product (service) -> type"** structure:

```
tencentcloud/framework/
├── provider.go              # framework Provider implementation (Schema / Configure / Resources / DataSources / Functions / EphemeralResources / ListResources / Actions)
├── registry.go              # six aggregator functions: frameworkResources / frameworkDataSources / frameworkFunctions / frameworkEphemeralResources / frameworkListResources / frameworkActions
├── provider_test.go         # in-package tests: mux startup sanity + no type-name collisions
├── testhelpers_test.go      # in-package test helpers
├── README.md                # this file
├── cvm/                     # product: CVM
│   └── actions/             # package cvmactions
│       ├── reboot_instance_action.go
│       └── reboot_instance_action_test.go
└── meta/                    # meta product: cross-product / not bound to any specific cloud product
    ├── resources/           # package metaresources
    │   ├── local_note_resource.go
    │   └── local_note_resource_test.go
    ├── datasources/         # package metadatasources
    │   ├── provider_runtime_data_source.go
    │   └── provider_runtime_data_source_test.go
    ├── functions/           # package metafunctions
    │   ├── parse_resource_id_function.go
    │   └── parse_resource_id_function_test.go
    ├── ephemerals/          # package metaephemerals
    │   ├── temp_credential_ephemeral_resource.go
    │   └── temp_credential_ephemeral_resource_test.go
    └── lists/               # package metalists (**L0 placeholder**, not wired into registry)
        ├── region_list_resource.go
        └── region_list_resource_test.go
```

### Package naming convention (mandatory)

- Top-level package name of `tencentcloud/framework/`: **`framework`**.
- Product directories (`cvm/` / `meta/` / future `vpc/` etc.) **do not**
  contain `.go` files of their own; they are pure namespace containers.
- The package name of a type subdirectory is **`<product><type-plural>`**
  (the product prefix disambiguates same-named subpackages across
  products), for example:
  - `framework/cvm/actions/` -> `package cvmactions`
  - `framework/meta/resources/` -> `package metaresources`
  - `framework/meta/datasources/` -> `package metadatasources`
  - `framework/meta/functions/` -> `package metafunctions`
  - `framework/meta/ephemerals/` -> `package metaephemerals`
  - `framework/meta/lists/` -> `package metalists`

This way `registry.go` can import same-typed subpackages from multiple
products without needing aliases.

## Wiring into main.go

`main.go` merges the SDKv2 and the framework providers into a single
provider binary via `tf5muxserver`:

```go
primary := tencentcloud.Provider()                  // SDKv2 entry (still under the tencentcloud package for now)
fw := framework.NewProvider(primary)                // entry of this package
providers := []func() tfprotov5.ProviderServer{
    primary.GRPCProvider,
    providerserver.NewProtocol5(fw),
}
muxServer, _ := tf5muxserver.NewMuxServer(ctx, providers...)
```

Callers **use the package name `framework` directly**; the legacy
`fwprovider` alias is no longer used.

## Adding a new framework resource / data source / type

Land directly into the two-level "product / type" layout. There is no
longer a `services/<service>/framework.go` middle layer.

### Example: add a framework resource under CVM

```bash
mkdir -p tencentcloud/framework/cvm/resources
```

```go
// tencentcloud/framework/cvm/resources/instance_resource.go
package cvmresources

func NewInstanceResource() resource.Resource { return &instanceResource{} }
// ... implement resource.Resource + resource.ResourceWithConfigure ...
```

Then wire it up inside `tencentcloud/framework/registry.go`:

```go
import (
    cvmresources "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/cvm/resources"
    // ... other existing imports ...
)

func frameworkResources() []func() resource.Resource {
    out := make([]func() resource.Resource, 0)
    out = append(out, metaresources.NewLocalNoteResource)
    out = append(out, cvmresources.NewInstanceResource)  // <- newly wired in
    return out
}
```

`provider.go` does not need to change.

### Product-ownership rules

- Resources that **clearly belong to a real cloud product** MUST land in
  the corresponding product directory (e.g. `framework/vpc/resources/...`).
- Resources that are **cross-product or not bound to any specific cloud
  product** MAY land under `framework/meta/` (for example, the provider's
  own runtime metadata, or local-only helper functions).

## Six-type reference implementations cheat sheet

| Type | Package path | Type name | Level | Notes |
|---|---|---|---|---|
| resource | `framework/meta/resources/` | `tencentcloud_local_note` | L2 | Pure local in-memory resource (CRUD only operates on a sync.Map). |
| datasource | `framework/meta/datasources/` | `tencentcloud_provider_runtime` | L2 | Provider runtime metadata (no IO). |
| function | `framework/meta/functions/` | `parse_resource_id` | L2 | `(id, sep) -> list[string]`, plain `strings.Split`. |
| ephemeral | `framework/meta/ephemerals/` | `tencentcloud_temp_credential` | L2 | Locally constructs a 5-minute placeholder credential. |
| list | `framework/meta/lists/` | `tencentcloud_region` | **L0** | Static region data + helper only; `list.ListResource` interface is **not** implemented (see the in-package doc comment). |
| action | `framework/cvm/actions/` | `tencentcloud_reboot_instance` | L2 | Validates instance_id with a regex and logs; **does NOT call** the CVM API. The execution method is **`Invoke`**. |

> **L0 downgrade note for the list type**: framework v1.19's
> `list.ListResource` requires the list type name to match an
> **already-registered managed resource** and demands `ResourceIdentity`
> plus a Go 1.23 `iter.Seq[ListResult]` iterator. A full integration
> requires first implementing the same-named `tencentcloud_region`
> resource together with its IdentitySchema, which is beyond the scope of
> this directory's initial drop. It will be addressed in a separate
> follow-up change.

## Relationship with SDKv2

- Credentials, the SDK client, UA and retry are **resolved and built only
  by the SDKv2 provider**. This provider reuses the same
  `*connectivity.TencentCloudClient` via
  `internal/sharedmeta.GetSharedMeta()`.
- This provider's Schema **must mirror SDKv2** (same names, same
  semantics, same nested structure); otherwise mux will reject the user's
  HCL fields when merging the two schemas.
- During Configure, `*sharedmeta.ProviderMeta` is written into all four
  fields `resp.{ResourceData, DataSourceData, EphemeralResourceData,
  ActionData}`.
- Each resource / data source / action retrieves the shared client by
  type-asserting `*sharedmeta.ProviderMeta` inside its own `Configure`
  method.
