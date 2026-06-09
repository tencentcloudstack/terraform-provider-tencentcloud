## 1. Service Layer Implementation

- [x] 1.1 Add `CreateCvmResourcePoolPacks()` method to `service_tencentcloud_cvm.go` that wraps `PurchaseResourcePoolPacks` API with retry logic
- [x] 1.2 Add `DescribeCvmResourcePoolPackById()` method to `service_tencentcloud_cvm.go` that wraps `DescribeResourcePoolPacks` API with filtering by pack ID and retry logic
- [x] 1.3 Add `DeleteCvmResourcePoolPacks()` method to `service_tencentcloud_cvm.go` that wraps `TerminateResourcePoolPacks` API
- [x] 1.4 Implement `ratelimit.Check()` calls before each API invocation
- [x] 1.5 Add proper error handling with `defer tccommon.LogElapsed()` and `defer tccommon.InconsistentCheck()`

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_cvm_resource_pool_packs.go` with resource definition function `ResourceTencentCloudCvmResourcePoolPacks()`
- [x] 2.2 Define resource schema with all required and optional fields, marking all fields with `ForceNew: true`
- [x] 2.3 Implement `resourceTencentCloudCvmResourcePoolPacksCreate()` function that calls service layer Create method
- [x] 2.4 Implement `resourceTencentCloudCvmResourcePoolPacksRead()` function with retry logic for eventual consistency
- [x] 2.5 Implement `resourceTencentCloudCvmResourcePoolPacksDelete()` function that calls service layer Delete method
- [x] 2.6 Implement `resourceTencentCloudCvmResourcePoolPacksImporter()` for import support using pack ID
- [x] 2.7 Add proper context handling for timeout support in CRUD operations

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_cvm_resource_pool_packs` in `provider.go` ResourcesMap
- [x] 3.2 Add resource declaration in `provider.md` under CVM resources section

## 4. Test Implementation

- [x] 4.1 Create `resource_tc_cvm_resource_pool_packs_test.go` with `TestAccTencentCloudCvmResourcePoolPacksResource_basic` test
- [x] 4.2 Implement test covering basic create/read/delete lifecycle
- [x] 4.3 Implement test for import functionality
- [x] 4.4 Implement test verifying ForceNew behavior on field changes
- [x] 4.5 Add test configuration with valid test parameters

## 5. Documentation

- [x] 5.1 Create `resource_tc_cvm_resource_pool_packs.md` with usage example showing basic resource creation
- [x] 5.2 Document all schema fields with types and descriptions
- [x] 5.3 Add import example in documentation
- [x] 5.4 Document ForceNew behavior and limitations (no update support)
- [x] 5.5 Add notes about prerequisites for deletion (no active instances)

## 6. Validation and Testing

- [x] 6.1 Run `go build ./tencentcloud/services/cvm/...` to verify compilation
- [x] 6.2 Run `make doc` to generate website documentation
- [x] 6.3 Run `TF_ACC=1 go test ./tencentcloud/services/cvm -run TestAccTencentCloudCvmResourcePoolPacks -v` to verify acceptance tests
- [x] 6.4 Verify generated documentation in `website/docs/r/cvm_resource_pool_packs.html.markdown`
- [x] 6.5 Test manual terraform apply/destroy with the new resource
