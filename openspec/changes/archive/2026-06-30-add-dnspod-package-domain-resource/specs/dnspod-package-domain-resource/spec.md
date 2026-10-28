## ADDED Requirements

### Requirement: Create package domain binding
The system SHALL allow users to bind a DNSPod package to a domain using the `tencentcloud_dnspod_package_domain` resource.

#### Scenario: Successful bind
- **WHEN** user specifies a valid `resource_id` and `domain_id`
- **THEN** the system calls `ModifyPackageDomain` with Operation="bind", ResourceId=resource_id, NewDomainId=domain_id
- **THEN** the resource ID is set to `{resource_id}#{domain_id}`
- **THEN** the system reads back and populates all Computed fields from `DescribeDomainVipList`

#### Scenario: Invalid parameters
- **WHEN** user specifies a non-existent `resource_id` or `domain_id`
- **THEN** the system returns an error from the API and does not set the resource ID

### Requirement: Read package domain binding status
The system SHALL query the package domain binding status using `DescribeDomainVipList` API.

#### Scenario: Binding exists
- **WHEN** `DescribeDomainVipList` returns a `PackageListItem` with matching `DomainId`
- **THEN** the system populates `domain`, `grade`, `grade_title`, `vip_start_at`, `vip_end_at`, `vip_auto_renew`, `remain_times`, `grade_level`, `status`, `is_grace_period`, and `downgrade` fields

#### Scenario: Binding does not exist
- **WHEN** `DescribeDomainVipList` returns no matching `PackageListItem`
- **THEN** the system clears the resource ID via `d.SetId("")` and logs a warning

#### Scenario: API returns nil response
- **WHEN** `DescribeDomainVipList` returns a nil response
- **THEN** the system logs the resource ID for debugging and clears the resource ID

### Requirement: Update package domain binding (change domain)
The system SHALL support changing the domain bound to a package via `ModifyPackageDomain` with Operation="change".

#### Scenario: Successful domain change
- **WHEN** user changes the `domain_id` field
- **THEN** the system calls `ModifyPackageDomain` with Operation="change", ResourceId=resource_id, DomainId=old_domain_id, NewDomainId=new_domain_id
- **THEN** the resource ID is updated to reflect the new domain_id
- **THEN** the system reads back the updated state

#### Scenario: resource_id change triggers ForceNew
- **WHEN** user changes the `resource_id` field
- **THEN** Terraform destroys the existing resource and creates a new one (ForceNew behavior)

### Requirement: Delete package domain binding (unbind)
The system SHALL support unbinding a domain from a package via `ModifyPackageDomain` with Operation="unbind".

#### Scenario: Successful unbind
- **WHEN** user deletes the resource
- **THEN** the system calls `ModifyPackageDomain` with Operation="unbind", ResourceId=resource_id, DomainId=domain_id
- **THEN** the resource is removed from Terraform state

### Requirement: Import existing package domain binding
The system SHALL support importing existing package domain bindings using the format `{resource_id}#{domain_id}`.

#### Scenario: Successful import
- **WHEN** user runs `terraform import tencentcloud_dnspod_package_domain.example "res-xxxxx#12345"`
- **THEN** the system parses the ID, calls `DescribeDomainVipList` to read the binding state, and populates all fields
