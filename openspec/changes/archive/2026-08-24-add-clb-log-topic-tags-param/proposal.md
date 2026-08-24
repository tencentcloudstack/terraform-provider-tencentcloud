## Why

The `tencentcloud_clb_log_topic` resource currently does not support configuring tags, even though the underlying cloud APIs (`CreateTopic`, `ModifyTopic`, `DescribeTopics`) all support tag parameters. Users cannot attach tags to CLB log topics through Terraform, which limits cost allocation, access control, and resource organization capabilities. This change adds tag support to achieve feature parity with the cloud API.

## What Changes

- Add a new optional `tags` parameter (TypeMap of string key/value pairs) to the `tencentcloud_clb_log_topic` resource schema.
- Update the `Create` method to pass tags to the `CreateTopic` API (`clb/v20180317`) using the `TagInfo` structure (`TagKey`/`TagValue` fields).
- Update the `Update` method to pass tags to the `ModifyTopic` API (`cls/v20201016`) using the `Tag` structure (`Key`/`Value` fields) when tags change.
- Update the `Read` method to read tags from the `DescribeTopics` API response (`cls/v20201016`) `TopicInfo.Tags` (`Tag` structure with `Key`/`Value` fields) and set them into state.

## Capabilities

### New Capabilities
- `clb-log-topic-tags`: Adds tag support to the `tencentcloud_clb_log_topic` resource, allowing users to configure and manage tags on CLB log topics through Terraform.

### Modified Capabilities
<!-- None - no existing capability spec is being modified. -->

## Impact

- **Files Modified**:
  - `tencentcloud/services/clb/resource_tc_clb_log_topic.go` - Add `tags` schema field and update CRUD logic
  - `tencentcloud/services/clb/resource_tc_clb_log_topic_test.go` - Add test cases for the new `tags` parameter
  - `tencentcloud/services/clb/resource_tc_clb_log_topic.md` - Update documentation with tags example
- **APIs Involved**:
  - `CreateTopic` (`clb/v20180317`) - Tags via `[]*TagInfo` (`TagKey`/`TagValue`)
  - `ModifyTopic` (`cls/v20201016`) - Tags via `[]*Tag` (`Key`/`Value`)
  - `DescribeTopics` (`cls/v20201016`) - Tags returned in `TopicInfo.Tags` (`Tag` with `Key`/`Value`)
- **Backward Compatibility**: The new `tags` parameter is optional; existing configurations continue working without changes.
