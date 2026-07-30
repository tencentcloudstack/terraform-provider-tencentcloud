## ADDED Requirements

### Requirement: Bind acceleration domains to shared CNAME
The system SHALL provide a Terraform resource `tencentcloud_teo_domain_shared_cname_attachment` that binds one or more acceleration domains to a shared CNAME within a TEO zone by calling the `BindSharedCNAMEWithContext` API with `BindType = "bind"`.

The resource represents a single `shared_cname` binding within a zone. The schema exposes `zone_id`, `shared_cname`, and `domain_names` as top-level fields (not nested under `bind_shared_cname_maps`).

#### Scenario: Successful binding creation
- **WHEN** user applies a `tencentcloud_teo_domain_shared_cname_attachment` resource with valid `zone_id`, `shared_cname`, and `domain_names`
- **THEN** the system SHALL call `BindSharedCNAMEWithContext` with `BindType = "bind"` and the provided parameters, and set the resource ID as `{zone_id}#{shared_cname}`

#### Scenario: API returns error on create
- **WHEN** the `BindSharedCNAMEWithContext` API returns an error during resource creation
- **THEN** the system SHALL wrap the error with `tccommon.RetryError()` and retry within `tccommon.WriteRetryTimeout`

### Requirement: Read binding status via DescribeSharedCNAME
The system SHALL read the current binding state by calling `DescribeSharedCNAMEWithContext` with a filter on the shared CNAME and populating `domain_names` from the `AccelerationDomains` field of the response.

#### Scenario: Binding exists
- **WHEN** the Read function is called and `DescribeSharedCNAMEWithContext` returns a `SharedCNAMEInfo` entry matching the shared CNAME
- **THEN** the system SHALL set `zone_id`, `shared_cname`, and `domain_names` from the API response

#### Scenario: Binding does not exist
- **WHEN** `DescribeSharedCNAMEWithContext` returns an empty response or the shared CNAME is not found
- **THEN** the system SHALL log the resource ID, then call `d.SetId("")` to remove the resource from state

### Requirement: Update domain_names by unbinding removed and binding added domains
The system SHALL support updating `domain_names` in-place without recreating the resource. When `domain_names` changes, the system SHALL:
1. Unbind domains that were in the old list but not in the new list (call `BindSharedCNAMEWithContext` with `BindType = "unbind"`)
2. Bind domains that are in the new list but not in the old list (call `BindSharedCNAMEWithContext` with `BindType = "bind"`)

#### Scenario: Domains removed from list
- **WHEN** user removes one or more domain names from `domain_names`
- **THEN** the system SHALL call `BindSharedCNAMEWithContext` with `BindType = "unbind"` for the removed domains

#### Scenario: Domains added to list
- **WHEN** user adds one or more domain names to `domain_names`
- **THEN** the system SHALL call `BindSharedCNAMEWithContext` with `BindType = "bind"` for the added domains

#### Scenario: API returns error on update
- **WHEN** the `BindSharedCNAMEWithContext` API returns an error during update
- **THEN** the system SHALL wrap the error with `tccommon.RetryError()` and retry within `tccommon.WriteRetryTimeout`

### Requirement: Unbind all acceleration domains from shared CNAME on delete
The system SHALL unbind all acceleration domains from the shared CNAME by calling `BindSharedCNAMEWithContext` with `BindType = "unbind"` when the resource is destroyed.

#### Scenario: Successful unbinding
- **WHEN** user destroys the `tencentcloud_teo_domain_shared_cname_attachment` resource
- **THEN** the system SHALL parse the composite ID to extract `zone_id` and `shared_cname`, read `domain_names` from state, then call `BindSharedCNAMEWithContext` with `BindType = "unbind"`

#### Scenario: API returns error on delete
- **WHEN** the `BindSharedCNAMEWithContext` API returns an error during resource deletion
- **THEN** the system SHALL wrap the error with `tccommon.RetryError()` and retry within `tccommon.WriteRetryTimeout`

### Requirement: Resource import support
The system SHALL support importing existing bindings using the composite ID format `{zone_id}#{shared_cname}`.

#### Scenario: Import existing binding
- **WHEN** user runs `terraform import tencentcloud_teo_domain_shared_cname_attachment.example zone-xxx#shared.example.com`
- **THEN** the system SHALL parse the composite ID and populate the resource state by calling the Read function, which reads `domain_names` from the API

### Requirement: zone_id and shared_cname are ForceNew; domain_names supports in-place update
`zone_id` and `shared_cname` SHALL be marked as `ForceNew`. Changes to `domain_names` SHALL trigger an in-place update (not recreation).

#### Scenario: zone_id or shared_cname change triggers recreation
- **WHEN** user modifies `zone_id` or `shared_cname`
- **THEN** Terraform SHALL plan a destroy-and-recreate operation

#### Scenario: domain_names change triggers update
- **WHEN** user modifies `domain_names`
- **THEN** Terraform SHALL plan an in-place update operation

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Resource is available in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_domain_shared_cname_attachment` SHALL be available as a resource type
