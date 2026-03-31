## 1. Resource Implementation

- [x] 1.1 Create resource schema definition in `tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config.go`
- [x] 1.2 Implement Create function to call ExportZoneConfig API
- [x] 1.3 Implement Read function to re-export and update state
- [x] 1.4 Implement Update function to handle parameter changes
- [x] 1.5 Implement Delete function to remove from state only
- [x] 1.6 Implement service layer function in `tencentcloud/services/teo/service_tencentcloud_teo.go`
- [x] 1.7 Register the resource in the provider

## 2. Testing

- [x] 2.1 Create unit test file `tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config_test.go`
- [x] 2.2 Implement unit tests with mocked API responses
- [x] 2.3 Create acceptance test cases in the unit test file
- [x] 2.4 Run acceptance tests with TF_ACC=1 to verify integration with real API

## 3. Documentation

- [x] 3.1 Create resource example file `examples/resources/teo_export_zone_config/README.md`
- [x] 3.2 Run `make doc` to generate documentation in `website/docs/r/teo_export_zone_config.html.markdown`
- [x] 3.3 Verify generated documentation contains all parameters and examples

## 4. Validation

- [x] 4.1 Run `make build` to verify the code compiles successfully
- [x] 4.2 Run `make lint` to check code style and potential issues
- [x] 4.3 Run `make test` to execute all unit tests
- [x] 4.4 Run acceptance tests to ensure the resource works with real TEO zones

## 5. Final Review

- [x] 5.1 Verify all CRUD operations work correctly
- [x] 5.2 Verify error handling for invalid zone IDs
- [x] 5.3 Verify error handling for unsupported configuration types
- [x] 5.4 Verify state management after terraform refresh
- [x] 5.5 Verify large configuration content handling
