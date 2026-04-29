# teo-security-api-service-resource Specification

## Purpose
TBD - created by archiving change add-teo-security-api-service-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource Schema Definition
The `tencentcloud_teo_security_api_service` resource SHALL define the following schema fields:
- `zone_id` (TypeString, Required, ForceNew): Zone ID of the TEO site
- `api_services` (TypeList, Required): List of API services to create, each containing:
  - `name` (TypeString, Required): API service name
  - `base_path` (TypeString, Required): Base path of the API service
- `api_service_ids` (TypeList of TypeString, Computed): API service IDs returned by the create operation
- `api_resources` (TypeList, Optional): List of API resources for update, each containing:
  - `id` (TypeString, Optional): Resource ID
  - `name` (TypeString, Optional): Resource name
  - `api_service_ids` (TypeList of TypeString, Optional): Associated API service IDs
  - `path` (TypeString, Optional): Resource path
  - `methods` (TypeList of TypeString, Optional): HTTP methods (GET, POST, PUT, HEAD, PATCH, OPTIONS, DELETE)
  - `request_constraint` (TypeString, Optional): Request content matching rule

#### Scenario: Schema validation on create
- **WHEN** a user creates a `tencentcloud_teo_security_api_service` resource with required fields `zone_id` and `api_services`
- **THEN** the resource SHALL accept the configuration and proceed with creation

#### Scenario: Missing required fields
- **WHEN** a user creates the resource without `zone_id` or `api_services`
- **THEN** Terraform SHALL report a validation error

### Requirement: Resource Create Operation
The resource SHALL call `CreateSecurityAPIService` API to create API services. The request SHALL include `ZoneId` and `APIServices` fields. Upon success, the response `APIServiceIds` SHALL be stored in the `api_service_ids` computed field. The resource ID SHALL be composed of `zone_id` and comma-joined `api_service_ids` separated by `FILED_SP`.

#### Scenario: Successful creation
- **WHEN** the create operation is called with valid `zone_id` and `api_services`
- **THEN** the `CreateSecurityAPIService` API SHALL be called with the corresponding parameters
- **AND** the resource ID SHALL be set to `zone_id:api_service_ids` (colon as `FILED_SP`)
- **AND** `api_service_ids` SHALL be populated from the response

#### Scenario: Create with retry on transient error
- **WHEN** the `CreateSecurityAPIService` API call fails with a transient error
- **THEN** the operation SHALL retry with `tccommon.WriteRetryTimeout`

### Requirement: Resource Read Operation
The resource SHALL call `DescribeSecurityAPIService` API to query API services. The request SHALL include `ZoneId` and pagination parameters (`Limit=100`, `Offset=0`). The response `APIServices` SHALL be mapped to the `api_services` field, filtered by `api_service_ids` from the resource ID.

#### Scenario: Successful read
- **WHEN** the read operation is called and the resource exists
- **THEN** the `DescribeSecurityAPIService` API SHALL be called with `ZoneId`
- **AND** `api_services` SHALL be populated from the response, filtered by `api_service_ids`

#### Scenario: Resource not found
- **WHEN** the read operation is called and the API services are not found
- **THEN** the resource ID SHALL be cleared (`d.SetId("")`)

#### Scenario: Read with retry on transient error
- **WHEN** the `DescribeSecurityAPIService` API call fails with a transient error
- **THEN** the operation SHALL retry with `tccommon.ReadRetryTimeout`

### Requirement: Resource Update Operation
The resource SHALL call `ModifySecurityAPIResource` API to update API resources when `api_resources` field changes. The `zone_id` field SHALL be immutable (ForceNew). The `api_services` field changes SHALL trigger a ForceNew since there is no Modify API for APIService itself.

#### Scenario: Update api_resources
- **WHEN** the `api_resources` field changes
- **THEN** the `ModifySecurityAPIResource` API SHALL be called with `ZoneId` and `APIResources`

#### Scenario: Immutable zone_id
- **WHEN** the `zone_id` field changes
- **THEN** Terraform SHALL force resource recreation

#### Scenario: Update with retry on transient error
- **WHEN** the `ModifySecurityAPIResource` API call fails with a transient error
- **THEN** the operation SHALL retry with `tccommon.WriteRetryTimeout`

### Requirement: Resource Delete Operation
The resource SHALL call `DeleteSecurityAPIService` API to delete API services. The request SHALL include `ZoneId` and `APIServiceIds` obtained from `d.Get()` rather than parsing `d.Id()`.

#### Scenario: Successful deletion
- **WHEN** the delete operation is called
- **THEN** the `DeleteSecurityAPIService` API SHALL be called with `ZoneId` and `APIServiceIds`

#### Scenario: Delete with retry on transient error
- **WHEN** the `DeleteSecurityAPIService` API call fails with a transient error
- **THEN** the operation SHALL retry with `tccommon.WriteRetryTimeout`

### Requirement: Resource Import
The resource SHALL support Terraform import via `schema.ImportStatePassthrough`. The imported ID SHALL be in the format `zone_id:api_service_ids`.

#### Scenario: Import existing resource
- **WHEN** a user imports an existing `tencentcloud_teo_security_api_service` resource
- **THEN** Terraform SHALL call the Read operation to populate the resource state

### Requirement: Provider Registration
The resource SHALL be registered in `provider.go` under `ResourcesMap` with key `tencentcloud_teo_security_api_service` and value `teo.ResourceTencentCloudTeoSecurityApiService()`. The `provider.md` file SHALL be updated to include the resource in the TEO section.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_security_api_service` SHALL be available as a resource type

### Requirement: Unit Tests
The resource SHALL have unit tests using gomonkey mock approach. Tests SHALL cover create, read, update, and delete operations. Tests SHALL be run with `go test -gcflags=all=-l`.

#### Scenario: Unit test for create operation
- **WHEN** the create function is tested
- **THEN** it SHALL mock the `CreateSecurityAPIService` API call and verify the resource state

#### Scenario: Unit test for read operation
- **WHEN** the read function is tested
- **THEN** it SHALL mock the `DescribeSecurityAPIService` API call and verify the resource state

#### Scenario: Unit test for update operation
- **WHEN** the update function is tested
- **THEN** it SHALL mock the `ModifySecurityAPIResource` API call and verify the resource state

#### Scenario: Unit test for delete operation
- **WHEN** the delete function is tested
- **THEN** it SHALL mock the `DeleteSecurityAPIService` API call and verify the resource is removed

### Requirement: Documentation
The resource SHALL have a `.md` documentation file following the gendoc/README.md format, including:
- A one-sentence description mentioning TEO product name
- Example Usage section with HCL configuration
- Import section (as RESOURCE_KIND_GENERAL resource)

#### Scenario: Documentation file exists
- **WHEN** the resource is created
- **THEN** a `resource_tc_teo_security_api_service.md` file SHALL exist in the `tencentcloud/services/teo/` directory
