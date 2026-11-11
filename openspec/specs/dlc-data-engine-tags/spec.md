# dlc-data-engine-tags Specification

## Purpose
TBD - created by archiving change add-dlc-data-engine-tags. Update Purpose after archive.
## Requirements
### Requirement: dlc_data_engine tags schema definition
The system SHALL add a `tags` parameter to the `tencentcloud_dlc_data_engine` resource schema to bind tags to a DLC data engine at creation time.

#### Scenario: tags block schema fields
- **WHEN** the `tencentcloud_dlc_data_engine` resource schema is defined
- **THEN** it includes a `tags` field of type `schema.TypeList`
- **AND** the `tags` field is `Optional: true`
- **AND** the `tags` field is `ForceNew: true` because the DLC `UpdateDataEngine` API does not support modifying tags
- **AND** the `tags` element is a `schema.Resource` with two string sub-fields: `tag_key` (Required) and `tag_value` (Optional)
- **AND** the `Description` for `tags` states that changing this parameter will force a new resource

#### Scenario: tag_key sub-field
- **WHEN** the `tags` element schema is defined
- **THEN** it contains `tag_key` of type `schema.TypeString` with `Required: true`
- **AND** the description documents it as the tag key

#### Scenario: tag_value sub-field
- **WHEN** the `tags` element schema is defined
- **THEN** it contains `tag_value` of type `schema.TypeString` with `Optional: true`
- **AND** the description documents it as the tag value

### Requirement: dlc_data_engine tags population on Create
The system SHALL map the `tags` block to `CreateDataEngineRequest.Tags` when creating a DLC data engine.

#### Scenario: tags passed to CreateDataEngine API
- **WHEN** the user configures the `tags` block on `tencentcloud_dlc_data_engine` and a create is performed
- **THEN** for each element in the `tags` list, the system constructs a `dlc.TagInfo` with `TagKey` set from `tag_key` and `TagValue` set from `tag_value`
- **AND** the constructed `[]*dlc.TagInfo` is assigned to `request.Tags`
- **AND** when `tag_value` is empty it is still passed as an empty string, matching the cloud API `TagInfo` structure

#### Scenario: tags omitted on create
- **WHEN** the user does not configure the `tags` block
- **THEN** `request.Tags` is left as its zero value (nil) so the cloud API applies defaults
- **AND** resource creation proceeds normally without error

### Requirement: dlc_data_engine tags synchronization on Read
The system SHALL read `TagList` from the DLC describe response and flatten it into the `tags` block during `Read`.

#### Scenario: tags flattened from TagList
- **WHEN** the `Read` handler obtains a `DataEngineInfo` with a non-nil `TagList`
- **THEN** the system builds a `[]map[string]interface{}` where each map carries `tag_key` and `tag_value`
- **AND** `tag_key` is populated only when `TagInfo.TagKey` is non-nil
- **AND** `tag_value` is populated only when `TagInfo.TagValue` is non-nil
- **AND** the list is written to state via `d.Set("tags", ...)`

#### Scenario: TagList is nil or empty
- **WHEN** the described `DataEngineInfo.TagList` is nil or has length zero
- **THEN** the system does not populate `tag_key`/`tag_value` entries
- **AND** the `tags` state is set to an empty list (or left unset) without error
- **AND** other fields continue to be read normally

### Requirement: dlc_data_engine tags are immutable after creation
The system SHALL prevent in-place updates of the `tags` parameter.

#### Scenario: ForceNew triggers recreation
- **WHEN** the user changes the `tags` block in an existing `tencentcloud_dlc_data_engine` configuration
- **THEN** Terraform plans a destroy-and-recreate because `tags` is `ForceNew: true`
- **AND** the `Update` handler is not invoked for a `tags`-only change

#### Scenario: tags listed as immutable in Update
- **WHEN** the `resourceTencentCloudDlcDataEngineUpdate` function builds its `immutableArgs` guard list
- **THEN** the list includes `"tags"` alongside the other create-only arguments
- **AND** if a `tags` change somehow reaches `Update`, it returns an error stating the argument cannot be changed

### Requirement: dlc_data_engine tags documentation
The system SHALL document the `tags` parameter in the resource example markdown.

#### Scenario: tags documented in resource md
- **WHEN** the resource documentation is generated/updated
- **THEN** `resource_tc_dlc_data_engine.md` includes the `tags` block in the Example Usage
- **AND** the example shows a `tags` block with `tag_key` and `tag_value` entries
- **AND** the `Argument Reference` / `Attribute Reference` sections are left for auto-generation (per project convention, these are not hand-written)

### Requirement: dlc_data_engine tags unit tests
The system SHALL include unit tests covering the tags marshalling and flattening logic using gomonkey mocks.

#### Scenario: create marshals tags to request
- **WHEN** a unit test exercises `resourceTencentCloudDlcDataEngineCreate` with a `tags` block configured
- **THEN** the test mocks the DLC client `CreateDataEngine` and `DescribeDataEngine` calls via gomonkey
- **AND** the test asserts that `request.Tags` contains the expected `TagInfo` entries with correct `TagKey` and `TagValue`

#### Scenario: read flattens TagList to state
- **WHEN** a unit test exercises `resourceTencentCloudDlcDataEngineRead` with a `DataEngineInfo` whose `TagList` is populated
- **THEN** the test asserts that `d.Get("tags")` returns the expected list of `tag_key`/`tag_value` maps
- **AND** the test covers the nil/empty `TagList` case without error

