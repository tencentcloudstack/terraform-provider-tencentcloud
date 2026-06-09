## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_secret_key_config.go` with schema definition (zone_id: Required/ForceNew/TypeString, secret_key: Required/TypeString/Sensitive) and Import support
- [x] 1.2 Implement Create function: set zone_id as resource ID, then call Update function (which calls DescribeMultiPathGatewaySecretKey first to check key existence, then calls CreateMultiPathGatewaySecretKey API if key exists, or ModifyMultiPathGatewaySecretKey API if key doesn't exist)
- [x] 1.3 Implement Read function: call DescribeMultiPathGatewaySecretKey API with ZoneId, populate secret_key from response, handle resource-not-found by removing from state
- [x] 1.4 Implement Update function: first call DescribeMultiPathGatewaySecretKey to check if key exists; if key exists, call CreateMultiPathGatewaySecretKey API to replace the key; if key doesn't exist, call ModifyMultiPathGatewaySecretKey API to set the key. Include retry logic with tccommon.WriteRetryTimeout
- [x] 1.5 Implement Delete function: remove resource from Terraform state only (no API call needed)

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_multi_path_gateway_secret_key` resource in `tencentcloud/provider.go` with factory function `ResourceTencentCloudTeoMultiPathGatewaySecretKeyConfig`
- [x] 2.2 Add resource entry in `tencentcloud/provider.md`

## 3. Unit Tests

- [x] 3.1 Create test file `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_secret_key_config_test.go` with gomonkey mocks for Create, Read, Update, Delete operations
- [x] 3.2 Cover both key-exists and key-not-exists scenarios for Create and Update operations
- [x] 3.3 Run unit tests with `go test -gcflags=all=-l` and ensure all tests pass

## 4. Documentation

- [x] 4.1 Create documentation file `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_secret_key_config.md` with one-line description, example usage, and import section
