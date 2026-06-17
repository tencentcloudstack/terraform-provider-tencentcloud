## Context

The `tencentcloud_tdmq_topic` resource manages TDMQ Pulsar topics. It currently supports basic parameters (environ_id, topic_name, partitions, topic_type, cluster_id, pulsar_topic_type, remark) but does not support the `Tags` parameter.

The TDMQ SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217`) supports Tags in:
- `CreateTopicRequest.Tags` — `[]*Tag` (TagKey, TagValue) for setting tags at creation time.
- `Topic.Tags` — `[]*Tag` returned by `DescribeTopics` for reading tags.
- `ModifyTopicRequest` — does **not** have a Tags field, so tags cannot be updated via this API.
- `DeleteTopicsRequest` — does not involve tags.

The TencentCloud Tag service SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813`) provides:
- `TagResources` — binds tags to cloud resources using the resource six-segment format.
- `UnTagResources` — unbinds tag keys from cloud resources using the resource six-segment format.

The resource code is in `tencentcloud/services/tpulsar/resource_tc_tdmq_topic.go`, and the service layer is in `tencentcloud/services/tdmq/service_tencentcloud_tdmq.go`.

## Goals / Non-Goals

**Goals:**
- Add an `Optional` `tags` parameter (type `map[string]string`) to the `tencentcloud_tdmq_topic` resource schema.
- Pass tags to the `CreateTopic` API during resource creation by converting the map to `[]*tdmq.Tag`.
- Read tags from the `Topic` struct returned by `DescribeTopics` and set them back into Terraform state.
- Support in-place tag updates in the Update function using the Tag service's `UnTagResources` and `TagResources` APIs.
- Add unit tests using gomonkey mocks for the new tags functionality.
- Update the resource documentation (`.md` file) with example usage including tags.

**Non-Goals:**
- Tag-based filtering in data sources.

## Decisions

### 1. Use Tag service APIs (`TagResources`/`UnTagResources`) for tag updates

**Decision**: Use the TencentCloud Tag service's `TagResources` and `UnTagResources` APIs to handle tag updates in the Update function, since the `ModifyTopic` API does not support a `Tags` field.

**Rationale**: The Tag service provides a universal mechanism for managing tags on any TencentCloud resource. By using `UnTagResources` to remove deleted tag keys and `TagResources` to add/update tags, we can support in-place tag modifications without requiring resource recreation.

**Alternative considered**: Marking `tags` as `ForceNew` to force resource recreation on tag changes. Rejected because using the Tag service APIs provides a better user experience by avoiding unnecessary resource destruction and recreation.

### 2. Use `map[string]string` type for tags

**Decision**: Use `schema.TypeMap` with `schema.TypeString` elements, following the standard Terraform tags pattern.

**Rationale**: This is the conventional Terraform pattern for tags. The conversion from `map[string]string` to `[]*tdmq.Tag` (TagKey/TagValue) is straightforward.

### 3. Modify the service layer `CreateTdmqTopic` function signature

**Decision**: Add a `tags []*tdmq.Tag` parameter to the `CreateTdmqTopic` service function.

**Rationale**: The service layer function directly constructs the API request. Adding the tags parameter keeps the pattern consistent with how other parameters are passed.

### 4. Use `svctag.DiffTags` for computing tag differences

**Decision**: Use the existing `svctag.DiffTags` utility to compute which tags to add/update (`replaceTags`) and which to delete (`deleteTags`).

**Rationale**: This utility is already used by other resources in the project (e.g., `resource_tc_tdmq_instance.go`) and correctly handles the diff logic.

## Risks / Trade-offs

- **[Risk] Breaking existing configurations** → Mitigated: `tags` is `Optional`, so existing configurations without tags continue to work unchanged.
- **[Risk] Tag service API consistency** → The Tag service APIs operate on the resource six-segment format (`qcs::tdmq:{region}:uin/{account}:topic/{topicName}`). The resource name is constructed using `tccommon.BuildTagResourceName("tdmq", "topic", region, topicName)`.
- **[Trade-off] Tag service vs ModifyTopic** → Using the Tag service adds a dependency on the tag SDK but enables in-place updates, which is a better UX than ForceNew.
