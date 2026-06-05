## ADDED Requirements

### Requirement: Apply free certificate operation resource

The system SHALL provide a Terraform resource `tencentcloud_teo_apply_free_certificate` of type RESOURCE_KIND_OPERATION that calls the TEO ApplyFreeCertificate API to initiate a free certificate application for a specified domain.

#### Scenario: Successful certificate application with DNS verification
- **WHEN** user creates a `tencentcloud_teo_apply_free_certificate` resource with `zone_id`, `domain`, and `verification_method` set to `dns_challenge`
- **THEN** the system SHALL call the ApplyFreeCertificate API and store the returned `dns_verification` (subdomain, record_type, record_value) in the resource state

#### Scenario: Successful certificate application with HTTP file verification
- **WHEN** user creates a `tencentcloud_teo_apply_free_certificate` resource with `zone_id`, `domain`, and `verification_method` set to `http_challenge`
- **THEN** the system SHALL call the ApplyFreeCertificate API and store the returned `file_verification` (path, content) in the resource state

#### Scenario: API call failure with retry
- **WHEN** the ApplyFreeCertificate API call fails due to a transient error
- **THEN** the system SHALL retry the API call using `tccommon.ReadRetryTimeout` and `tccommon.RetryError` wrapping

### Requirement: Resource schema definition

The resource SHALL define the following schema fields:
- `zone_id`: Required, ForceNew, String - The zone ID
- `domain`: Required, ForceNew, String - The target domain for the free certificate
- `verification_method`: Required, ForceNew, String - The verification method (http_challenge or dns_challenge)
- `dns_verification`: Computed, List(MaxItems:1) - DNS verification info with subdomain, record_type, record_value
- `file_verification`: Computed, List(MaxItems:1) - File verification info with path, content

#### Scenario: Schema validation
- **WHEN** user provides all required fields (zone_id, domain, verification_method)
- **THEN** the resource SHALL accept the configuration and proceed with creation

### Requirement: OPERATION resource lifecycle

The resource SHALL implement OPERATION lifecycle:
- Create: Call ApplyFreeCertificate API, set resource ID, store verification results
- Read: No-op (return nil)
- Delete: No-op (return nil)

#### Scenario: Resource deletion
- **WHEN** user destroys the resource
- **THEN** the system SHALL perform no API calls and simply remove the resource from state

### Requirement: Provider registration

The resource SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Resource available in provider
- **WHEN** user references `tencentcloud_teo_apply_free_certificate` in their Terraform configuration
- **THEN** the provider SHALL recognize and handle the resource type
