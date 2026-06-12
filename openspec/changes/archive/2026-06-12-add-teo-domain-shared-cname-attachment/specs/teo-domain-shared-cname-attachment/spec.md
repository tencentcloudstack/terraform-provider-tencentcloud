## ADDED Requirements

### Requirement: Bind acceleration domains to shared CNAME
The system SHALL provide a Terraform resource `tencentcloud_teo_domain_shared_cname_attachment` that binds one or more acceleration domains to a shared CNAME within a TEO zone by calling the `BindSharedCNAME` API with `BindType = "bind"`.

#### Scenario: Successful binding creation
- **WHEN** user applies a `tencentcloud_teo_domain_shared_cname_attachment` resource with valid `zone_id` and `bind_shared_cname_maps` containing a `shared_cname` and `domain_names`
- **THEN** the system SHALL call `BindSharedCNAME` with `BindType = "bind"` and the provided parameters, and set the resource ID as a composite of `zone_id`, `shared_cname`, and `domain_names` joined by `tccommon.FILED_SP`

#### Scenario: API returns error on create
- **WHEN** the `BindSharedCNAME` API returns an error during resource creation
- **THEN** the system SHALL wrap the error with `tccommon.RetryError()` and retry within `tccommon.ReadRetryTimeout`

### Requirement: Read binding status via DescribeSharedCNAME
The system SHALL verify the binding relationship exists by calling `DescribeSharedCNAME` with a filter on the shared CNAME and checking that the expected domains appear in the response's `AccelerationDomains` field.

#### Scenario: Binding exists
- **WHEN** the Read function is called and `DescribeSharedCNAME` returns a `SharedCNAMEInfo` entry whose `AccelerationDomains` contains the expected domain names
- **THEN** the system SHALL set the resource attributes from the stored state

#### Scenario: Binding does not exist
- **WHEN** the Read function is called and the expected domains are not found in the `AccelerationDomains` of the shared CNAME
- **THEN** the system SHALL log the resource ID, then call `d.SetId("")` to remove the resource from state

#### Scenario: API returns empty response on read
- **WHEN** `DescribeSharedCNAME` returns an empty response or the shared CNAME is not found
- **THEN** the system SHALL log the resource ID, then call `d.SetId("")` to remove the resource from state

### Requirement: Unbind acceleration domains from shared CNAME on delete
The system SHALL unbind the acceleration domains from the shared CNAME by calling `BindSharedCNAME` with `BindType = "unbind"` when the resource is destroyed.

#### Scenario: Successful unbinding
- **WHEN** user destroys the `tencentcloud_teo_domain_shared_cname_attachment` resource
- **THEN** the system SHALL parse the composite ID to extract `zone_id`, `shared_cname`, and `domain_names`, then call `BindSharedCNAME` with `BindType = "unbind"`

#### Scenario: API returns error on delete
- **WHEN** the `BindSharedCNAME` API returns an error during resource deletion
- **THEN** the system SHALL wrap the error with `tccommon.RetryError()` and retry within `tccommon.ReadRetryTimeout`

### Requirement: Resource import support
The system SHALL support importing existing bindings using the composite ID format `{zone_id}#{shared_cname}#{domain_name1,domain_name2,...}`.

#### Scenario: Import existing binding
- **WHEN** user runs `terraform import tencentcloud_teo_domain_shared_cname_attachment.example zone-xxx#shared.example.com#domain1.example.com,domain2.example.com`
- **THEN** the system SHALL parse the composite ID and populate the resource state by calling the Read function

### Requirement: All fields are ForceNew
Since this is a CRD-only resource (no Update API), all top-level fields SHALL be marked as `ForceNew`. Any change to `zone_id` or `bind_shared_cname_maps` SHALL trigger resource recreation (destroy + create).

#### Scenario: Field change triggers recreation
- **WHEN** user modifies any field in the resource configuration
- **THEN** Terraform SHALL plan a destroy-and-recreate operation

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Resource is available in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_domain_shared_cname_attachment` SHALL be available as a resource type
