# bdrc-disaster-recovery-site-pair-resource Specification

## Purpose
TBD - created by archiving change add-bdrc-disaster-recovery-site-pair-resource. Update Purpose after archive.
## Requirements
### Requirement: Create a BDRC disaster recovery site pair

The resource `tencentcloud_bdrc_disaster_recovery_site_pair` SHALL create a site pair by calling `CreateDisasterRecoverySitePair` with the required parameters (`disaster_recovery_type`, `source_region`, `source_zone`, `target_region`, `target_zone`, `source_vpc`, `target_vpc`, `site_pair_product_type`) and optional parameters (`site_pair_name`, `copy_type`). After a successful creation, the resource SHALL set `d.Id()` to the returned `SitePairId` and call Read to populate state.

#### Scenario: Successful creation

- **WHEN** the user applies a configuration with all required fields populated
- **THEN** the provider calls `CreateDisasterRecoverySitePair`, sets the resource ID to the returned `SitePairId`, and reads back all attributes

#### Scenario: Creation returns an empty SitePairId

- **WHEN** `CreateDisasterRecoverySitePair` returns successfully but `SitePairId` is nil or empty
- **THEN** the provider returns a `NonRetryableError` to prevent writing an empty ID to state

### Requirement: Read a BDRC disaster recovery site pair

The resource SHALL read the current state of a site pair by calling `DescribeDisasterRecoverySitePairs` with `SitePairIds` set to the resource ID and `SitePairType` derived from the `site_pair_product_type` field in state. The response `SitePairSet[0]` (first matching element) SHALL be used to populate all schema fields. If the site pair no longer exists, the provider SHALL log `[CRUD] tencentcloud_bdrc_disaster_recovery_site_pair id=<id>` and then set `d.SetId("")`.

#### Scenario: Site pair exists

- **WHEN** the provider reads the resource and the API returns a matching `SitePair`
- **THEN** all schema fields (including computed nested fields like `protected_resource_set`, `protected_resource_status_set`, `cross_cloud_details`) are populated from the response

#### Scenario: Site pair not found

- **WHEN** the API returns an empty `SitePairSet` for the given `SitePairId`
- **THEN** the provider logs the ID for diagnostics and removes the resource from state via `d.SetId("")`

### Requirement: Update a BDRC disaster recovery site pair

The resource SHALL support updating `site_pair_name` by calling `ModifySitePairAttribute` with `SitePairId` (from `d.Id()`) and the new `SitePairName`. All other parameters are immutable (`ForceNew`) and cannot be changed without recreating the resource.

#### Scenario: Updating the site pair name

- **WHEN** the user changes `site_pair_name` in the configuration
- **THEN** the provider calls `ModifySitePairAttribute` with the new name and re-reads the resource

#### Scenario: Changing an immutable field

- **WHEN** the user changes any `ForceNew` field (e.g., `disaster_recovery_type`, `source_region`, `target_vpc`)
- **THEN** Terraform recreates the resource (destroy + create) because the field is marked `ForceNew`

### Requirement: Delete a BDRC disaster recovery site pair

The resource SHALL delete a site pair by calling `DeleteDisasterRecoverySitePairs` with `SitePairIds` set to a single-element list containing the resource ID.

#### Scenario: Successful deletion

- **WHEN** the user destroys the resource
- **THEN** the provider calls `DeleteDisasterRecoverySitePairs` with the `SitePairId` and the resource is removed from state

### Requirement: Import a BDRC disaster recovery site pair

The resource SHALL support Terraform import using the `SitePairId` as the import ID.

#### Scenario: Import by SitePairId

- **WHEN** the user runs `terraform import tencentcloud_bdrc_disaster_recovery_site_pair.example <site_pair_id>`
- **THEN** the provider sets `d.Id()` to the imported ID and calls Read to populate the full state

### Requirement: Register the resource in the provider

The provider SHALL register `tencentcloud_bdrc_disaster_recovery_site_pair` in the `ResourcesMap` of `provider.go` and list it under the BDRC section of `provider.md`.

#### Scenario: Provider registration

- **WHEN** the provider is initialized
- **THEN** `tencentcloud_bdrc_disaster_recovery_site_pair` appears in the resource map and is usable in Terraform configurations

### Requirement: BDRC SDK client connectivity

The provider SHALL add a `UseBdrcV20260330Client()` method to `connectivity.TencentCloudClient` that lazily initializes and returns the BDRC SDK client (`bdrc/v20260330.Client`), mirroring the existing `UseIgtmV20231024Client` pattern.

#### Scenario: Client initialization

- **WHEN** a BDRC resource CRUD operation executes
- **THEN** the provider obtains the BDRC client via `UseBdrcV20260330Client()` and uses it to call BDRC cloud APIs with retry and error handling

