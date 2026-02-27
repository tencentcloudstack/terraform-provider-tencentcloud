# Tasks: Add Parquet Format Support to CLS COS Shipper Content

**Change ID**: `add-cls-cos-shipper-parquet`  
**Total Tasks**: 26  
**Estimated Time**: 2.5 hours

---

## Phase 1: Schema Definition (30 min)

### Task 1.1: Add `parquet` Nested Block to Content Schema
**File**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go`

Add after the `json` block definition (around line 150):

```go
"parquet": {
	Type:     schema.TypeList,
	Optional: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"parquet_key_info": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "Array of Parquet column definitions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Column name in the Parquet file.",
						},
						"key_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"string", "boolean", "int32", "int64", "float", "double"}),
							Description:  "Data type of the column. Valid values: string, boolean, int32, int64, float, double.",
						},
					},
				},
			},
		},
	},
	Description: "Parquet format content description.Note: this field may return null, indicating that no valid values can be obtained.",
},
```

**Validation**:
- ✅ Parquet block is optional
- ✅ MaxItems: 1 ensures only one parquet block
- ✅ parquet_key_info is required with MinItems: 1
- ✅ key_type has validation for 6 allowed types

---

### Task 1.2: Update Content Format Description
**File**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go`

Update the `format` field description (line ~110):

```go
"format": {
	Type:        schema.TypeString,
	Required:    true,
	Description: "Content format. Valid values: json, csv, parquet.",
},
```

**Validation**:
- ✅ Description mentions all three formats

---

### Task 1.3: Verify Schema Compilation
**Command**: `go build ./tencentcloud/services/cls/...`

**Validation**:
- ✅ No compilation errors
- ✅ Schema syntax is correct

---

## Phase 2: Create Operation (20 min)

### Task 2.1: Handle Parquet Block in Create
**File**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go`

Add after JSON handling (around line 310):

```go
if v, ok := dMap["parquet"]; ok {
	if len(v.([]interface{})) == 1 {
		parquet := v.([]interface{})[0].(map[string]interface{})
		parquetInfo := cls.ParquetInfo{}
		
		if keyInfos, ok := parquet["parquet_key_info"]; ok {
			parquetKeyInfoList := keyInfos.([]interface{})
			parquetInfo.ParquetKeyInfo = make([]*cls.ParquetKeyInfo, 0, len(parquetKeyInfoList))
			
			for _, keyInfo := range parquetKeyInfoList {
				keyInfoMap := keyInfo.(map[string]interface{})
				parquetKeyInfo := &cls.ParquetKeyInfo{
					KeyName: helper.String(keyInfoMap["key_name"].(string)),
					KeyType: helper.String(keyInfoMap["key_type"].(string)),
				}
				parquetInfo.ParquetKeyInfo = append(parquetInfo.ParquetKeyInfo, parquetKeyInfo)
			}
		}
		content.Parquet = &parquetInfo
	}
}
```

**Validation**:
- ✅ Parquet block is read from schema
- ✅ ParquetKeyInfo array is populated
- ✅ Proper type conversion with helper functions
- ✅ Handles missing parquet_key_info gracefully

---

### Task 2.2: Verify Create Logic
**Verification**:
- Check that `content.Parquet` is assigned when format is "parquet"
- Ensure no conflicts with CSV/JSON handling

**Validation**:
- ✅ Logic flow is correct
- ✅ Mutual exclusivity with csv/json maintained

---

### Task 2.3: Add Logging for Parquet Creation
**File**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go`

The existing logging at line 336 will automatically log Parquet content.

**Validation**:
- ✅ Request body logging includes Parquet data

---

### Task 2.4: Test Create Compilation
**Command**: `go build ./tencentcloud/services/cls/...`

**Validation**:
- ✅ No compilation errors in Create logic

---

## Phase 3: Read Operation (20 min)

### Task 3.1: Handle Parquet in Read Response
**File**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go`

Add after JSON handling in Read (around line 438):

```go
if shipper.Content.Parquet != nil {
	parquetKeyInfoList := make([]interface{}, 0, len(shipper.Content.Parquet.ParquetKeyInfo))
	
	for _, keyInfo := range shipper.Content.Parquet.ParquetKeyInfo {
		parquetKeyInfoMap := map[string]interface{}{
			"key_name": keyInfo.KeyName,
			"key_type": keyInfo.KeyType,
		}
		parquetKeyInfoList = append(parquetKeyInfoList, parquetKeyInfoMap)
	}
	
	parquet := map[string]interface{}{
		"parquet_key_info": parquetKeyInfoList,
	}
	content["parquet"] = []interface{}{parquet}
}
```

**Validation**:
- ✅ Nil check for Parquet pointer
- ✅ ParquetKeyInfo array is converted to interface slice
- ✅ Nested structure matches schema definition
- ✅ Added to content map correctly

---

### Task 3.2: Verify Read State Management
**Verification**:
- Ensure `parquet` is added to content map
- Verify the final `d.Set("content", ...)` includes parquet data

**Validation**:
- ✅ State is set correctly
- ✅ No data loss in round-trip (create → read)

---

### Task 3.3: Handle Nil Pointers Safely
**File**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go`

Review all pointer accesses in Read operation:

```go
// Already has nil checks:
if shipper.Content != nil { ... }
if shipper.Content.Parquet != nil { ... }
```

**Validation**:
- ✅ All pointers checked before dereference
- ✅ No potential panics

---

### Task 3.4: Test Read Compilation
**Command**: `go build ./tencentcloud/services/cls/...`

**Validation**:
- ✅ No compilation errors in Read logic

---

## Phase 4: Update Operation (20 min)

### Task 4.1: Handle Parquet in Update
**File**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go`

Add after JSON handling in Update (around line 589):

```go
if v, ok := dMap["parquet"]; ok {
	if len(v.([]interface{})) == 1 {
		parquet := v.([]interface{})[0].(map[string]interface{})
		parquetInfo := cls.ParquetInfo{}
		
		if keyInfos, ok := parquet["parquet_key_info"]; ok {
			parquetKeyInfoList := keyInfos.([]interface{})
			parquetInfo.ParquetKeyInfo = make([]*cls.ParquetKeyInfo, 0, len(parquetKeyInfoList))
			
			for _, keyInfo := range parquetKeyInfoList {
				keyInfoMap := keyInfo.(map[string]interface{})
				parquetKeyInfo := &cls.ParquetKeyInfo{
					KeyName: helper.String(keyInfoMap["key_name"].(string)),
					KeyType: helper.String(keyInfoMap["key_type"].(string)),
				}
				parquetInfo.ParquetKeyInfo = append(parquetInfo.ParquetKeyInfo, parquetKeyInfo)
			}
		}
		content.Parquet = &parquetInfo
	}
}
```

**Validation**:
- ✅ Logic mirrors Create operation
- ✅ Wrapped in `d.HasChange("content")` check
- ✅ Proper error handling

---

### Task 4.2: Verify Update Detection
**Verification**:
- Update operation already checks `d.HasChange("content")`
- Parquet changes will trigger update correctly

**Validation**:
- ✅ Change detection works for parquet modifications

---

### Task 4.3: Test Update Compilation
**Command**: `go build ./tencentcloud/services/cls/...`

**Validation**:
- ✅ No compilation errors in Update logic

---

## Phase 5: Testing (30 min)

### Task 5.1: Add Test Case - Create with Parquet
**File**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper_test.go` (if exists)

Add test configuration:

```go
const testAccClsCosShipperParquet = `
resource "tencentcloud_cls_cos_shipper" "parquet" {
  bucket       = tencentcloud_cos_bucket.example.id
  topic_id     = tencentcloud_cls_topic.example.id
  shipper_name = "test-parquet-shipper"
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
`
```

**Validation**:
- ✅ Test compiles
- ✅ Configuration is valid HCL

---

### Task 5.2: Add Test Case - All Data Types
**File**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper_test.go`

```go
const testAccClsCosShipperParquetAllTypes = `
resource "tencentcloud_cls_cos_shipper" "all_types" {
  bucket       = tencentcloud_cos_bucket.example.id
  topic_id     = tencentcloud_cls_topic.example.id
  shipper_name = "test-all-types"
  prefix       = "logs/types/"

  content {
    format = "parquet"

    parquet {
      parquet_key_info {
        key_name = "str_field"
        key_type = "string"
      }
      parquet_key_info {
        key_name = "bool_field"
        key_type = "boolean"
      }
      parquet_key_info {
        key_name = "int32_field"
        key_type = "int32"
      }
      parquet_key_info {
        key_name = "int64_field"
        key_type = "int64"
      }
      parquet_key_info {
        key_name = "float_field"
        key_type = "float"
      }
      parquet_key_info {
        key_name = "double_field"
        key_type = "double"
      }
    }
  }
}
`
```

**Validation**:
- ✅ All 6 data types tested
- ✅ Configuration is valid

---

### Task 5.3: Test Invalid Type Validation
**Verification**:
- Attempt to create with invalid key_type (e.g., "varchar")
- Should fail at plan time with validation error

**Validation**:
- ✅ Validation catches invalid types
- ✅ Error message is clear

---

### Task 5.4: Compile Tests
**Command**: `go test -c ./tencentcloud/services/cls/... -o /dev/null`

**Validation**:
- ✅ Tests compile successfully
- ✅ No syntax errors

---

## Phase 6: Documentation (20 min)

### Task 6.1: Add Parquet Example to Resource Documentation
**File**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.md`

Add after existing example (around line 62):

```markdown
Example with Parquet format:

```hcl
resource "tencentcloud_cls_cos_shipper" "parquet_example" {
  bucket       = tencentcloud_cos_bucket.example.id
  topic_id     = tencentcloud_cls_topic.example.id
  interval     = 300
  max_size     = 256
  partition    = "/%Y/%m/%d/%H/"
  prefix       = "logs/parquet/"
  shipper_name = "parquet-shipper"

  compress {
    format = "gzip"
  }

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
```

**Validation**:
- ✅ Example is syntactically correct
- ✅ Shows realistic use case

---

### Task 6.2: Document Parquet Data Types
**File**: `tencentcloud/services/cls/resource_tc_cls_cos_shipper.md`

Add after examples:

```markdown
## Parquet Format

When using `format = "parquet"`, you must define the Parquet column schema using `parquet_key_info` blocks.

### Supported Data Types

| Type | Description | Example Use Case |
|------|-------------|------------------|
| `string` | Text data | Log messages, usernames, IDs |
| `boolean` | True/false values | Flags, success/failure states |
| `int32` | 32-bit integers | Status codes, counters (< 2B) |
| `int64` | 64-bit integers | Timestamps, large counters, IDs |
| `float` | 32-bit floating point | CPU usage, percentages |
| `double` | 64-bit floating point | High-precision metrics, latencies |

### Best Practices

1. **Column Names**: Must match the field names in your CLS logs
2. **Data Types**: Choose appropriate types for your data to optimize storage
3. **Compression**: Parquet with gzip compression provides excellent storage efficiency
4. **Analytics**: Parquet format is optimized for analytical queries in data warehouses
```

**Validation**:
- ✅ All types documented
- ✅ Clear guidance for users

---

### Task 6.3: Generate Website Documentation
**Command**: `make doc`

**Validation**:
- ✅ Command runs successfully
- ✅ File `website/docs/r/cls_cos_shipper.html.markdown` updated
- ✅ Parquet fields appear in generated docs

---

### Task 6.4: Verify Generated Documentation
**File**: `website/docs/r/cls_cos_shipper.html.markdown`

Check that generated docs include:
- `parquet` block in content
- `parquet_key_info` nested block
- All fields with descriptions

**Validation**:
- ✅ Documentation is complete
- ✅ All new fields documented

---

## Phase 7: Validation (20 min)

### Task 7.1: Format Code
**Command**: `gofmt -w tencentcloud/services/cls/resource_tc_cls_cos_shipper.go`

**Validation**:
- ✅ Code is properly formatted
- ✅ Consistent indentation

---

### Task 7.2: Compile Provider
**Command**: `go build -o /tmp/terraform-provider-tencentcloud .`

**Validation**:
- ✅ Provider compiles successfully
- ✅ No compilation errors
- ✅ Binary created

---

### Task 7.3: Check Linter (if available)
**Command**: `make lint` or `golangci-lint run ./tencentcloud/services/cls/...`

**Validation**:
- ✅ No new linter errors
- ✅ Code follows project style

---

### Task 7.4: Final Review Checklist

Review the implementation against requirements:

- ✅ Schema fields added correctly
- ✅ Validation for key_type implemented
- ✅ Create operation handles parquet
- ✅ Read operation retrieves parquet
- ✅ Update operation modifies parquet
- ✅ Nil pointers handled safely
- ✅ Error handling in place
- ✅ Documentation complete
- ✅ Examples provided for all types
- ✅ Code formatted
- ✅ No linter errors
- ✅ Backward compatible

**Validation**:
- ✅ All items checked
- ✅ Implementation complete
- ✅ Ready for PR

---

## Summary

**Total Tasks**: 26  
**Phases**: 7  
**Estimated Time**: 2.5 hours

### Task Breakdown by Phase:
1. Schema Definition: 3 tasks (30 min)
2. Create Operation: 4 tasks (20 min)
3. Read Operation: 4 tasks (20 min)
4. Update Operation: 3 tasks (20 min)
5. Testing: 4 tasks (30 min)
6. Documentation: 4 tasks (20 min)
7. Validation: 4 tasks (20 min)

### Key Implementation Points:
1. Parquet block is optional in content
2. parquet_key_info is a repeatable block (list)
3. key_type has validation for 6 allowed types
4. Same handling pattern in Create/Read/Update
5. Backward compatible (optional field)

### Success Criteria:
- ✅ Users can configure Parquet format
- ✅ All 6 data types supported and validated
- ✅ Create/Read/Update work correctly
- ✅ Import works for Parquet shippers
- ✅ Tests pass
- ✅ Documentation complete
- ✅ No breaking changes

---

**Ready for Implementation!**

Run `openspec apply add-cls-cos-shipper-parquet` to begin.
