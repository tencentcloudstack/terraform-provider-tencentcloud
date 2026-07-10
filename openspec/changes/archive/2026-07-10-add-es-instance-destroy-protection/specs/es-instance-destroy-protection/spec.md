# Capability: ES Instance Destroy Protection

## Overview

This specification defines the requirements for managing the destroy protection (`EnableDestroyProtection`) setting of a Tencent Cloud Elasticsearch Service (ES) instance through the `tencentcloud_elasticsearch_instance` Terraform resource.

## ADDED Requirements

### Requirement: Resource schema exposes destroy protection

The `tencentcloud_elasticsearch_instance` resource SHALL expose an `enable_destroy_protection` schema field of type string that is optional and computed, accepting only the values `OPEN` and `CLOSE`.

**Rationale**: Users need to manage cluster destroy protection through Terraform to prevent accidental deletion of production Elasticsearch clusters. The cloud API field `EnableDestroyProtection` is a string enum (`OPEN`/`CLOSE`), so the Terraform field mirrors it.

#### Scenario: Schema field is defined as optional + computed string

- **WHEN** the `tencentcloud_elasticsearch_instance` resource schema is inspected
- **THEN** it SHALL contain an `enable_destroy_protection` field of type `schema.TypeString`
- **AND** the field SHALL be `Optional: true` and `Computed: true`
- **AND** the field SHALL validate that the value is one of `OPEN` or `CLOSE`

#### Scenario: Existing configurations remain valid

- **WHEN** a user applies an existing `tencentcloud_elasticsearch_instance` configuration that does not specify `enable_destroy_protection`
- **THEN** the plan SHALL succeed without requiring the new field
- **AND** no existing schema fields SHALL be removed or changed

### Requirement: Destroy protection is applied on create

The resource create flow SHALL enable or disable destroy protection after the instance is created and reaches normal status, by calling the ES `UpdateInstance` API with the `EnableDestroyProtection` request field set to the configured value.

**Rationale**: The ES `CreateInstance` API does not accept an `EnableDestroyProtection` parameter, so destroy protection cannot be set at creation time directly. The resource applies it as a post-create update step, consistent with how other attributes (e.g. `es_acl`, `public_access`, `cos_backup`) are applied after creation.

#### Scenario: Create with destroy protection enabled

- **WHEN** a user creates a `tencentcloud_elasticsearch_instance` with `enable_destroy_protection = "OPEN"`
- **THEN** the resource SHALL create the instance via `CreateInstance`
- **AND** after the instance reaches normal status, the resource SHALL call `UpdateInstance` with `EnableDestroyProtection` set to `OPEN`
- **AND** the resource SHALL wait for the instance upgrade to complete

#### Scenario: Create without destroy protection specified

- **WHEN** a user creates a `tencentcloud_elasticsearch_instance` without specifying `enable_destroy_protection`
- **THEN** the resource SHALL create the instance via `CreateInstance`
- **AND** the resource SHALL NOT call `UpdateInstance` solely for destroy protection
- **AND** the create SHALL succeed

### Requirement: Destroy protection is updated on change

The resource update flow SHALL detect changes to `enable_destroy_protection` and call the ES `UpdateInstance` API with the `EnableDestroyProtection` request field, using retry logic and waiting for the instance upgrade to complete.

**Rationale**: Destroy protection is a modifiable attribute of an existing ES instance, exposed via the `UpdateInstance` API input field `EnableDestroyProtection`.

#### Scenario: Update destroy protection from CLOSE to OPEN

- **WHEN** a user changes `enable_destroy_protection` from `CLOSE` to `OPEN` on an existing `tencentcloud_elasticsearch_instance`
- **THEN** the resource SHALL call `UpdateInstance` with `EnableDestroyProtection` set to `OPEN`
- **AND** the call SHALL be wrapped in retry logic using `tccommon.WriteRetryTimeout`
- **AND** the resource SHALL wait for the instance upgrade to complete via the upgrade-wait helper

#### Scenario: Update destroy protection from OPEN to CLOSE

- **WHEN** a user changes `enable_destroy_protection` from `OPEN` to `CLOSE` on an existing `tencentcloud_elasticsearch_instance`
- **THEN** the resource SHALL call `UpdateInstance` with `EnableDestroyProtection` set to `CLOSE`
- **AND** the resource SHALL wait for the instance upgrade to complete

#### Scenario: No update when destroy protection is unchanged

- **WHEN** a user updates other fields of `tencentcloud_elasticsearch_instance` without changing `enable_destroy_protection`
- **THEN** the resource SHALL NOT call `UpdateInstance` for destroy protection

### Requirement: Destroy protection is read from the API

The resource read flow SHALL read the `EnableDestroyProtection` field from the `DescribeInstances` response (`InstanceInfo`) and set it into Terraform state, only when the API returns a non-null value.

**Rationale**: The ES `DescribeInstances` API returns `EnableDestroyProtection` on `InstanceInfo` (marked as possibly null), allowing the current protection state to be reconciled into state.

#### Scenario: Read populates state when API returns a value

- **WHEN** the resource read flow queries `DescribeInstances` for an instance
- **AND** the response `InstanceInfo.EnableDestroyProtection` is non-null
- **THEN** the resource SHALL set `enable_destroy_protection` in state to the returned value

#### Scenario: Read is nil-safe when API returns null

- **WHEN** the resource read flow queries `DescribeInstances` for an instance
- **AND** the response `InstanceInfo.EnableDestroyProtection` is null
- **THEN** the resource SHALL NOT overwrite `enable_destroy_protection` in state
- **AND** the read SHALL succeed without error

### Requirement: Destroy protection blocks deletion when OPEN

When destroy protection is `OPEN`, attempting to destroy the `tencentcloud_elasticsearch_instance` resource SHALL fail at the cloud API `DeleteInstance` call until the user sets `enable_destroy_protection` to `CLOSE`.

**Rationale**: This is the intended security behavior of the cloud API; destroy protection exists to prevent accidental deletion. The Terraform provider surfaces the cloud API error rather than silently disabling protection.

#### Scenario: Destroy fails while protection is OPEN

- **WHEN** a user runs `terraform destroy` on a `tencentcloud_elasticsearch_instance` whose `enable_destroy_protection` is `OPEN`
- **THEN** the `DeleteInstance` API call SHALL fail
- **AND** the cloud API error SHALL be surfaced to the user

#### Scenario: Destroy succeeds after disabling protection

- **WHEN** a user sets `enable_destroy_protection = "CLOSE"`, applies the change, and then runs `terraform destroy`
- **THEN** the `DeleteInstance` API call SHALL succeed
- **AND** the resource SHALL be destroyed

### Requirement: Documentation and changelog are updated

The source documentation file and changelog SHALL be updated to reflect the new `enable_destroy_protection` parameter.

**Rationale**: All resource changes require updated documentation and a changelog entry per provider conventions.

#### Scenario: Resource documentation includes the new parameter

- **WHEN** the source documentation file `resource_tc_elasticsearch_instance.md` is updated
- **THEN** it SHALL include an example usage demonstrating `enable_destroy_protection`
- **AND** the generated `website/docs/r/elasticsearch_instance.html.markdown` SHALL be produced via `make doc`

#### Scenario: Changelog entry is created

- **WHEN** the change is finalized
- **THEN** a changelog file SHALL be created under `.changelog/` describing the new `enable_destroy_protection` parameter for `tencentcloud_elasticsearch_instance`
