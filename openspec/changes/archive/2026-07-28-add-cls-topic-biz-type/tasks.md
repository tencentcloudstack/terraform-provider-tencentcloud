## 1. Schema & CRUD Implementation

- [x] 1.1 Add `biz_type` parameter to the `tencentcloud_cls_topic` resource schema in `tencentcloud/services/cls/resource_tc_cls_topic.go` (TypeInt, Optional, Computed, ForceNew, description: "Topic type. 0: log topic (default), 1: metric topic.")
- [x] 1.2 Add `biz_type` to the Create method: read `biz_type` from resource data and set `request.BizType = helper.IntUint64(v.(int))` before the API call
- [x] 1.3 Add `biz_type` to the Read method: set `biz_type` from `topic.BizType` in the DescribeTopics response, with nil check before setting
- [x] 1.4 Add `biz_type` to the `immutableArgs` array in the Update method to prevent update attempts
- [x] 1.5 Modify `DescribeClsTopicById` in `service_tencentcloud_cls.go` to accept an optional `bizType *uint64` parameter; when non-nil, set `request.BizType`
- [x] 1.6 Update the Read method to pass `biz_type` from state to `DescribeClsTopicById`; update all other callers to pass `nil`

## 2. Unit Tests

- [x] 2.1 Add unit test in `tencentcloud/services/cls/resource_tc_cls_topic_test.go` to verify `biz_type` is correctly passed to CreateTopic API request using gomonkey mock
- [x] 2.2 Add unit test to verify `biz_type` is correctly read from TopicInfo response and set in resource data

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/cls/resource_tc_cls_topic.md` to include `biz_type` parameter in the example usage
