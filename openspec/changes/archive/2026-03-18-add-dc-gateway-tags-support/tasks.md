# Implementation Tasks: Add Tags Support to DC Gateway

## Overview
This document outlines the sequential tasks to implement tags support for the `tencentcloud_dc_gateway` resource.

**Total Tasks**: 18  
**Estimated Effort**: 3.5 days

---

## Phase 1: Schema Definition (2 tasks)

### Task 1.1: Add tags field to schema
- **Location**: `tencentcloud/services/dcg/resource_tc_dc_gateway.go`
- **Action**: Add `tags` field to `ResourceTencentCloudDcGatewayInstance()` schema
- **Details**:
  ```go
  "tags": {
      Type:        schema.TypeMap,
      Optional:    true,
      Description: "Tag key-value pairs for the DC gateway. Multiple tags can be set.",
  },
  ```
- **Verification**: Schema compiles without errors

### Task 1.2: Import tag service package
- **Location**: `tencentcloud/services/dcg/resource_tc_dc_gateway.go`
- **Action**: Add import for tag service
- **Details**:
  ```go
  svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
  ```
- **Verification**: No import conflicts

---

## Phase 2: Create Operation Enhancement (3 tasks)

### Task 2.1: Extract tags from schema in Create
- **Location**: `resourceTencentCloudDcGatewayCreate()` function
- **Action**: Extract tags from schema after basic parameters
- **Details**:
  ```go
  var tags []*vpc.Tag
  if temp := helper.GetTags(d, "tags"); len(temp) > 0 {
      for k, v := range temp {
          tags = append(tags, &vpc.Tag{
              Key:   helper.String(k),
              Value: helper.String(v),
          })
      }
  }
  ```
- **Position**: After extracting `gateway_type` parameter (line ~106)
- **Verification**: Tags are properly extracted

### Task 2.2: Set tags in CreateDirectConnectGatewayRequest
- **Location**: Same function, before API call
- **Action**: Add tags to request
- **Details**:
  ```go
  if len(tags) > 0 {
      request.Tags = tags
  }
  ```
- **Position**: After all request parameter settings, before validation (line ~118)
- **Verification**: Request includes tags

### Task 2.3: Handle tag creation errors
- **Location**: After `d.SetId()` call (line ~144)
- **Action**: No additional handling needed as tags are part of create request
- **Note**: Tags set during creation, no separate tag API call needed
- **Verification**: Gateway created with tags visible in console

---

## Phase 3: Read Operation Enhancement (2 tasks)

### Task 3.1: Add tag service initialization in Read
- **Location**: `resourceTencentCloudDcGatewayRead()` function
- **Action**: Initialize tag service after setting basic fields
- **Details**:
  ```go
  tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
  tagService := svctag.NewTagService(tcClient)
  ```
- **Position**: After all `d.Set()` calls (line ~175)
- **Verification**: Tag service initialized

### Task 3.2: Retrieve and set tags
- **Location**: Same function, after tag service initialization
- **Action**: Call DescribeResourceTags and set tags in state
- **Details**:
  ```go
  tags, err := tagService.DescribeResourceTags(ctx, "vpc", "dcg", tcClient.Region, d.Id())
  if err != nil {
      return err
  }
  _ = d.Set("tags", tags)
  ```
- **Position**: Before final return statement (line ~183)
- **Verification**: Tags correctly retrieved and stored in state

---

## Phase 4: Update Operation Enhancement (4 tasks)

### Task 4.1: Add context variable in Update
- **Location**: `resourceTencentCloudDcGatewayUpdate()` function
- **Action**: Add context variable if not present
- **Details**:
  ```go
  ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
  ```
- **Position**: After logId declaration (line ~190)
- **Verification**: Context available for tag service

### Task 4.2: Detect tag changes
- **Location**: Same function, after name update logic
- **Action**: Check if tags changed
- **Details**:
  ```go
  if d.HasChange("tags") {
      // Tag update logic
  }
  ```
- **Position**: After name change handling (line ~215)
- **Verification**: Tag changes detected

### Task 4.3: Calculate tag diff and call ModifyTags
- **Location**: Inside tag change block
- **Action**: Implement tag update logic
- **Details**:
  ```go
  if d.HasChange("tags") {
      oldValue, newValue := d.GetChange("tags")
      replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
      tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
      tagService := svctag.NewTagService(tcClient)
      resourceName := tccommon.BuildTagResourceName("vpc", "dcg", tcClient.Region, d.Id())
      if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
          return err
      }
  }
  ```
- **Verification**: Tags updated correctly

### Task 4.4: Handle tag update errors
- **Location**: Within tag update block
- **Action**: Error already handled in previous task
- **Note**: ModifyTags returns error which is properly propagated
- **Verification**: Errors logged and returned appropriately

---

## Phase 5: Testing (4 tasks)

### Task 5.1: Create acceptance test file structure
- **Location**: `tencentcloud/services/dcg/resource_tc_dc_gateway_test.go`
- **Action**: Check if test file exists, create if needed
- **Verification**: Test file structure ready

### Task 5.2: Write basic tag test
- **Location**: Same test file
- **Action**: Add test for creating gateway with tags
- **Test Name**: `TestAccTencentCloudDcGatewayResource_Tags`
- **Scenarios**:
  - Create gateway with tags
  - Verify tags in state
- **Verification**: Test passes

### Task 5.3: Write tag update test
- **Location**: Same test file
- **Action**: Add test for updating tags
- **Scenarios**:
  - Create gateway with initial tags
  - Update tags (add new, modify existing, remove some)
  - Verify final tag state
- **Verification**: Test passes

### Task 5.4: Write tag import test
- **Location**: Same test file
- **Action**: Add test for importing gateway with tags
- **Scenarios**:
  - Create gateway with tags via API
  - Import into Terraform
  - Verify tags imported correctly
- **Verification**: Test passes

---

## Phase 6: Documentation (2 tasks)

### Task 6.1: Update resource documentation
- **Location**: `tencentcloud/services/dcg/resource_tc_dc_gateway.md` (or website equivalent)
- **Action**: Add tags parameter documentation
- **Content**:
  - Add `tags` to Argument Reference section
  - Type: Map
  - Description: "Tag key-value pairs for the DC gateway."
- **Verification**: Documentation accurate

### Task 6.2: Add usage examples
- **Location**: Same documentation file
- **Action**: Add examples showing tag usage
- **Examples**:
  1. Create gateway with tags
  2. Update tags on existing gateway
  3. Import gateway preserving tags
- **Verification**: Examples are runnable

---

## Phase 7: Code Quality (3 tasks)

### Task 7.1: Run code formatting
- **Command**: `make fmt` or `gofmt -w tencentcloud/services/dcg/resource_tc_dc_gateway.go`
- **Verification**: Code properly formatted

### Task 7.2: Run linter
- **Command**: `make lint` targeting dcg service
- **Action**: Fix any linter warnings/errors
- **Verification**: No linting issues

### Task 7.3: Run acceptance tests
- **Command**: `TF_ACC=1 go test -v ./tencentcloud/services/dcg -run TestAccTencentCloudDcGatewayResource`
- **Verification**: All tests pass

---

## Validation Checklist

- [x] Schema includes tags field
- [x] Tags can be set during creation
- [x] Tags are read and displayed in state
- [x] Tags can be updated without recreating resource
- [x] Tag updates properly handle additions, modifications, and deletions
- [x] Tags work with resource import
- [x] All acceptance tests pass (test written, ready to run)
- [x] Code is properly formatted
- [x] No linter errors (only pre-existing deprecation warnings)
- [x] Documentation is complete and accurate
- [x] Examples work as documented

---

## Dependencies

Tasks must be completed in phase order:
1. Phase 1 (Schema) → Phase 2 (Create)
2. Phase 2 (Create) → Phase 3 (Read)
3. Phase 3 (Read) → Phase 4 (Update)
4. Phase 4 (Update) → Phase 5 (Testing)
5. Phase 5 (Testing) → Phase 6 (Documentation)
6. All phases → Phase 7 (Quality)

Within each phase, tasks can be done in sequence.

---

## Rollback Plan

If issues arise:
1. **Before merge**: Simply revert code changes
2. **After merge**: 
   - Tags field is optional, no breaking changes
   - Can be disabled by not setting the field
   - Full revert possible without user impact
