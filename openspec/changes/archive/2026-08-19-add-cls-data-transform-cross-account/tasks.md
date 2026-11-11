# Tasks: Add Cross-Account Fields to CLS Data Transform DstResources

**Change ID**: `add-cls-data-transform-cross-account`
**Total Tasks**: 12
**Estimated Time**: 1 hour
**Status**: Completed

---

## Phase 1: Schema Definition (15 min)

### Task 1.1: Add Cross-Account Fields to DstResources Schema
**File**: `tencentcloud/services/cls/resource_tc_cls_data_transform.go`

Add five new fields inside the `dst_resources` block schema (the `Elem.Schema` map, around line 68-80), after the existing `alias` field:

```go
"dst_resources": {
    Optional:    true,
    Type:        schema.TypeList,
    Description: "Data transform des resources. If `func_type` is `1`, this parameter is required. If `func_type` is `2`, this parameter does not need to be filled in.",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "topic_id": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "Dst topic ID.",
            },
            "alias": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "Alias.",
            },
            "is_cross_account": {
                Type:        schema.TypeBool,
                Optional:    true,
                Description: "Whether the destination topic is in another TencentCloud account. `false`: same account (default); `true`: cross-account.",
            },
            "role_arn": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "In the cross-account scenario, the Role ARN value. The target account creates a role for the delivering account. Find this in the target account role list.",
            },
            "external_id": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "External ID value, used for cross-account role assumption. Find this in the target account role trust policy.",
            },
            "topic_name": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "Name of the destination topic. Used in cross-account scenario.",
            },
            "logset_name": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "Name of the logset that contains the destination topic. Used in cross-account scenario.",
            },
        },
    },
},
```

**Validation**:
- All five new fields are Optional
- `is_cross_account` is TypeBool (matches SDK `*bool`)
- The other four are TypeString (matches SDK `*string`)
- Existing `topic_id` and `alias` fields unchanged

---

## Phase 2: Create Operation (10 min)

### Task 2.1: Handle Cross-Account Fields in Create
**File**: `tencentcloud/services/cls/resource_tc_cls_data_transform.go`

In `resourceTencentCloudClsDataTransformCreate`, inside the `dst_resources` loop (around line 228-242), after setting `Alias`, add handling for the five new fields:

```go
if v, ok := dMap["dst_resources"]; ok {
    for _, item := range v.([]interface{}) {
        dMap := item.(map[string]interface{})
        dataTransformResouceInfo := cls.DataTransformResouceInfo{}
        if v, ok := dMap["topic_id"]; ok {
            dataTransformResouceInfo.TopicId = helper.String(v.(string))
        }

        if v, ok := dMap["alias"]; ok {
            dataTransformResouceInfo.Alias = helper.String(v.(string))
        }

        if v, ok := dMap["is_cross_account"]; ok {
            dataTransformResouceInfo.IsCrossAccount = helper.Bool(v.(bool))
        }

        if v, ok := dMap["role_arn"]; ok {
            dataTransformResouceInfo.RoleARN = helper.String(v.(string))
        }

        if v, ok := dMap["external_id"]; ok {
            dataTransformResouceInfo.ExternalId = helper.String(v.(string))
        }

        if v, ok := dMap["topic_name"]; ok {
            dataTransformResouceInfo.TopicName = helper.String(v.(string))
        }

        if v, ok := dMap["logset_name"]; ok {
            dataTransformResouceInfo.LogsetName = helper.String(v.(string))
        }

        request.DstResources = append(request.DstResources, &dataTransformResouceInfo)
    }
}
```

**Validation**:
- `is_cross_account` uses `helper.Bool()` (matches SDK `*bool`)
- All string fields use `helper.String()`
- Field names map correctly: `role_arn` → `RoleARN`, `external_id` → `ExternalId`, `topic_name` → `TopicName`, `logset_name` → `LogsetName`, `is_cross_account` → `IsCrossAccount`

---

## Phase 3: Read Operation (10 min)

### Task 3.1: Handle Cross-Account Fields in Read
**File**: `tencentcloud/services/cls/resource_tc_cls_data_transform.go`

In `resourceTencentCloudClsDataTransformRead`, inside the `dst_resources` loop (around line 381-398), after setting `alias`, add handling for the five new fields:

```go
if dataTransform.DstResources != nil {
    var dstResourcesList []interface{}
    for _, dstResources := range dataTransform.DstResources {
        dstResourcesMap := map[string]interface{}{}

        if dstResources.TopicId != nil {
            dstResourcesMap["topic_id"] = dstResources.TopicId
        }

        if dstResources.Alias != nil {
            dstResourcesMap["alias"] = dstResources.Alias
        }

        if dstResources.IsCrossAccount != nil {
            dstResourcesMap["is_cross_account"] = dstResources.IsCrossAccount
        }

        if dstResources.RoleARN != nil {
            dstResourcesMap["role_arn"] = dstResources.RoleARN
        }

        if dstResources.ExternalId != nil {
            dstResourcesMap["external_id"] = dstResources.ExternalId
        }

        if dstResources.TopicName != nil {
            dstResourcesMap["topic_name"] = dstResources.TopicName
        }

        if dstResources.LogsetName != nil {
            dstResourcesMap["logset_name"] = dstResources.LogsetName
        }

        dstResourcesList = append(dstResourcesList, dstResourcesMap)
    }

    _ = d.Set("dst_resources", dstResourcesList)
}
```

**Validation**:
- Each field has a nil check before assignment to the map
- Field name mapping is correct (snake_case in Terraform, PascalCase in SDK)
- `RoleARN` SDK field maps to `role_arn` Terraform field

---

## Phase 4: Update Operation (10 min)

### Task 4.1: Handle Cross-Account Fields in Update
**File**: `tencentcloud/services/cls/resource_tc_cls_data_transform.go`

In `resourceTencentCloudClsDataTransformUpdate`, inside the `dst_resources` change handling (around line 522-538), after setting `Alias`, add handling for the five new fields (mirrors Create):

```go
if d.HasChange("dst_resources") {
    if v, ok := d.GetOk("dst_resources"); ok {
        for _, item := range v.([]interface{}) {
            dataTransformResouceInfo := cls.DataTransformResouceInfo{}
            dMap := item.(map[string]interface{})
            if v, ok := dMap["topic_id"]; ok {
                dataTransformResouceInfo.TopicId = helper.String(v.(string))
            }

            if v, ok := dMap["alias"]; ok {
                dataTransformResouceInfo.Alias = helper.String(v.(string))
            }

            if v, ok := dMap["is_cross_account"]; ok {
                dataTransformResouceInfo.IsCrossAccount = helper.Bool(v.(bool))
            }

            if v, ok := dMap["role_arn"]; ok {
                dataTransformResouceInfo.RoleARN = helper.String(v.(string))
            }

            if v, ok := dMap["external_id"]; ok {
                dataTransformResouceInfo.ExternalId = helper.String(v.(string))
            }

            if v, ok := dMap["topic_name"]; ok {
                dataTransformResouceInfo.TopicName = helper.String(v.(string))
            }

            if v, ok := dMap["logset_name"]; ok {
                dataTransformResouceInfo.LogsetName = helper.String(v.(string))
            }

            request.DstResources = append(request.DstResources, &dataTransformResouceInfo)
        }
    }
}
```

**Validation**:
- Logic mirrors Create operation
- Wrapped in `d.HasChange("dst_resources")` check
- Field mapping is identical to Create

---

## Phase 5: Testing (10 min)

### Task 5.1: Add Unit Test Cases for Cross-Account Fields
**File**: `tencentcloud/services/cls/resource_tc_cls_data_transform_test.go`

Add unit test cases using gomonkey mocks (not the Terraform test suite) to cover the new cross-account fields. The tests should mock the CLS API client methods (`CreateDataTransform`, `DescribeDataTransformInfo`, `ModifyDataTransform`) and verify that:
- The Create operation correctly populates the cross-account fields in the request
- The Read operation correctly sets the cross-account fields in state
- The Update operation correctly populates the cross-account fields in the modify request

Refer to the existing test patterns in the file and other mock-based tests in the cls service.

**Validation**:
- Tests use gomonkey mocks (not `resource.Test` / Terraform acceptance suite)
- Tests cover Create, Read, and Update for the cross-account fields
- Tests compile correctly

---

## Phase 6: Documentation (5 min)

### Task 6.1: Add Cross-Account Example to Resource Documentation
**File**: `tencentcloud/services/cls/resource_tc_cls_data_transform.md`

Add a cross-account usage example to the documentation, after the existing example:

```markdown
Cross-account data transform example:

```hcl
resource "tencentcloud_cls_data_transform" "example_cross_account" {
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
```

**Validation**:
- Example is syntactically valid HCL
- Shows all five new cross-account fields
- Follows the existing documentation format

---

## Phase 7: Validation (finalization phase)

### Task 7.1: Format Code
**Command**: `gofmt -w tencentcloud/services/cls/resource_tc_cls_data_transform.go`

Handled by the tfpacer-finalize skill.

**Validation**:
- Code is properly formatted

### Task 7.2: Generate Website Documentation
**Command**: `make doc`

Handled by the tfpacer-finalize skill.

**Validation**:
- `website/docs/r/cls_data_transform.html.markdown` updated with new fields

### Task 7.3: Generate Changelog
Handled by the tfpacer-finalize skill via `.changelog/` directory.

**Validation**:
- Changelog entry created

---

## Summary

**Total Tasks**: 12
**Phases**: 7
**Estimated Time**: 1 hour

### Task Breakdown by Phase:
1. Schema Definition: 1 task (15 min)
2. Create Operation: 1 task (10 min)
3. Read Operation: 1 task (10 min)
4. Update Operation: 1 task (10 min)
5. Testing: 1 task (10 min)
6. Documentation: 1 task (5 min)
7. Validation: 3 tasks (finalization phase)

### Key Implementation Points:
1. Five new optional fields added to the `dst_resources` block
2. `is_cross_account` is TypeBool; the other four are TypeString
3. Same handling pattern in Create/Read/Update
4. Backward compatible (all fields optional)
5. Field name mapping: `is_cross_account`→`IsCrossAccount`, `role_arn`→`RoleARN`, `external_id`→`ExternalId`, `topic_name`→`TopicName`, `logset_name`→`LogsetName`

### Success Criteria:
- Users can configure cross-account destination topics
- Create/Read/Update work correctly for cross-account transforms
- Import of existing cross-account transforms populates all fields
- Same-account transforms remain unaffected
- Documentation includes cross-account example
- No breaking changes

---

**Ready for Implementation!**
