# ES Instance Cerebro Config Specification

## Purpose

This specification defines the requirements for managing Cerebro service settings (enable/disable, public network access, private network access, custom private domain) of a Tencent Cloud Elasticsearch Service (ES) instance through the `tencentcloud_elasticsearch_instance` Terraform resource.

## Requirements

### Requirement: Resource schema exposes Cerebro configuration fields

The `tencentcloud_elasticsearch_instance` resource SHALL expose four new optional computed schema fields for Cerebro configuration: `enable_cerebro` (bool), `cerebro_public_access` (string), `cerebro_private_access` (string), and `cerebro_private_domain` (string).

**Rationale**: The ES `UpdateInstance` API accepts `EnableCerebro`, `CerebroPublicAccess`, `CerebroPrivateAccess`, and `CerebroPrivateDomain` as request fields, allowing users to manage the Cerebro monitoring service through Terraform. Since the `DescribeInstances` API response (`InstanceInfo`) does not return these fields, they are write-only parameters and must be `Computed: true` to preserve user-configured values in state.

#### Scenario: Schema fields are defined as optional + computed

- **WHEN** the `tencentcloud_elasticsearch_instance` resource schema is inspected
- **THEN** it SHALL contain an `enable_cerebro` field of type `schema.TypeBool` with `Optional: true` and `Computed: true`
- **AND** it SHALL contain a `cerebro_public_access` field of type `schema.TypeString` with `Optional: true` and `Computed: true`
- **AND** it SHALL contain a `cerebro_private_access` field of type `schema.TypeString` with `Optional: true` and `Computed: true`
- **AND** it SHALL contain a `cerebro_private_domain` field of type `schema.TypeString` with `Optional: true` and `Computed: true`
- **AND** the access fields SHALL validate that the value is one of `OPEN` or `CLOSE`

#### Scenario: Existing configurations remain valid

- **WHEN** a user applies an existing `tencentcloud_elasticsearch_instance` configuration that does not specify Cerebro fields
- **THEN** the plan SHALL succeed without requiring the new fields
- **AND** no existing schema fields SHALL be removed or changed

### Requirement: Cerebro configuration is applied on update

The resource update flow SHALL detect changes to any Cerebro field and call the ES `UpdateInstance` API with the corresponding request field, using retry logic (`tccommon.WriteRetryTimeout`) and waiting for the instance upgrade to complete.

**Rationale**: Cerebro is a modifiable attribute of an existing ES instance, exposed via the `UpdateInstance` API. Cerebro cannot be configured at creation time, so the update flow is the only path to apply these settings.

#### Scenario: Enable Cerebro

- **WHEN** a user changes `enable_cerebro` from `false` to `true` on an existing `tencentcloud_elasticsearch_instance`
- **THEN** the resource SHALL call `UpdateInstance` with `EnableCerebro` set to `true`
- **AND** the call SHALL be wrapped in retry logic using `tccommon.WriteRetryTimeout`
- **AND** the resource SHALL wait for the instance upgrade to complete

#### Scenario: Update Cerebro public access

- **WHEN** a user changes `cerebro_public_access` from `CLOSE` to `OPEN` on an existing `tencentcloud_elasticsearch_instance`
- **THEN** the resource SHALL call `UpdateInstance` with `CerebroPublicAccess` set to `OPEN`
- **AND** the resource SHALL wait for the instance upgrade to complete

#### Scenario: Update Cerebro private access

- **WHEN** a user changes `cerebro_private_access` from `CLOSE` to `OPEN` on an existing `tencentcloud_elasticsearch_instance`
- **THEN** the resource SHALL call `UpdateInstance` with `CerebroPrivateAccess` set to `OPEN`
- **AND** the resource SHALL wait for the instance upgrade to complete

#### Scenario: Update Cerebro private domain

- **WHEN** a user sets `cerebro_private_domain` to a custom domain on an existing `tencentcloud_elasticsearch_instance`
- **THEN** the resource SHALL call `UpdateInstance` with `CerebroPrivateDomain` set to the configured value
- **AND** the resource SHALL wait for the instance upgrade to complete

#### Scenario: No update when Cerebro fields are unchanged

- **WHEN** a user updates other fields of `tencentcloud_elasticsearch_instance` without changing Cerebro fields
- **THEN** the resource SHALL NOT call `UpdateInstance` for Cerebro configuration

### Requirement: Cerebro fields are not read back from the API

The resource read flow SHALL NOT attempt to set Cerebro fields from the `DescribeInstances` response, because the `InstanceInfo` struct does not contain Cerebro fields. User-configured Cerebro values SHALL be preserved in Terraform state due to the `Computed` nature of the fields.

**Rationale**: The ES `DescribeInstances` API does not return Cerebro configuration on `InstanceInfo`. Attempting to set these fields to nil/empty would cause Terraform to detect drift on every plan. Marking the fields as `Computed` preserves user values.

#### Scenario: Read preserves user-configured Cerebro values

- **WHEN** the resource read flow queries `DescribeInstances` for an instance with user-configured Cerebro fields
- **AND** the response `InstanceInfo` has no Cerebro fields
- **THEN** the resource SHALL NOT overwrite the Cerebro fields in state
- **AND** the read SHALL succeed without error

#### Scenario: Cerebro drift outside Terraform is not detected

- **WHEN** Cerebro state is changed outside of Terraform (e.g., via console)
- **THEN** the next `terraform plan` SHALL NOT detect the drift
- **AND** this limitation SHALL be documented in the field descriptions

### Requirement: Service layer passes Cerebro parameters to UpdateInstance API

The `UpdateInstance` service layer function SHALL accept Cerebro parameters and set the corresponding `UpdateInstanceRequest` fields (`EnableCerebro`, `CerebroPublicAccess`, `CerebroPrivateAccess`, `CerebroPrivateDomain`) when the parameters are non-empty/non-nil.

**Rationale**: The service layer is the single place where the `UpdateInstanceRequest` is constructed. Centralizing field mapping here keeps the resource layer clean.

#### Scenario: Service layer sets Cerebro request fields

- **WHEN** the `UpdateInstance` service function is called with `enableCerebro = true`, `cerebroPublicAccess = "OPEN"`, `cerebroPrivateAccess = "CLOSE"`, and `cerebroPrivateDomain = "example.com"`
- **THEN** the constructed `UpdateInstanceRequest` SHALL have `EnableCerebro` set to `true`
- **AND** `CerebroPublicAccess` set to `"OPEN"`
- **AND** `CerebroPrivateAccess` set to `"CLOSE"`
- **AND** `CerebroPrivateDomain` set to `"example.com"`

#### Scenario: Service layer skips empty Cerebro request fields

- **WHEN** the `UpdateInstance` service function is called with empty Cerebro parameters
- **THEN** the constructed `UpdateInstanceRequest` SHALL NOT have Cerebro fields set
- **AND** the update SHALL succeed for non-Cerebro changes

### Requirement: Documentation and changelog are updated

The source documentation file and changelog SHALL be updated to reflect the new Cerebro parameters.

**Rationale**: All resource changes require updated documentation and a changelog entry per provider conventions.

#### Scenario: Resource documentation includes the new parameters

- **WHEN** the source documentation file `resource_tc_elasticsearch_instance.md` is updated
- **THEN** it SHALL include an example usage demonstrating the Cerebro parameters
- **AND** the generated `website/docs/r/elasticsearch_instance.html.markdown` SHALL be produced via `make doc`

#### Scenario: Changelog entry is created

- **WHEN** the change is finalized
- **THEN** a changelog file SHALL be created under `.changelog/` describing the new Cerebro parameters for `tencentcloud_elasticsearch_instance`
