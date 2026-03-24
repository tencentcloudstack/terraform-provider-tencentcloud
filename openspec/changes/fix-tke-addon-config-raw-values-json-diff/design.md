# Technical Design: JSON Order Diff Suppression for TKE Addon Config

## Architecture Overview

This change adds a diff suppression layer to the Terraform resource schema that normalizes JSON comparison by parsing and deep-comparing the semantic content instead of string bytes.

```
┌─────────────────────────────────────────────────────────────┐
│                    Terraform Plan Phase                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Schema Diff Detection (raw_values)              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  DiffSuppressFunc: suppressJSONOrderDiff               │ │
│  │                                                          │ │
│  │  1. Receive: old (from state), new (from config)       │ │
│  │  2. Parse: json.Unmarshal both strings                 │ │
│  │  3. Compare: reflect.DeepEqual(oldJSON, newJSON)       │ │
│  │  4. Return: true (suppress) or false (show diff)       │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│         Terraform Shows Diff Only If Content Changed         │
└─────────────────────────────────────────────────────────────┘
```

## Component Design

### 1. Diff Suppression Function

**Function Signature:**
```go
func suppressJSONOrderDiff(k, old, new string, d *schema.ResourceData) bool
```

**Parameters:**
- `k` (string): Schema key being compared (e.g., "raw_values")
- `old` (string): Current value from Terraform state (API response)
- `new` (string): Desired value from user configuration
- `d` (*schema.ResourceData): Access to full resource data (unused in this implementation)

**Return Value:**
- `true`: Suppress diff (values are equivalent)
- `false`: Show diff (values differ)

**Algorithm:**

```
┌─────────────────────────────────┐
│ Start: Compare old and new      │
└────────────┬────────────────────┘
             │
             ▼
      ┌──────────────┐
      │ Both empty?  │───Yes──→ Return true (no diff)
      └──────┬───────┘
             No
             │
             ▼
      ┌──────────────┐
      │ One empty?   │───Yes──→ Return false (diff)
      └──────┬───────┘
             No
             │
             ▼
      ┌──────────────────────┐
      │ Parse old as JSON    │
      └──────┬───────────────┘
             │
             ├─Error→ Return old == new (string compare)
             │
             ▼
      ┌──────────────────────┐
      │ Parse new as JSON    │
      └──────┬───────────────┘
             │
             ├─Error→ Return old == new (string compare)
             │
             ▼
      ┌──────────────────────────────┐
      │ reflect.DeepEqual(old, new)  │
      └──────────┬───────────────────┘
                 │
                 ▼
         ┌─────────────┐
         │ Return bool │
         └─────────────┘
```

**Error Handling:**
- JSON parse failures gracefully fall back to string comparison
- Logged as warnings for debugging
- Ensures safe behavior even with malformed data

### 2. Schema Modification

**Location:** `resource_tc_kubernetes_addon_config.go`, Schema definition

**Before:**
```go
"raw_values": {
    Type:        schema.TypeString,
    Optional:    true,
    Computed:    true,
    Description: "Params of addon, base64 encoded json format.",
},
```

**After:**
```go
"raw_values": {
    Type:             schema.TypeString,
    Optional:         true,
    Computed:         true,
    Description:      "Params of addon, base64 encoded json format.",
    DiffSuppressFunc: suppressJSONOrderDiff,
},
```

**Rationale:**
- `DiffSuppressFunc` is a standard Terraform SDK v2 feature
- Called automatically during plan/diff operations
- Non-invasive: doesn't affect Create/Read/Update/Delete logic

### 3. Import Additions

**New imports required:**

```go
import (
    // ... existing imports ...
    "encoding/json"  // For JSON parsing
    "reflect"        // For deep equality comparison
)
```

**Why these packages:**
- `encoding/json`: Standard library, zero external dependencies
- `reflect`: Standard library, provides `DeepEqual` for structural comparison

## Data Flow

### Plan/Diff Operation

```
User Config (HCL)                    Terraform State (JSON)
     │                                        │
     │  raw_values:                          │  raw_values:
     │  '{"a":1,"b":2}'                      │  '{"b":2,"a":1}'
     │                                        │
     └────────────┬───────────────────────────┘
                  │
                  ▼
      ┌──────────────────────────┐
      │ Terraform Diff Engine    │
      │ Calls DiffSuppressFunc   │
      └──────────┬───────────────┘
                 │
                 ▼
      ┌─────────────────────────────┐
      │ suppressJSONOrderDiff        │
      │                              │
      │ old = '{"b":2,"a":1}'       │
      │ new = '{"a":1,"b":2}'       │
      └──────────┬──────────────────┘
                 │
                 ▼
      ┌─────────────────────────────┐
      │ json.Unmarshal              │
      │ old → map[a:1 b:2]          │
      │ new → map[a:1 b:2]          │
      └──────────┬──────────────────┘
                 │
                 ▼
      ┌─────────────────────────────┐
      │ reflect.DeepEqual           │
      │ Returns: true               │
      └──────────┬──────────────────┘
                 │
                 ▼
      ┌─────────────────────────────┐
      │ Terraform: No diff shown    │
      └─────────────────────────────┘
```

### Create/Update Operation

**No Changes Required** - The diff suppression only affects plan/diff logic, not actual CRUD operations.

## Edge Cases and Handling

### 1. Empty Strings
| Old    | New    | Result | Reason                |
|--------|--------|--------|-----------------------|
| ""     | ""     | true   | Both empty = no diff  |
| ""     | "{}"   | false  | Different values      |
| "{}"   | ""     | false  | Different values      |

### 2. Invalid JSON
| Scenario                    | Behavior                      |
|-----------------------------|-------------------------------|
| Old is invalid JSON         | Fallback to string comparison |
| New is invalid JSON         | Fallback to string comparison |
| Both are invalid JSON       | String comparison decides     |

**Example:**
```go
old = "not-json"
new = "not-json"
// Result: true (strings match)

old = "not-json"
new = "{valid}"
// Result: false (strings don't match, new is valid JSON)
```

### 3. Nested JSON
```json
old: {"outer": {"inner": "value", "other": 1}}
new: {"outer": {"other": 1, "inner": "value"}}
// Result: true (deep equality handles nested objects)
```

### 4. JSON Arrays
**Important:** Arrays ARE order-sensitive in JSON semantics.

```json
old: {"list": [1, 2, 3]}
new: {"list": [3, 2, 1]}
// Result: false (array order matters)
```

This is correct behavior - array ordering is semantically significant.

### 5. Base64 Encoded Data
The field description mentions "base64 encoded json format", but in the code:

- **Read function** (line 115-117): Decodes base64 → stores decoded JSON string in state
- **Update function** (line 167-169): Encodes JSON string → sends base64 to API

**Implication:** The diff suppression operates on the **decoded JSON string** in Terraform state, which is correct.

## Performance Considerations

### Time Complexity
- JSON parsing: O(n) where n = string length
- Deep equality: O(m) where m = number of JSON elements
- Overall: O(n + m), linear and acceptable

### Space Complexity
- Allocates temporary objects during parsing
- For typical addon configs (< 10KB), memory overhead is negligible

### Benchmark Estimates
Typical addon config size: 1-5 KB

| Operation          | Time         | Notes                    |
|--------------------|--------------|--------------------------|
| JSON parse (1KB)   | ~50-100µs    | Standard library         |
| Deep compare       | ~10-50µs     | Few dozen keys typically |
| Total overhead     | ~100-200µs   | Per plan operation       |

**Conclusion:** Performance impact is negligible compared to network API calls (typically 100-500ms).

## Testing Strategy

### Unit Tests

**Test File:** `resource_tc_kubernetes_addon_config_test.go`

**Test Cases:**

```go
func TestSuppressJSONOrderDiff(t *testing.T) {
    cases := []struct {
        name     string
        old      string
        new      string
        expected bool
    }{
        {
            name:     "empty strings",
            old:      "",
            new:      "",
            expected: true,
        },
        {
            name:     "one empty",
            old:      "",
            new:      `{"a":1}`,
            expected: false,
        },
        {
            name:     "same json different order",
            old:      `{"a":1,"b":2,"c":3}`,
            new:      `{"c":3,"a":1,"b":2}`,
            expected: true,
        },
        {
            name:     "different json content",
            old:      `{"a":1}`,
            new:      `{"a":2}`,
            expected: false,
        },
        {
            name:     "nested objects",
            old:      `{"outer":{"a":1,"b":2}}`,
            new:      `{"outer":{"b":2,"a":1}}`,
            expected: true,
        },
        {
            name:     "arrays order matters",
            old:      `{"list":[1,2,3]}`,
            new:      `{"list":[3,2,1]}`,
            expected: false,
        },
        {
            name:     "invalid json fallback",
            old:      "not-json",
            new:      "not-json",
            expected: true,
        },
        {
            name:     "one invalid json",
            old:      "not-json",
            new:      `{"a":1}`,
            expected: false,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            result := suppressJSONOrderDiff("raw_values", tc.old, tc.new, nil)
            if result != tc.expected {
                t.Errorf("Expected %v, got %v", tc.expected, result)
            }
        })
    }
}
```

### Integration Tests

**Acceptance Test:**

```go
func TestAccTencentCloudKubernetesAddonConfig_jsonOrderDiff(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:     func() { testAccPreCheck(t) },
        Providers:    testAccProviders,
        CheckDestroy: testAccCheckKubernetesAddonConfigDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccKubernetesAddonConfig_basic(`{"image":"nginx","port":80}`),
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckKubernetesAddonConfigExists("tencentcloud_kubernetes_addon_config.test"),
                ),
            },
            {
                // Assume API returns different order: {"port":80,"image":"nginx"}
                // This plan should show no changes
                Config:             testAccKubernetesAddonConfig_basic(`{"image":"nginx","port":80}`),
                PlanOnly:           true,
                ExpectNonEmptyPlan: false, // No diff expected
            },
        },
    })
}
```

## Rollback Plan

If issues arise:

1. **Immediate Rollback:**
   ```go
   "raw_values": {
       Type:        schema.TypeString,
       Optional:    true,
       Computed:    true,
       Description: "Params of addon, base64 encoded json format.",
       // Remove: DiffSuppressFunc: suppressJSONOrderDiff,
   },
   ```

2. **Remove Helper Function:**
   Delete `suppressJSONOrderDiff` function

3. **Remove Imports:**
   Remove `encoding/json` and `reflect` if not used elsewhere

**Risk:** Rollback is trivial; change is fully encapsulated.

## Security Considerations

- **No security impact**: Only affects local diff comparison
- **No API changes**: No changes to network requests or authentication
- **No data exposure**: Function doesn't log sensitive data
- **Standard library only**: No external dependencies with vulnerabilities

## Compliance and Best Practices

✅ **Follows Terraform Plugin SDK v2 patterns**
✅ **Uses standard library only (no external deps)**
✅ **Backward compatible**
✅ **Graceful error handling**
✅ **Minimal performance overhead**
✅ **Well-tested approach**
✅ **Follows project code organization (helper function at end of file)**

## Future Enhancements

Potential future improvements (not in scope):

1. **Logging:** Add DEBUG-level logs for troubleshooting
2. **Metrics:** Track JSON parse failures
3. **Configuration:** Allow users to disable suppression (likely unnecessary)
4. **Apply to other fields:** If similar issues arise in other resources

## References

- [Terraform Plugin SDK v2 Schema Documentation](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas)
- [Go encoding/json Package](https://pkg.go.dev/encoding/json)
- [Go reflect.DeepEqual](https://pkg.go.dev/reflect#DeepEqual)
- [TKE UpdateAddon API](https://cloud.tencent.com/document/product/457/43259)
