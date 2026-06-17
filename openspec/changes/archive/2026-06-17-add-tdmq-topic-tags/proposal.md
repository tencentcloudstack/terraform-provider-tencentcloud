## Why

The `tencentcloud_tdmq_topic` resource currently does not support the `tags` parameter, which prevents users from assigning tags to TDMQ topics at creation time. The TDMQ CreateTopic API already supports a `Tags` field, and the DescribeTopics API returns tags in the `Topic` struct. Adding this parameter enables users to manage topic-level tags through Terraform.

## What Changes

- Add a new `tags` parameter (type: `map[string]string`) to the `tencentcloud_tdmq_topic` resource schema, marked as `Optional` and `ForceNew` (since the ModifyTopic API does not support updating tags).
- Pass the `tags` parameter to the `CreateTopic` API request as `request.Tags` during resource creation.
- Read the `Tags` field from the `Topic` struct returned by the `DescribeTopics` API and set it back into state during resource read.

## Capabilities

### New Capabilities
- `topic-tags`: Support setting and reading tags on TDMQ topics via the `tencentcloud_tdmq_topic` resource.

### Modified Capabilities
(none)

## Impact

- **Code**: `tencentcloud/services/tpulsar/resource_tc_tdmq_topic.go` — schema, create, and read functions need modification.
- **Tests**: `tencentcloud/services/tpulsar/resource_tc_tdmq_topic_test.go` — add unit tests for the tags parameter using gomonkey mocks.
- **Documentation**: `tencentcloud/services/tpulsar/resource_tc_tdmq_topic.md` — update example usage and description.
- **APIs**: Uses existing `CreateTopic` (Tags input) and `DescribeTopics` (Tags in Topic response) from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217`. No new API dependencies.
- **Constraint**: The `ModifyTopic` API does not support a `Tags` field, so the `tags` parameter must be `ForceNew` — changing tags requires resource recreation.
