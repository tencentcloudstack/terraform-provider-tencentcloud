# Spec: MongoDB Backups Data Source Enhancement

**Capability**: `mongodb-backups-enhancement`  
**Related Change**: `enhance-mongodb-backups-datasource`

## Overview
This spec defines the enhanced behavior of the `tencentcloud_mongodb_instance_backups` data source to include all fields from the MongoDB DescribeDBBackups API and support pagination parameters.

---

## MODIFIED Requirements

### Requirement: Data Source Schema Must Include All API Response Fields

**Priority**: High  
**Category**: Schema Definition

The `tencentcloud_mongodb_instance_backups` data source must expose all fields available in the MongoDB `DescribeDBBackups` API response, not just a subset.

#### Scenario: User queries backup details with all available fields

**Given** a MongoDB instance with existing backups  
**When** the user queries the `tencentcloud_mongodb_instance_backups` data source  
**Then** the response must include:
- `instance_id` (existing)
- `backup_type` (existing)
- `backup_name` (existing)
- `backup_desc` (existing)
- `backup_size` (existing)
- `start_time` (existing)
- `end_time` (existing)
- `status` (existing)
- `backup_method` (existing)
- `back_id` (NEW - backup record ID)
- `delete_time` (NEW - backup deletion time)
- `backup_region` (NEW - backup location)
- `restore_time` (NEW - supported restore time)

**Expected Behavior**:
```hcl
data "tencentcloud_mongodb_instance_backups" "example" {
  instance_id = "cmgo-xxxxxxxx"
}

output "backup_details" {
  value = data.tencentcloud_mongodb_instance_backups.example.backup_list[0]
  # Returns all 13 fields including the 4 new ones
}
```

**Implementation Notes**:
- All new fields must be `Computed: true`
- Field types must match SDK types (back_id: Int, others: String)
- Handle nil values gracefully (API may not always return optional fields)

---

### Requirement: Data Source Must Support Pagination Parameters

**Priority**: High  
**Category**: Query Control

The data source must accept `limit` and `offset` parameters to control pagination, matching the API capabilities.

#### Scenario: User queries first page of backups

**Given** a MongoDB instance with 50 backups  
**When** the user sets `limit = 10` and `offset = 0`  
**Then** the data source must return exactly 10 backups (backups 1-10)

**Expected Behavior**:
```hcl
data "tencentcloud_mongodb_instance_backups" "first_page" {
  instance_id = "cmgo-xxxxxxxx"
  limit       = 10
  offset      = 0
}

# Returns backups[0:10]
```

#### Scenario: User queries second page of backups

**Given** a MongoDB instance with 50 backups  
**When** the user sets `limit = 10` and `offset = 10`  
**Then** the data source must return backups 11-20

**Expected Behavior**:
```hcl
data "tencentcloud_mongodb_instance_backups" "second_page" {
  instance_id = "cmgo-xxxxxxxx"
  limit       = 10
  offset      = 10
}

# Returns backups[10:20]
```

#### Scenario: User queries all backups without pagination

**Given** a MongoDB instance with backups  
**When** the user does NOT specify `limit` or `offset`  
**Then** the data source must return ALL backups (existing behavior)

**Expected Behavior**:
```hcl
data "tencentcloud_mongodb_instance_backups" "all" {
  instance_id = "cmgo-xxxxxxxx"
  # No limit/offset specified
}

# Returns all backups (backward compatible)
```

**Implementation Notes**:
- `limit` must be Optional, default behavior = query all
- `offset` must be Optional, default = 0
- `limit` validation: must be > 0 and <= 100 (API limit)
- `offset` validation: must be >= 0

---

## MODIFIED Requirements

### Requirement: Service Layer Must Support Caller-Controlled Pagination

**Priority**: High  
**Category**: Service Layer

The `DescribeMongodbInstanceBackupsByFilter` function MUST accept pagination parameters from the caller instead of always using hardcoded values.

#### Scenario: Service layer uses caller's pagination values

**Given** the data source passes `limit=10` and `offset=5` via paramMap  
**When** the service layer calls the API  
**Then** the request MUST include:
- `Limit = 10`
- `Offset = 5`

**And** MUST NOT loop for additional pages

#### Scenario: Service layer fetches all when no pagination specified

**Given** the data source does NOT pass `limit` or `offset`  
**When** the service layer calls the API  
**Then** the function MUST:
- Use internal pagination (existing behavior: limit=20, offset increments)
- Loop until all records are fetched
- Return complete list

**Implementation Notes**:
- Check paramMap for `limit` and `offset` keys
- If present, use caller's values and make single API call
- If absent, maintain existing loop behavior
- Properly handle last page (fewer results than limit)

---

### Requirement: Field Mapping Must Handle All API Response Fields

**Priority**: High  
**Category**: Data Transformation

The data source read function MUST map all API response fields to the schema output, including new fields.

#### Scenario: API returns all fields including new ones

**Given** an API response with BackId, DeleteTime, BackupRegion, RestoreTime  
**When** the data source processes the response  
**Then** the output map MUST include:
- `back_id` ← BackId
- `delete_time` ← DeleteTime
- `backup_region` ← BackupRegion
- `restore_time` ← RestoreTime

**And** existing field mappings MUST remain unchanged

#### Scenario: API returns nil for optional fields

**Given** an API response where DeleteTime or RestoreTime is nil  
**When** the data source processes the response  
**Then** the output map MUST:
- Include the field key
- Set the value to nil (not omit the key)
- Not cause errors or panics

**Implementation Notes**:
- Use nil-safe checks: `if field != nil { map["key"] = field }`
- Match SDK types (BackId is *int64, others are *string)
- Follow existing patterns in the codebase

---

## ADDED Requirements

### Requirement: Documentation Must Clearly Explain Pagination Behavior

**Priority**: Medium  
**Category**: Documentation

The markdown documentation MUST explain pagination behavior, including when to use it and default behavior.

#### Scenario: User reads documentation for pagination guidance

**Given** a user wants to query a large number of backups  
**When** the user reads the documentation  
**Then** the documentation MUST explain:
1. Default behavior (queries all backups)
2. When to use pagination (large result sets, specific page needs)
3. Valid ranges for `limit` (1-100) and `offset` (>=0)
4. Pagination example code

**Expected Content**:
```markdown
## Argument Reference

* `limit` - (Optional, Int) Number of backups to return per page. 
  Valid range: 1-100. Default: all backups are returned.
* `offset` - (Optional, Int) Offset for pagination. Default: 0.

## Example Usage

### Query with pagination
```hcl
data "tencentcloud_mongodb_instance_backups" "page" {
  instance_id = "cmgo-xxxxxxxx"
  limit       = 20
  offset      = 0
}
```
```

---

### Requirement: Documentation Must Describe New Output Fields

**Priority**: Medium  
**Category**: Documentation

The markdown documentation MUST describe all new output fields with their purpose and format.

#### Scenario: User reads field descriptions

**Given** a user wants to understand new fields  
**When** the user reads the Attributes Reference  
**Then** each new field MUST have:
- Clear description of purpose
- Data type/format
- When the field is populated (always/conditional)

**Expected Content**:
```markdown
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backup_list` - Backup list.
  * `back_id` - Backup record ID.
  * `delete_time` - Scheduled deletion time for the backup.
  * `backup_region` - Region where the backup is stored (for cross-region backups).
  * `restore_time` - Time point supported for backup restore.
  * (... existing fields ...)
```

---

## Validation Criteria

### Schema Validation
- [ ] All 4 new fields are defined in backup_list schema
- [ ] All new fields are `Computed: true`
- [ ] Pagination parameters are `Optional: true`
- [ ] Field types match SDK (back_id: TypeInt, others: TypeString)

### Behavior Validation
- [ ] Query without pagination returns all backups (backward compatible)
- [ ] Query with limit=10 returns max 10 backups
- [ ] Query with offset=5 skips first 5 backups
- [ ] New fields are populated when API provides them
- [ ] Nil fields don't cause errors

### Code Quality
- [ ] No breaking changes to existing schema
- [ ] Follows existing code patterns in mongodb service
- [ ] Proper error handling for nil values
- [ ] Comments explain pagination logic

### Documentation Quality
- [ ] All parameters documented with valid ranges
- [ ] All new fields documented with descriptions
- [ ] Pagination example included
- [ ] Generated website docs include all changes

---

## Related APIs

**Primary API**: `DescribeDBBackups`
- **Endpoint**: mongodb.tencentcloudapi.com
- **Version**: 2019-07-25
- **Documentation**: https://cloud.tencent.com/document/api/240/38574

**Request Parameters**:
- InstanceId (Required, String)
- BackupMethod (Optional, Integer): 0=logical, 1=physical, 3=snapshot
- Limit (Optional, Integer): max 100
- Offset (Optional, Integer): min 0

**Response Structure**:
```go
type BackupInfo struct {
    InstanceId   *string  // Existing
    BackupType   *uint64  // Existing
    BackupName   *string  // Existing
    BackupDesc   *string  // Existing
    BackupSize   *uint64  // Existing
    StartTime    *string  // Existing
    EndTime      *string  // Existing
    Status       *uint64  // Existing
    BackupMethod *uint64  // Existing
    BackId       *int64   // NEW
    DeleteTime   *string  // NEW
    BackupRegion *string  // NEW
    RestoreTime  *string  // NEW
}
```

---

## Dependencies

**SDK Package**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725`
- No version change required
- All fields already present in SDK models

**Related Files**:
- `data_source_tc_mongodb_instance_backups.go` (PRIMARY)
- `service_tencentcloud_mongodb.go` (service layer)
- `data_source_tc_mongodb_instance_backups.md` (docs)

---

## Success Metrics

1. **Completeness**: All 13 API fields exposed (9 existing + 4 new)
2. **Backward Compatibility**: Existing configs work without changes
3. **Pagination Control**: Users can limit query size when needed
4. **Documentation**: Clear guidance on all parameters and fields
5. **Test Coverage**: Tests verify new functionality works correctly
