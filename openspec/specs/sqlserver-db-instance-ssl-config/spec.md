### Requirement: Resource schema definition
The resource `tencentcloud_sqlserver_db_instance_ssl_config` SHALL define the following schema fields:

Required fields:
- `instance_id` (TypeString, ForceNew): SQL Server 实例 ID
- `type` (TypeString): SSL 操作类型，可选值 enable/disable/renew

Optional fields:
- `wait_switch` (TypeInt): 执行时机，0-立即执行，1-维护时间内执行，默认 0
- `is_kms` (TypeInt): 是否 KMS 加密保护，0-否，1-是，默认 0
- `key_id` (TypeString): KMS 中购买的用户主密钥 ID，IsKMS 为 1 时必填
- `key_region` (TypeString): CMK 所属地域，IsKMS 为 1 时必填

Computed fields:
- `encryption` (TypeString): SSL 加密状态（enable/disable/enable_doing/disable_doing/renew_doing/wait_doing）
- `ssl_validity_period` (TypeString): SSL 证书有效期
- `ssl_validity` (TypeInt): SSL 证书有效性，0-无效，1-有效

#### Scenario: Resource created with required fields only
- **WHEN** user creates a `tencentcloud_sqlserver_db_instance_ssl_config` resource with `instance_id` and `type`
- **THEN** the resource SHALL be created and the SSL configuration SHALL be updated accordingly

#### Scenario: Resource created with KMS encryption
- **WHEN** user creates a `tencentcloud_sqlserver_db_instance_ssl_config` resource with `instance_id`, `type`, `is_kms=1`, `key_id`, and `key_region`
- **THEN** the SSL configuration SHALL be updated with KMS encryption protection

### Requirement: Resource Create operation
The resource Create operation SHALL set the resource ID to the `instance_id` value, then delegate to the Update operation. This follows the RESOURCE_KIND_CONFIG pattern where the underlying resource always exists.

#### Scenario: Create triggers SSL enable
- **WHEN** user creates the resource with `type = "enable"`
- **THEN** the resource SHALL set ID to instance_id and call ModifyDBInstanceSSL with Type="enable"

#### Scenario: Create triggers SSL disable
- **WHEN** user creates the resource with `type = "disable"`
- **THEN** the resource SHALL set ID to instance_id and call ModifyDBInstanceSSL with Type="disable"

### Requirement: Resource Read operation
The resource Read operation SHALL call `DescribeDBInstancesAttribute` via `SqlserverService.DescribeSqlserverInstanceSslById` and read the `SSLConfig` field from the response. If `SSLConfig` is not nil, its sub-fields SHALL be set into the Terraform state.

#### Scenario: Read with SSL enabled
- **WHEN** the resource reads the current state and SSLConfig.Encryption is "enable"
- **THEN** `encryption` SHALL be set to "enable", and `ssl_validity_period` and `ssl_validity` SHALL be populated from the response

#### Scenario: Read with SSL disabled
- **WHEN** the resource reads the current state and SSLConfig.Encryption is "disable"
- **THEN** `encryption` SHALL be set to "disable"

#### Scenario: Read with instance not found
- **WHEN** the DescribeDBInstancesAttribute API returns an error indicating the instance does not exist
- **THEN** the resource SHALL set ID to empty string to indicate resource removal from state

### Requirement: Resource Update operation
The resource Update operation SHALL call `ModifyDBInstanceSSL` with the configured parameters. Since this is an async API returning `FlowId`, the operation SHALL poll `DescribeFlowStatus` until the flow status equals 0 (SQLSERVER_TASK_SUCCESS).

#### Scenario: Update SSL type to enable
- **WHEN** user updates `type` to "enable"
- **THEN** ModifyDBInstanceSSL SHALL be called with Type="enable", and the operation SHALL wait for the async flow to complete

#### Scenario: Update SSL type to renew
- **WHEN** user updates `type` to "renew"
- **THEN** ModifyDBInstanceSSL SHALL be called with Type="renew", and the operation SHALL wait for the async flow to complete

#### Scenario: Update with wait_switch
- **WHEN** user sets `wait_switch` to 1
- **THEN** ModifyDBInstanceSSL SHALL be called with WaitSwitch=1, indicating execution during maintenance window

#### Scenario: Async operation polling
- **WHEN** ModifyDBInstanceSSL returns a FlowId
- **THEN** the operation SHALL poll DescribeFlowStatus with the FlowId until Status == 0

### Requirement: Resource Delete operation
The resource Delete operation SHALL be a no-op, following the RESOURCE_KIND_CONFIG pattern. The underlying SQL Server instance continues to exist with its current SSL configuration.

#### Scenario: Delete resource
- **WHEN** user destroys the `tencentcloud_sqlserver_db_instance_ssl_config` resource
- **THEN** the resource SHALL be removed from Terraform state without making any API call

### Requirement: Resource Import support
The resource SHALL support Terraform import, allowing users to import an existing SSL configuration by specifying the `instance_id`.

#### Scenario: Import existing SSL config
- **WHEN** user runs `terraform import tencentcloud_sqlserver_db_instance_ssl_config.example mssql-instance-id`
- **THEN** the resource SHALL be imported with the current SSL configuration read from the API

### Requirement: Error handling and retry
The resource SHALL use `tccommon.WriteRetryTimeout` with retry for the ModifyDBInstanceSSL API call. If the API call fails, the error SHALL be wrapped using `tccommon.RetryError()`. The response SHALL be checked for nil values.

#### Scenario: API call retry on transient failure
- **WHEN** the ModifyDBInstanceSSL API call fails with a transient error
- **THEN** the operation SHALL retry within WriteRetryTimeout

#### Scenario: API returns nil response
- **WHEN** ModifyDBInstanceSSL returns nil response
- **THEN** a NonRetryableError SHALL be returned

#### Scenario: FlowId is nil in response
- **WHEN** ModifyDBInstanceSSL response has nil FlowId
- **THEN** an error SHALL be returned indicating FlowId is nil

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_sqlserver_db_instance_ssl_config` SHALL be available as a resource

### Requirement: Unit tests
The resource SHALL have unit tests using gomonkey mock approach (not Terraform test framework) covering Create, Read, and Update operations.

#### Scenario: Unit test for Create operation
- **WHEN** the Create function is called with valid parameters
- **THEN** it SHALL successfully create the resource by calling Update

#### Scenario: Unit test for Read operation
- **WHEN** the Read function is called
- **THEN** it SHALL correctly populate the state from the API response

#### Scenario: Unit test for Update operation
- **WHEN** the Update function is called
- **THEN** it SHALL correctly call ModifyDBInstanceSSL and wait for async completion
