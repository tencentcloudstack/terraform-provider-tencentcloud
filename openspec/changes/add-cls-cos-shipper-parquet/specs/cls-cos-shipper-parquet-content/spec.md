# Spec: CLS COS Shipper Parquet Content Format Support

**Capability**: cls-cos-shipper-parquet-content  
**Related Change**: add-cls-cos-shipper-parquet  
**Status**: Draft

---

## Overview

This spec defines the requirements for adding Parquet format support to the `tencentcloud_cls_cos_shipper` resource's `content` configuration. Parquet is a columnar storage format optimized for analytics workloads, and the TencentCloud CLS API already supports it in shipper configurations.

---

## ADDED Requirements

### Requirement 1: Parquet Content Format Configuration

**ID**: CLS-SHIPPER-PARQUET-001  
**Priority**: High  
**Type**: Feature

Users must be able to configure the content format as "parquet" and define Parquet column schemas with data type specifications.

#### Scenario: Configure basic Parquet format with string and integer columns

**Given**: A user wants to ship CLS logs to COS in Parquet format  
**When**: The user configures content with format "parquet" and defines columns  
**Then**: 
- The shipper is created with Parquet content format
- The Parquet schema includes the defined columns
- Each column has a name and data type
- The API accepts the configuration successfully

**Example**:
```hcl
resource "tencentcloud_cls_cos_shipper" "example" {
  bucket       = "my-bucket"
  topic_id     = "topic-123"
  shipper_name = "parquet-shipper"
  prefix       = "logs/"

  content {
    format = "parquet"

    parquet {
      parquet_key_info {
        key_name = "timestamp"
        key_type = "int64"
      }
      parquet_key_info {
        key_name = "message"
        key_type = "string"
      }
    }
  }
}
```

**Acceptance Criteria**:
- ✅ `content.format` can be set to "parquet"
- ✅ `content.parquet` block is available
- ✅ Multiple `parquet_key_info` blocks can be defined
- ✅ Each `parquet_key_info` has required `key_name` and `key_type` fields
- ✅ Terraform plan succeeds without errors
- ✅ Provider sends correct ParquetInfo structure to API

---

#### Scenario: Support all six Parquet data types

**Given**: A user needs to define columns with various data types  
**When**: The user specifies key_type as one of the supported types  
**Then**: 
- All six data types are accepted: string, boolean, int32, int64, float, double
- Invalid types are rejected at plan time
- Type validation provides clear error messages

**Example**:
```hcl
content {
  format = "parquet"

  parquet {
    parquet_key_info {
      key_name = "app_name"
      key_type = "string"
    }
    parquet_key_info {
      key_name = "is_error"
      key_type = "boolean"
    }
    parquet_key_info {
      key_name = "status_code"
      key_type = "int32"
    }
    parquet_key_info {
      key_name = "trace_id"
      key_type = "int64"
    }
    parquet_key_info {
      key_name = "cpu_usage"
      key_type = "float"
    }
    parquet_key_info {
      key_name = "latency"
      key_type = "double"
    }
  }
}
```

**Acceptance Criteria**:
- ✅ `string` type is accepted and works
- ✅ `boolean` type is accepted and works
- ✅ `int32` type is accepted and works
- ✅ `int64` type is accepted and works
- ✅ `float` type is accepted and works
- ✅ `double` type is accepted and works
- ✅ Invalid type like "varchar" fails validation
- ✅ Validation error message mentions valid types

---

#### Scenario: Validate minimum column requirement

**Given**: A user configures Parquet format  
**When**: The user tries to define parquet block with no columns  
**Then**: 
- Validation fails with clear error
- Error indicates at least one parquet_key_info is required

**Example (should fail)**:
```hcl
content {
  format = "parquet"
  
  parquet {
    # No parquet_key_info defined
  }
}
```

**Acceptance Criteria**:
- ✅ Empty parquet_key_info array fails validation
- ✅ Error message indicates minimum 1 column required
- ✅ Validation happens at plan time

---

### Requirement 2: Create Operation Integration

**ID**: CLS-SHIPPER-PARQUET-002  
**Priority**: High  
**Type**: Feature

The Create operation must correctly convert Terraform schema to API request structure for Parquet content.

#### Scenario: Create shipper with Parquet content successfully

**Given**: Valid Parquet content configuration  
**When**: User runs `terraform apply` to create the shipper  
**Then**: 
- Provider reads parquet block from schema
- Provider constructs ParquetInfo structure
- Provider populates ParquetKeyInfo array correctly
- API request includes Content.Parquet field
- Shipper is created successfully

**API Mapping**:
```
Schema                              → API Request
------                              -----------
content.format = "parquet"          → request.Content.Format = "parquet"
content.parquet.parquet_key_info    → request.Content.Parquet.ParquetKeyInfo[]
  .key_name = "timestamp"           →   .KeyName = "timestamp"
  .key_type = "int64"               →   .KeyType = "int64"
```

**Acceptance Criteria**:
- ✅ Parquet block is correctly parsed from schema
- ✅ ParquetKeyInfo array is built with correct length
- ✅ Each ParquetKeyInfo has KeyName and KeyType set
- ✅ String pointers are created with helper.String()
- ✅ Content.Parquet is assigned to request
- ✅ API accepts the request
- ✅ Shipper ID is returned and set in state

---

#### Scenario: Create with empty parquet_key_info list

**Given**: User configures parquet with empty parquet_key_info  
**When**: Terraform attempts to create the shipper  
**Then**: 
- Validation catches the error at plan time
- Clear error message is displayed
- No API call is made

**Acceptance Criteria**:
- ✅ Schema validation with MinItems: 1 catches this
- ✅ Error is shown during plan, not apply
- ✅ Error message is user-friendly

---

### Requirement 3: Read Operation Integration

**ID**: CLS-SHIPPER-PARQUET-003  
**Priority**: High  
**Type**: Feature

The Read operation must correctly parse API response and populate Terraform state with Parquet content configuration.

#### Scenario: Read shipper with Parquet content

**Given**: An existing shipper with Parquet content format  
**When**: Provider reads the shipper from API  
**Then**: 
- Provider retrieves Content.Parquet from response
- Provider converts ParquetKeyInfo array to schema format
- Terraform state includes parquet block with all columns
- State matches the original configuration

**API Mapping**:
```
API Response                        → Terraform State
------------                        ---------------
response.Content.Format = "parquet" → content[0].format = "parquet"
response.Content.Parquet.ParquetKeyInfo[] → content[0].parquet[0].parquet_key_info[]
  .KeyName = "timestamp"            →   [0].key_name = "timestamp"
  .KeyType = "int64"                →   [0].key_type = "int64"
```

**Acceptance Criteria**:
- ✅ Content.Parquet is checked for nil before access
- ✅ ParquetKeyInfo array is converted to interface slice
- ✅ Each element is converted to map[string]interface{}
- ✅ key_name and key_type are correctly populated
- ✅ parquet block is added to content map
- ✅ State is set via d.Set("content", ...)
- ✅ No data loss in round-trip (create → read)

---

#### Scenario: Read shipper with null Parquet field

**Given**: API returns Content with nil Parquet field  
**When**: Provider reads the response  
**Then**: 
- Provider handles nil pointer safely
- No panic occurs
- Parquet block is not added to state (or added as empty)

**Acceptance Criteria**:
- ✅ Nil check prevents panic: `if shipper.Content.Parquet != nil`
- ✅ Read operation completes successfully
- ✅ Other content fields (csv, json) are still processed

---

#### Scenario: Import existing Parquet shipper

**Given**: A Parquet shipper created via console or API  
**When**: User imports the shipper into Terraform: `terraform import tencentcloud_cls_cos_shipper.example shipper-id`  
**Then**: 
- Provider reads the shipper configuration
- Parquet columns are populated in state
- Subsequent `terraform plan` shows no changes

**Acceptance Criteria**:
- ✅ Import command succeeds
- ✅ State includes complete parquet configuration
- ✅ `terraform plan` shows 0 changes
- ✅ All column names and types are preserved

---

### Requirement 4: Update Operation Integration

**ID**: CLS-SHIPPER-PARQUET-004  
**Priority**: High  
**Type**: Feature

The Update operation must support modifying Parquet content configuration when content block changes.

#### Scenario: Update Parquet column definitions

**Given**: An existing shipper with Parquet format  
**When**: User modifies the parquet_key_info blocks (add/remove/change columns)  
**Then**: 
- Terraform detects the change in content
- Provider sends ModifyShipper request with updated Parquet schema
- Shipper is updated successfully
- State reflects new configuration

**Example**:
```hcl
# Before
content {
  format = "parquet"
  parquet {
    parquet_key_info {
      key_name = "timestamp"
      key_type = "int64"
    }
  }
}

# After (add column)
content {
  format = "parquet"
  parquet {
    parquet_key_info {
      key_name = "timestamp"
      key_type = "int64"
    }
    parquet_key_info {
      key_name = "level"
      key_type = "string"
    }
  }
}
```

**Acceptance Criteria**:
- ✅ `d.HasChange("content")` detects the modification
- ✅ Update logic reads new parquet configuration
- ✅ ParquetInfo structure is rebuilt with new columns
- ✅ API request includes updated Content.Parquet
- ✅ Update succeeds
- ✅ Subsequent read shows new configuration

---

#### Scenario: Change format from JSON to Parquet

**Given**: A shipper with JSON format  
**When**: User changes format to "parquet" and adds parquet block  
**Then**: 
- Update operation handles format change
- JSON configuration is removed
- Parquet configuration is applied
- Update succeeds

**Example**:
```hcl
# Before
content {
  format = "json"
  json {
    enable_tag = true
    meta_fields = ["__TIMESTAMP__"]
  }
}

# After
content {
  format = "parquet"
  parquet {
    parquet_key_info {
      key_name = "timestamp"
      key_type = "int64"
    }
  }
}
```

**Acceptance Criteria**:
- ✅ Format change is detected
- ✅ JSON configuration is cleared from request
- ✅ Parquet configuration is applied
- ✅ Update succeeds without errors
- ✅ State reflects new format

---

### Requirement 5: Schema Validation

**ID**: CLS-SHIPPER-PARQUET-005  
**Priority**: Medium  
**Type**: Quality

Schema validation must catch configuration errors early at plan time with clear error messages.

#### Scenario: Reject invalid data types

**Given**: User specifies an unsupported key_type  
**When**: User runs `terraform plan`  
**Then**: 
- Validation fails immediately
- Error message lists valid types
- No API call is made

**Example (should fail)**:
```hcl
content {
  format = "parquet"
  parquet {
    parquet_key_info {
      key_name = "id"
      key_type = "uuid"  # Invalid type
    }
  }
}
```

**Expected Error**:
```
Error: expected key_type to be one of [string boolean int32 int64 float double], got uuid
```

**Acceptance Criteria**:
- ✅ ValidateFunc catches invalid types
- ✅ Error is shown during plan phase
- ✅ Error message lists all valid types
- ✅ No API request is sent

---

#### Scenario: Enforce required fields

**Given**: User omits required fields in parquet_key_info  
**When**: User runs `terraform plan`  
**Then**: 
- Schema validation fails
- Error indicates which field is missing
- Clear remediation guidance

**Example (should fail)**:
```hcl
content {
  format = "parquet"
  parquet {
    parquet_key_info {
      key_name = "timestamp"
      # Missing key_type
    }
  }
}
```

**Acceptance Criteria**:
- ✅ Missing key_type triggers validation error
- ✅ Error identifies the missing required field
- ✅ Validation happens at plan time

---

### Requirement 6: Backward Compatibility

**ID**: CLS-SHIPPER-PARQUET-006  
**Priority**: Critical  
**Type**: Quality

The addition of Parquet support must not break existing JSON and CSV configurations.

#### Scenario: Existing JSON shipper remains unchanged

**Given**: A shipper configured with JSON format before Parquet support  
**When**: User upgrades provider to version with Parquet support  
**Then**: 
- JSON configuration continues to work
- No changes detected in `terraform plan`
- No migration required

**Acceptance Criteria**:
- ✅ JSON shippers show 0 changes after upgrade
- ✅ JSON block structure unchanged
- ✅ JSON shipper CRUD operations work as before

---

#### Scenario: Existing CSV shipper remains unchanged

**Given**: A shipper configured with CSV format before Parquet support  
**When**: User upgrades provider to version with Parquet support  
**Then**: 
- CSV configuration continues to work
- No changes detected in `terraform plan`
- No migration required

**Acceptance Criteria**:
- ✅ CSV shippers show 0 changes after upgrade
- ✅ CSV block structure unchanged
- ✅ CSV shipper CRUD operations work as before

---

#### Scenario: Mutual exclusivity maintained

**Given**: Content format can be json, csv, or parquet  
**When**: User configures one format  
**Then**: 
- Only one format block is populated
- API receives only the configured format
- No conflicts between formats

**Acceptance Criteria**:
- ✅ Format string determines which block is used
- ✅ Only one of csv/json/parquet is populated in API request
- ✅ Validation prevents multiple format blocks if needed

---

### Requirement 7: Documentation

**ID**: CLS-SHIPPER-PARQUET-007  
**Priority**: High  
**Type**: Quality

Complete and accurate documentation must be provided for Parquet format usage.

#### Scenario: Usage examples cover common use cases

**Given**: A user wants to configure Parquet format  
**When**: User reads the resource documentation  
**Then**: 
- Examples show basic Parquet configuration
- Examples demonstrate all 6 data types
- Examples show best practices (e.g., with compression)

**Acceptance Criteria**:
- ✅ At least one basic Parquet example
- ✅ Example showing all data types
- ✅ Example with compression configuration
- ✅ Examples are syntactically valid HCL

---

#### Scenario: Data type reference is clear and complete

**Given**: User needs to choose appropriate data types  
**When**: User consults documentation  
**Then**: 
- All 6 data types are documented
- Description explains when to use each type
- Examples show typical use cases for each type

**Acceptance Criteria**:
- ✅ Table of data types with descriptions
- ✅ Guidance on type selection
- ✅ Performance implications noted (if any)
- ✅ Storage efficiency tips

---

### Requirement 8: Error Handling

**ID**: CLS-SHIPPER-PARQUET-008  
**Priority**: Medium  
**Type**: Quality

Error handling must be robust and provide helpful feedback for common failure scenarios.

#### Scenario: Handle API rejection of Parquet configuration

**Given**: User configures Parquet with schema that doesn't match log structure  
**When**: API rejects the configuration  
**Then**: 
- Provider captures API error
- Error message is passed to user
- User can identify the issue

**Acceptance Criteria**:
- ✅ API errors are caught and returned
- ✅ Error message includes API response details
- ✅ Request body is logged for debugging

---

#### Scenario: Safe handling of nil pointers in responses

**Given**: API returns incomplete or null Parquet data  
**When**: Provider processes the response  
**Then**: 
- No panic occurs
- Nil checks prevent crashes
- Operation completes gracefully

**Acceptance Criteria**:
- ✅ All pointer fields checked before dereference
- ✅ `if shipper.Content != nil` guard
- ✅ `if shipper.Content.Parquet != nil` guard
- ✅ No runtime panics in Read operation

---

## API Contract

### CreateShipper API

**Endpoint**: `CreateShipper`

**Request Structure**:
```json
{
  "TopicId": "topic-123",
  "Bucket": "bucket-123",
  "Prefix": "logs/",
  "ShipperName": "test-shipper",
  "Content": {
    "Format": "parquet",
    "Parquet": {
      "ParquetKeyInfo": [
        {
          "KeyName": "timestamp",
          "KeyType": "int64"
        },
        {
          "KeyName": "message",
          "KeyType": "string"
        }
      ]
    }
  }
}
```

### ModifyShipper API

**Endpoint**: `ModifyShipper`

**Request Structure**: Same as CreateShipper, with ShipperId included

### DescribeShippers API

**Endpoint**: `DescribeShippers`

**Response Structure**:
```json
{
  "Shippers": [
    {
      "ShipperId": "shipper-123",
      "TopicId": "topic-123",
      "Content": {
        "Format": "parquet",
        "Parquet": {
          "ParquetKeyInfo": [
            {
              "KeyName": "timestamp",
              "KeyType": "int64"
            }
          ]
        }
      }
    }
  ]
}
```

---

## Non-Functional Requirements

### Performance
- **NFR-1**: Parquet configuration parsing must not add > 50ms overhead
- **NFR-2**: Read operation with Parquet data must complete within existing timeout

### Code Quality
- **NFR-3**: Code must pass golangci-lint with no new errors
- **NFR-4**: Code must follow existing patterns (csv/json handling)
- **NFR-5**: All pointer dereferences must have nil checks

### Maintainability
- **NFR-6**: Parquet handling code must mirror csv/json structure
- **NFR-7**: Variable names must be consistent with SDK field names
- **NFR-8**: Comments must explain non-obvious logic

---

## Success Metrics

1. **Feature Completeness**: 100% of Parquet API capabilities exposed
2. **Backward Compatibility**: 0 breaking changes to existing JSON/CSV configs
3. **Validation Coverage**: 100% of invalid inputs caught at plan time
4. **Documentation Quality**: All data types documented with examples
5. **Test Coverage**: All CRUD operations tested for Parquet format

---

## Related Specs

- **CLS COS Shipper Base**: Existing resource spec (if exists)
- **CLS Content Formats**: Specification for JSON and CSV formats

---

## Change History

- **2026-02-11**: Initial spec created for Parquet support addition

---

**END OF SPEC**
