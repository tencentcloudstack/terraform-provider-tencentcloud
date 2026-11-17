## Why

The `tencentcloud_clb_log_topic` resource manages CLB (Cloud Load Balancer) log topics but does not expose the log retention period (`Period`). The underlying cloud APIs (`CreateTopic`, `ModifyTopic`, `DescribeTopics`) all support a `Period` field that controls how long (in days) logs are stored, yet users currently cannot configure log retention through Terraform and must fall back to the console or API. Exposing `period` enables full Infrastructure-as-Code lifecycle management of CLB log topics, including retention policy.

## What Changes

- Add an optional `period` (TypeInt) parameter to the `tencentcloud_clb_log_topic` resource schema. The value represents the log storage lifecycle in days (range 1-3600; 3640 means permanent retention). The parameter is updatable and not `ForceNew`.
- In the Create flow, pass `period` to the `CreateTopic` API (CLB SDK `clb/v20180317`, `CreateTopicRequest.Period` is `*uint64`).
- In the Update flow, when `period` changes, call the `ModifyTopic` API (CLS SDK `cls/v20201016`, `ModifyTopicRequest.Period` is `*int64`) with the new value.
- In the Read flow, read `period` back from the `DescribeTopics` API response (`TopicInfo.Period` is `*int64`) and set it into state so refresh and import work correctly.
- Update the resource `.md` example file with a `period` usage example.
- Add unit tests (mock-based, using gomonkey) covering the new `period` parameter in Create, Read, and Update.

## Capabilities

### New Capabilities
- `clb-log-topic-period`: Enable the optional `period` parameter on the `tencentcloud_clb_log_topic` resource to allow users to configure the log retention lifecycle (in days) when creating and updating CLB log topics.

### Modified Capabilities
<!-- No existing spec requirements require modification; period is a new, independent parameter on the resource. -->

## Impact

- **Affected files:**
  - `tencentcloud/services/clb/resource_tc_clb_log_topic.go` — add `period` schema field, wire through Create/Read/Update flows
  - `tencentcloud/services/clb/service_tencentcloud_clb.go` — extend `CreateTopic` service function to accept and pass `Period` to the CLB `CreateTopic` API request
  - `tencentcloud/services/clb/resource_tc_clb_log_topic.md` — update documentation example with `period` usage
  - `tencentcloud/services/clb/resource_tc_clb_log_topic_test.go` — add mock-based unit tests for the `period` parameter
- **SDK dependency:** No SDK upgrade required. The vendored SDKs already contain `Period` on `CreateTopicRequest` (`clb/v20180317`, `*uint64`), `ModifyTopicRequest` (`cls/v20201016`, `*int64`), and `TopicInfo` (`cls/v20201016`, `*int64`).
- **Backward compatibility:** Fully backward compatible — `period` is Optional and defaults to not being set, so existing configurations continue to work unchanged. The API applies its own default (30 days) when `period` is not provided.
- **API constraints:** `Period` is accepted by both `CreateTopic` (CLB SDK) and `ModifyTopic` (CLS SDK), so the parameter is updatable in-place (not `ForceNew`). The Read path uses `DescribeTopics` (via the existing `DescribeClsTopicById` service helper) which returns `TopicInfo.Period`.
