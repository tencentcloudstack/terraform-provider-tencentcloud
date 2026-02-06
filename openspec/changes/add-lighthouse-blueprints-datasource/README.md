# Add Lighthouse Blueprints DataSource

**Change ID**: `add-lighthouse-blueprints-datasource`  
**Status**: Proposal Phase  
**Type**: New DataSource  
**Estimated Time**: 3.5 hours

---

## Quick Summary

Add a new datasource `tencentcloud_lighthouse_blueprints` to query Lighthouse blueprint (镜像) information using the `DescribeBlueprints` API.

**Key Features**:
- Query all blueprints or filter by type, platform, state, name, scene
- Query specific blueprints by ID
- Automatic pagination (offset/limit hidden from users)
- All 18+ Blueprint fields mapped
- Follows project conventions

---

## What's Being Added

### New Datasource: `tencentcloud_lighthouse_blueprints`

```hcl
# Query all blueprints
data "tencentcloud_lighthouse_blueprints" "all" {
}

# Filter by platform type
data "tencentcloud_lighthouse_blueprints" "linux" {
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
}

# Query specific IDs
data "tencentcloud_lighthouse_blueprints" "specific" {
  blueprint_ids = ["lhbp-xxx", "lhbp-yyy"]
}
```

---

## Files to Create

1. **Data Source**: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_blueprints.go`
2. **Service Method**: Add to `tencentcloud/services/lighthouse/service_tencentcloud_lighthouse.go`
3. **Test**: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_blueprints_test.go`
4. **Documentation**: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_blueprints.md`
5. **Website Docs**: Auto-generated via `make doc`

---

## Design Highlights

### 1. Hidden Pagination
**User doesn't specify offset/limit** - the service layer automatically loops through all pages:

```go
// Service layer handles this internally
var offset int64 = 0
var pageSize int64 = 100
for {
    request.Offset = &offset
    request.Limit = &pageSize
    // Fetch page, append results, check if more...
}
```

### 2. Comprehensive Filtering
Supports all API filters:
- `blueprint-id`: Filter by ID
- `blueprint-type`: APP_OS, PURE_OS, DOCKER, PRIVATE, SHARED
- `platform-type`: LINUX_UNIX, WINDOWS
- `blueprint-name`: Filter by name
- `blueprint-state`: Filter by state
- `scene-id`: Filter by scene

### 3. Complete Field Mapping
All 18+ Blueprint fields from SDK are mapped:
- Basic info: ID, name, title, version, description
- OS info: os_name, platform, platform_type
- Type info: blueprint_type, blueprint_state
- Requirements: required_system_disk_size, required_memory_size
- Metadata: image_url, created_time, support_automation_tools
- Optional: image_id, community_url, guide_url, scene_id_set, docker_version

---

## API Details

**API**: `DescribeBlueprints`  
**Documentation**: https://cloud.tencent.com/document/product/1207/47689  
**SDK Package**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324`

**Request Parameters**:
- `BlueprintIds` (optional): List of blueprint IDs
- `Filters` (optional): Filter list (cannot use with BlueprintIds)
- `Offset` (internal): Pagination offset
- `Limit` (internal): Page size

**Response**:
- `TotalCount`: Total number of blueprints
- `BlueprintSet`: List of Blueprint objects

---

## Implementation Plan

### Phase 1: Service Layer (30 min)
- Add `DescribeLighthouseBlueprintsByFilter()` method
- Implement pagination loop
- Handle BlueprintIds and Filters parameters

### Phase 2: DataSource Schema (45 min)
- Define input arguments (blueprint_ids, filters)
- Define output schema (blueprint_set with all fields)
- Add result_output_file support

### Phase 3: Read Function (45 min)
- Parse inputs
- Call service layer
- Map response to schema
- Set resource ID

### Phase 4: Testing (1 hour)
- Basic query test
- Filter tests (platform, type, multiple)
- ID query test
- Compile and run tests

### Phase 5: Documentation (30 min)
- Write usage examples
- Document arguments and attributes
- Generate website docs

### Phase 6: Quality & Integration (15 min)
- Format code
- Run linters
- Compile provider
- Register datasource

### Phase 7: Validation (15 min)
- Run acceptance tests
- Manual testing with Terraform CLI
- Verify pagination and filters

---

## Breaking Changes

**None** - This is a new datasource, no changes to existing resources.

---

## Testing Strategy

### Acceptance Tests
1. ✅ Query all blueprints
2. ✅ Filter by platform-type
3. ✅ Filter by blueprint-type
4. ✅ Multiple filters
5. ✅ Query by IDs

### Manual Tests
- Large result sets (>100 blueprints) - verify pagination
- Invalid filters - verify error handling
- Non-existent IDs - verify graceful handling

---

## Documentation Examples

### Example 1: Query All Blueprints
```hcl
data "tencentcloud_lighthouse_blueprints" "all" {
}

output "blueprint_count" {
  value = length(data.tencentcloud_lighthouse_blueprints.all.blueprint_set)
}
```

### Example 2: Linux Application Blueprints
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

output "linux_apps" {
  value = [
    for bp in data.tencentcloud_lighthouse_blueprints.linux_apps.blueprint_set :
    {
      id    = bp.blueprint_id
      name  = bp.blueprint_name
      type  = bp.blueprint_type
    }
  ]
}
```

### Example 3: Specific Blueprints
```hcl
data "tencentcloud_lighthouse_blueprints" "my_blueprints" {
  blueprint_ids = ["lhbp-f1lkcd41", "lhbp-ab123456"]
}
```

---

## Success Criteria

### Functional
- ✅ Query all blueprints without filters
- ✅ Filter by type, platform, state, name, scene
- ✅ Query specific blueprint IDs
- ✅ All Blueprint fields returned
- ✅ Pagination transparent to users

### Technical
- ✅ Follows project conventions
- ✅ All tests pass
- ✅ Code formatted and linted
- ✅ Documentation complete
- ✅ No nil pointer issues

---

## References

- **Proposal**: `proposal.md` - Detailed design and rationale
- **Tasks**: `tasks.md` - 62 implementation tasks
- **Spec**: `specs/lighthouse-blueprints-datasource/spec.md` - Detailed requirements
- **API Docs**: https://cloud.tencent.com/document/product/1207/47689
- **SDK Source**: `vendor/.../lighthouse/v20200324/models.go:422-488`
- **Similar Datasource**: `data_source_tc_lighthouse_bundle.go`

---

## Next Steps

1. **Review Proposal**: Read `proposal.md` for design details
2. **Review Tasks**: Check `tasks.md` for implementation steps
3. **Review Spec**: Read spec for detailed requirements
4. **Validate**: Run `openspec validate add-lighthouse-blueprints-datasource --strict`
5. **Implement**: Run `openspec apply add-lighthouse-blueprints-datasource` when ready

---

**Status**: ✅ Proposal Complete - Ready for Review
