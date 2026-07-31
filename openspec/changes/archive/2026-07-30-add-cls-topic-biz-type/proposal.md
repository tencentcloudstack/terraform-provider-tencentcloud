## Why

The `tencentcloud_cls_topic` resource currently does not support the `BizType` parameter, which determines the topic type (0: log topic, 1: metric topic). Users need this parameter to create metric topics through Terraform, as the default behavior only creates log topics. This is a missing capability that prevents full management of CLS topic types via Terraform.

## What Changes

- Add a new `biz_type` parameter (TypeInt, Optional, Computed, ForceNew) to the `tencentcloud_cls_topic` resource schema
- The `biz_type` parameter maps to `BizType` in the CLS CreateTopic API (input, uint64 type) and is returned in the DescribeTopics API response (TopicInfo.BizType, uint64 type)
- Since `ModifyTopic` API does not support the `BizType` parameter, it is immutable after creation (ForceNew)
- Add `biz_type` to the `immutableArgs` array in the Update method to prevent update attempts
- Set `biz_type` in the Read method from the DescribeTopics response
- Modify `DescribeClsTopicById` service method to accept an optional `bizType *uint64` parameter; when non-nil, set `BizType` on the `DescribeTopicsRequest` for targeted API filtering
- **Resource ID format**: When `biz_type=0` or not specified, the resource ID is the plain `topic_id`. When `biz_type=1`, the resource ID is `"topic_id#1"` to encode the topic type in the ID
- **ID parsing**: Add a `parseClsTopicId` helper function to extract `topic_id` and `biz_type` from the resource ID. Used in Read, Update, and Delete methods
- **Read method**: Parse the resource ID to extract `topicId` and `bizType`, then pass them to `DescribeClsTopicById`. The API response's `BizType` is also set in state for validation
- **Import support**: When importing a metric topic (biz_type=1), use the `"topic_id#1"` format. Standard log topics use the plain `topic_id`

## Capabilities

### New Capabilities
- `cls-topic-biz-type`: Adds BizType (topic type) parameter support to the tencentcloud_cls_topic resource, allowing users to specify whether a topic is a log topic (0) or metric topic (1)

### Modified Capabilities
<!-- No existing capability requirements are changing -->

## Impact

- **Resource file**: `tencentcloud/services/cls/resource_tc_cls_topic.go` — schema definition, Create, Read, Update, Delete methods; added `parseClsTopicId` helper
- **Service file**: `tencentcloud/services/cls/service_tencentcloud_cls.go` — `DescribeClsTopicById` modified to accept optional `bizType` parameter
- **Test file**: `tencentcloud/services/cls/resource_tc_cls_topic_test.go` — unit tests for the new parameter
- **Documentation**: `tencentcloud/services/cls/resource_tc_cls_topic.md` — example usage and import examples for both log and metric topics
- **Other callers updated**: `resource_tc_clb_log_topic.go`, `resource_tc_clb_log_topic_test.go`, `resource_tc_redis_log_delivery.go`, `resource_tc_cls_cloud_product_log_task_v2.go` — pass `nil` for new `bizType` parameter
- **Cloud API**: CreateTopic (input: BizType), DescribeTopics (input: BizType filter, output: TopicInfo.BizType)
- **No breaking changes**: The new parameter is Optional with Computed default, existing configurations remain valid. Plain `topic_id` format for log topics is fully backward compatible
