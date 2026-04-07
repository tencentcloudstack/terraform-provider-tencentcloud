## Why

CVM (Cloud Virtual Machine) needs to support resource pool pack management to allow users to purchase, query, and terminate resource pool packs through Terraform. Currently, there is no Terraform resource to manage CVM resource pool packs, requiring users to manage these resources manually through the console or API, which is inefficient and error-prone for infrastructure-as-code workflows.

## What Changes

- Add a new Terraform resource `tencentcloud_cvm_resource_pool_packs` to manage CVM resource pool packs
- Support creating resource pool packs via `PurchaseResourcePoolPacks` API
- Support querying resource pool packs via `DescribeResourcePoolPacks` API  
- Support deleting resource pool packs via `TerminateResourcePoolPacks` API
- Add corresponding test files for the new resource
- Add documentation (`.md` files) for the new resource
- Register the new resource in `provider.go` and `provider.md`

## Capabilities

### New Capabilities

- `cvm-resource-pool-packs-resource`: Terraform resource for managing CVM resource pool packs lifecycle (create, read, delete)

### Modified Capabilities

<!-- No existing capabilities are being modified -->

## Impact

- **New files**: 
  - `tencentcloud/services/cvm/resource_tc_cvm_resource_pool_packs.go`
  - `tencentcloud/services/cvm/resource_tc_cvm_resource_pool_packs_test.go`
  - `tencentcloud/services/cvm/resource_tc_cvm_resource_pool_packs.md`
  - `openspec/changes/add-cvm-resource-pool-packs/specs/cvm-resource-pool-packs-resource/spec.md`
  
- **Modified files**:
  - `tencentcloud/provider.go` - register new resource
  - `tencentcloud/provider.md` - add resource declaration
  - `tencentcloud/services/cvm/service_tencentcloud_cvm.go` - add service layer methods

- **Dependencies**: 
  - Requires `tencentcloud-sdk-go` CVM package (already included)
  - No new external dependencies needed

- **Breaking changes**: None - this is a new resource addition
