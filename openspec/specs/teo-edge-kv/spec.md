## ADDED Requirements

### Requirement: Resource schema definition

The resource `tencentcloud_teo_edge_kv` SHALL define the following schema:

- `zone_id` (Required, ForceNew, String): 站点 ID
- `namespace` (Required, ForceNew, String): 命名空间名称
- `key` (Required, ForceNew, String): 键名，长度为 1-512 个字符，允许的字符为字母、数字、中划线和下划线
- `value` (Required, String): 键值，不能为空，最大支持 1 MB
- `expiration` (Optional, Int): 键值对的过期时间，绝对时间，Unix 时间戳，取值必须大于等于当前时间 + 60
- `expiration_ttl` (Optional, Int): 键值对的存活时长，相对时间，单位为秒，取值范围大于等于 60

The resource ID SHALL be composed of `zone_id`, `namespace`, and `key` joined by `tccommon.FILED_SP` separator.

#### Scenario: Schema fields are correctly defined
- **WHEN** user defines a `tencentcloud_teo_edge_kv` resource in HCL
- **THEN** Terraform SHALL validate that `zone_id`, `namespace`, `key`, and `value` are provided as required fields

#### Scenario: ForceNew triggers resource recreation
- **WHEN** user changes `zone_id`, `namespace`, or `key` in the configuration
- **THEN** Terraform SHALL destroy the existing resource and create a new one

### Requirement: Create operation writes KV data

The resource Create method SHALL call the `EdgeKVPut` API with the following parameter mapping:
- `zone_id` → `request.ZoneId`
- `namespace` → `request.Namespace`
- `key` → `request.Key`
- `value` → `request.Value`
- `expiration` → `request.Expiration` (if set)
- `expiration_ttl` → `request.ExpirationTTL` (if set)

The Create method SHALL use `tccommon.ReadRetryTimeout` for retry logic. After successful creation, the resource ID SHALL be set to `zone_id#namespace#key`.

#### Scenario: Successful KV data write
- **WHEN** user applies a configuration with valid `zone_id`, `namespace`, `key`, and `value`
- **THEN** the resource SHALL call EdgeKVPut API and set the composite ID

#### Scenario: Create with expiration parameters
- **WHEN** user specifies `expiration` or `expiration_ttl`
- **THEN** the resource SHALL pass these values to the EdgeKVPut API

#### Scenario: Create API returns error
- **WHEN** EdgeKVPut API returns an error
- **THEN** the resource SHALL wrap the error with `tccommon.RetryError()` and return it

### Requirement: Read operation queries KV data

The resource Read method SHALL call the `EdgeKVGet` API with:
- `zone_id` → `request.ZoneId`
- `namespace` → `request.Namespace`
- `[key]` → `request.Keys` (single-element array)

The Read method SHALL parse the composite ID to extract `zone_id`, `namespace`, and `key`. If the returned Data is empty or the key's Value is empty string, the resource SHALL be marked as deleted by calling `d.SetId("")`.

#### Scenario: Successful KV data read
- **WHEN** the KV key exists
- **THEN** the resource SHALL set `value` from the response Data

#### Scenario: KV key does not exist
- **WHEN** EdgeKVGet returns empty Data or empty Value for the key
- **THEN** the resource SHALL call `d.SetId("")` to mark it as deleted

#### Scenario: Read parses composite ID
- **WHEN** Read method is called
- **THEN** it SHALL split `d.Id()` by `tccommon.FILED_SP` to extract `zone_id`, `namespace`, and `key`

### Requirement: Delete operation removes KV data

The resource Delete method SHALL call the `EdgeKVDelete` API with:
- `zone_id` → `request.ZoneId`
- `namespace` → `request.Namespace`
- `[key]` → `request.Keys` (single-element array)

The Delete method SHALL use `tccommon.ReadRetryTimeout` for retry logic.

#### Scenario: Successful KV data deletion
- **WHEN** user destroys the resource
- **THEN** the resource SHALL call EdgeKVDelete API to remove the KV data

#### Scenario: Delete API returns error
- **WHEN** EdgeKVDelete API returns an error
- **THEN** the resource SHALL wrap the error with `tccommon.RetryError()` and return it

### Requirement: Update method enforces immutability

The resource Update method SHALL check if any of `value`, `expiration`, `expiration_ttl` have changed. If any of these fields are in the change set, the Update method SHALL return an error indicating that the field cannot be modified in-place.

#### Scenario: Attempt to modify immutable field
- **WHEN** user changes `value`, `expiration`, or `expiration_ttl` without changing ForceNew fields
- **THEN** the resource SHALL return an error stating the field is immutable

### Requirement: Import support

The resource SHALL support import using the composite ID format `zone_id#namespace#key`.

#### Scenario: Import existing KV data
- **WHEN** user runs `terraform import tencentcloud_teo_edge_kv.example zone_id#namespace#key`
- **THEN** the resource SHALL read the KV data and populate the state

### Requirement: Provider registration

The resource SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Resource is available in provider
- **WHEN** user references `tencentcloud_teo_edge_kv` in their configuration
- **THEN** Terraform SHALL recognize it as a valid resource type

### Requirement: Unit tests with gomonkey mock

The resource SHALL have unit tests in `resource_tc_teo_edge_kv_test.go` that use gomonkey to mock the cloud API calls. Tests SHALL cover Create, Read, and Delete operations.

#### Scenario: Unit test for Create
- **WHEN** unit test for Create is executed
- **THEN** it SHALL mock EdgeKVPut API and verify the resource is created correctly

#### Scenario: Unit test for Read
- **WHEN** unit test for Read is executed
- **THEN** it SHALL mock EdgeKVGet API and verify state is populated correctly

#### Scenario: Unit test for Delete
- **WHEN** unit test for Delete is executed
- **THEN** it SHALL mock EdgeKVDelete API and verify the resource is deleted correctly
