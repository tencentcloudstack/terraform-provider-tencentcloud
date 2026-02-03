# Spec: VOD Sub Applications Data Source

**Capability ID**: `datasource-vod-sub-applications`  
**Type**: Data Source  
**Status**: Draft

## Overview
This spec defines the `tencentcloud_vod_sub_applications` data source for querying Tencent Cloud VOD (Video on Demand) sub-application lists.

## ADDED Requirements

### Requirement 1: Query Sub-Applications by Name Filter
**ID**: `DS-VOD-SUBAPPS-001`  
**Priority**: High  
**Category**: Query Filtering

The data source MUST support filtering sub-applications by application name.

#### Scenario: Query by exact name match
**Given** a VOD account with multiple sub-applications  
**When** user specifies a `name` parameter  
**Then** only sub-applications matching that exact name are returned  
**And** the result set includes all fields from the matched applications

```hcl
data "tencentcloud_vod_sub_applications" "example" {
  name = "MyVideoApp"
}

output "app_id" {
  value = data.tencentcloud_vod_sub_applications.example.sub_application_info_set[0].sub_app_id
}
```

---

### Requirement 2: Query Sub-Applications by Tags
**ID**: `DS-VOD-SUBAPPS-002`  
**Priority**: High  
**Category**: Query Filtering

The data source MUST support filtering sub-applications by resource tags.

#### Scenario: Query applications with specific tags
**Given** sub-applications with various tag assignments  
**When** user specifies tags filter as a map  
**Then** only sub-applications with all matching tags are returned  
**And** returned applications include their complete tag information

```hcl
data "tencentcloud_vod_sub_applications" "by_tags" {
  tags = {
    Environment = "Production"
    Team        = "VideoTeam"
  }
}

output "production_apps" {
  value = data.tencentcloud_vod_sub_applications.by_tags.sub_application_info_set[*].sub_app_id_name
}
```

---

### Requirement 3: Support Pagination
**ID**: `DS-VOD-SUBAPPS-003`  
**Priority**: Medium  
**Category**: Query Control

The data source MUST support pagination for handling large result sets efficiently.

#### Scenario: Query with pagination parameters
**Given** an account with many sub-applications  
**When** user specifies `offset` and `limit` parameters  
**Then** results are returned starting from the offset position  
**And** at most `limit` number of results are returned  
**And** pagination follows the API's limit constraints (max 200 per page)

```hcl
data "tencentcloud_vod_sub_applications" "page_one" {
  offset = 0
  limit  = 50
}

data "tencentcloud_vod_sub_applications" "page_two" {
  offset = 50
  limit  = 50
}
```

---

### Requirement 4: Return Complete Sub-Application Information
**ID**: `DS-VOD-SUBAPPS-004`  
**Priority**: High  
**Category**: Data Completeness

The data source MUST return all available fields from the `SubAppIdInfo` API structure.

#### Scenario: Access all sub-application fields
**Given** a query that returns sub-applications  
**When** results are read from the data source  
**Then** each application includes:
- `sub_app_id` (Integer): Application ID
- `sub_app_id_name` (String): Application name
- `name` (String): Legacy name field
- `description` (String): Application description
- `create_time` (String): Creation time in ISO 8601
- `status` (String): Status (On/Off/Destroying/Destroyed)
- `mode` (String): Mode (fileid or fileid+path)
- `storage_regions` (List): Enabled storage regions
- `tags` (Map): Resource tags

```hcl
data "tencentcloud_vod_sub_applications" "all" {}

output "first_app_details" {
  value = {
    id          = data.tencentcloud_vod_sub_applications.all.sub_application_info_set[0].sub_app_id
    name        = data.tencentcloud_vod_sub_applications.all.sub_application_info_set[0].sub_app_id_name
    description = data.tencentcloud_vod_sub_applications.all.sub_application_info_set[0].description
    status      = data.tencentcloud_vod_sub_applications.all.sub_application_info_set[0].status
    created     = data.tencentcloud_vod_sub_applications.all.sub_application_info_set[0].create_time
    mode        = data.tencentcloud_vod_sub_applications.all.sub_application_info_set[0].mode
    regions     = data.tencentcloud_vod_sub_applications.all.sub_application_info_set[0].storage_regions
    tags        = data.tencentcloud_vod_sub_applications.all.sub_application_info_set[0].tags
  }
}
```

---

### Requirement 5: Export Results to JSON File
**ID**: `DS-VOD-SUBAPPS-005`  
**Priority**: Low  
**Category**: Data Export

The data source SHOULD support exporting query results to a JSON file.

#### Scenario: Export query results
**Given** a successful query returning sub-applications  
**When** user specifies `result_output_file` parameter  
**Then** query results are written to the specified file path in JSON format  
**And** the file contains the complete application information

```hcl
data "tencentcloud_vod_sub_applications" "export" {
  result_output_file = "/tmp/vod_apps.json"
}
```

---

### Requirement 6: Handle Empty Results Gracefully
**ID**: `DS-VOD-SUBAPPS-006`  
**Priority**: Medium  
**Category**: Error Handling

The data source MUST handle empty result sets without error.

#### Scenario: Query returns no results
**Given** filter parameters that match no sub-applications  
**When** the query is executed  
**Then** an empty list is returned for `sub_application_info_set`  
**And** no error is raised  
**And** the data source read completes successfully

```hcl
data "tencentcloud_vod_sub_applications" "nonexistent" {
  name = "NonExistentApp123456"
}

output "count" {
  value = length(data.tencentcloud_vod_sub_applications.nonexistent.sub_application_info_set)
  # Expected: 0
}
```

---

### Requirement 7: Support API Retry Logic
**ID**: `DS-VOD-SUBAPPS-007`  
**Priority**: High  
**Category**: Reliability

The data source MUST implement retry logic for API calls to handle transient failures.

#### Scenario: Retry on temporary API failure
**Given** the VOD API experiences temporary unavailability  
**When** a query is executed  
**Then** the request is retried with exponential backoff  
**And** retries continue until success or timeout  
**And** appropriate error messages are logged

**Implementation notes**:
- Use `resource.Retry` with `tccommon.ReadRetryTimeout`
- Use `tccommon.RetryError` for retryable errors
- Follow patterns from `resource_tc_vod_sub_application.go`

---

### Requirement 8: Implement Proper Rate Limiting
**ID**: `DS-VOD-SUBAPPS-008`  
**Priority**: Medium  
**Category**: API Management

The data source MUST respect API rate limits to prevent throttling.

#### Scenario: Rate limit compliance
**Given** multiple API requests in pagination loop  
**When** each request is about to be sent  
**Then** `ratelimit.Check(request.GetAction())` is called  
**And** rate limiting delays are applied as needed

---

### Requirement 9: Reference Sub-Applications in Other Resources
**ID**: `DS-VOD-SUBAPPS-009`  
**Priority**: High  
**Category**: Integration

The data source MUST enable referencing sub-application IDs in other resources.

#### Scenario: Use queried sub-app in resource configuration
**Given** a query that returns sub-applications  
**When** the sub-app ID is referenced in another resource  
**Then** the reference resolves correctly  
**And** dependencies are properly tracked by Terraform

```hcl
data "tencentcloud_vod_sub_applications" "existing" {
  name = "ProductionApp"
}

resource "tencentcloud_vod_super_player_config" "example" {
  sub_app_id = data.tencentcloud_vod_sub_applications.existing.sub_application_info_set[0].sub_app_id
  name       = "player-config"
  # ... other configuration
}
```

---

## Schema Definition

### Input Arguments

| Argument | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `name` | String | No | - | Application name for exact match filtering |
| `tags` | Map(String) | No | - | Tag key-value pairs for filtering applications |
| `offset` | Integer | No | 0 | Pagination offset for result set |
| `limit` | Integer | No | 200 | Maximum number of results per page (max: 200) |
| `result_output_file` | String | No | - | File path for exporting results in JSON format |

### Output Attributes

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | Data source identifier (hash of result IDs) |
| `sub_application_info_set` | List(Object) | List of sub-application information |

### sub_application_info_set Object Schema

| Field | Type | Description |
|-------|------|-------------|
| `sub_app_id` | Integer | Sub-application unique ID |
| `sub_app_id_name` | String | Sub-application name (current field) |
| `name` | String | Legacy name field (for backward compatibility) |
| `description` | String | Application description |
| `create_time` | String | Creation time in ISO 8601 format |
| `status` | String | Application status: On, Off, Destroying, Destroyed |
| `mode` | String | Application mode: fileid or fileid+path |
| `storage_regions` | List(String) | List of enabled storage regions |
| `tags` | Map(String) | Resource tags as key-value pairs |

---

## API Mapping

### Source API
- **Name**: `DescribeSubAppIds`
- **Version**: `2018-07-17`
- **Package**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717`

### Request Mapping

| Terraform Input | API Parameter | Transformation |
|-----------------|---------------|----------------|
| `name` | `Name` | Direct mapping |
| `tags` | `Tags` | Convert map to `[]*ResourceTag` |
| `offset` | `Offset` | Direct mapping as `*uint64` |
| `limit` | `Limit` | Direct mapping as `*uint64` |

### Response Mapping

| API Response Field | Terraform Output | Transformation |
|-------------------|------------------|----------------|
| `SubAppIdInfoSet[].SubAppId` | `sub_app_id` | `*uint64` to `int` |
| `SubAppIdInfoSet[].SubAppIdName` | `sub_app_id_name` | `*string` to `string` |
| `SubAppIdInfoSet[].Name` | `name` | `*string` to `string` |
| `SubAppIdInfoSet[].Description` | `description` | `*string` to `string` |
| `SubAppIdInfoSet[].CreateTime` | `create_time` | `*string` to `string` |
| `SubAppIdInfoSet[].Status` | `status` | `*string` to `string` |
| `SubAppIdInfoSet[].Mode` | `mode` | `*string` to `string` |
| `SubAppIdInfoSet[].StorageRegions` | `storage_regions` | `[]*string` to `[]string` |
| `SubAppIdInfoSet[].Tags` | `tags` | `[]*ResourceTag` to `map[string]string` |

---

## Implementation Notes

1. **Service Layer**: Add `DescribeSubApplicationsByFilter` method to `VodService` in `service_tencentcloud_vod.go`
2. **Pagination**: Implement full pagination to collect all results across multiple API calls
3. **Error Handling**: Use project-standard retry patterns and error logging
4. **Testing**: Cover basic query, name filter, tag filter, and pagination scenarios
5. **Documentation**: Follow existing VOD data source documentation patterns

---

## Related Changes
- Reference: `align-vod-sub-application-params` (understanding of SubAppIdInfo structure)
- Pattern: `data_source_tc_vod_super_player_configs.go` (similar data source implementation)

---

## Backward Compatibility
This is a new data source with no impact on existing resources or data sources.

---

## Security Considerations
- No sensitive data exposed beyond standard cloud resource metadata
- Tags may contain sensitive information but are already handled by VOD API
- No special authentication beyond standard provider credentials

---

## Performance Considerations
- Pagination ensures efficient handling of large result sets
- Rate limiting prevents API throttling
- Retry logic handles transient failures without manual intervention
