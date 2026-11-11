## Context

The `tencentcloud_cls_topic` resource manages CLS (Cloud Log Service) log topics. Currently, the resource only creates log topics (default BizType=0) and does not support creating metric topics (BizType=1). The CLS CreateTopic API supports a `BizType` parameter (uint64) to specify the topic type, and the DescribeTopics API returns `BizType` in the `TopicInfo` response. However, the ModifyTopic API does not support `BizType`, meaning the topic type is immutable after creation.

Current resource file: `tencentcloud/services/cls/resource_tc_cls_topic.go`

Cloud API details (from vendor):
- `CreateTopicRequest.BizType`: `*uint64` — Topic type. 0: log topic (default), 1: metric topic
- `TopicInfo.BizType`: `*uint64` — Returned by DescribeTopics, same semantics
- `ModifyTopicRequest`: Does NOT contain `BizType` field — immutable after creation

## Goals / Non-Goals

**Goals:**
- Add `biz_type` parameter to `tencentcloud_cls_topic` resource schema as Optional, Computed, ForceNew
- Pass `biz_type` to the CreateTopic API during resource creation
- Encode `biz_type` in the resource ID: plain `topic_id` for biz_type=0, `"topic_id#1"` for biz_type=1
- Parse and validate `biz_type` from the resource ID in Read, Update, and Delete methods
- Read `biz_type` from the DescribeTopics API response during resource read
- Prevent updates to `biz_type` by adding it to `immutableArgs` in the Update method (belt-and-suspenders with ForceNew)
- Update unit tests for the new parameter
- Update documentation with import examples for both log and metric topics

**Non-Goals:**
- This change does NOT add support for updating `biz_type` (ModifyTopic API does not support it)
- This change does NOT modify any other CLS resources or data sources
- This change does NOT add new resources or data sources

## Decisions

1. **Schema type: TypeInt with ForceNew**
   - `BizType` in the cloud API is `*uint64`, mapped to `TypeInt` in Terraform schema (standard pattern for integer enum fields)
   - ForceNew is required because ModifyTopic does not support `BizType` — changing it requires recreation
   - Computed is set so existing resources without `biz_type` specified will still read the value from the API

2. **Resource ID encoding**
   - When `biz_type=0` or not specified: resource ID is the plain `topic_id` (e.g., `"2f5764c1-c833-44c5-84c7-950979b2a278"`)
   - When `biz_type=1`: resource ID is `"topic_id#1"` (e.g., `"2f5764c1-c833-44c5-84c7-950979b2a278#1"`)
   - This follows the existing pattern in the codebase for encoding multiple identifiers in a single resource ID (e.g., CLB redirection, CAM policy attachments)
   - A `parseClsTopicId` helper function extracts `topicId` and `bizType` from the ID using `strings.Split`

3. **ID parsing in Read/Update/Delete**
   - The Read method parses `d.Id()` to extract `topicId` and `bizType`, then passes `topicId` and a `*uint64` bizType filter to `DescribeClsTopicById`
   - The Update method parses `d.Id()` to get `topicId` for the `ModifyTopicRequest.TopicId`
   - The Delete method parses `d.Id()` to get `topicId` for the `DeleteTopicRequest.TopicId`
   - The `biz_type` field in state is set from the parsed ID value, and also read from the API response for validation

4. **Add to immutableArgs in Update method**
   - Although ForceNew handles this at the schema level, adding `biz_type` to the `immutableArgs` array in the Update method provides an additional safety net and follows the existing pattern used by `partition_count` and `storage_type`

5. **Default value handling**
   - If `biz_type` is not specified by the user, the cloud API defaults to 0 (log topic). No explicit default is set in the Terraform schema; the Computed attribute handles reading back the actual value

6. **Pass BizType to DescribeTopics via DescribeClsTopicById**
   - The `DescribeClsTopicById` service method is modified to accept an optional `bizType *uint64` parameter
   - When `bizType` is non-nil, it is set on the `DescribeTopicsRequest.BizType` field, enabling the API to filter results by topic type
   - The Read method in `resource_tc_cls_topic.go` passes the parsed `bizType` to the service method
   - All other callers of `DescribeClsTopicById` (CLB log topic, Redis log delivery, cloud product log task) pass `nil`, preserving existing behavior without BizType filtering

## Risks / Trade-offs

- **[Backward compatibility]** → Adding an Optional+Computed+ForceNew parameter is backward compatible. Existing configurations without `biz_type` will continue to work unchanged. The ForceNew attribute only triggers recreation if the user explicitly changes the `biz_type` value after initial creation.
- **[ForceNew behavior]** → If a user adds `biz_type` to an existing resource configuration where it was previously unset, Terraform will see a diff (0 vs unset) and may trigger recreation. This is mitigated by using Computed, which allows the state to be populated from the API read without requiring explicit user input.
- **[ID format change]** → Resource IDs for metric topics (biz_type=1) now have a `"#1"` suffix. This requires the import command to use the `"topic_id#1"` format. Standard log topics use the plain `topic_id` format (backward compatible).
