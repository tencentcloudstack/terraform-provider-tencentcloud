## ADDED Requirements

### Requirement: Resource schema definition
The `tencentcloud_teo_security_api_resource` resource SHALL define the following schema fields:
- `zone_id` (string, required, ForceNew): TEO 站点 ID
- `api_resources` (TypeList of nested blocks, required): API 资源列表
  - `id` (string, computed): API 资源 ID，创建时由服务端生成
  - `name` (string, required): API 资源名称
  - `api_service_ids` (TypeList of string, optional): API 资源关联的 API 服务 ID 列表
  - `path` (string, optional): 资源路径
  - `methods` (TypeList of string, optional): 请求方法列表，支持 GET, POST, PUT, HEAD, PATCH, OPTIONS, DELETE
  - `request_constraint` (string, optional): 请求内容匹配规则
- `api_resource_ids` (TypeList of string, computed): 创建后由服务端返回的 API 资源 ID 列表

The resource SHALL support import via `schema.ImportStatePassthrough`.

#### Scenario: Schema fields are correctly defined
- **WHEN** the resource schema is registered
- **THEN** all fields with correct types, required/computed status, and ForceNew attributes are defined

### Requirement: Create API resource
The resource SHALL call `CreateSecurityAPIResource` API to create API resources under the specified zone.

#### Scenario: Successful creation
- **WHEN** a `tencentcloud_teo_security_api_resource` resource is created with valid `zone_id` and `api_resources`
- **THEN** the system SHALL call CreateSecurityAPIResource with ZoneId and APIResources parameters
- **AND** the returned APIResourceIds SHALL be stored in `api_resource_ids` field
- **AND** the resource ID SHALL be set to the `zone_id` value
- **AND** a Read operation SHALL be performed after creation to sync state

#### Scenario: Creation with retry
- **WHEN** the CreateSecurityAPIResource API call fails with a retryable error
- **THEN** the system SHALL retry the operation with `tccommon.WriteRetryTimeout`

### Requirement: Read API resource
The resource SHALL call `DescribeSecurityAPIResource` API to query all API resources under the specified zone with pagination (Limit=100, looping until all records are fetched).

#### Scenario: Successful read
- **WHEN** a Read operation is performed for an existing resource
- **THEN** the system SHALL call DescribeSecurityAPIResource with ZoneId
- **AND** the system SHALL paginate with Limit=100 until TotalCount records are retrieved
- **AND** all API resource fields SHALL be populated from the response
- **AND** `api_resource_ids` SHALL be populated from the Id field of each APIResource in the response

#### Scenario: Resource not found
- **WHEN** a Read operation is performed and the DescribeSecurityAPIResource returns no API resources for the zone
- **THEN** the resource ID SHALL be cleared (`d.SetId("")`) and nil error SHALL be returned

#### Scenario: Read with retry
- **WHEN** the DescribeSecurityAPIResource API call fails with a retryable error
- **THEN** the system SHALL retry the operation with `tccommon.ReadRetryTimeout`

### Requirement: Update API resource
The resource SHALL call `ModifySecurityAPIResource` API to update API resources under the specified zone.

#### Scenario: Successful update
- **WHEN** a `tencentcloud_teo_security_api_resource` resource is updated and `api_resources` has changed
- **THEN** the system SHALL call ModifySecurityAPIResource with ZoneId and the updated APIResources list
- **AND** each APIResource in the request SHALL include the Id field from the current state
- **AND** a Read operation SHALL be performed after update to sync state

#### Scenario: Update with retry
- **WHEN** the ModifySecurityAPIResource API call fails with a retryable error
- **THEN** the system SHALL retry the operation with `tccommon.WriteRetryTimeout`

### Requirement: Delete API resource
The resource SHALL call `DeleteSecurityAPIResource` API to delete all API resources under the specified zone.

#### Scenario: Successful deletion
- **WHEN** a `tencentcloud_teo_security_api_resource` resource is destroyed
- **THEN** the system SHALL call DeleteSecurityAPIResource with ZoneId and APIResourceIds obtained from `d.Get("api_resource_ids")`
- **AND** the system SHALL NOT use d.Id() to obtain api_resource_ids

#### Scenario: Delete with retry
- **WHEN** the DeleteSecurityAPIResource API call fails with a retryable error
- **THEN** the system SHALL retry the operation with `tccommon.WriteRetryTimeout`

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` under ResourcesMap with key `tencentcloud_teo_security_api_resource` and value `teo.ResourceTencentCloudTeoSecurityApiResource()`.

The resource SHALL be listed in `tencentcloud/provider.md` under the TEO section.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_security_api_resource` SHALL be available as a resource type

### Requirement: Service layer method
A `DescribeSecurityAPIResourceById` method SHALL be added to the TeoService in `service_tencentcloud_teo.go`.

#### Scenario: Service method returns API resources
- **WHEN** `DescribeSecurityAPIResourceById` is called with a valid zoneId
- **THEN** the method SHALL paginate through DescribeSecurityAPIResource with Limit=100
- **AND** return all API resources for the zone
- **AND** return nil if no API resources exist

### Requirement: Unit tests
Unit tests SHALL be created using gomonkey mock approach (not Terraform test suite) for the resource.

#### Scenario: Unit test coverage
- **WHEN** unit tests are executed with `go test -gcflags=all=-l`
- **THEN** all CRUD operations SHALL be tested with mocked cloud API responses
