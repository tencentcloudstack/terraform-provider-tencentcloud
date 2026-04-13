## 1. Code Implementation

- [x] 1.1 Create resource file `tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config.go` with basic structure
- [x] 1.2 Define resource schema with all parameters from CAPI interface `iacpres-i2etM5NTBN` (version `iacpresv-k90qSqRTp9`), ensuring Required/Optional attributes match exactly
- [x] 1.3 Add Timeouts block to schema with default values for Create, Read, Update, and Delete operations
- [x] 1.4 Implement resource Create function with CAPI call and error handling using `defer tccommon.LogElapsed()` and `defer tccommon.InconsistentCheck()`
- [x] 1.5 Implement resource Read function with CAPI call and state refresh logic
- [x] 1.6 Implement resource Update function with CAPI call and state update logic
- [x] 1.7 Implement resource Delete function with CAPI call and state cleanup logic
- [x] 1.8 Add retry logic using `helper.Retry()` in all CRUD operations for eventual consistency
- [x] 1.9 Create resource example file `tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config.md` with complete usage examples
- [x] 1.10 Update service layer in `service_tencentcloud_teo.go` if additional helper functions are needed

## 2. Test Implementation

- [x] 2.1 Create unit test file `tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config_test.go`
- [x] 2.2 Implement unit tests for Create operation with mock data
- [x] 2.3 Implement unit tests for Read operation with mock data
- [x] 2.4 Implement unit tests for Update operation with mock data
- [x] 2.5 Implement unit tests for Delete operation with mock data
- [x] 2.6 Implement edge case tests (missing required parameters, invalid parameters, timeout handling)
- [x] 2.7 Create acceptance test cases in the test file with real API calls

## 3. Documentation

- [x] 3.1 Run `make doc` to generate `website/docs/r/teo_export_zone_config.html.md` automatically from the resource schema

## 4. Verification

- [x] 4.1 Run `make build` to verify code compiles without errors
- [x] 4.2 Run `make lint` to verify code follows project style guidelines
- [x] 4.3 Run unit tests for the resource to ensure all tests pass: `go test -v ./tencentcloud/services/teo -run TestAccTencentCloudTeoExportZoneConfig`
- [x] 4.4 Run acceptance tests with `TF_ACC=1` to verify real API integration: `TF_ACC=1 go test -v ./tencentcloud/services/teo -run TestAccTencentCloudTeoExportZoneConfig -timeout 60m`
