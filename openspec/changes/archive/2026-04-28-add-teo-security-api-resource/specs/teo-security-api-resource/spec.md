## ADDED Requirements

### Requirement: Resource schema definition
The `tencentcloud_teo_security_api_resource` resource SHALL define the following schema fields:
- `zone_id` (string, required, ForceNew): TEO 站点 ID
- `api_resources` (TypeList of nested blocks, required, MaxItems: 1): API 资源配置，每次仅允许一个
  - `name` (string, required): API 资源名称
  - `path` (string, required): API 资源路径，如 /api/v1/orders
  - `api_service_ids` (TypeList of string, optional): API 资源关联的 API 服务 ID 列表
  - `methods` (TypeList of string, optional): 请求方法列表，支持 GET, POST, PUT, HEAD, PATCH, OPTIONS, DELETE
  - `request_constraint` (string, optional): 请求内容匹配规则表达式
  - `id` (string, computed): API 资源 ID，创建时由服务端生成，如 apires-xxxxxxxx

The resource SHALL support import via `schema.ImportStatePassthrough`.
The resource ID SHALL use composite ID format `zoneId#apiResourceId` (using tccommon.FILED_SP as separator).

#### Scenario: Schema fields are correctly defined
- **WHEN** the resource schema is registered
- **THEN** all fields with correct types, required/computed status, and ForceNew attributes are defined
- **AND** api_resources has MaxItems: 1
- **AND** path is required

### Requirement: Create API resource
The resource SHALL call `CreateSecurityAPIResource` API to create an API resource under the specified zone.

#### Scenario: Successful creation
- **WHEN** a `tencentcloud_teo_security_api_resource` resource is created with valid `zone_id` and `api_resources`
- **THEN** the system SHALL call CreateSecurityAPIResource with ZoneId and a single APIResource (using buildSecurityAPIResourceFromMap with id="")
- **AND** the returned APIResourceIds[0] SHALL be combined with zoneId as composite ID using FILED_SP separator
- **AND** a Read operation SHALL be performed after creation to sync state

#### Scenario: Creation with retry
- **WHEN** the CreateSecurityAPIResource API call fails with a retryable error
- **THEN** the system SHALL retry the operation with `tccommon.WriteRetryTimeout`

#### Scenario: Empty APIResourceIds
- **WHEN** the CreateSecurityAPIResource response has empty APIResourceIds
- **THEN** the system SHALL return an error indicating APIResourceIds is empty

### Requirement: Read API resource
The resource SHALL call `DescribeTeoSecurityAPIResourceById` service method to query a specific API resource.

#### Scenario: Successful read
- **WHEN** a Read operation is performed for an existing resource
- **THEN** the system SHALL parse the composite ID to extract zoneId and apiResourceId
- **AND** call DescribeTeoSecurityAPIResourceById with zoneId and apiResourceId
- **AND** populate the api_resources nested block from the returned APIResource
- **AND** set zone_id from the parsed zoneId

#### Scenario: Resource not found
- **WHEN** a Read operation is performed and DescribeTeoSecurityAPIResourceById returns nil
- **THEN** the resource ID SHALL be cleared (`d.SetId("")`) and nil error SHALL be returned

#### Scenario: Invalid composite ID
- **WHEN** the composite ID format is invalid (not exactly 2 parts)
- **THEN** the system SHALL return an error indicating the ID is broken

### Requirement: Update API resource
The resource SHALL call `ModifySecurityAPIResource` API to update the API resource.

#### Scenario: Successful update
- **WHEN** a `tencentcloud_teo_security_api_resource` resource is updated
- **THEN** the system SHALL parse the composite ID to extract zoneId and apiResourceId
- **AND** call ModifySecurityAPIResource with ZoneId and a single APIResource (using buildSecurityAPIResourceFromMap with the apiResourceId)
- **AND** a Read operation SHALL be performed after update to sync state

#### Scenario: Update with retry
- **WHEN** the ModifySecurityAPIResource API call fails with a retryable error
- **THEN** the system SHALL retry the operation with `tccommon.WriteRetryTimeout`

### Requirement: Delete API resource
The resource SHALL call `DeleteSecurityAPIResource` API to delete the API resource.

#### Scenario: Successful deletion
- **WHEN** a `tencentcloud_teo_security_api_resource` resource is destroyed
- **THEN** the system SHALL parse the composite ID to extract zoneId and apiResourceId
- **AND** call DeleteSecurityAPIResource with ZoneId and APIResourceIds containing the single apiResourceId

#### Scenario: Delete with retry
- **WHEN** the DeleteSecurityAPIResource API call fails with a retryable error
- **THEN** the system SHALL retry the operation with `tccommon.WriteRetryTimeout`

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` under ResourcesMap with key `tencentcloud_teo_security_api_resource` and value `teo.ResourceTencentCloudTeoSecurityAPIResource()`.

The resource SHALL be listed in `tencentcloud/provider.md` under the TEO section.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_security_api_resource` SHALL be available as a resource type

### Requirement: Service layer method
A `DescribeTeoSecurityAPIResourceById` method SHALL be added to the TeoService in `service_tencentcloud_teo.go`.

#### Scenario: Service method returns API resource
- **WHEN** `DescribeTeoSecurityAPIResourceById` is called with a valid zoneId and apiResourceId
- **THEN** the method SHALL paginate through DescribeSecurityAPIResource with Limit=100
- **AND** match the apiResourceId against each APIResource's Id field
- **AND** return the matching `*APIResource`
- **AND** return nil if no matching API resource is found

### Requirement: Helper function
A `buildSecurityAPIResourceFromMap` helper function SHALL be implemented to convert schema map to *teo.APIResource.

#### Scenario: Helper function with empty id
- **WHEN** called with id=""
- **THEN** the resulting APIResource SHALL NOT have the Id field set

#### Scenario: Helper function with actual id
- **WHEN** called with a non-empty id string
- **THEN** the resulting APIResource SHALL have the Id field set to the given id

### Requirement: Unit tests
Unit tests SHALL be created using gomonkey mock approach (not Terraform test suite) for the resource.

#### Scenario: Unit test coverage
- **WHEN** unit tests are executed with `go test -gcflags=all=-l`
- **THEN** all CRUD operations SHALL be tested with mocked cloud API responses
- **AND** schema definition SHALL be verified including MaxItems, Required fields, and Computed fields
