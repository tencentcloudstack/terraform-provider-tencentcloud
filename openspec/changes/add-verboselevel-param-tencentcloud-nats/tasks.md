## 1. Schema Implementation

- [x] 1.1 Add `verbose_level` parameter to DataSourceTencentCloudNats schema with TypeString and Optional attributes
- [x] 1.2 Create validateVerboseLevel function to accept only "DETAIL", "COMPACT", "SIMPLE" values
- [x] 1.3 Add validation function to `verbose_level` parameter definition

## 2. Read Function Implementation

- [x] 2.1 Retrieve `verbose_level` parameter value in dataSourceTencentCloudNatsRead function
- [x] 2.2 Set request.VerboseLevel when verbose_level parameter is provided
- [x] 2.3 Ensure default behavior (no VerboseLevel) when parameter is not set

## 3. Documentation

- [x] 3.1 Update data source documentation file if needed
- [x] 3.2 Add usage example showing verbose_level parameter with different values
- [x] 3.3 Run `make doc` to generate website documentation

## 4. Testing

- [x] 4.1 Add acceptance test case for verbose_level parameter with "DETAIL" value
- [x] 4.2 Add acceptance test case for verbose_level parameter with "COMPACT" value
- [x] 4.3 Add acceptance test case for verbose_level parameter with "SIMPLE" value
- [x] 4.4 Add acceptance test case for invalid verbose_level value (should fail validation)
- [x] 4.5 Add acceptance test case for omitted verbose_level parameter (default behavior)
- [x] 4.6 Run acceptance tests with TF_ACC=1 to verify implementation

## 5. Build and Verification

- [x] 5.1 Build provider to ensure no compilation errors
- [x] 5.2 Run go fmt to ensure code formatting
- [x] 5.3 Run go vet to check for potential issues
- [x] 5.4 Verify all acceptance tests pass
