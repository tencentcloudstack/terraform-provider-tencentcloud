## ADDED Requirements

### Requirement: Support dedicated resource pack placement parameters

The `tencentcloud_instance` resource SHALL support specifying dedicated resource pack placement parameters to enable instance creation using pre-purchased resource pool packs.

#### Scenario: Create instance with dedicated resource pack

- **WHEN** user specifies both `dedicated_resource_pack_tenancy` and `dedicated_resource_pack_ids` in the instance configuration
- **THEN** the provider SHALL pass `Placement.DedicatedResourcePackTenancy` and `Placement.DedicatedResourcePackIds` to the `RunInstances` API
- **AND** the instance SHALL be created using the specified resource pool pack

#### Scenario: Validation when only one parameter is provided

- **WHEN** user specifies `dedicated_resource_pack_tenancy` without `dedicated_resource_pack_ids` (or vice versa)
- **THEN** Terraform SHALL return a validation error during plan phase
- **AND** the error message SHALL indicate that both parameters must be specified together

#### Scenario: ForceNew behavior on parameter change

- **WHEN** user modifies `dedicated_resource_pack_tenancy` or `dedicated_resource_pack_ids` in an existing instance configuration
- **THEN** Terraform plan SHALL show the instance will be destroyed and recreated
- **AND** the new instance SHALL use the new resource pack parameters

### Requirement: Schema field definitions

The resource schema SHALL include two new optional fields for dedicated resource pack placement.

#### Scenario: dedicated_resource_pack_tenancy field specification

- **WHEN** defining the `dedicated_resource_pack_tenancy` schema field
- **THEN** it SHALL be of type String
- **AND** it SHALL be Optional
- **AND** it SHALL have ForceNew set to true
- **AND** it SHALL be RequiredWith `dedicated_resource_pack_ids`
- **AND** it SHALL have a description explaining it specifies the tenancy strategy (e.g., "ResourcePool")

#### Scenario: dedicated_resource_pack_ids field specification

- **WHEN** defining the `dedicated_resource_pack_ids` schema field
- **THEN** it SHALL be of type List of Strings
- **AND** it SHALL be Optional
- **AND** it SHALL have ForceNew set to true
- **AND** it SHALL be RequiredWith `dedicated_resource_pack_tenancy`
- **AND** it SHALL have a description explaining it specifies the resource pack IDs to use

### Requirement: Backward compatibility

The addition of these parameters SHALL NOT break existing instance configurations.

#### Scenario: Existing configurations without new parameters

- **WHEN** user applies an existing instance configuration that does not specify dedicated resource pack parameters
- **THEN** the instance SHALL be created/updated successfully without using resource pool packs
- **AND** no changes SHALL be required to existing Terraform configurations

#### Scenario: State file compatibility

- **WHEN** provider reads state from an instance created before this feature was added
- **THEN** the provider SHALL handle the missing fields gracefully
- **AND** no spurious diffs SHALL be generated

### Requirement: Documentation

The resource documentation SHALL include examples and guidance for the new parameters.

#### Scenario: Basic usage example

- **WHEN** user views the resource documentation
- **THEN** it SHALL include an example showing how to create an instance with dedicated resource pack parameters
- **AND** the example SHALL demonstrate correct usage of both `dedicated_resource_pack_tenancy` and `dedicated_resource_pack_ids`

#### Scenario: Parameter descriptions

- **WHEN** user views the Argument Reference section
- **THEN** both new fields SHALL be documented with their type, optionality, ForceNew behavior, and purpose
- **AND** the documentation SHALL note that both parameters must be specified together

---

## ADDED Requirements - Disaster Recover Group IDs

### Requirement: disaster_recover_group_ids takes priority over placement_group_id

When `disaster_recover_group_ids` is set, it SHALL take priority over `placement_group_id` for all CRUD operations. The `ConflictsWith` constraint between these fields SHALL be removed, and they SHALL be allowed to coexist with `disaster_recover_group_ids` always winning.

#### Scenario: Create with disaster_recover_group_ids

- **WHEN** user specifies `disaster_recover_group_ids = ["pg-a", "pg-b"]` in the instance configuration
- **THEN** the provider SHALL pass these IDs to `request.DisasterRecoverGroupIds` in `RunInstances`
- **AND** `placement_group_id` SHALL be ignored in the Create request even if also specified
- **AND** the `rpgFlag` force_replace post-create call SHALL be skipped

#### Scenario: Create with both disaster_recover_group_ids and placement_group_id

- **WHEN** user specifies both `disaster_recover_group_ids = ["pg-a"]` and `placement_group_id = "pg-b"`
- **THEN** the provider SHALL use `disaster_recover_group_ids` in the Create request
- **AND** `placement_group_id` SHALL be ignored (not passed to API)

#### Scenario: Create with only placement_group_id (backward compatibility)

- **WHEN** user specifies only `placement_group_id = "pg-a"` (no `disaster_recover_group_ids`)
- **THEN** the provider SHALL pass it to `request.DisasterRecoverGroupIds` as a single-element list
- **AND** behavior SHALL be identical to previous releases

#### Scenario: Read with disaster_recover_group_ids in state

- **WHEN** the provider reads an instance that has `disaster_recover_group_ids` in state
- **THEN** `disaster_recover_group_ids` SHALL be updated from `instance.DisasterRecoverGroupIds` API response
- **AND** `placement_group_id` SHALL NOT be overwritten by `instance.DisasterRecoverGroupId` API response
- **AND** `placement_group_id` SHALL retain its state value to avoid plan diffs

#### Scenario: Read with only placement_group_id (backward compatibility)

- **WHEN** the provider reads an instance that has only `placement_group_id` in state (no `disaster_recover_group_ids`)
- **THEN** `placement_group_id` SHALL be updated from `instance.DisasterRecoverGroupId` API response
- **AND** behavior SHALL be identical to previous releases

#### Scenario: Update placement_group_id when disaster_recover_group_ids is set

- **WHEN** user attempts to change `placement_group_id` while `disaster_recover_group_ids` is in state
- **THEN** the provider SHALL return an error indicating `disaster_recover_group_ids` is set and `placement_group_id` changes are ignored

### Requirement: disaster_recover_group_ids schema constraints

#### Scenario: MaxItems validation

- **WHEN** user specifies more than 3 group IDs in `disaster_recover_group_ids`
- **THEN** Terraform SHALL return a validation error during plan phase

#### Scenario: disaster_recover_group_ids is Computed

- **WHEN** the provider reads an instance from API
- **THEN** `disaster_recover_group_ids` SHALL be populated from `instance.DisasterRecoverGroupIds`
- **AND** the field SHALL work correctly with import scenarios
