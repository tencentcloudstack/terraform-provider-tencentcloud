## Why

The `tencentcloud_cls_topic` resource currently does not support the `BizType` parameter, which determines the topic type (0: log topic, 1: metric topic). Users need this parameter to create metric topics through Terraform, as the default behavior only creates log topics. This is a missing capability that prevents full management of CLS topic types via Terraform.

## What Changes

- Add a new `biz_type` parameter (TypeInt, Optional, Computed, ForceNew) to the `tencentcloud_cls_topic` resource schema
- The `biz_type` parameter maps to `BizType` in the CLS CreateTopic API (input, uint64 type) and is returned in the DescribeTopics API response (TopicInfo.BizType, uint64 type)
- Since `ModifyTopic` API does not support the `BizType` parameter, it is immutable after creation (ForceNew)
- Add `biz_type` to the `immutableArgs` array in the Update method to prevent update attempts
- Set `biz_type` in the Read method from the DescribeTopics response
- Modify `DescribeClsTopicById` service method to accept an optional `bizType *uint64` parameter; when non-nil, set `BizType` on the `DescribeTopicsRequest` for targeted API filtering
- Update the Read method to pass `biz_type` from state to `DescribeClsTopicById`; other callers pass `nil` (no BizType filter)

## Capabilities

### New Capabilities
- `cls-topic-biz-type`: Adds BizType (topic type) parameter support to the tencentcloud_cls_topic resource, allowing users to specify whether a topic is a log topic (0) or metric topic (1)

### Modified Capabilities
<!-- No existing capability requirements are changing -->

## Impact

- **Resource file**: `tencentcloud/services/cls/resource_tc_cls_topic.go` — schema definition, Create, Read, Update methods
- **Service file**: `tencentcloud/services/cls/service_tencentcloud_cls.go` — `DescribeClsTopicById` modified to accept optional `bizType` parameter
- **Test file**: `tencentcloud/services/cls/resource_tc_cls_topic_test.go` — unit tests for the new parameter
- **Documentation**: `tencentcloud/services/cls/resource_tc_cls_topic.md` — example usage update
- **Other callers updated**: `resource_tc_clb_log_topic.go`, `resource_tc_clb_log_topic_test.go`, `resource_tc_redis_log_delivery.go`, `resource_tc_cls_cloud_product_log_task_v2.go` — pass `nil` for new `bizType` parameter
- **Cloud API**: CreateTopic (input: BizType), DescribeTopics (input: BizType filter, output: TopicInfo.BizType)
- **No breaking changes**: The new parameter is Optional with Computed default, existing configurations remain valid
