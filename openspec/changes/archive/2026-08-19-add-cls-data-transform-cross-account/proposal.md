# Proposal: Add Cross-Account Fields to CLS Data Transform DstResources

## Metadata
- **Change ID**: `add-cls-data-transform-cross-account`
- **Status**: Implemented
- **Created**: 2026-05-21
- **Author**: AI Assistant
- **Type**: Enhancement
- **Estimated Effort**: 1 hour

---

## Problem Statement

The `tencentcloud_cls_data_transform` resource manages CLS (Cloud Log Service) data transform tasks. Each transform task can deliver processed logs to one or more destination topics via the `dst_resources` block.

Currently, the `dst_resources` block schema only exposes two fields:
- `topic_id` (Required) - Destination topic ID
- `alias` (Required) - Alias for the destination

This only supports **same-account** delivery: the destination topic must be in the same TencentCloud account as the source topic.

However, the TencentCloud CLS API's `DataTransformResouceInfo` struct already supports **cross-account** delivery with five additional fields:
- `IsCrossAccount` - Whether the destination topic is in another account
- `RoleARN` - Role ARN for cross-account role assumption
- `ExternalId` - External ID for the role assumption
- `TopicName` - Name of the destination topic
- `LogsetName` - Name of the logset containing the destination topic

### Current Limitations

Users who want to deliver transformed logs to a **different TencentCloud account** cannot configure this through Terraform, even though:
- The TencentCloud SDK already includes all five fields in `DataTransformResouceInfo`
- The CLS API supports cross-account delivery in CreateDataTransform / ModifyDataTransform
- DescribeDataTransformInfo returns these fields for existing cross-account transforms

### User Impact

- **Incomplete feature coverage**: Terraform does not expose all available API capabilities
- **Manual configuration required**: Cross-account transforms must be created via console or API directly
- **State drift**: Importing an existing cross-account transform loses the cross-account configuration

---

## Proposed Solution

Add five new optional fields to the `dst_resources` block schema in `tencentcloud_cls_data_transform`:

1. `is_cross_account` (Optional, bool)
2. `role_arn` (Optional, string)
3. `external_id` (Optional, string)
4. `topic_name` (Optional, string)
5. `logset_name` (Optional, string)

### Schema Addition

```hcl
resource "tencentcloud_cls_data_transform" "example" {
  func_type    = 1
  src_topic_id = tencentcloud_cls_topic.topic_src.id
  name         = "tf-example-cross-account"
  etl_content  = "ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")"
  task_type    = 3
  enable_flag  = 1

  dst_resources {
    topic_id         = "topic-id-in-target-account"
    alias            = "cross-account-dst"
    is_cross_account = true
    role_arn         = "qcs::cam::uin/123456789:roleName/cls-cross-account-role"
    external_id      = "external-id-value"
    topic_name       = "target-topic-name"
    logset_name      = "target-logset-name"
  }
}
```

### Field Specification

#### `dst_resources.is_cross_account` (Optional, bool)
- **Type**: Bool
- **Description**: Whether the destination topic is in another TencentCloud account. `false`: same account (default); `true`: cross-account.
- **Default**: `false`

#### `dst_resources.role_arn` (Optional, string)
- **Type**: String
- **Description**: In the cross-account scenario, the Role ARN value. The target (delivering-to) account creates a role for the delivering account. Find this in the target account's role list.
- **Required when**: `is_cross_account = true`

#### `dst_resources.external_id` (Optional, string)
- **Type**: String
- **Description**: External ID value, used for cross-account role assumption. Find this in the target account's role trust policy.
- **Required when**: `is_cross_account = true`

#### `dst_resources.topic_name` (Optional, string)
- **Type**: String
- **Description**: Name of the destination topic (used in cross-account scenario where the topic_id may not be directly accessible).

#### `dst_resources.logset_name` (Optional, string)
- **Type**: String
- **Description**: Name of the logset that contains the destination topic (cross-account scenario).

---

## Technical Design

### API Mapping

| Operation | API Endpoint | Field Mapping |
|-----------|-------------|---------------|
| **Create** | CreateDataTransform | `dst_resources.*` → `request.DstResources[].*` |
| **Read** | DescribeDataTransformInfo | `response.DataTransformTaskInfos[].DstResources[].*` → `dst_resources.*` |
| **Update** | ModifyDataTransform | `dst_resources.*` → `request.DstResources[].*` |

### SDK Structures (Already Available)

```go
// From vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016/models.go (line 6884)
type DataTransformResouceInfo struct {
    // Existing fields
    TopicId *string `json:"TopicId,omitnil,omitempty" name:"TopicId"`
    Alias   *string `json:"Alias,omitnil,omitempty" name:"Alias"`

    // New fields (already in SDK)
    IsCrossAccount *bool   `json:"IsCrossAccount,omitnil,omitempty" name:"IsCrossAccount"`
    RoleARN        *string `json:"RoleARN,omitnil,omitempty" name:"RoleARN"`
    ExternalId     *string `json:"ExternalId,omitnil,omitempty" name:"ExternalId"`
    TopicName      *string `json:"TopicName,omitnil,omitempty" name:"TopicName"`
    LogsetName     *string `json:"LogsetName,omitnil,omitempty" name:"LogsetName"`
}
```

The `DstResources` field is `[]*DataTransformResouceInfo` and is present in:
- `CreateDataTransformRequest` (line 3705)
- `ModifyDataTransformRequest` (line 18754)
- `DataTransformTaskInfo` (returned by DescribeDataTransformInfo, line 6966)

### Implementation Files

1. **Resource File**: `tencentcloud/services/cls/resource_tc_cls_data_transform.go`
   - Add 5 fields to the `dst_resources` schema block (around line 67-81)
   - Handle new fields in Create operation (around line 228-242)
   - Handle new fields in Read operation (around line 381-398)
   - Handle new fields in Update operation (around line 522-538)

2. **Documentation**: `tencentcloud/services/cls/resource_tc_cls_data_transform.md`
   - Add cross-account usage example

3. **Website Documentation**: Auto-generated via `make doc` (finalization phase)

---

## Implementation Plan

### Phase 1: Schema Definition (15 min)
1. Add `is_cross_account` field to `dst_resources` schema
2. Add `role_arn` field to `dst_resources` schema
3. Add `external_id` field to `dst_resources` schema
4. Add `topic_name` field to `dst_resources` schema
5. Add `logset_name` field to `dst_resources` schema

### Phase 2: Create Operation (10 min)
1. Read the new fields from each `dst_resources` map in Create
2. Assign to the corresponding `DataTransformResouceInfo` fields
3. Use `helper.Bool()` for `is_cross_account` and `helper.String()` for the rest

### Phase 3: Read Operation (10 min)
1. Populate the new fields from `dstResources` (API response) into the `dstResourcesMap`
2. Add nil checks before each `d.Set`/map assignment

### Phase 4: Update Operation (10 min)
1. Read the new fields from each `dst_resources` map in Update (mirrors Create)
2. Assign to the corresponding `DataTransformResouceInfo` fields

### Phase 5: Testing (10 min)
1. Add unit test cases using gomonkey mocks for the cross-account scenario

### Phase 6: Documentation (5 min)
1. Add cross-account example to `resource_tc_cls_data_transform.md`

---

## Design Considerations

### 1. Field Placement

**Decision**: Add the new fields inside the existing `dst_resources` block (as nested schema fields).

**Rationale**:
- The API models these fields as part of `DataTransformResouceInfo`, which is exactly the structure represented by the `dst_resources` block.
- Placing them at the resource top level would break the 1:1 mapping with the API struct.
- Each `dst_resources` entry can independently be cross-account or same-account.

### 2. Optionality

**Decision**: All five new fields are `Optional`.

**Rationale**:
- The API marks all five fields as `omitnil,omitempty`, meaning they are optional.
- `is_cross_account` defaults to `false` on the API side when omitted.
- Same-account configs do not need to specify these fields.
- This preserves full backward compatibility.

### 3. No Validation Constraints

**Decision**: Do not add `ConflictsWith` or `RequiredWhen` constraints between `is_cross_account` and `role_arn`/`external_id`.

**Rationale**:
- While logically `role_arn` and `external_id` are needed when `is_cross_account = true`, enforcing this at the Terraform schema level is complex for nested blocks.
- The API will reject invalid combinations with a clear error message.
- Keeping the schema simple mirrors the API's own optionality.

### 4. Backward Compatibility

**Impact**: Fully backward compatible
- New optional fields only
- Existing same-account `dst_resources` configs unaffected
- No breaking changes to schema structure
- No migration required

---

## Example Usage

### Example: Cross-Account Data Transform

```hcl
resource "tencentcloud_cls_logset" "logset_src" {
  logset_name = "tf-example-src"
  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_cls_topic" "topic_src" {
  topic_name           = "tf-example_src"
  logset_id            = tencentcloud_cls_logset.logset_src.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_cls_data_transform" "example" {
  func_type    = 1
  src_topic_id = tencentcloud_cls_topic.topic_src.id
  name         = "tf-example-cross-account"
  etl_content  = "ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")"
  task_type    = 3
  enable_flag  = 1

  dst_resources {
    topic_id         = "topic-id-in-target-account"
    alias            = "cross-account-dst"
    is_cross_account = true
    role_arn         = "qcs::cam::uin/123456789:roleName/cls-cross-account-role"
    external_id      = "external-id-value"
    topic_name       = "target-topic-name"
    logset_name      = "target-logset-name"
  }
}
```

---

## Impact Analysis

### Breaking Changes
**None** - This is a purely additive change.

### Backward Compatibility
Fully backward compatible:
- New optional fields
- Existing same-account `dst_resources` configs continue to work
- No changes to existing schema fields
- Import of existing cross-account transforms will now populate the cross-account fields

### API Dependencies
- SDK structures already available in vendor
- API supports cross-account fields in Create/Modify/Describe operations
- No version constraints

---

## Risks and Mitigation

### Risk 1: Cross-Account Role Misconfiguration
**Risk**: User provides an invalid `role_arn` or `external_id`
**Impact**: Medium - API will reject the transform creation
**Mitigation**:
- Clear documentation describing the required role setup
- API provides detailed error messages
- Let the API validate the role assumption

### Risk 2: Missing Cross-Account Fields
**Risk**: User sets `is_cross_account = true` but omits `role_arn`/`external_id`
**Impact**: Low - API will reject with a clear error
**Mitigation**:
- Document that `role_arn` and `external_id` are required for cross-account
- Rely on API validation

### Risk 3: Null Handling in Read
**Risk**: API returns `nil` for cross-account fields on same-account transforms
**Impact**: Low - handled with nil checks
**Mitigation**:
- Add `if dstResources.IsCrossAccount != nil` checks before map assignment
- Follow the existing pattern used for `TopicId` and `Alias`

---

## Success Criteria

### Functional Requirements
- Users can configure cross-account destination topics
- Create/Read/Update/Delete operations work correctly for cross-account transforms
- Import of existing cross-account transforms populates all fields
- Same-account transforms remain unaffected

### Non-Functional Requirements
- Code follows project conventions (naming, error handling, nil checks)
- Documentation is complete and includes cross-account example
- Backward compatible with existing configurations

---

## Alternatives Considered

### Alternative 1: Separate Resource for Cross-Account Config
**Approach**: Split cross-account config into a separate resource
**Pros**: Modularity
**Cons**: Over-engineering, breaks existing pattern, complex for users
**Verdict**: Rejected - Cross-account fields are integral to `dst_resources`

### Alternative 2: Top-Level Fields Instead of Nested
**Approach**: Add `is_cross_account`, `role_arn`, etc. as top-level resource fields
**Pros**: Simpler schema
**Cons**: Breaks 1:1 mapping with API struct, doesn't support multiple dst_resources with mixed account types
**Verdict**: Rejected - Each dst_resources entry can independently be cross-account

---

**END OF PROPOSAL**
