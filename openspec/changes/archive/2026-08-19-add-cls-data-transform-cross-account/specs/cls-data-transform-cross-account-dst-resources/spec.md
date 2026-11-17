# Capability: CLS Data Transform Cross-Account DstResources Fields

## Overview

This specification defines the requirements for adding cross-account destination topic configuration fields to the `dst_resources` block in the `tencentcloud_cls_data_transform` resource. The TencentCloud CLS API already supports cross-account delivery through the `DataTransformResouceInfo` struct, but the Terraform provider does not currently expose these fields.

## ADDED Requirements

### Requirement: Support Cross-Account Destination Configuration

The resource MUST support configuring cross-account destination topic delivery by setting five new optional fields inside the `dst_resources` block: `is_cross_account`, `role_arn`, `external_id`, `topic_name`, and `logset_name`.

**Rationale**: The TencentCloud CLS API's `DataTransformResouceInfo` struct already includes `IsCrossAccount`, `RoleARN`, `ExternalId`, `TopicName`, and `LogsetName` fields. These enable delivering transformed logs to topics in other TencentCloud accounts. Without exposing them in Terraform, users cannot manage cross-account data transforms via Infrastructure as Code.

**Acceptance Criteria**:
- `dst_resources.is_cross_account` is an optional bool field
- `dst_resources.role_arn` is an optional string field
- `dst_resources.external_id` is an optional string field
- `dst_resources.topic_name` is an optional string field
- `dst_resources.logset_name` is an optional string field
- All five fields map to the corresponding `DataTransformResouceInfo` SDK fields
- Existing same-account configs without these fields continue to work

#### Scenario: Configure basic cross-account destination

**Given** a user wants to deliver transformed logs to a topic in another TencentCloud account  
**When** they specify:
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
**Then** the data transform is created with cross-account delivery configuration  
**And** the API accepts the configuration successfully  
**And** the state reflects all cross-account fields

#### Scenario: Same-account config remains backward compatible

**Given** a user has an existing same-account data transform config without cross-account fields  
**When** they upgrade the provider to the version with cross-account support  
**Then** the existing config continues to work  
**And** no changes are detected in `terraform plan`  
**And** no migration is required

---

### Requirement: Create Operation Integration

The Create operation MUST correctly convert the Terraform schema cross-account fields to the API request structure for `CreateDataTransform`.

**Acceptance Criteria**:
- Provider reads the five cross-account fields from each `dst_resources` block
- `is_cross_account` uses `helper.Bool()` to create `*bool` (maps to `DataTransformResouceInfo.IsCrossAccount`)
- `role_arn` uses `helper.String()` (maps to `DataTransformResouceInfo.RoleARN`)
- `external_id` uses `helper.String()` (maps to `DataTransformResouceInfo.ExternalId`)
- `topic_name` uses `helper.String()` (maps to `DataTransformResouceInfo.TopicName`)
- `logset_name` uses `helper.String()` (maps to `DataTransformResouceInfo.LogsetName`)
- API request includes the cross-account fields in `request.DstResources`

#### Scenario: Create data transform with cross-account dst_resources successfully

**Given** valid cross-account `dst_resources` configuration  
**When** the user runs `terraform apply` to create the data transform  
**Then** the provider reads the five cross-account fields from each `dst_resources` block  
**And** constructs `DataTransformResouceInfo` with the cross-account fields set  
**And** the API request includes the cross-account fields in `request.DstResources`  
**And** the data transform is created successfully  
**And** the Task ID is returned and set in state

---

### Requirement: Read Operation Integration

The Read operation MUST correctly parse the API response and populate Terraform state with the cross-account fields from `DescribeDataTransformInfo`.

**Acceptance Criteria**:
- Provider retrieves `DstResources` from `DataTransformTaskInfo`
- Each cross-account field has a nil check before assignment to the state map
- Fields are correctly mapped from PascalCase (SDK) to snake_case (Terraform)
- `RoleARN` SDK field maps to `role_arn` Terraform field
- State is set via `d.Set("dst_resources", ...)`

#### Scenario: Read data transform with cross-account dst_resources

**Given** an existing data transform with cross-account destination configuration  
**When** the provider reads the transform from API via `DescribeClsDataTransformById`  
**Then** the provider retrieves `DstResources` from `DataTransformTaskInfo`  
**And** converts each `DataTransformResouceInfo` to a schema map  
**And** the cross-account fields are populated in state  
**And** state matches the original configuration

#### Scenario: Read data transform with null cross-account fields

**Given** the API returns `DstResources` with nil cross-account fields (same-account transform)  
**When** the provider reads the response  
**Then** the provider handles nil pointers safely  
**And** no panic occurs  
**And** cross-account fields are not added to the state map  
**And** other fields (`topic_id`, `alias`) are still processed

#### Scenario: Import existing cross-account data transform

**Given** a cross-account data transform created via console or API  
**When** the user imports the transform into Terraform  
**Then** the provider reads the transform configuration  
**And** cross-account fields are populated in state  
**And** subsequent `terraform plan` shows no changes

---

### Requirement: Update Operation Integration

The Update operation MUST support modifying cross-account fields when the `dst_resources` block changes, via `ModifyDataTransform`.

**Acceptance Criteria**:
- `d.HasChange("dst_resources")` detects modifications to cross-account fields
- Update logic reads the new cross-account configuration from schema
- `DataTransformResouceInfo` is rebuilt with the new cross-account fields
- API request includes updated `DstResources` in `ModifyDataTransform`
- Update succeeds and subsequent read shows the new configuration

#### Scenario: Update cross-account configuration

**Given** an existing data transform with cross-account destination  
**When** the user modifies the `dst_resources` cross-account fields (e.g., changes `role_arn`)  
**Then** Terraform detects the change in `dst_resources`  
**And** the provider sends `ModifyDataTransform` request with updated `DataTransformResouceInfo`  
**And** the data transform is updated successfully  
**And** state reflects the new configuration

#### Scenario: Switch from same-account to cross-account

**Given** a same-account data transform  
**When** the user adds `is_cross_account = true` and cross-account fields to `dst_resources`  
**Then** the update operation handles the change  
**And** cross-account configuration is applied  
**And** the update succeeds without errors

---

### Requirement: Backward Compatibility

The addition of cross-account fields MUST NOT break existing same-account data transform configurations.

**Acceptance Criteria**:
- All five new fields are optional
- Existing same-account configs show 0 changes after provider upgrade
- `dst_resources` block structure unchanged for existing fields (`topic_id`, `alias`)
- CRUD operations work as before for same-account transforms

#### Scenario: Existing same-account transform remains unchanged

**Given** a data transform configured without cross-account fields  
**When** the user upgrades the provider to the version with cross-account support  
**Then** the same-account configuration continues to work  
**And** no changes are detected in `terraform plan`  
**And** no migration is required

---

### Requirement: Error Handling

Error handling MUST be robust and provide helpful feedback for common failure scenarios.

**Acceptance Criteria**:
- API errors are caught and wrapped with `tccommon.RetryError()`
- Error messages include API response details
- All pointer fields checked before dereference in Read operation
- Nil checks prevent runtime panics

#### Scenario: Handle API rejection of cross-account configuration

**Given** the user provides an invalid `role_arn` or `external_id`  
**When** the API rejects the configuration  
**Then** the provider captures the API error  
**And** the error message is passed to the user  
**And** the request body is logged for debugging

#### Scenario: Safe handling of nil pointers in responses

**Given** the API returns incomplete or null cross-account fields  
**When** the provider processes the response  
**Then** no panic occurs  
**And** nil checks (`if dstResources.IsCrossAccount != nil`, etc.) prevent crashes  
**And** the operation completes gracefully

---

### Requirement: Documentation

Complete and accurate documentation MUST be provided for cross-account field usage.

**Acceptance Criteria**:
- At least one cross-account example in `resource_tc_cls_data_transform.md`
- Example shows all five new fields
- Example is syntactically valid HCL
- Website docs are auto-generated via `make doc` in the finalization phase

#### Scenario: Usage example covers cross-account scenario

**Given** a user wants to configure cross-account data transform  
**When** the user reads the resource documentation  
**Then** examples show cross-account configuration  
**And** all five new fields are demonstrated  
**And** the example is syntactically valid HCL

---

## Implementation Notes

### API Mapping

The new parameters map to the following `DataTransformResouceInfo` SDK fields:

| Terraform Field | SDK Field | Type |
|-----------------|-----------|------|
| `is_cross_account` | `IsCrossAccount` | `*bool` |
| `role_arn` | `RoleARN` | `*string` |
| `external_id` | `ExternalId` | `*string` |
| `topic_name` | `TopicName` | `*string` |
| `logset_name` | `LogsetName` | `*string` |

### SDK Availability

All five fields are already present in the vendored SDK (`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016/models.go`, line 6884):

- `CreateDataTransformRequest.DstResources []*DataTransformResouceInfo` (line 3705)
- `ModifyDataTransformRequest.DstResources []*DataTransformResouceInfo` (line 18754)
- `DataTransformTaskInfo.DstResources []*DataTransformResouceInfo` (returned by `DescribeDataTransformInfo`, line 6966)

### Backward Compatibility

All new parameters are optional. Existing configurations without these parameters will continue to work without modification. The API defaults `IsCrossAccount` to `false` when omitted.
