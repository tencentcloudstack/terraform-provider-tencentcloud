## 1. Service Layer Changes

- [x] 1.1 Modify `CreateTdmqTopic` function in `tencentcloud/services/tdmq/service_tencentcloud_tdmq.go` to accept a `tags []*tdmq.Tag` parameter and set `request.Tags` when tags is non-nil and non-empty

## 2. Resource Schema and CRUD Changes

- [x] 2.1 Add `tags` parameter to the `tencentcloud_tdmq_topic` resource schema in `tencentcloud/services/tpulsar/resource_tc_tdmq_topic.go` as `Optional`, `ForceNew`, type `map[string]string`
- [x] 2.2 Update `resourceTencentCloudTdmqTopicCreate` to read `tags` from config, convert `map[string]string` to `[]*tdmq.Tag`, and pass to `CreateTdmqTopic`
- [x] 2.3 Update `resourceTencentCloudTdmqTopicRead` to read `Tags` from the `Topic` struct returned by `DescribeTopics` and set into state as `map[string]string`

## 3. Unit Tests

- [x] 3.1 Add unit tests in `tencentcloud/services/tpulsar/resource_tc_tdmq_topic_test.go` using gomonkey mocks to test create with tags, create without tags, and read with tags scenarios

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/tpulsar/resource_tc_tdmq_topic.md` with example usage including the `tags` parameter
