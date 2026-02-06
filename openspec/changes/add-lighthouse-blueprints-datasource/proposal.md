# Proposal: Add Lighthouse Blueprints DataSource

## Change ID
`add-lighthouse-blueprints-datasource`

## Status
- **Phase**: Proposal
- **Created**: 2026-02-06
- **Type**: New DataSource

---

## Problem Statement

### Current Situation
The terraform-provider-tencentcloud currently lacks a direct datasource for querying Lighthouse blueprint (镜像) information using the `DescribeBlueprints` API. While there is `tencentcloud_lighthouse_instance_blueprint` datasource, it queries blueprints for specific instances only, not the general blueprint catalog.

### User Pain Points
1. **Cannot browse blueprint catalog**: Users cannot query available blueprints to discover what images are available for creating Lighthouse instances
2. **No filtering capability**: Users cannot filter blueprints by type (APP_OS, PURE_OS, DOCKER, PRIVATE, SHARED), platform (LINUX_UNIX, WINDOWS), or state
3. **Missing blueprint details**: Cannot get comprehensive blueprint information including OS details, required resources, and scene associations
4. **Workaround required**: Users must use the instance-specific datasource or directly query the API outside Terraform

### Business Impact
- Reduced user experience when working with Lighthouse resources
- Incomplete Terraform coverage of Lighthouse APIs
- Users cannot implement blueprint discovery workflows in Terraform

---

## Proposed Solution

### Overview
Add a new datasource `tencentcloud_lighthouse_blueprints` that wraps the TencentCloud `DescribeBlueprints` API to allow users to query and filter available Lighthouse blueprints.

### Core Design Decision
**Hide offset/limit from users** - The datasource will automatically handle pagination internally and return all matching blueprints, following the pattern used in other datasources like `tencentcloud_lighthouse_bundle`.

### Key Features
1. **Blueprint querying**: Query blueprints by IDs or filters
2. **Comprehensive filtering**: Support all API filters (blueprint-id, blueprint-type, platform-type, blueprint-name, blueprint-state, scene-id)
3. **Complete field mapping**: Expose all Blueprint struct fields from the SDK
4. **Automatic pagination**: Handle offset/limit internally, return all results
5. **Standard patterns**: Follow existing project conventions for datasources

---

## Technical Design

### API Integration
- **API**: `DescribeBlueprints` ([documentation](https://cloud.tencent.com/document/product/1207/47689))
- **SDK Package**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324`
- **Request**: `DescribeBlueprintsRequest`
- **Response**: `DescribeBlueprintsResponse`

### Schema Design

#### Input Arguments (Optional)
```hcl
data "tencentcloud_lighthouse_blueprints" "example" {
  blueprint_ids = ["lhbp-xxx", "lhbp-yyy"]  # Optional: Blueprint ID list
  
  filters {                                   # Optional: Filter list
    name   = "blueprint-type"
    values = ["APP_OS"]
  }
  
  result_output_file = "output.json"         # Optional: Save results to file
}
```

#### Output Attributes (Computed)
```hcl
# Returns all blueprint details
output "blueprints" {
  value = data.tencentcloud_lighthouse_blueprints.example.blueprint_set
}
```

**Note**: `offset` and `limit` are NOT exposed to users - pagination is handled automatically.

### Field Mapping

| Blueprint Field | Schema Field | Type | Description |
|----------------|--------------|------|-------------|
| BlueprintId | blueprint_id | string | Blueprint ID (unique identifier) |
| DisplayTitle | display_title | string | Display title |
| DisplayVersion | display_version | string | Display version |
| Description | description | string | Description (nullable) |
| OsName | os_name | string | OS name |
| Platform | platform | string | OS platform |
| PlatformType | platform_type | string | Platform type (LINUX_UNIX/WINDOWS) |
| BlueprintType | blueprint_type | string | Blueprint type (APP_OS/PURE_OS/DOCKER/PRIVATE/SHARED) |
| ImageUrl | image_url | string | Blueprint image URL |
| RequiredSystemDiskSize | required_system_disk_size | int | Required system disk size (GB) |
| BlueprintState | blueprint_state | string | Blueprint state |
| CreatedTime | created_time | string | Creation time (ISO 8601, nullable) |
| BlueprintName | blueprint_name | string | Blueprint name |
| SupportAutomationTools | support_automation_tools | bool | Supports automation tools |
| RequiredMemorySize | required_memory_size | int | Required memory size (GB) |
| ImageId | image_id | string | CVM image ID (nullable) |
| CommunityUrl | community_url | string | Community URL |
| GuideUrl | guide_url | string | Guide URL |
| SceneIdSet | scene_id_set | []string | Associated scene IDs (nullable) |
| DockerVersion | docker_version | string | Docker version (nullable) |

### Filter Support

The datasource will support all API filters:

| Filter Name | Description | Type | Example Values |
|------------|-------------|------|----------------|
| blueprint-id | Blueprint ID | String | lhbp-xxx |
| blueprint-type | Blueprint type | String | APP_OS, PURE_OS, DOCKER, PRIVATE, SHARED |
| platform-type | Platform type | String | LINUX_UNIX, WINDOWS |
| blueprint-name | Blueprint name | String | Ubuntu Server 20.04 LTS |
| blueprint-state | Blueprint state | String | ONLINE, OFFLINE |
| scene-id | Scene ID | String | scene-xxx |

**Constraint**: Cannot use both `blueprint_ids` and `filters` simultaneously (API limitation).

---

## Implementation Plan

### Files to Create

1. **Data Source**: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_blueprints.go`
   - Schema definition with all Blueprint fields
   - Read function with pagination logic
   - Filter and ID list support

2. **Service Layer**: Add method to `tencentcloud/services/lighthouse/service_tencentcloud_lighthouse.go`
   - `DescribeLighthouseBlueprintsByFilter()` method
   - Automatic pagination handling (offset/limit loop)

3. **Test**: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_blueprints_test.go`
   - Basic query test
   - Filter test (by type, platform)
   - ID list test

4. **Documentation**: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_blueprints.md`
   - Usage examples
   - Argument reference
   - Attribute reference

5. **Website Documentation**: Auto-generated via `make doc`
   - `website/docs/d/lighthouse_blueprints.html.markdown`

### Code Structure

```go
// Service layer method
func (me *LightHouseService) DescribeLighthouseBlueprintsByFilter(
    ctx context.Context, 
    param map[string]interface{},
) (blueprints []*lighthouse.Blueprint, errRet error) {
    // Pagination loop with offset/limit
    // Returns all matching blueprints
}

// DataSource read function
func dataSourceTencentCloudLighthouseBlueprintsRead(
    d *schema.ResourceData, 
    meta interface{},
) error {
    // Build paramMap from schema
    // Call service layer
    // Map SDK response to schema
    // Set ID as hash of blueprint IDs
}
```

---

## Alternatives Considered

### Alternative 1: Expose offset/limit to users
**Rejected** because:
- Inconsistent with other datasources (bundle, zone, region don't expose them)
- Adds complexity for users who just want all results
- Pagination is an implementation detail, not a feature

### Alternative 2: Use existing instance_blueprint datasource
**Rejected** because:
- That datasource requires instance IDs (different API)
- Cannot query the general blueprint catalog
- Different use case (instance-specific vs. catalog browsing)

### Alternative 3: Don't implement filters, only IDs
**Rejected** because:
- API provides rich filtering capabilities
- Users need to discover blueprints by type/platform
- Reduces usability significantly

---

## Impact Analysis

### Breaking Changes
**None** - This is a new datasource with no effect on existing resources.

### Compatibility
- ✅ Backward compatible (new feature only)
- ✅ No changes to existing datasources or resources
- ✅ No state migration required

### Dependencies
- Requires SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324`
- No new external dependencies

### Testing Requirements
1. **Acceptance Tests**: Query blueprints with different filters
2. **Manual Testing**: Verify pagination works correctly with large result sets
3. **Documentation Review**: Ensure examples are clear and accurate

---

## Risk Assessment

### Technical Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| API rate limiting with large queries | Low | Medium | Pagination already handles this |
| Null fields causing crashes | Low | High | All fields have nil checks |
| Filter validation issues | Low | Medium | Use API's built-in validation |

### Operational Risks
- **None identified** - Read-only datasource with no state

---

## Success Criteria

### Functional Requirements
- ✅ Users can query all blueprints
- ✅ Users can filter by type, platform, state, name, scene
- ✅ Users can query specific blueprint IDs
- ✅ All Blueprint fields are returned correctly
- ✅ Pagination is transparent to users

### Non-Functional Requirements
- ✅ Code follows project conventions
- ✅ Tests pass (unit and acceptance)
- ✅ Documentation is complete and clear
- ✅ No linter errors or warnings

### Acceptance Criteria
1. Datasource can query all blueprints without errors
2. Filters work correctly (type, platform, state)
3. All Blueprint fields are mapped and returned
4. Tests pass with TF_ACC=1
5. Documentation includes working examples
6. Code passes `make fmt`, `make lint`, `make doc`

---

## Timeline Estimate

| Phase | Tasks | Estimated Time |
|-------|-------|---------------|
| **Implementation** | Service layer + datasource code | 1.5 hours |
| **Testing** | Write and run tests | 1 hour |
| **Documentation** | Write docs and generate website docs | 0.5 hours |
| **Review & Polish** | Code review, fix issues | 0.5 hours |
| **Total** | | **3.5 hours** |

---

## References

- **API Documentation**: https://cloud.tencent.com/document/product/1207/47689
- **SDK Models**: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324/models.go:422-488`
- **Similar Datasource**: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_bundle.go`
- **Project Conventions**: `openspec/project.md`

---

## Approval

- [ ] Design reviewed and approved
- [ ] No security concerns identified
- [ ] Impact assessment complete
- [ ] Ready for implementation

---

**Next Steps**: Create detailed tasks.md and spec.md files for implementation.
