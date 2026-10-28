## ADDED Requirements

### Requirement: Tags parameter in resource schema

The `tencentcloud_tdmq_topic` resource SHALL support an optional `tags` parameter of type `map[string]string`. Each map entry represents a tag where the key is the tag key and the value is the tag value.

#### Scenario: Create topic with tags
- **WHEN** a user specifies `tags` in the `tencentcloud_tdmq_topic` resource configuration
- **THEN** the resource SHALL convert the map entries to `[]*tdmq.Tag` (with `TagKey` and `TagValue` fields) and pass them in the `CreateTopic` API request's `Tags` field

#### Scenario: Create topic without tags
- **WHEN** a user does not specify `tags` in the `tencentcloud_tdmq_topic` resource configuration
- **THEN** the resource SHALL not set the `Tags` field in the `CreateTopic` API request, and the resource SHALL be created successfully without tags

#### Scenario: Update tags in-place
- **WHEN** a user modifies the `tags` parameter in an existing `tencentcloud_tdmq_topic` resource
- **THEN** the resource Update function SHALL:
  1. Compute the diff between old and new tags using `svctag.DiffTags`
  2. Call `UnTagResources` API with the deleted tag keys to remove them from the resource
  3. Call `TagResources` API with the new/updated tag key-value pairs to bind them to the resource
  4. The resource six-segment name SHALL be constructed as `tccommon.BuildTagResourceName("tdmq", "topic", region, topicName)`

### Requirement: Read tags from API response

The resource read function SHALL read the `Tags` field from the `Topic` struct returned by the `DescribeTopics` API and set it into the Terraform state as a `map[string]string`.

#### Scenario: Read topic with tags
- **WHEN** the `DescribeTopics` API returns a `Topic` with non-nil `Tags` field containing tag entries
- **THEN** the resource read function SHALL convert `[]*tdmq.Tag` to `map[string]string` and set it into state via `d.Set("tags", tagsMap)`

#### Scenario: Read topic without tags
- **WHEN** the `DescribeTopics` API returns a `Topic` with nil or empty `Tags` field
- **THEN** the resource read function SHALL not call `d.Set("tags", ...)` to avoid overwriting user configuration with empty values
