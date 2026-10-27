# clb-log-topic-period Specification

## Purpose
TBD - created by archiving change add-clb-log-topic-period. Update Purpose after archive.
## Requirements
### Requirement: Support period parameter on CLB log topic

The `tencentcloud_clb_log_topic` resource SHALL support an optional `period` (TypeInt) parameter that controls the log storage lifecycle in days. The parameter SHALL be updatable in-place (not `ForceNew`). When set, the provider SHALL pass `Period` to the `CreateTopic` API (CLB SDK `clb/v20180317`) on creation and to the `ModifyTopic` API (CLS SDK `cls/v20201016`) on update. When the user does not set `period`, the provider SHALL NOT send `Period` and the cloud API default (30 days) SHALL apply.

**Rationale**: The underlying `CreateTopic`, `ModifyTopic`, and `DescribeTopics` APIs all support a `Period` field, but the Terraform resource does not expose it. Exposing `period` enables users to configure log retention through Infrastructure as Code.

#### Scenario: Create a CLB log topic with period

- **WHEN** a user creates a `tencentcloud_clb_log_topic` resource with the `period` parameter set (e.g. `period = 30`)
- **THEN** the provider SHALL pass the value to the `CreateTopic` API request as `Period` (cast to `*uint64`)
- **AND** a subsequent `terraform plan` shows no changes

#### Scenario: Create a CLB log topic without period

- **WHEN** a user creates a `tencentcloud_clb_log_topic` resource without specifying the `period` parameter
- **THEN** the provider SHALL NOT set `Period` in the `CreateTopic` API request
- **AND** the cloud API default retention (30 days) SHALL apply

#### Scenario: Read period back from the DescribeTopics API

- **WHEN** the `Read` method calls `DescribeTopics` and the response `TopicInfo.Period` is non-nil
- **THEN** the provider SHALL set the `period` field in state to the value of `TopicInfo.Period`
- **AND** state refresh and import SHALL populate `period` correctly

#### Scenario: Read when period is not returned

- **WHEN** the `Read` method calls `DescribeTopics` and the response `TopicInfo.Period` is nil
- **THEN** the provider SHALL NOT set the `period` field in state
- **AND** no error is returned

#### Scenario: Update period on an existing CLB log topic

- **WHEN** a user updates the `period` parameter on an existing `tencentcloud_clb_log_topic` resource
- **THEN** the provider SHALL call the `ModifyTopic` API with the new value as `Period` (cast to `*int64`)
- **AND** the updated value SHALL be reflected in state after the subsequent Read

#### Scenario: Update other fields without changing period

- **WHEN** a user updates `status` or `tags` but does not change `period`
- **THEN** the provider SHALL NOT include `Period` in the `ModifyTopic` request
- **AND** the existing retention value SHALL remain unchanged

