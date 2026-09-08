# tencentcloud/framework

Home of the terraform-plugin-framework Provider entry point. Business
references (resources, data sources, functions, ephemeral resources, list
resources and actions) are **co-located with the SDKv2 implementations**
under `tencentcloud/services/<product>/` — there is no separate
framework subtree for business code.

> Current state: the framework Provider is wired into mux but **no
> business references are registered yet**. The six factory slices in
> `registry.go` are all empty. Add new references following the workflow
> in [Adding a new framework reference](#adding-a-new-framework-reference).

## Directory layout

```
tencentcloud/
├── framework/                # framework entry only — no business code lives here
│   ├── provider.go           # framework Provider implementation (Schema / Configure / Resources / DataSources / Functions / EphemeralResources / ListResources / Actions)
│   ├── registry.go           # SDKv2-style central manifest: six append-only factory slices
│   ├── provider_test.go      # in-package tests: mux startup sanity + no type-name collisions
│   ├── testhelpers_test.go   # in-package test helpers
│   ├── acctest/              # ProtoV5 test factories shared by acceptance tests
│   ├── internal/             # framework-only helpers (Go-internal-visible to framework/...)
│   └── README.md             # this file
└── services/
    ├── common/               # cross-product / provider-meta references (package common)
    ├── cvm/                  # CVM product (package cvm; SDKv2 + framework can mix)
    └── <product>/            # any other product follows the same pattern
```

> The single SDKv2 reference index `tencentcloud/provider.md` is the
> only documentation index of the provider. Framework references are
> listed under their corresponding product node alongside SDKv2
> entries (using the new sub-section headers `Function`,
> `Ephemeral Resource`, `List Resource`, `Action`). The previously
> dedicated `framework/provider.md` has been removed.

### Generated docs directory layout

`make doc` writes one HTML markdown per reference into
`website/docs/<short>/<resName>.html.markdown`, where `<short>` is a
single letter mirroring the SDKv2 `d` / `r` shorthand:

| Framework type | Output dir |
|---|---|
| Resource | `r/` |
| Data Source | `d/` |
| Function | `f/` |
| Ephemeral Resource | `e/` |
| List Resource | `l/` |
| Action | `a/` |

### File naming convention (mandatory)

Framework references follow the same `<dtype>_tc_<name>.go` template
used by SDKv2 resources / data sources, with one new prefix per
framework type:

| Framework type | File prefix | Example |
|---|---|---|
| Resource | `resource_tc_` | `resource_tc_local_note.go` |
| Data Source | `data_source_tc_` | `data_source_tc_provider_runtime.go` |
| Function | `function_tc_` | `function_tc_parse_resource_id.go` |
| Ephemeral Resource | `ephemeral_tc_` | `ephemeral_tc_temp_credential.go` |
| List Resource | `list_tc_` | `list_tc_region.go` |
| Action | `action_tc_` | `action_tc_cvm_reboot_instance.go` |

When a reference is product-specific, include the product segment (e.g.
`action_tc_cvm_reboot_instance`); when it lives under `services/common/`,
the product segment is omitted (e.g. `resource_tc_local_note`).

Each reference ships with a sibling Markdown file (`<stem>.md`)
containing the description and `Example Usage` block consumed by `make
doc`. `gendoc` accepts both `<dtype>_tc_<resName>.md` and
`<dtype>_tc_<product>_<resName>.md` — the same naming the file uses.

### Package naming convention

The Go package of every framework reference is the **same as the SDKv2
package of that product**:

- `services/common/` → `package common`
- `services/cvm/` → `package cvm`

There is no per-type subpackage (the older `cvmactions` / `metaresources`
split has been removed): the SDKv2 and framework code share one package
per service directory and import paths stay simple.

## Wiring into main.go

`main.go` merges the SDKv2 and the framework providers into a single
provider binary via `tf5muxserver`:

```go
primary := tencentcloud.Provider()                  // SDKv2 entry
fw := framework.NewProvider(primary)                // entry of this package
providers := []func() tfprotov5.ProviderServer{
    primary.GRPCProvider,
    providerserver.NewProtocol5(fw),
}
muxServer, _ := tf5muxserver.NewMuxServer(ctx, providers...)
```

## Adding a new framework reference

The central manifest `tencentcloud/framework/registry.go` mirrors the
SDKv2 `provider.go` style: each framework reference type owns an
append-only factory slice (`resourceFactories`, `dataSourceFactories`,
`functionFactories`, `ephemeralResourceFactories`,
`listResourceFactories`, `actionFactories`). The six framework Provider
callbacks return these slices verbatim.

### Workflow

1. **Implement the factory** in the matching `services/<product>/`
   package (e.g. `services/cvm/resource_tc_cvm_my_thing.go`):

   ```go
   package cvm

   func NewMyThingResource() resource.Resource { return &myThingResource{} }
   // ... implement resource.Resource + resource.ResourceWithConfigure ...
   ```

2. **Register the factory** by adding **one line** to the matching
   slice in `framework/registry.go`:

   ```go
   var resourceFactories = []func() resource.Resource{
       cvm.NewMyThingResource, // <- one line, alphabetically grouped by product
   }
   ```

   If the product subpackage is not yet imported in `registry.go`, add
   one import line as well. `provider.go` does not need to change.

3. **Add the reference to the unified index** in
   `tencentcloud/provider.md` so `make doc` picks it up. Use the
   product node that already exists for that cloud product (or add a
   new product node if the product is new), then list the type under
   the matching sub-section header. Six headers are recognised:
   `Data Source`, `Resource`, `Function`, `Ephemeral Resource`,
   `List Resource`, `Action`. For example:

   ```
   Cloud Virtual Machine(CVM)
   Resource
   tencentcloud_my_thing

   Action
   tencentcloud_cvm_reboot_instance
   ```

   The same provider.md is consumed by the SDKv2 stack — there is no
   separate framework index file.

4. **Ship a sibling Markdown** next to the Go file (e.g.
   `resource_tc_cvm_my_thing.md`) containing a short description and a
   required `Example Usage` HCL block.

That is the complete workflow — no edits to `provider.go`, no per-service
register file, no init() magic.

### Product-ownership rules

- References that **clearly belong to a real cloud product** MUST land
  in the corresponding product directory (e.g.
  `services/vpc/resource_tc_vpc_xxx.go`).
- References that are **cross-product or not bound to any specific cloud
  product** MAY land in `services/common/` (for example, the provider's
  own runtime metadata, or local-only helper functions).

## Relationship with SDKv2

- Credentials, the SDK client, UA and retry are **resolved and built
  only by the SDKv2 provider**. This provider reuses the same
  `*connectivity.TencentCloudClient` via
  `internal/sharedmeta.GetSharedMeta()`.
- This provider's Schema **must mirror SDKv2** (same names, same
  semantics, same nested structure); otherwise mux will reject the
  user's HCL fields when merging the two schemas. The mux startup
  invariants are exercised by `make check-mux`.
- During Configure, `*sharedmeta.ProviderMeta` is written into all four
  fields `resp.{ResourceData, DataSourceData, EphemeralResourceData,
  ActionData}`.
- Each resource / data source / action retrieves the shared client by
  type-asserting `*sharedmeta.ProviderMeta` inside its own `Configure`
  method.
