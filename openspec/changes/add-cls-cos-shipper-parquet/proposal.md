# Proposal: Add Parquet Format Support to CLS COS Shipper Content

## Metadata
- **Change ID**: `add-cls-cos-shipper-parquet`
- **Status**: Proposal
- **Created**: 2026-02-11
- **Author**: AI Assistant
- **Type**: Enhancement
- **Estimated Effort**: 2-3 hours

---

## Problem Statement

The `tencentcloud_cls_cos_shipper` resource currently supports only two content formats for shipping logs to COS:
1. **JSON format** - with metadata fields configuration
2. **CSV format** - with delimiter, keys, and field options

However, the TencentCloud CLS API already supports a third format - **Parquet**, which is:
- A columnar storage format optimized for analytics workloads
- More efficient in terms of storage and query performance
- Widely used in big data ecosystems (Spark, Hive, Presto)

### Current Limitations

Users who want to ship CLS logs in Parquet format to COS for:
- Integration with data warehouses (e.g., DLC, EMR)
- Cost-effective long-term storage with compression
- Efficient analytical queries with columnar access patterns

**Cannot** configure this through Terraform, even though:
- The TencentCloud SDK already includes `ParquetInfo` structure
- The CLS API supports Parquet in CreateShipper/ModifyShipper/DescribeShippers
- The `ContentInfo` struct has a `Parquet` field

### User Impact

- **Incomplete feature coverage**: Terraform doesn't expose all available API capabilities
- **Manual configuration required**: Users must use console or API directly
- **Infrastructure drift**: Parquet shippers created outside Terraform cannot be fully managed

---

## Proposed Solution

Add a new optional `parquet` nested block to the `content` schema field in `tencentcloud_cls_cos_shipper` resource.

### Schema Addition

```hcl
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
    parquet_key_info {
      key_name = "message"
      key_type = "string"
    }
  }
}
```

### Field Specification

#### `content.parquet` (Optional, Block)
- **Type**: List with MaxItems: 1
- **Description**: Parquet format content description
- **Required when**: `content.format = "parquet"`

#### `content.parquet.parquet_key_info` (Required, Block List)
- **Type**: List (repeatable block)
- **MinItems**: 1
- **Description**: Array of Parquet column definitions

#### `content.parquet.parquet_key_info.key_name` (Required, String)
- **Type**: String
- **Description**: Column name in the Parquet file

#### `content.parquet.parquet_key_info.key_type` (Required, String)
- **Description**: Data type of the column
- **Valid Values**: `string`, `boolean`, `int32`, `int64`, `float`, `double`
- **Validation**: ValidateAllowedStringValue

---

## Technical Design

### API Mapping

| Operation | API Endpoint | Field Mapping |
|-----------|-------------|---------------|
| **Create** | CreateShipper | `content.parquet` → `Content.Parquet` |
| **Read** | DescribeShippers | `Content.Parquet` → `content.parquet` |
| **Update** | ModifyShipper | `content.parquet` → `Content.Parquet` |

### SDK Structures (Already Available)

```go
// From tencentcloud-sdk-go/tencentcloud/cls/v20201016/models.go
type ContentInfo struct {
    Format  *string      `json:"Format,omitnil,omitempty" name:"Format"`
    Csv     *CsvInfo     `json:"Csv,omitnil,omitempty" name:"Csv"`
    Json    *JsonInfo    `json:"Json,omitnil,omitempty" name:"Json"`
    Parquet *ParquetInfo `json:"Parquet,omitnil,omitempty" name:"Parquet"` // ✅ Already exists
}

type ParquetInfo struct {
    ParquetKeyInfo []*ParquetKeyInfo `json:"ParquetKeyInfo,omitnil,omitempty" name:"ParquetKeyInfo"`
}

type ParquetKeyInfo struct {
    KeyName *string `json:"KeyName,omitnil,omitempty" name:"KeyName"`
    KeyType *string `json:"KeyType,omitnil,omitempty" name:"KeyType"`
}
```

### Implementation Files

1. **Resource File**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go`
   - Add `parquet` block to `content` schema (line ~107)
   - Handle `parquet` in Create operation (line ~272-314)
   - Handle `parquet` in Read operation (line ~418-440)
   - Handle `parquet` in Update operation (line ~550-594)

2. **Documentation**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.md`
   - Add Parquet usage example
   - Document supported data types
   - Add best practices for column definitions

3. **Website Documentation**: Auto-generated via `make doc`

---

## Implementation Plan

### Phase 1: Schema Definition (30 min)
1. Add `parquet` nested block to `content` schema
2. Add `parquet_key_info` repeatable block with `key_name` and `key_type`
3. Add validation for `key_type` (6 supported types)

### Phase 2: Create Operation (20 min)
1. Read `parquet` block from schema
2. Build `ParquetInfo` structure
3. Populate `ParquetKeyInfo` array from `parquet_key_info` blocks
4. Assign to `request.Content.Parquet`

### Phase 3: Read Operation (20 min)
1. Check if `shipper.Content.Parquet` exists
2. Build `parquet` map from `ParquetInfo`
3. Build `parquet_key_info` array from `ParquetKeyInfo` slice
4. Set state via `d.Set("content", ...)`

### Phase 4: Update Operation (20 min)
1. Detect change with `d.HasChange("content")`
2. Read `parquet` block if format is "parquet"
3. Build and assign `ParquetInfo` to `request.Content.Parquet`

### Phase 5: Testing (30 min)
1. Add test case: Create shipper with Parquet format
2. Add test case: Update from JSON to Parquet
3. Add test case: Modify Parquet column definitions
4. Verify Import works correctly

### Phase 6: Documentation (20 min)
1. Add Parquet example to resource documentation
2. Document all 6 supported data types
3. Add notes about column name requirements
4. Generate website docs via `make doc`

### Phase 7: Validation (20 min)
1. Run `gofmt` and `goimports`
2. Compile provider: `go build`
3. Check linter: `make lint` (if available)
4. Verify documentation generation

---

## Design Considerations

### 1. Schema Structure Choice

**Decision**: Use nested block `parquet_key_info` (repeatable) instead of List[Map]

**Rationale**:
- **Consistent with existing pattern**: `csv` and `json` blocks follow this pattern
- **Better type safety**: Block schema provides explicit field validation
- **Better UX**: HCL block syntax is more readable than list of maps
- **Terraform best practice**: Blocks for structural data, attributes for simple values

**Example**:
```hcl
# ✅ Chosen approach - cleaner syntax
parquet {
  parquet_key_info {
    key_name = "id"
    key_type = "int64"
  }
  parquet_key_info {
    key_name = "name"
    key_type = "string"
  }
}

# ❌ Alternative - less readable
parquet {
  parquet_key_info = [
    {
      key_name = "id"
      key_type = "int64"
    },
    {
      key_name = "name"
      key_type = "string"
    }
  ]
}
```

### 2. Validation Strategy

**key_type Validation**:
- Use `tccommon.ValidateAllowedStringValue([]string{"string", "boolean", "int32", "int64", "float", "double"})`
- Fail fast at plan time (not at apply time)
- Clear error messages for invalid types

**No key_name Validation**:
- Allow any valid identifier (alphanumeric + underscore)
- Let API validate against actual log field names
- API knows the actual log schema, Terraform doesn't

### 3. Mutual Exclusivity

**Current behavior**: `content` can have only one of `csv`, `json`, or `parquet`

**Implementation**: 
- No explicit ConflictsWith needed
- Natural exclusivity: format string determines which block is used
- Validation happens at API level if multiple formats provided

### 4. Backward Compatibility

**Impact**: ✅ **Fully backward compatible**
- New optional field
- Existing `json` and `csv` configurations unaffected
- No breaking changes to schema structure

---

## Example Usage

### Example 1: Basic Parquet Format

```hcl
resource "tencentcloud_cls_cos_shipper" "parquet_example" {
  bucket       = tencentcloud_cos_bucket.example.id
  topic_id     = tencentcloud_cls_topic.example.id
  shipper_name = "parquet-shipper"
  prefix       = "logs/parquet/"
  interval     = 300
  max_size     = 256

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
      parquet_key_info {
        key_name = "message"
        key_type = "string"
      }
      parquet_key_info {
        key_name = "user_id"
        key_type = "int64"
      }
      parquet_key_info {
        key_name = "success"
        key_type = "boolean"
      }
      parquet_key_info {
        key_name = "duration_ms"
        key_type = "float"
      }
    }
  }
}
```

### Example 2: All Data Types

```hcl
resource "tencentcloud_cls_cos_shipper" "all_types" {
  bucket       = tencentcloud_cos_bucket.example.id
  topic_id     = tencentcloud_cls_topic.example.id
  shipper_name = "all-types-shipper"
  prefix       = "logs/types/"

  content {
    format = "parquet"

    parquet {
      # String type
      parquet_key_info {
        key_name = "app_name"
        key_type = "string"
      }
      
      # Boolean type
      parquet_key_info {
        key_name = "is_error"
        key_type = "boolean"
      }
      
      # Int32 type
      parquet_key_info {
        key_name = "status_code"
        key_type = "int32"
      }
      
      # Int64 type
      parquet_key_info {
        key_name = "trace_id"
        key_type = "int64"
      }
      
      # Float type
      parquet_key_info {
        key_name = "cpu_usage"
        key_type = "float"
      }
      
      # Double type
      parquet_key_info {
        key_name = "response_time"
        key_type = "double"
      }
    }
  }
}
```

### Example 3: With Compression (Recommended)

```hcl
resource "tencentcloud_cls_cos_shipper" "compressed_parquet" {
  bucket       = tencentcloud_cos_bucket.example.id
  topic_id     = tencentcloud_cls_topic.example.id
  shipper_name = "compressed-parquet"
  prefix       = "logs/compressed/"

  compress {
    format = "gzip"  # Parquet + gzip = excellent compression
  }

  content {
    format = "parquet"

    parquet {
      parquet_key_info {
        key_name = "event_time"
        key_type = "int64"
      }
      parquet_key_info {
        key_name = "event_name"
        key_type = "string"
      }
      parquet_key_info {
        key_name = "value"
        key_type = "double"
      }
    }
  }
}
```

---

## Impact Analysis

### Breaking Changes
**None** - This is a purely additive change.

### Backward Compatibility
✅ **Fully backward compatible**:
- New optional field
- Existing JSON/CSV shippers continue to work
- No changes to existing schema fields
- Import of existing Parquet shippers (created via console) will now work

### API Dependencies
- ✅ SDK structures already available
- ✅ API supports Parquet in all CRUD operations
- ✅ No version constraints

### Performance Impact
- **Minimal**: Same API call pattern as CSV/JSON
- **Read operation**: No extra API calls needed
- **User benefit**: Parquet files are more efficient for storage and queries

---

## Risks and Mitigation

### Risk 1: Type Mismatch
**Risk**: User specifies `int64` but log field contains strings  
**Impact**: Medium - API will reject invalid type mappings  
**Mitigation**: 
- Clear documentation with type descriptions
- Validation at Terraform plan time for valid type names
- API will provide detailed error messages

### Risk 2: Missing Fields
**Risk**: User defines Parquet columns that don't exist in logs  
**Impact**: Low - API handles gracefully (likely skips or nulls)  
**Mitigation**: 
- Document that column names must match log field names
- Recommend testing with actual log data

### Risk 3: Parquet-specific API Limitations
**Risk**: Undocumented API constraints for Parquet format  
**Impact**: Low - Parquet is a standard feature in CLS  
**Mitigation**: 
- Rely on API validation
- Add detailed error context if API rejects

---

## Success Criteria

### Functional Requirements
- ✅ Users can configure Parquet format with column definitions
- ✅ All 6 data types are supported and validated
- ✅ Create/Read/Update/Delete operations work correctly
- ✅ Import of existing Parquet shippers works
- ✅ Validation catches invalid type names at plan time

### Non-Functional Requirements
- ✅ Code follows project conventions (naming, error handling)
- ✅ No linter errors or warnings
- ✅ Documentation is complete and includes examples
- ✅ Backward compatible with existing configurations

---

## Alternatives Considered

### Alternative 1: List of Maps for parquet_key_info
**Approach**: Use `List[Map]` instead of repeatable blocks  
**Pros**: Slightly less verbose in HCL  
**Cons**: Less type-safe, inconsistent with csv/json patterns, harder to validate  
**Verdict**: ❌ Rejected - Block pattern is Terraform best practice

### Alternative 2: Single Map with Name→Type Mapping
**Approach**: 
```hcl
parquet {
  key_types = {
    "timestamp" = "int64"
    "level"     = "string"
  }
}
```
**Pros**: Very concise  
**Cons**: Loses ordering (maps are unordered), inconsistent with SDK structure  
**Verdict**: ❌ Rejected - Order may matter for columnar formats

### Alternative 3: Separate Resource for Content Configuration
**Approach**: Split content configuration into separate resource  
**Pros**: Modularity  
**Cons**: Over-engineering, breaks existing pattern, complex for users  
**Verdict**: ❌ Rejected - Content is integral to shipper configuration

---

## Timeline

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| 1. Schema Definition | 3 tasks | 30 min |
| 2. Create Operation | 4 tasks | 20 min |
| 3. Read Operation | 4 tasks | 20 min |
| 4. Update Operation | 3 tasks | 20 min |
| 5. Testing | 4 tasks | 30 min |
| 6. Documentation | 4 tasks | 20 min |
| 7. Validation | 4 tasks | 20 min |
| **Total** | **26 tasks** | **2.5 hours** |

---

## Related Documentation

- **TencentCloud CLS API**:
  - CreateShipper: https://cloud.tencent.com/document/api/614/31574
  - ModifyShipper: https://cloud.tencent.com/document/api/614/31575
  - DescribeShippers: https://cloud.tencent.com/document/api/614/31576

- **Parquet Format**:
  - Apache Parquet: https://parquet.apache.org/
  - Data Types: https://github.com/apache/parquet-format/blob/master/LogicalTypes.md

- **SDK Source**:
  - Models: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016/models.go`
  - ContentInfo (line 1571), ParquetInfo (line 18630), ParquetKeyInfo (line 18635)

- **Existing Implementation**:
  - Resource: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go`
  - Documentation: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.md`

---

## Approval

This proposal is ready for review and implementation.

**Next Steps**:
1. Review and approve this proposal
2. Create `tasks.md` with detailed implementation tasks
3. Create spec delta in `specs/cls-cos-shipper-parquet-content/spec.md`
4. Run `openspec validate add-cls-cos-shipper-parquet --strict`
5. Begin implementation after validation passes

---

**END OF PROPOSAL**
