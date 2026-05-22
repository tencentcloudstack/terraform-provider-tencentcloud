## ADDED Requirements

### Requirement: Resource schema definition
The resource `tencentcloud_lighthouseShareBlueprintAcrossAccount` SHALL define the following schema fields:
- `blueprint_id`: Required, ForceNew, TypeString - 镜像ID，可以通过DescribeBlueprints接口返回的BlueprintId获取
- `account_ids`: Required, ForceNew, TypeList of TypeString - 接收共享镜像的账号ID列表，最多10个

#### Scenario: Valid schema with all required fields
- **WHEN** user provides both `blueprint_id` and `account_ids` in the resource configuration
- **THEN** the resource SHALL accept the configuration and proceed with creation

#### Scenario: Missing required field blueprint_id
- **WHEN** user does not provide `blueprint_id` in the resource configuration
- **THEN** Terraform SHALL report a required field missing error

#### Scenario: Missing required field account_ids
- **WHEN** user does not provide `account_ids` in the resource configuration
- **THEN** Terraform SHALL report a required field missing error

### Requirement: Create operation
The resource Create function SHALL call the `ShareBlueprintAcrossAccounts` API with the provided `blueprint_id` and `account_ids` parameters. The API call SHALL be wrapped with `resource.Retry(tccommon.WriteRetryTimeout, ...)`, and errors SHALL be wrapped with `tccommon.RetryError(e)`.

#### Scenario: Successful share blueprint across accounts
- **WHEN** the resource is created with valid `blueprint_id` and `account_ids`
- **THEN** the Create function SHALL call `ShareBlueprintAcrossAccounts` API and set the resource ID to a generated token

#### Scenario: API call failure with retryable error
- **WHEN** the `ShareBlueprintAcrossAccounts` API call fails with a retryable error
- **THEN** the Create function SHALL retry the API call within WriteRetryTimeout

#### Scenario: API call failure with non-retryable error
- **WHEN** the `ShareBlueprintAcrossAccounts` API call fails with a non-retryable error
- **THEN** the Create function SHALL return the error without retrying

### Requirement: Read operation
The resource Read function SHALL be a no-op and return nil, as this is a one-time operation that does not require state tracking.

#### Scenario: Read operation
- **WHEN** Terraform refreshes the resource state
- **THEN** the Read function SHALL return nil without performing any API calls

### Requirement: Delete operation
The resource Delete function SHALL be a no-op and return nil, as this is a one-time operation that does not require cleanup.

#### Scenario: Delete operation
- **WHEN** the resource is destroyed
- **THEN** the Delete function SHALL return nil without performing any API calls

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` with the key `tencentcloud_lighthouse_share_blueprint_across_account` and the factory function `lighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountOperation()`.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_lighthouse_share_blueprint_across_account` SHALL be available for use in Terraform configurations

### Requirement: Resource documentation
The resource SHALL have a corresponding markdown documentation file at `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_operation.md` with example usage.

#### Scenario: Documentation exists
- **WHEN** the resource is implemented
- **THEN** a markdown file with example usage SHALL exist in the lighthouse service directory

### Requirement: Unit tests
The resource SHALL have unit tests using gomonkey mock for the `ShareBlueprintAcrossAccounts` API call, verifying the Create function behavior.

#### Scenario: Unit test for successful creation
- **WHEN** the unit test runs the Create function with mocked API response
- **THEN** the test SHALL verify the resource ID is set and no error is returned

#### Scenario: Unit test for API error
- **WHEN** the unit test runs the Create function with mocked API error
- **THEN** the test SHALL verify the error is properly returned
