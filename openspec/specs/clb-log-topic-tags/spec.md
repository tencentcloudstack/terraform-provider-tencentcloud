# clb-log-topic-tags Specification

## Purpose
TBD - created by archiving change add-clb-log-topic-tags-param. Update Purpose after archive.
## Requirements
### Requirement: Support tags on CLB log topic

The `tencentcloud_clb_log_topic` resource SHALL support an optional `tags` parameter that allows users to configure tags on a CLB log topic. Each tag element SHALL have a `key` and `value` string field.

**Rationale**: The underlying cloud APIs (`CreateTopic`, `ModifyTopic`, `DescribeTopics`) all support tags, but the Terraform resource does not expose them. Adding tag support enables cost allocation, access control, and resource organization.

#### Scenario: Create a CLB log topic with tags

- **WHEN** a user creates a `tencentcloud_clb_log_topic` resource with the `tags` parameter set to a list of `{ key, value }` objects
- **THEN** the resource is created successfully
- **AND** the tags are passed to the `CreateTopic` API as `[]*TagInfo` with `TagKey`/`TagValue` mapped from the schema `key`/`value`
- **AND** a subsequent `terraform plan` shows no changes

#### Scenario: Create a CLB log topic without tags

- **WHEN** a user creates a `tencentcloud_clb_log_topic` resource without specifying the `tags` parameter
- **THEN** the resource is created successfully
- **AND** no tags are sent to the `CreateTopic` API

#### Scenario: Read tags back from the DescribeTopics API

- **WHEN** the `Read` method calls `DescribeTopics` and the response `TopicInfo.Tags` contains a list of `Tag` objects with `Key`/`Value`
- **THEN** the tags are flattened into the `tags` schema list with `key`/`value` fields
- **AND** the state reflects the tags returned by the API

#### Scenario: Read when no tags are returned

- **WHEN** the `Read` method calls `DescribeTopics` and the response `TopicInfo.Tags` is nil or empty
- **THEN** the `tags` field in state is set to an empty list
- **AND** no error is returned

#### Scenario: Update tags on an existing CLB log topic

- **WHEN** a user updates the `tags` parameter on an existing `tencentcloud_clb_log_topic` resource
- **THEN** the `ModifyTopic` API is called with the new tags as `[]*Tag` with `Key`/`Value` mapped from the schema `key`/`value`
- **AND** the resource is not recreated
- **AND** a subsequent read reflects the updated tags

#### Scenario: Tags parameter is optional and backward compatible

- **WHEN** an existing `tencentcloud_clb_log_topic` configuration without the `tags` parameter is applied after the provider upgrade
- **THEN** the resource continues to work without any changes
- **AND** no tags are added or removed

