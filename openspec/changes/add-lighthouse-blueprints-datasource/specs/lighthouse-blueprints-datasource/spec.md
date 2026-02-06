# Spec: Lighthouse Blueprints DataSource

**Capability**: `lighthouse-blueprints-datasource`  
**Type**: New DataSource  
**Service**: Lighthouse (轻量应用服务器)

---

## Overview

This spec defines a new Terraform datasource `tencentcloud_lighthouse_blueprints` that allows users to query Lighthouse blueprint (镜像) information using the TencentCloud `DescribeBlueprints` API.

**Key Principles**:
1. Hide pagination (offset/limit) from users - handle automatically
2. Support all API filters for maximum flexibility
3. Map all Blueprint fields from SDK response
4. Follow existing project patterns (similar to `tencentcloud_lighthouse_bundle`)
5. Maintain nil safety for all pointer fields

---

## ADDED Requirements

### Requirement 1: Query All Blueprints
**Priority**: P0 (Critical)  
**Category**: Core Functionality

Users must be able to query all available Lighthouse blueprints without any filters.

#### Scenario: Query all blueprints

**Given** the user has valid TencentCloud credentials  
**When** the user defines:
```hcl
data "tencentcloud_lighthouse_blueprints" "all" {
}
```
**And** runs `terraform plan` or `terraform apply`  
**Then** the datasource should:
- Call the `DescribeBlueprints` API
- Automatically paginate through all results (hide offset/limit)
- Return all available blueprints in `blueprint_set`
- Set a composite ID based on returned blueprint IDs

**Acceptance Criteria**:
- ✅ No required arguments
- ✅ Returns all blueprints (potentially hundreds)
- ✅ Pagination is transparent to user
- ✅ ID is set as hash of blueprint IDs

---

### Requirement 2: Filter Blueprints by Criteria
**Priority**: P0 (Critical)  
**Category**: Core Functionality

Users must be able to filter blueprints using the API's filter capabilities.

#### Scenario: Filter by platform type

**Given** the user wants only Linux blueprints  
**When** the user defines:
```hcl
data "tencentcloud_lighthouse_blueprints" "linux" {
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
}
```
**Then** the datasource should:
- Pass the filter to the API
- Return only blueprints matching the platform type
- Include all fields for matching blueprints

**Acceptance Criteria**:
- ✅ Filter by `platform-type` (LINUX_UNIX, WINDOWS)
- ✅ Only matching blueprints returned
- ✅ All blueprint fields populated

---

#### Scenario: Filter by blueprint type

**Given** the user wants only application blueprints  
**When** the user defines:
```hcl
data "tencentcloud_lighthouse_blueprints" "apps" {
  filters {
    name   = "blueprint-type"
    values = ["APP_OS"]
  }
}
```
**Then** the datasource should:
- Filter by blueprint type
- Return only APP_OS blueprints

**Acceptance Criteria**:
- ✅ Supports values: APP_OS, PURE_OS, DOCKER, PRIVATE, SHARED
- ✅ Filtering works correctly

---

#### Scenario: Multiple filters

**Given** the user wants Linux application blueprints  
**When** the user defines:
```hcl
data "tencentcloud_lighthouse_blueprints" "linux_apps" {
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
  
  filters {
    name   = "blueprint-type"
    values = ["APP_OS"]
  }
}
```
**Then** the datasource should:
- Apply both filters (AND logic)
- Return only blueprints matching both criteria

**Acceptance Criteria**:
- ✅ Multiple filters supported
- ✅ Filters combined with AND logic

---

#### Scenario: Filter by name

**Given** the user wants blueprints with specific name pattern  
**When** the user filters by blueprint-name  
**Then** matching blueprints are returned

**Acceptance Criteria**:
- ✅ `blueprint-name` filter works
- ✅ Partial name matching (if supported by API)

---

#### Scenario: Filter by state

**Given** the user wants only online blueprints  
**When** the user defines:
```hcl
data "tencentcloud_lighthouse_blueprints" "online" {
  filters {
    name   = "blueprint-state"
    values = ["ONLINE"]
  }
}
```
**Then** only ONLINE blueprints are returned

**Acceptance Criteria**:
- ✅ Filter by blueprint state
- ✅ State filtering works correctly

---

#### Scenario: Filter by scene

**Given** the user wants blueprints for a specific scene  
**When** the user filters by scene-id  
**Then** blueprints associated with that scene are returned

**Acceptance Criteria**:
- ✅ `scene-id` filter supported
- ✅ Scene-based filtering works

---

### Requirement 3: Query Specific Blueprint IDs
**Priority**: P0 (Critical)  
**Category**: Core Functionality

Users must be able to query specific blueprints by ID.

#### Scenario: Query by IDs

**Given** the user knows specific blueprint IDs  
**When** the user defines:
```hcl
data "tencentcloud_lighthouse_blueprints" "specific" {
  blueprint_ids = ["lhbp-f1lkcd41", "lhbp-ab123456"]
}
```
**Then** the datasource should:
- Query only those blueprint IDs
- Return matching blueprints
- Handle case where some IDs don't exist

**Acceptance Criteria**:
- ✅ Accepts list of blueprint IDs
- ✅ Returns only specified blueprints
- ✅ Gracefully handles non-existent IDs

---

#### Scenario: Cannot use IDs and filters together

**Given** the API constraint that IDs and filters cannot be used together  
**When** the user defines both:
```hcl
data "tencentcloud_lighthouse_blueprints" "invalid" {
  blueprint_ids = ["lhbp-xxx"]
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
}
```
**Then** the API should return an error  
**And** the error should be propagated to the user

**Acceptance Criteria**:
- ✅ API handles mutual exclusivity
- ✅ Clear error message returned
- ✅ Documentation warns about this constraint

---

### Requirement 4: Return Complete Blueprint Information
**Priority**: P0 (Critical)  
**Category**: Data Mapping

All Blueprint fields from the SDK must be mapped to the schema output.

#### Scenario: All fields populated

**Given** a blueprint with all fields present  
**When** the datasource fetches the blueprint  
**Then** the output should include all 18+ fields:

Required fields:
- `blueprint_id`: String, unique identifier
- `display_title`: String, display title
- `display_version`: String, version
- `os_name`: String, OS name
- `platform`: String, OS platform
- `platform_type`: String (LINUX_UNIX/WINDOWS)
- `blueprint_type`: String (APP_OS/PURE_OS/DOCKER/PRIVATE/SHARED)
- `image_url`: String, image URL
- `required_system_disk_size`: Int, required disk size (GB)
- `blueprint_state`: String, state
- `blueprint_name`: String, name
- `support_automation_tools`: Bool, automation support
- `required_memory_size`: Int, required memory (GB)
- `community_url`: String, community URL
- `guide_url`: String, guide URL

Optional/nullable fields:
- `description`: String, may be null
- `created_time`: String (ISO 8601), may be null
- `image_id`: String (CVM image ID), may be null
- `scene_id_set`: List of strings, may be null
- `docker_version`: String, may be null

**Acceptance Criteria**:
- ✅ All fields defined in schema
- ✅ All fields have correct types
- ✅ Nil checks for all pointer fields
- ✅ No nil pointer panics

---

#### Scenario: Handle null fields gracefully

**Given** a blueprint with some null fields  
**When** the datasource processes the response  
**Then** null fields should:
- Not cause crashes or panics
- Be omitted from output (Terraform behavior)
- Have proper nil checks in code

**Acceptance Criteria**:
- ✅ Nil checks for all pointer fields
- ✅ No crashes on null fields
- ✅ Terraform handles omitted fields correctly

---

### Requirement 5: Automatic Pagination
**Priority**: P0 (Critical)  
**Category**: Implementation Detail

The datasource must handle pagination automatically without exposing offset/limit to users.

#### Scenario: Pagination with many results

**Given** the API returns >100 blueprints  
**When** the user queries all blueprints  
**Then** the service layer should:
- Start with offset=0, limit=100
- Fetch first page
- Check if more results exist
- Increment offset and fetch next page
- Repeat until all results fetched
- Return complete list to datasource

**Acceptance Criteria**:
- ✅ offset/limit not in schema (hidden)
- ✅ Service layer loops through pages
- ✅ All results returned to user
- ✅ No duplicate results
- ✅ Correct termination when no more results

**Implementation**:
```go
var offset int64 = 0
var pageSize int64 = 100
blueprints := make([]*lighthouse.Blueprint, 0)

for {
    request.Offset = &offset
    request.Limit = &pageSize
    
    response, err := me.client.UseLighthouseClient().DescribeBlueprints(request)
    // ... error handling ...
    
    if len(response.Response.BlueprintSet) < 1 {
        break
    }
    
    blueprints = append(blueprints, response.Response.BlueprintSet...)
    
    if len(response.Response.BlueprintSet) < int(pageSize) {
        break
    }
    
    offset += pageSize
}
```

---

### Requirement 6: Result Output File
**Priority**: P1 (Nice to have)  
**Category**: User Convenience

Users should be able to save results to a JSON file.

#### Scenario: Save results to file

**Given** the user wants to save results  
**When** the user defines:
```hcl
data "tencentcloud_lighthouse_blueprints" "all" {
  result_output_file = "blueprints.json"
}
```
**Then** the datasource should:
- Write the blueprint_set to the specified file
- Format as JSON
- Create file if not exists
- Overwrite if exists

**Acceptance Criteria**:
- ✅ File written with valid JSON
- ✅ Contains blueprint_set data
- ✅ Standard feature across datasources

---

### Requirement 7: Follow Project Conventions
**Priority**: P0 (Critical)  
**Category**: Code Quality

The implementation must follow all project conventions.

#### Scenario: Code structure

**Given** the terraform-provider-tencentcloud project conventions  
**When** implementing the datasource  
**Then** it must:

**File Naming**:
- ✅ Data source file: `data_source_tc_lighthouse_blueprints.go`
- ✅ Test file: `data_source_tc_lighthouse_blueprints_test.go`
- ✅ Doc file: `data_source_tc_lighthouse_blueprints.md`

**Function Naming**:
- ✅ DataSource constructor: `DataSourceTencentCloudLighthouseBlueprints()`
- ✅ Read function: `dataSourceTencentCloudLighthouseBlueprintsRead()`
- ✅ Service method: `DescribeLighthouseBlueprintsByFilter()`

**Resource Name**:
- ✅ Terraform name: `tencentcloud_lighthouse_blueprints`

**Imports**:
- ✅ Use `tccommon` for common package
- ✅ Use `lighthouse` for SDK package

**Error Handling**:
- ✅ Use `defer tccommon.LogElapsed()`
- ✅ Use `defer tccommon.InconsistentCheck()`
- ✅ Use `resource.Retry()` for read operations
- ✅ Use `tccommon.RetryError()` for retryable errors

**Logging**:
- ✅ Log all API requests/responses
- ✅ Log errors with context

---

#### Scenario: Service layer pattern

**Given** existing service layer methods  
**When** adding new method  
**Then** it must follow the pattern:

```go
func (me *LightHouseService) DescribeLighthouseBlueprintsByFilter(
    ctx context.Context,
    param map[string]interface{},
) (blueprints []*lighthouse.Blueprint, errRet error) {
    var (
        logId   = tccommon.GetLogId(ctx)
        request = lighthouse.NewDescribeBlueprintsRequest()
    )

    defer func() {
        if errRet != nil {
            log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
                logId, request.GetAction(), request.ToJsonString(), errRet.Error())
        }
    }()
    
    // ... implementation ...
}
```

**Acceptance Criteria**:
- ✅ Context passed as first parameter
- ✅ param map for flexibility
- ✅ Proper error handling with defer
- ✅ Logging for debugging
- ✅ ratelimit.Check() before API calls

---

### Requirement 8: Testing
**Priority**: P0 (Critical)  
**Category**: Quality Assurance

Comprehensive tests must be provided.

#### Scenario: Basic acceptance test

**Given** acceptance test infrastructure  
**When** running tests with TF_ACC=1  
**Then** tests must verify:

```go
func TestAccTencentCloudLighthouseBlueprintsDataSource_basic(t *testing.T) {
    t.Parallel()
    resource.Test(t, resource.TestCase{
        PreCheck:  func() { tcacctest.AccPreCheck(t) },
        Providers: tcacctest.AccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccLighthouseBlueprintsDataSource,
                Check: resource.ComposeTestCheckFunc(
                    tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_blueprints.blueprints"),
                    resource.TestCheckResourceAttrSet("data.tencentcloud_lighthouse_blueprints.blueprints", "blueprint_set.#"),
                ),
            },
        },
    })
}
```

**Acceptance Criteria**:
- ✅ Test is parallel (`t.Parallel()`)
- ✅ Test name follows convention
- ✅ Checks datasource ID is set
- ✅ Checks blueprint_set has results

---

#### Scenario: Filter tests

**Given** filter functionality  
**When** running filter tests  
**Then** tests must cover:
- Filter by platform-type
- Filter by blueprint-type
- Multiple filters
- Filter validation

**Acceptance Criteria**:
- ✅ Each filter type tested
- ✅ Tests verify filtered results
- ✅ Tests use realistic filter values

---

### Requirement 9: Documentation
**Priority**: P0 (Critical)  
**Category**: User Experience

Complete documentation must be provided.

#### Scenario: Usage examples

**Given** the datasource documentation  
**When** users read the docs  
**Then** examples must cover:

1. **Query all blueprints**:
```hcl
data "tencentcloud_lighthouse_blueprints" "all" {
}
```

2. **Filter by platform**:
```hcl
data "tencentcloud_lighthouse_blueprints" "linux" {
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
}
```

3. **Filter by type**:
```hcl
data "tencentcloud_lighthouse_blueprints" "apps" {
  filters {
    name   = "blueprint-type"
    values = ["APP_OS"]
  }
}
```

4. **Query by IDs**:
```hcl
data "tencentcloud_lighthouse_blueprints" "specific" {
  blueprint_ids = ["lhbp-xxx", "lhbp-yyy"]
}
```

5. **Multiple filters**:
```hcl
data "tencentcloud_lighthouse_blueprints" "linux_apps" {
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
  filters {
    name   = "blueprint-type"
    values = ["APP_OS"]
  }
}
```

**Acceptance Criteria**:
- ✅ All examples are valid HCL
- ✅ Examples demonstrate key use cases
- ✅ Examples are copy-paste ready

---

#### Scenario: Argument reference

**Given** the documentation  
**Then** it must document:

**Arguments** (all optional):
- `blueprint_ids` - (Optional, Set: [`String`]) Blueprint ID list.
- `filters` - (Optional, List) Filter list. Each element contains:
  - `name` - (Required, String) Filter name. Valid values: `blueprint-id`, `blueprint-type`, `platform-type`, `blueprint-name`, `blueprint-state`, `scene-id`.
  - `values` - (Required, Set: [`String`]) Filter values.
- `result_output_file` - (Optional, String) Used to save results.

**Note**: Cannot use `blueprint_ids` and `filters` together.

**Acceptance Criteria**:
- ✅ All arguments documented
- ✅ Types specified correctly
- ✅ Constraints noted (mutual exclusivity)

---

#### Scenario: Attributes reference

**Given** the documentation  
**Then** it must document all output attributes:

- `blueprint_set` - (List) List of blueprint details. Each element contains:
  - All 18+ Blueprint fields with types and descriptions
  - Nullable fields marked

**Acceptance Criteria**:
- ✅ All output fields documented
- ✅ Descriptions are clear and helpful
- ✅ Types match schema definition

---

## API Mapping

### Request Mapping

| Schema Field | API Parameter | Transform |
|-------------|--------------|-----------|
| blueprint_ids | BlueprintIds | Convert Set to []*string |
| filters.name | Filters[].Name | Direct map |
| filters.values | Filters[].Values | Convert Set to []*string |
| (internal) | Offset | Handled in service layer |
| (internal) | Limit | Handled in service layer |

### Response Mapping

| API Field | Schema Field | Type | Nullable |
|-----------|--------------|------|----------|
| BlueprintId | blueprint_id | string | No |
| DisplayTitle | display_title | string | No |
| DisplayVersion | display_version | string | No |
| Description | description | string | Yes |
| OsName | os_name | string | No |
| Platform | platform | string | No |
| PlatformType | platform_type | string | No |
| BlueprintType | blueprint_type | string | No |
| ImageUrl | image_url | string | No |
| RequiredSystemDiskSize | required_system_disk_size | int | No |
| BlueprintState | blueprint_state | string | No |
| CreatedTime | created_time | string | Yes |
| BlueprintName | blueprint_name | string | No |
| SupportAutomationTools | support_automation_tools | bool | No |
| RequiredMemorySize | required_memory_size | int | No |
| ImageId | image_id | string | Yes |
| CommunityUrl | community_url | string | No |
| GuideUrl | guide_url | string | No |
| SceneIdSet | scene_id_set | []string | Yes |
| DockerVersion | docker_version | string | Yes |

---

## Schema Definition

### Input Schema

```go
Schema: map[string]*schema.Schema{
    "blueprint_ids": {
        Optional: true,
        Type:     schema.TypeSet,
        Elem: &schema.Schema{
            Type: schema.TypeString,
        },
        Description: "Blueprint ID list.",
    },
    "filters": {
        Optional: true,
        Type:     schema.TypeList,
        Description: "Filter list. ...",
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "name": {
                    Type:        schema.TypeString,
                    Required:    true,
                    Description: "Field to be filtered.",
                },
                "values": {
                    Type: schema.TypeSet,
                    Elem: &schema.Schema{
                        Type: schema.TypeString,
                    },
                    Required:    true,
                    Description: "Filter value of field.",
                },
            },
        },
    },
    "blueprint_set": {
        Computed:    true,
        Type:        schema.TypeList,
        Description: "List of blueprint details.",
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                // ... all Blueprint fields ...
            },
        },
    },
    "result_output_file": {
        Type:        schema.TypeString,
        Optional:    true,
        Description: "Used to save results.",
    },
}
```

---

## Non-Functional Requirements

### Performance
- Pagination must be efficient (batch 100 at a time)
- API calls must use rate limiting (`ratelimit.Check()`)
- Large result sets (>500 blueprints) should complete in reasonable time

### Reliability
- Must handle API errors gracefully
- Must retry on transient failures
- Must handle nil fields without crashing

### Maintainability
- Code must follow project conventions
- Clear separation of concerns (datasource vs. service layer)
- Comprehensive error messages

### Security
- No security concerns (read-only datasource)
- API credentials handled by provider

---

## Testing Strategy

### Unit Tests
- Schema validation
- Parameter conversion logic
- Filter conversion logic

### Acceptance Tests (with TF_ACC=1)
1. Basic query (no filters)
2. Filter by platform-type
3. Filter by blueprint-type
4. Multiple filters
5. Query by IDs
6. result_output_file

### Manual Tests
- Large result sets (pagination)
- Error scenarios (invalid filters, non-existent IDs)
- Performance with many filters

---

## Success Metrics

- ✅ All acceptance tests pass
- ✅ Code passes linting
- ✅ Documentation complete and accurate
- ✅ No nil pointer panics
- ✅ Pagination transparent to users
- ✅ All Blueprint fields mapped correctly
- ✅ Follows project conventions

---

## References

- **API**: https://cloud.tencent.com/document/product/1207/47689
- **SDK**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324`
- **Blueprint struct**: `models.go:422-488`
- **Similar datasource**: `data_source_tc_lighthouse_bundle.go`
