# Add Parquet Format Support to CLS COS Shipper

**Change ID**: `add-cls-cos-shipper-parquet`  
**Status**: Proposal  
**Type**: Enhancement  
**Estimated Effort**: 2.5 hours

---

## Quick Summary

Add support for **Parquet format** to the `content` configuration in the `tencentcloud_cls_cos_shipper` resource. This enables users to ship CLS logs to COS in columnar Parquet format, which is optimized for analytics and data warehouse integration.

---

## Problem

Currently, users can only ship logs in **JSON** or **CSV** formats. The TencentCloud CLS API already supports **Parquet** format, but the Terraform provider doesn't expose this capability.

**Impact**:
- Users cannot leverage Parquet's superior compression and query performance
- Integration with data warehouses (DLC, EMR) requires manual configuration
- Terraform state cannot fully manage Parquet shippers

---

## Solution

Add a new `parquet` block to the `content` schema:

```hcl
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
```

### Supported Data Types

| Type | Description |
|------|-------------|
| `string` | Text data (messages, IDs) |
| `boolean` | True/false values |
| `int32` | 32-bit integers (status codes) |
| `int64` | 64-bit integers (timestamps, IDs) |
| `float` | 32-bit floating point (percentages) |
| `double` | 64-bit floating point (high-precision metrics) |

---

## Implementation Overview

### Files to Modify

1. **`resource_tc_cls_cos_shipper.go`** (main changes)
   - Add `parquet` block to schema (~line 107)
   - Handle `parquet` in Create (~line 272)
   - Handle `parquet` in Read (~line 418)
   - Handle `parquet` in Update (~line 550)

2. **`resource_tc_cls_cos_shipper.md`** (documentation)
   - Add Parquet usage examples
   - Document data types
   - Add best practices

3. **Auto-generated**: `website/docs/r/cls_cos_shipper.html.markdown`

### SDK Structures (Already Available)

```go
// From TencentCloud SDK
type ContentInfo struct {
    Format  *string      // "json", "csv", "parquet"
    Csv     *CsvInfo
    Json    *JsonInfo
    Parquet *ParquetInfo  // ✅ Already exists!
}

type ParquetInfo struct {
    ParquetKeyInfo []*ParquetKeyInfo
}

type ParquetKeyInfo struct {
    KeyName *string  // Column name
    KeyType *string  // Data type (string, boolean, int32, int64, float, double)
}
```

---

## Key Features

### 1. Full CRUD Support
- ✅ **Create**: Configure Parquet format during shipper creation
- ✅ **Read**: Import existing Parquet shippers
- ✅ **Update**: Modify column definitions
- ✅ **Delete**: Standard deletion (format-agnostic)

### 2. Validation
- ✅ **Type validation**: Only 6 supported types allowed
- ✅ **Required fields**: key_name and key_type must be provided
- ✅ **Minimum columns**: At least one column required
- ✅ **Plan-time errors**: Invalid configs fail before API call

### 3. Backward Compatibility
- ✅ **Additive change**: New optional field
- ✅ **No breaking changes**: Existing JSON/CSV configs unaffected
- ✅ **Graceful upgrade**: Provider upgrade requires no migration

---

## Usage Examples

### Example 1: Basic Parquet Configuration

```hcl
resource "tencentcloud_cls_cos_shipper" "example" {
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
    }
  }
}
```

### Example 2: All Data Types

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

### Example 3: With Compression (Recommended)

```hcl
resource "tencentcloud_cls_cos_shipper" "optimized" {
  bucket       = tencentcloud_cos_bucket.example.id
  topic_id     = tencentcloud_cls_topic.example.id
  shipper_name = "optimized-shipper"
  prefix       = "logs/"

  compress {
    format = "gzip"  # Parquet + gzip = excellent compression!
  }

  content {
    format = "parquet"

    parquet {
      parquet_key_info {
        key_name = "timestamp"
        key_type = "int64"
      }
      parquet_key_info {
        key_name = "event"
        key_type = "string"
      }
    }
  }
}
```

---

## Implementation Plan

### Phase 1: Schema (30 min)
- Add `parquet` nested block to `content` schema
- Add `parquet_key_info` repeatable block
- Add validation for `key_type`

### Phase 2-4: CRUD Operations (60 min)
- **Create**: Build ParquetInfo from schema
- **Read**: Parse ParquetInfo to state
- **Update**: Rebuild ParquetInfo on change

### Phase 5: Testing (30 min)
- Test basic Parquet configuration
- Test all 6 data types
- Test validation (invalid types)

### Phase 6: Documentation (20 min)
- Add usage examples
- Document data types
- Generate website docs

### Phase 7: Validation (20 min)
- Format code
- Compile provider
- Run linter

**Total**: 2.5 hours

---

## Documentation Structure

### In `resource_tc_cls_cos_shipper.md`

1. **Usage Example**: Basic Parquet configuration
2. **Data Types Table**: All 6 types with descriptions
3. **Best Practices**: 
   - Column names must match log fields
   - Use appropriate types for storage optimization
   - Combine with compression for best results
4. **Import Example**: How to import existing Parquet shippers

### Auto-Generated Docs

- **Argument Reference**: Includes `parquet` block and `parquet_key_info` fields
- **Field Descriptions**: Clear explanations for key_name and key_type
- **Validation Rules**: Documents valid key_type values

---

## Testing Strategy

### Unit Tests
1. **Create with Parquet**: Basic Parquet configuration
2. **All data types**: Test each of 6 types
3. **Validation**: Invalid type rejection
4. **Update**: Modify column definitions
5. **Import**: Import existing Parquet shipper

### Validation Tests
- Schema compilation
- Validation function for key_type
- Minimum columns requirement

### Integration Tests (with credentials)
- End-to-end create → read → update → delete
- Format conversion (JSON → Parquet)

---

## Success Criteria

### Functional
- ✅ Users can configure Parquet format
- ✅ All 6 data types work correctly
- ✅ Create/Read/Update/Delete operations succeed
- ✅ Import works for Parquet shippers
- ✅ Validation catches errors at plan time

### Quality
- ✅ No linter errors
- ✅ Code follows project conventions
- ✅ Documentation is complete
- ✅ Backward compatible
- ✅ Tests pass

---

## Risk Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Type mismatch in logs | Medium | Document clearly, let API validate |
| Missing log fields | Low | API handles gracefully (nulls/skips) |
| Undocumented API limits | Low | Rely on API validation |

---

## Benefits

### For Users
- 🚀 **Better performance**: Columnar format optimized for analytics
- 💾 **Storage efficiency**: Better compression ratios
- 🔗 **Integration**: Direct compatibility with data warehouses
- 📊 **Query speed**: Faster analytical queries

### For Project
- ✅ **Feature parity**: Matches API capabilities
- 📦 **No dependencies**: Uses existing SDK structures
- 🛡️ **Low risk**: Additive, backward-compatible change

---

## Files

```
openspec/changes/add-cls-cos-shipper-parquet/
├── README.md          (this file)
├── proposal.md        (detailed design - 13 KB)
├── tasks.md           (26 implementation tasks - 10 KB)
└── specs/
    └── cls-cos-shipper-parquet-content/
        └── spec.md    (8 requirements, 24 scenarios - 15 KB)
```

---

## Next Steps

1. ✅ **Review** proposal and spec documents
2. ⏳ **Validate** with `openspec validate add-cls-cos-shipper-parquet --strict`
3. ⏳ **Approve** the proposal
4. ⏳ **Implement** using `openspec apply add-cls-cos-shipper-parquet`

---

## References

- **API Docs**: 
  - CreateShipper: https://cloud.tencent.com/document/api/614/31574
  - ModifyShipper: https://cloud.tencent.com/document/api/614/31575
  - DescribeShippers: https://cloud.tencent.com/document/api/614/31576

- **Parquet Format**: https://parquet.apache.org/

- **SDK**: `vendor/.../cls/v20201016/models.go` (lines 1571, 18630-18640)

- **Current Implementation**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go`

---

**Contact**: AI Assistant  
**Created**: 2026-02-11  
**Status**: Ready for Review ✅
