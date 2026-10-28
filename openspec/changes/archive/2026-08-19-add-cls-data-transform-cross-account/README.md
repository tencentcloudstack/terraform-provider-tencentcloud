# Add Cross-Account Fields to CLS Data Transform DstResources

**Change ID**: `add-cls-data-transform-cross-account`
**Status**: Implemented
**Type**: Enhancement
**Estimated Effort**: 1 hour

---

## Quick Summary

Add support for **cross-account destination topic** configuration to the `dst_resources` block in the `tencentcloud_cls_data_transform` resource. This enables users to perform data transform tasks that deliver logs to topics located in **other TencentCloud accounts**.

---

## Problem

Currently, the `dst_resources` block of the `tencentcloud_cls_data_transform` resource only supports configuring the destination `topic_id` and `alias` for **same-account** log topics. The TencentCloud CLS API already supports cross-account delivery through the `DataTransformResouceInfo` struct, but the Terraform provider does not expose these capabilities.

**Impact**:
- Users cannot configure cross-account data transform delivery through Terraform
- Cross-account transforms created via console/API cannot be fully managed by Terraform
- State drift occurs when importing cross-account transforms (the cross-account fields are dropped)

---

## Solution

Add five new optional fields to the `dst_resources` block schema:

```hcl
resource "tencentcloud_cls_data_transform" "example" {
  # ... other fields ...
  dst_resources {
    topic_id        = tencentcloud_cls_topic.topic_dst.id
    alias           = "iac-test-dst"
    is_cross_account = true
    role_arn        = "qcs::cam::uin/123456789:roleName/cls-cross-account-role"
    external_id     = "external-id-value"
    topic_name      = "topic-name-in-target-account"
    logset_name     = "logset-name-in-target-account"
  }
}
```

### New Fields

| Field | Type | Description |
|------|------|-------------|
| `is_cross_account` | bool | Whether the destination topic is in another account. Default: `false`. |
| `role_arn` | string | Role ARN granted by the target account to the delivering account. |
| `external_id` | string | External ID used for the cross-account role assumption. |
| `topic_name` | string | Name of the destination topic (cross-account scenario). |
| `logset_name` | string | Name of the logset that contains the destination topic (cross-account scenario). |

---

## Implementation Overview

### Files to Modify

1. **`tencentcloud/services/cls/resource_tc_cls_data_transform.go`** (main changes)
   - Add 5 new fields to the `dst_resources` schema block
   - Handle the new fields in Create operation
   - Handle the new fields in Read operation
   - Handle the new fields in Update operation

2. **`tencentcloud/services/cls/resource_tc_cls_data_transform_test.go`** (tests)
   - Add unit test cases covering the new cross-account fields

3. **`tencentcloud/services/cls/resource_tc_cls_data_transform.md`** (documentation)
   - Add cross-account usage example

4. **Auto-generated**: `website/docs/r/cls_data_transform.html.markdown` (via `make doc`)

### SDK Structures (Already Available)

```go
// From vendor/.../cls/v20201016/models.go (line 6884)
type DataTransformResouceInfo struct {
    TopicId        *string `json:"TopicId,omitnil,omitempty"`
    Alias          *string `json:"Alias,omitnil,omitempty"`
    IsCrossAccount *bool   `json:"IsCrossAccount,omitnil,omitempty"`   // NEW
    RoleARN        *string `json:"RoleARN,omitnil,omitempty"`           // NEW
    ExternalId     *string `json:"ExternalId,omitnil,omitempty"`        // NEW
    TopicName      *string `json:"TopicName,omitnil,omitempty"`         // NEW
    LogsetName     *string `json:"LogsetName,omitnil,omitempty"`        // NEW
}
```

These fields are already present in:
- `CreateDataTransformRequest.DstResources []*DataTransformResouceInfo`
- `ModifyDataTransformRequest.DstResources []*DataTransformResouceInfo`
- `DescribeDataTransformInfoResponse.DataTransformTaskInfos[].DstResources []*DataTransformResouceInfo`

---

## Key Features

### 1. Full CRUD Support
- **Create**: Configure cross-account fields during transform creation
- **Read**: Populate cross-account fields from API response into state
- **Update**: Modify cross-account fields on transform update
- **Delete**: Standard deletion (field-agnostic)

### 2. Backward Compatibility
- **Additive change**: All new fields are optional
- **No breaking changes**: Existing same-account configs unaffected
- **Graceful upgrade**: Provider upgrade requires no migration

---

## Files

```
openspec/changes/add-cls-data-transform-cross-account/
├── README.md          (this file)
├── proposal.md        (detailed design)
├── tasks.md           (implementation tasks)
└── specs/
    └── cls-data-transform-cross-account-dst-resources/
        └── spec.md    (requirements and scenarios)
```

---

## Next Steps

1. **Review** proposal and spec documents
2. **Validate** with `pnpm dlx @fission-ai/openspec@1.2.0 validate add-cls-data-transform-cross-account --strict`
3. **Approve** the proposal
4. **Implement** the change

---

## References

- **API Docs**:
  - CreateDataTransform
  - ModifyDataTransform
  - DescribeDataTransformInfo

- **SDK**: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016/models.go` (line 6884)

- **Current Implementation**: `tencentcloud/services/cls/resource_tc_cls_data_transform.go`

---

**Created**: 2026-05-21
**Status**: Ready for Review
