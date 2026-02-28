# Tasks: Enhance MongoDB Backups Data Source

## Task List

### 1. Schema Enhancement - Add Pagination Parameters
**Status**: completed  
**Estimated Time**: 15 minutes  
**Dependencies**: None

Add `limit` and `offset` input parameters to the data source schema.

**Acceptance Criteria**:
- `limit` parameter is Optional, TypeInt, with description and validation (max 100)
- `offset` parameter is Optional, TypeInt, with description and validation (min 0)
- Parameters follow existing schema patterns

**Files Modified**:
- `tencentcloud/services/mongodb/data_source_tc_mongodb_instance_backups.go`

---

### 2. Schema Enhancement - Add Missing Output Fields
**Status**: completed  
**Estimated Time**: 20 minutes  
**Dependencies**: None

Add missing fields to the `backup_list` output schema.

**Acceptance Criteria**:
- `back_id` field is Computed, TypeInt
- `delete_time` field is Computed, TypeString
- `backup_region` field is Computed, TypeString
- `restore_time` field is Computed, TypeString
- All fields have appropriate descriptions

**Files Modified**:
- `tencentcloud/services/mongodb/data_source_tc_mongodb_instance_backups.go`

---

### 3. Data Source Logic - Handle Pagination Parameters
**Status**: completed  
**Estimated Time**: 15 minutes  
**Dependencies**: Task 1

Update data source read function to pass pagination parameters to service layer.

**Acceptance Criteria**:
- Read `limit` and `offset` from schema if present
- Pass parameters to service layer via paramMap
- Maintain backward compatibility (parameters are optional)

**Files Modified**:
- `tencentcloud/services/mongodb/data_source_tc_mongodb_instance_backups.go`

---

### 4. Data Source Logic - Map New Response Fields
**Status**: completed  
**Estimated Time**: 15 minutes  
**Dependencies**: Task 2

Map new API response fields to schema in the data source read function.

**Acceptance Criteria**:
- Map `BackId` to `back_id`
- Map `DeleteTime` to `delete_time`
- Map `BackupRegion` to `backup_region`
- Map `RestoreTime` to `restore_time`
- Handle nil values properly

**Files Modified**:
- `tencentcloud/services/mongodb/data_source_tc_mongodb_instance_backups.go`

---

### 5. Service Layer - Support Pagination Parameters
**Status**: completed  
**Estimated Time**: 30 minutes  
**Dependencies**: None

Update service layer to accept and use pagination parameters from caller.

**Acceptance Criteria**:
- `DescribeMongodbInstanceBackupsByFilter` accepts `limit` and `offset` from paramMap
- When provided, use caller's pagination values instead of hardcoded ones
- When not provided, maintain existing behavior (fetch all with internal pagination)
- Properly handle edge cases (limit=0, large offsets)

**Files Modified**:
- `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go`

---

### 6. Documentation - Update Parameter Descriptions
**Status**: completed  
**Estimated Time**: 15 minutes  
**Dependencies**: None

Document new input parameters with clear usage guidance.

**Acceptance Criteria**:
- Document `limit` parameter with valid range and default behavior
- Document `offset` parameter with valid range and default behavior
- Include notes about pagination behavior
- Add example showing pagination usage

**Files Modified**:
- `tencentcloud/services/mongodb/data_source_tc_mongodb_instance_backups.md`

---

### 7. Documentation - Document New Output Fields
**Status**: completed  
**Estimated Time**: 15 minutes  
**Dependencies**: None

Document new output fields with descriptions and use cases.

**Acceptance Criteria**:
- Document `back_id` with purpose and usage
- Document `delete_time` with format and purpose
- Document `backup_region` with cross-region context
- Document `restore_time` with point-in-time recovery context
- Update example output to show new fields

**Files Modified**:
- `tencentcloud/services/mongodb/data_source_tc_mongodb_instance_backups.md`

---

### 8. Website Documentation - Generate Updated Docs
**Status**: completed  
**Estimated Time**: 5 minutes  
**Dependencies**: Tasks 6, 7

Generate website documentation from markdown source.

**Acceptance Criteria**:
- Run `make doc` successfully
- Verify `website/docs/d/mongodb_instance_backups.html.markdown` is updated
- All new parameters and fields are documented

**Commands**:
```bash
make doc
```

**Files Generated**:
- `website/docs/d/mongodb_instance_backups.html.markdown`

---

### 9. Testing - Add Pagination Test Case
**Status**: completed  
**Estimated Time**: 20 minutes  
**Dependencies**: Tasks 1, 3, 5

Add test case to verify pagination parameters work correctly.

**Acceptance Criteria**:
- Test case uses `limit` and `offset` parameters
- Verify correct number of results returned
- Verify no errors with valid parameters
- Test case follows existing test patterns

**Files Modified**:
- `tencentcloud/services/mongodb/data_source_tc_mongodb_instance_backups_test.go`

---

### 10. Testing - Verify New Fields
**Status**: completed  
**Estimated Time**: 15 minutes  
**Dependencies**: Tasks 2, 4

Verify new fields are properly populated in test output.

**Acceptance Criteria**:
- Test queries backups and checks for new fields
- Verify fields are not nil when API provides them
- Handle cases where optional fields may be nil
- Test case follows existing test patterns

**Files Modified**:
- `tencentcloud/services/mongodb/data_source_tc_mongodb_instance_backups_test.go`

---

### 11. Build Validation - Compile and Lint
**Status**: completed  
**Estimated Time**: 10 minutes  
**Dependencies**: All code tasks

Verify code compiles and passes linting.

**Acceptance Criteria**:
- `go build` succeeds without errors
- `make lint` passes without new issues
- No formatting issues

**Commands**:
```bash
go build
make fmt
make lint
```

---

### 12. Integration Testing - Run Acceptance Tests
**Status**: completed  
**Estimated Time**: 15 minutes  
**Dependencies**: All tasks

Run acceptance tests to verify end-to-end functionality.

**Acceptance Criteria**:
- Existing test cases continue to pass
- New test cases pass
- No regressions in behavior

**Commands**:
```bash
TF_ACC=1 go test -v ./tencentcloud/services/mongodb -run TestAccTencentCloudMongodbInstanceBackupsDataSource
```

---

## Summary

**Total Tasks**: 12  
**Estimated Total Time**: ~3 hours  
**Parallelizable Tasks**: 1-2, 6-7 can be done in parallel  
**Critical Path**: 1 → 3 → 5 → 9 → 12

**Deliverables**:
- Enhanced data source with 4 new output fields
- 2 new pagination input parameters
- Updated service layer with pagination support
- Complete documentation with examples
- Test coverage for new functionality
