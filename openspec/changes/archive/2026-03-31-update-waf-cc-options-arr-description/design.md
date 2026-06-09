# Design: Update tencentcloud_waf_cc options_arr Description

## Architecture

This is a documentation-only change that updates the Schema description of the `options_arr` field in the `tencentcloud_waf_cc` resource. No functional logic or data structures will be modified.

### File Structure

```
tencentcloud/services/waf/resource_tc_waf_cc.go  # Schema definition (修改)
```

## Detailed Design

### 1. Current Description (Lines 82-103)

The current description has the following issues:
- Uses backticks (\`\`) instead of double quotes ("")
- Missing `URL` key type
- Missing `IPLocation` key type
- Incomplete match operator values for some key types
- Cookie and CustomHeader match operators are documented twice with inconsistencies

### 2. Updated Description

Based on the official API documentation (https://cloud.tencent.com/document/api/627/97646, updated 2026-03-06), the new description should include:

#### Supported Key Types
1. `URL`
2. `Method`
3. `Post`
4. `Referer`
5. `Cookie`
6. `User-Agent`
7. `CustomHeader`
8. `IPLocation`
9. `CaptchaRisk`
10. `CaptchaDeviceRisk`
11. `CaptchaScore`

#### Match Operators by Key Type

| Key | Match Operators |
|-----|-----------------|
| **URL** | 0 (equal), 3 (not equal), 1 (prefix), 6 (suffix), 2 (contains), 7 (not contains) |
| **Method** | 0 (equal), 3 (not equal) |
| **Post** | 0 (equal), 3 (not equal), 2 (contains), 7 (not contains) |
| **Cookie** | 0 (equal), 3 (not equal), 2 (contains), 7 (not contains) |
| **Referer** | 0 (equal), 3 (not equal), 1 (prefix), 6 (suffix), 2 (contains), 7 (not contains), 12 (exists), 5 (not exists), 4 (empty) |
| **User-Agent** | 0 (equal), 3 (not equal), 1 (prefix), 6 (suffix), 2 (contains), 7 (not contains), 12 (exists), 5 (not exists), 4 (empty) |
| **CustomHeader** | 0 (equal), 3 (not equal), 2 (contains), 7 (not contains), 4 (empty), 5 (not exists) |
| **IPLocation** | 13 (belongs to), 14 (not belongs to) |
| **CaptchaRisk** | 15 (numerically equal), 16 (numerically not equal), 13 (belongs to), 14 (not belongs to), 12 (exists), 5 (not exists) |
| **CaptchaDeviceRisk** | 13 (belongs to), 14 (not belongs to), 12 (exists), 5 (not exists) |
| **CaptchaScore** | 15 (numerically equal), 17 (numerically greater than), 18 (numerically less than), 19 (numerically greater than or equal), 20 (numerically less than or equal), 12 (exists), 5 (not exists) |

#### Encoding Rules

**For Post, Cookie, CustomHeader**:
- Base64 encode both parameter name and value (remove trailing `=`)
- Concatenate with `=` sign
- Format: `Base64(name)=Base64(value)`

**For Referer, User-Agent**:
- Base64 encode the value only (remove trailing `=`)
- Prefix with `=` sign
- Format: `=Base64(value)`

### 3. Implementation Details

**Location**: `tencentcloud/services/waf/resource_tc_waf_cc.go`, lines 82-103

**Change Type**: Replace the entire `Description` string for the `options_arr` field

**Formatting Requirements**:
- Use double quotes ("") instead of backticks (\`\`)
- Maintain proper line breaks for readability
- Keep JSON example at the beginning
- Document all key types and their match operators
- Include encoding rules at the end

### 4. Post-Modification Steps

After updating the description:
1. Run `go fmt` on the modified file to ensure proper formatting
2. Save the file
3. Verify no syntax errors

## Validation

### Success Criteria
- [ ] Description uses double quotes ("") instead of backticks (\`\`)
- [ ] All 11 key types are documented
- [ ] All match operators are correctly listed for each key type
- [ ] Encoding rules are clearly explained
- [ ] JSON example is included
- [ ] Code is properly formatted with `go fmt`
- [ ] No syntax errors

### Code Quality
- [ ] Follows existing code style
- [ ] Description is clear and accurate
- [ ] Matches official API documentation
