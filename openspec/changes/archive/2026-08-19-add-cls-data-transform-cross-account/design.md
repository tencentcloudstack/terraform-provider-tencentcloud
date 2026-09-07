# Design Document: Add Cross-Account Fields to CLS Data Transform DstResources

## Context

The `tencentcloud_cls_data_transform` resource manages CLS (Cloud Log Service) data transform tasks. Each transform task delivers processed logs to one or more destination topics via the `dst_resources` block.

The `DataTransformResouceInfo` struct in the vendored SDK (`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016/models.go`, line 6884) already includes five cross-account fields that are not yet exposed in the Terraform schema:

- `IsCrossAccount *bool` - Whether the destination topic is in another account
- `RoleARN *string` - Role ARN for cross-account role assumption
- `ExternalId *string` - External ID for cross-account role assumption
- `TopicName *string` - Name of the destination topic
- `LogsetName *string` - Name of the logset containing the destination topic

These fields are present in:
- `CreateDataTransformRequest.DstResources` (line 3705)
- `ModifyDataTransformRequest.DstResources` (line 18754)
- `DataTransformTaskInfo.DstResources` (returned by `DescribeDataTransformInfo`, line 6966)

No SDK upgrade is needed — all fields are already available in the vendor directory.

## Goals / Non-Goals

**Goals:**
- Expose the five cross-account fields as optional schema fields inside the `dst_resources` block
- Handle the new fields in Create, Read, and Update operations
- Add unit test cases using gomonkey mocks
- Add cross-account example to resource documentation

**Non-Goals:**
- No validation enforcing that `role_arn`/`external_id` must be set when `is_cross_account = true` (let the API validate)
- No changes to the Delete operation (field-agnostic)
- No changes to the service layer (`DescribeClsDataTransformById` already returns the full `DataTransformTaskInfo`)

## Decisions

### Field placement
- Add the five new fields inside the existing `dst_resources` block schema (as nested fields), matching the 1:1 mapping with the `DataTransformResouceInfo` SDK struct.
- Each `dst_resources` entry can independently be cross-account or same-account.

### Optionality
- All five new fields are `Optional`, matching the API's `omitnil,omitempty` tags.
- `is_cross_account` defaults to `false` on the API side when omitted.

### Type mapping
- `is_cross_account` → `schema.TypeBool` (SDK: `*bool`, use `helper.Bool()`)
- `role_arn` → `schema.TypeString` (SDK: `*string`, use `helper.String()`)
- `external_id` → `schema.TypeString` (SDK: `*string`, use `helper.String()`)
- `topic_name` → `schema.TypeString` (SDK: `*string`, use `helper.String()`)
- `logset_name` → `schema.TypeString` (SDK: `*string`, use `helper.String()`)

### Field name mapping
- `is_cross_account` → `IsCrossAccount`
- `role_arn` → `RoleARN` (note: SDK uses `RoleARN`, not `RoleArn`)
- `external_id` → `ExternalId` (note: SDK uses `ExternalId`, not `ExternalID`)
- `topic_name` → `TopicName`
- `logset_name` → `LogsetName`

### CRUD mapping
- **Create**: Read the five fields from each `dst_resources` map, assign to `DataTransformResouceInfo` fields, append to `request.DstResources`.
- **Read**: For each `DataTransformResouceInfo` in the API response, nil-check each cross-account field, then add to the `dstResourcesMap`. Set via `d.Set("dst_resources", ...)`.
- **Update**: Mirror Create logic inside `d.HasChange("dst_resources")`.

## Risks / Trade-offs

- [No schema-level validation for cross-account field combinations] → Mitigation: the API rejects invalid combinations with clear errors; document that `role_arn` and `external_id` are required for cross-account.
- [Null handling in Read] → Mitigation: add nil checks (`if dstResources.IsCrossAccount != nil`) before map assignment, following the existing pattern for `TopicId` and `Alias`.
- [Field name casing] → Mitigation: carefully map `RoleARN` and `ExternalId` (not `RoleArn`/`ExternalID`).

## Migration Plan

Purely additive; no migration. New optional fields added to existing `dst_resources` block. Rollback = revert the additive schema/CRUD changes and documentation.
