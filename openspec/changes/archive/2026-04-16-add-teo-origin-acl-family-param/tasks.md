## 1. Schema and Resource Implementation

- [x] 1.1 Add `origin_acl_family` parameter to resource schema in `resource_tc_teo_origin_acl.go`
  - Set Type to `schema.TypeString`
  - Set Optional to `true`
  - Set Computed to `true`
  - Add description explaining the parameter's purpose and available values

- [x] 1.2 Update `ResourceTencentCloudTeoOriginAclCreate` function to set `origin_acl_family`
  - Read `origin_acl_family` value from schema if provided
  - Set `request.OriginACLFamily` in `EnableOriginACLRequest`
  - Ensure the value is only set when explicitly provided by user

- [x] 1.3 Update `ResourceTencentCloudTeoOriginAclRead` function to read `origin_acl_family`
  - Extract `OriginACLFamily` from `respData.OriginACLFamily` (which comes from `OriginACLInfo.OriginACLFamily`)
  - Set the value in Terraform state using `d.Set("origin_acl_family", ...)`

- [x] 1.4 Update `ResourceTencentCloudTeoOriginAclUpdate` function to support `origin_acl_family` changes
  - Add check for `d.HasChange("origin_acl_family")`
  - If changed, call `ModifyOriginACL` API with the new `OriginACLFamily` value
  - Integrate this update into existing update flow or handle as separate update operation

## 2. Testing

- [x] 2.1 Add unit tests for `origin_acl_family` parameter in `resource_tc_teo_origin_acl_test.go`
  - Test creation with `origin_acl_family` parameter
  - Test creation without `origin_acl_family` parameter (verify computed behavior)
  - Test read operation to ensure `origin_acl_family` is correctly read from API response
  - Test update operation for `origin_acl_family` parameter changes
  - Use gomonkey to mock the cloud API calls (EnableOriginACL, DescribeOriginACL, ModifyOriginACL)

## 3. Documentation

- [x] 3.1 Update documentation example in `resource_tc_teo_origin_acl.md`
  - Add `origin_acl_family` parameter to the example configuration
  - Document available values: "gaz", "mlc", "emc", "plat-gaz", "plat-mlc", "plat-emc"
  - Explain the default behavior when parameter is not specified

## 4. Validation and Formatting

- [x] 4.1 Run gofmt to format the modified code
  - Ensure all code follows Go formatting standards
  - NOTE: This will be executed in the tfpacer-finalize stage

- [x] 4.2 Verify the implementation follows project patterns
  - Check that error handling matches existing patterns (defer tccommon.LogElapsed, defer tccommon.InconsistentCheck)
  - Verify parameter mapping is correct
  - Ensure backward compatibility is maintained

- [x] 4.3 Manual verification of implementation
  - Review that the parameter is correctly added to schema
  - Verify API parameter mappings (EnableOriginACL, DescribeOriginACL, ModifyOriginACL)
  - Check that the computed field behavior is correct

## 5. Final Documentation Generation

- [x] 5.1 Generate documentation using `make doc` command
  - This will automatically generate/update markdown files in `website/docs/` directory
  - Verify that the generated documentation includes the new `origin_acl_family` parameter
  - NOTE: This will be executed in the tfpacer-finalize stage
