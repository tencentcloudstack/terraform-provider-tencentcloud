## ADDED Requirements

### Requirement: Resource schema definition
The resource `tencentcloud_teo_multi_path_gateway_secret_key` SHALL define the following schema fields:
- `zone_id`: Required, ForceNew, TypeString - 站点 ID
- `secret_key`: Required, TypeString, Sensitive - 多通道安全加速网关接入密钥（base64 字符串，编码前字符串长度为 32-48 个字符）

The resource SHALL support import via `schema.ImportStatePassthrough`.

#### Scenario: Schema fields are correctly defined
- **WHEN** the resource schema is registered
- **THEN** `zone_id` is defined as Required, ForceNew, TypeString
- **AND** `secret_key` is defined as Required, TypeString, Sensitive

#### Scenario: Resource supports import
- **WHEN** a user imports an existing secret key config using `terraform import`
- **THEN** the resource SHALL be imported with `zone_id` as the resource ID
- **AND** the Read method SHALL populate `secret_key` from the cloud API

### Requirement: Create operation
The Create operation SHALL set `zone_id` as the resource ID and then invoke the Update operation to apply the secret key configuration via the `ModifyMultiPathGatewaySecretKey` API.

#### Scenario: Creating a new secret key config
- **WHEN** a user creates a `tencentcloud_teo_multi_path_gateway_secret_key` resource with `zone_id` and `secret_key`
- **THEN** the resource ID SHALL be set to the value of `zone_id`
- **AND** the `ModifyMultiPathGatewaySecretKey` API SHALL be called with `ZoneId` and `SecretKey` parameters

#### Scenario: Create with retry on API failure
- **WHEN** the `ModifyMultiPathGatewaySecretKey` API call fails during Create
- **THEN** the operation SHALL retry with `tccommon.ReadRetryTimeout`
- **AND** return a wrapped error via `tccommon.RetryError()`

### Requirement: Read operation
The Read operation SHALL call the `DescribeMultiPathGatewaySecretKey` API to retrieve the current secret key configuration and update the Terraform state accordingly.

#### Scenario: Reading existing secret key config
- **WHEN** the Read operation is invoked for an existing resource
- **THEN** the `DescribeMultiPathGatewaySecretKey` API SHALL be called with `ZoneId` from the resource ID
- **AND** `secret_key` SHALL be populated from `response.Response.SecretKey`

#### Scenario: Resource not found during Read
- **WHEN** the `DescribeMultiPathGatewaySecretKey` API returns an error indicating the resource does not exist
- **THEN** the resource SHALL be removed from the Terraform state

#### Scenario: Read with retry on API failure
- **WHEN** the `DescribeMultiPathGatewaySecretKey` API call fails
- **THEN** the operation SHALL retry with `tccommon.ReadRetryTimeout`
- **AND** return a wrapped error via `tccommon.RetryError()`

### Requirement: Update operation
The Update operation SHALL call the `ModifyMultiPathGatewaySecretKey` API to update the secret key configuration when `secret_key` changes.

#### Scenario: Updating the secret key
- **WHEN** a user updates the `secret_key` field of an existing resource
- **THEN** the `ModifyMultiPathGatewaySecretKey` API SHALL be called with the `ZoneId` and the new `SecretKey`

#### Scenario: Update with retry on API failure
- **WHEN** the `ModifyMultiPathGatewaySecretKey` API call fails during Update
- **THEN** the operation SHALL retry with `tccommon.ReadRetryTimeout`
- **AND** return a wrapped error via `tccommon.RetryError()`

### Requirement: Delete operation
The Delete operation SHALL only remove the resource from the Terraform state without making any API calls, as the secret key config is a site-level configuration that cannot be deleted.

#### Scenario: Deleting the resource
- **WHEN** a user destroys the `tencentcloud_teo_multi_path_gateway_secret_key` resource
- **THEN** the resource SHALL be removed from the Terraform state
- **AND** no cloud API call SHALL be made

### Requirement: Resource registration
The resource SHALL be registered in `tencentcloud/provider.go` and `tencentcloud/provider.md` following the pattern of existing TEO resources.

#### Scenario: Resource is registered in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_multi_path_gateway_secret_key` SHALL be available as a resource type
- **AND** the resource factory function `ResourceTencentCloudTeoMultiPathGatewaySecretKeyConfig` SHALL be registered

### Requirement: Documentation
The resource SHALL have a corresponding `.md` documentation file at `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_secret_key_config.md` following the project documentation standards.

#### Scenario: Documentation file exists
- **WHEN** the resource is created
- **THEN** a `.md` file SHALL exist with a one-line description, example usage, and import section

### Requirement: Unit tests
The resource SHALL have unit tests in `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_secret_key_config_test.go` using mock (gomonkey) approach.

#### Scenario: Unit tests cover CRUD operations
- **WHEN** unit tests are executed
- **THEN** tests SHALL cover the Create, Read, Update, and Delete operations using gomonkey mocks
- **AND** tests SHALL pass with `go test -gcflags=all=-l`
