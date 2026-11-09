# tencentcloud/framework

Home of the terraform-plugin-framework Provider entry point. Business
references (resources, data sources, functions, ephemeral resources, list
resources and actions) are **co-located with the SDKv2 implementations**
under `tencentcloud/services/<product>/` — there is no separate
framework subtree for business code.

> Current state: the framework Provider is wired into mux and references are
> registered in `registry.go`. Add new references following the workflow in
> [Adding a new framework reference](#adding-a-new-framework-reference).

## Directory layout

```
tencentcloud/
├── framework/                # framework entry only — no business code lives here
│   ├── provider.go           # framework Provider implementation (Schema / Configure / Resources / DataSources / Functions / EphemeralResources / ListResources / Actions)
│   ├── registry.go           # SDKv2-style central manifest: six type-name → factory maps
│   ├── registry_test.go      # registry invariants: key == Metadata type name, full exposure
│   ├── provider_test.go      # in-package tests: mux startup sanity + no type-name collisions
│   ├── testhelpers_test.go   # in-package test helpers
│   ├── acctest/              # ProtoV6 test factories shared by acceptance tests
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
provider binary via `tf6muxserver`; the SDKv2 half is upgraded to
protocol 6 with `tf5to6server` because nested attributes only exist in
tfplugin6:

```go
primary := tencentcloud.Provider()                  // SDKv2 entry
upgradedSdkv2, _ := tf5to6server.UpgradeServer(ctx, primary.GRPCProvider)
providers := []func() tfprotov6.ProviderServer{
    func() tfprotov6.ProviderServer { return upgradedSdkv2 },
    providerserver.NewProtocol6(framework.NewProvider(primary)),
}
muxServer, _ := tf6muxserver.NewMuxServer(ctx, providers...)
```

## Adding a new framework reference

The central manifest `tencentcloud/framework/registry.go` mirrors the
SDKv2 `provider.go` style: each framework reference type owns a map from
the name users write in HCL to its factory (`resourcesMap`,
`dataSourcesMap`, `functionsMap`, `ephemeralResourcesMap`,
`listResourcesMap`, `actionsMap`). The six framework Provider callbacks
return those factories through `registry.go`'s accessors. Registration
order carries no meaning: the framework immediately re-keys the
factories into maps by type name.

A map (rather than a slice) buys two guarantees: the file is searchable
by the Terraform type name, and a duplicate registration is a **compile
error** (duplicate map key) instead of the runtime "Duplicate Resource
Type Defined" diagnostic raised when the provider schema is first
requested.

### Workflow

1. **Implement the factory** in the matching `services/<product>/`
   package (e.g. `services/cvm/resource_tc_cvm_my_thing.go`):

   ```go
   package cvm

   func NewMyThingResource() resource.Resource { return &myThingResource{} }
   // ... implement resource.Resource + resource.ResourceWithConfigure ...
   ```

2. **Register the factory** by adding **one line** to the matching map
   in `framework/registry.go`. The key is the Terraform type name the
   resource reports from `Metadata`:

   ```go
   var resourcesMap = map[string]func() resource.Resource{
       "tencentcloud_cvm_my_thing": cvm.NewMyThingResource, // <- one line
   }
   ```

   If the product subpackage is not yet imported in `registry.go`, add
   one import line as well. `provider.go` does not need to change.

   `registry_test.go` then verifies that the key equals the type name the
   factory reports (so a typo fails CI, not `terraform plan`) and that the
   provider exposes every registry entry. Provider-defined functions are
   the exception: their key is the bare function name (e.g.
   `parse_resource_id`), because a function call already carries the
   provider name.

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
