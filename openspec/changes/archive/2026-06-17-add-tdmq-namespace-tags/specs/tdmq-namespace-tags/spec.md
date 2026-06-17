## ADDED Requirements

### Requirement: Tags parameter in tencentcloud_tdmq_namespace resource
The `tencentcloud_tdmq_namespace` resource SHALL support a `Tags` parameter that allows users to assign tags to a TDMQ namespace during creation. The `Tags` parameter SHALL be of type `TypeList`, with each element containing `tag_key` (Required, TypeString) and `tag_value` (Required, TypeString) sub-fields. The `Tags` parameter SHALL be Optional and ForceNew, since the `ModifyEnvironmentAttributes` API does not support updating tags.

#### Scenario: Create namespace with tags
- **WHEN** a user creates a `tencentcloud_tdmq_namespace` resource with the `Tags` parameter specified
- **THEN** the system SHALL pass the Tags to the `CreateEnvironment` API request as `request.Tags` (type `[]*Tag` with TagKey/TagValue)
- **THEN** the resource SHALL be created with the specified tags assigned

#### Scenario: Create namespace without tags
- **WHEN** a user creates a `tencentcloud_tdmq_namespace` resource without specifying the `Tags` parameter
- **THEN** the system SHALL create the namespace without tags, and the `Tags` field in the API request SHALL be omitted

#### Scenario: Read namespace with tags
- **WHEN** the system reads a `tencentcloud_tdmq_namespace` resource and the `DescribeEnvironments` API returns a non-nil `Tags` field in the `Environment` response
- **THEN** the system SHALL flatten the Tags into the schema by converting each `Tag` item's `TagKey` to `tag_key` and `TagValue` to `tag_value`
- **THEN** the system SHALL set the `tags` field in the Terraform state

#### Scenario: Read namespace with nil Tags
- **WHEN** the system reads a `tencentcloud_tdmq_namespace` resource and the `DescribeEnvironments` API returns nil `Tags` in the `Environment` response
- **THEN** the system SHALL skip setting the `tags` field to avoid writing nil data to the state

#### Scenario: Update namespace with changed tags
- **WHEN** a user attempts to update the `Tags` parameter on an existing `tencentcloud_tdmq_namespace` resource
- **THEN** since Tags is ForceNew, Terraform SHALL trigger a resource recreation (destroy + create)
- **THEN** the update method SHALL also include "tags" in the immutableArgs list to provide a clear error message

#### Scenario: Import namespace with tags
- **WHEN** a user imports an existing `tencentcloud_tdmq_namespace` that has tags assigned
- **THEN** the system SHALL read and populate the `tags` field in the Terraform state from the `DescribeEnvironments` response
