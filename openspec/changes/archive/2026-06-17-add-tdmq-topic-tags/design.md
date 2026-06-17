## Context

The `tencentcloud_tdmq_topic` resource manages TDMQ Pulsar topics. It currently supports basic parameters (environ_id, topic_name, partitions, topic_type, cluster_id, pulsar_topic_type, remark) but does not support the `Tags` parameter.

The TDMQ SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217`) supports Tags in:
- `CreateTopicRequest.Tags` — `[]*Tag` (TagKey, TagValue) for setting tags at creation time.
- `Topic.Tags` — `[]*Tag` returned by `DescribeTopics` for reading tags.
- `ModifyTopicRequest` — does **not** have a Tags field, so tags cannot be updated after creation.
- `DeleteTopicsRequest` — does not involve tags.

The resource code is in `tencentcloud/services/tpulsar/resource_tc_tdmq_topic.go`, and the service layer is in `tencentcloud/services/tdmq/service_tencentcloud_tdmq.go`.

## Goals / Non-Goals

**Goals:**
- Add an `Optional`, `ForceNew` `tags` parameter (type `map[string]string`) to the `tencentcloud_tdmq_topic` resource schema.
- Pass tags to the `CreateTopic` API during resource creation by converting the map to `[]*tdmq.Tag`.
- Read tags from the `Topic` struct returned by `DescribeTopics` and set them back into Terraform state.
- Add unit tests using gomonkey mocks for the new tags functionality.
- Update the resource documentation (`.md` file) with example usage including tags.

**Non-Goals:**
- Tag update support (ModifyTopic API does not support Tags).
- Tag-based filtering in data sources.
- Integration with the common tag service (this uses TDMQ-native tags, not the generic TencentCloud tag API).

## Decisions

### 1. Use `ForceNew` for the `tags` parameter

**Decision**: Mark `tags` as `ForceNew` so that changing tags forces resource recreation.

**Rationale**: The `ModifyTopic` API does not include a `Tags` field, so there is no way to update tags after creation. Using `ForceNew` is the standard Terraform pattern when an attribute cannot be updated in-place.

**Alternative considered**: Adding tags to the `immutableArgs` check in the update function. Rejected because `ForceNew` is the idiomatic Terraform approach and provides better plan output to users.

### 2. Use `map[string]string` type for tags

**Decision**: Use `schema.TypeMap` with `schema.TypeString` elements, following the standard Terraform tags pattern.

**Rationale**: This is the conventional Terraform pattern for tags. The conversion from `map[string]string` to `[]*tdmq.Tag` (TagKey/TagValue) is straightforward.

### 3. Modify the service layer `CreateTdmqTopic` function signature

**Decision**: Add a `tags []*tdmq.Tag` parameter to the `CreateTdmqTopic` service function.

**Rationale**: The service layer function directly constructs the API request. Adding the tags parameter keeps the pattern consistent with how other parameters are passed.

## Risks / Trade-offs

- **[Risk] Breaking existing configurations** → Mitigated: `tags` is `Optional`, so existing configurations without tags continue to work unchanged.
- **[Risk] Tags cannot be updated** → Mitigated: `ForceNew` clearly communicates to users that changing tags requires recreation. This is documented in the resource docs.
- **[Trade-off] ForceNew vs immutableArgs** → ForceNew provides better UX (Terraform shows "forces replacement" in plan) at the cost of resource recreation when tags change.
