## 1. Schema Definition

- [x] 1.1 Add `dedicated_resource_pack_tenancy` field to schema in `resource_tc_instance.go` with Type=String, Optional=true, ForceNew=true, RequiredWith=["dedicated_resource_pack_ids"]
- [x] 1.2 Add `dedicated_resource_pack_ids` field to schema in `resource_tc_instance.go` with Type=List, Elem=String, Optional=true, ForceNew=true, RequiredWith=["dedicated_resource_pack_tenancy"]
- [x] 1.3 Add descriptive descriptions for both fields explaining their purpose and usage with resource pool packs

## 2. Create Function Implementation

- [x] 2.1 In `resourceTencentCloudInstanceCreate()`, read `dedicated_resource_pack_tenancy` from resource data after line 604 (where Placement is created)
- [x] 2.2 If `dedicated_resource_pack_tenancy` is set, assign it to `request.Placement.DedicatedResourcePackTenancy`
- [x] 2.3 Read `dedicated_resource_pack_ids` from resource data and convert the list to []*string
- [x] 2.4 Assign the converted IDs to `request.Placement.DedicatedResourcePackIds`

## 2.5. Read Function Implementation

- [x] 2.5.1 In `resourceTencentCloudInstanceRead()`, read `instance.Placement.DedicatedResourcePackTenancy` from API response
- [x] 2.5.2 If not nil, set it to resource data via `d.Set("dedicated_resource_pack_tenancy", ...)`
- [x] 2.5.3 Read `instance.Placement.DedicatedResourcePackIds` from API response
- [x] 2.5.4 If the list is not empty, convert to interface slice and set to resource data via `d.Set("dedicated_resource_pack_ids", ...)`

## 3. Documentation

- [x] 3.1 Update `resource_tc_instance.md` to add a new example showing instance creation with dedicated resource pack parameters
- [x] 3.2 In the example, demonstrate both `dedicated_resource_pack_tenancy = "ResourcePool"` and `dedicated_resource_pack_ids = ["rpp-xxxxxxxx"]`
- [x] 3.3 Add both new fields to the Argument Reference section with descriptions including type, optionality, ForceNew behavior
- [x] 3.4 Add a note explaining that both parameters must be specified together and reference the `tencentcloud_cvm_resource_pool_packs` resource

## 4. Testing

- [x] 4.1 Create a test case `TestAccTencentCloudInstanceResource_dedicatedResourcePack` in `resource_tc_instance_test.go`
- [x] 4.2 The test should create an instance with both `dedicated_resource_pack_tenancy` and `dedicated_resource_pack_ids` specified
- [x] 4.3 Verify the instance is created successfully
- [ ] 4.4 Add a test step to verify ForceNew behavior by changing one of the parameters
- [ ] 4.5 (Optional) Add a negative test verifying validation fails when only one parameter is specified

## 5. Validation and Documentation Generation

- [x] 5.1 Run `go build ./tencentcloud/services/cvm/...` to verify compilation
- [x] 5.2 Run `make doc` to generate website documentation
- [x] 5.3 Verify generated documentation in `website/docs/r/instance.html.markdown` includes the new parameters
- [ ] 5.4 Run `go test ./tencentcloud/services/cvm -run TestAccTencentCloudInstanceResource_dedicatedResourcePack -v` to verify the new test passes (if acceptance test credentials are available)
- [x] 5.5 Manually review the code changes to ensure backward compatibility (no existing configurations are broken)
