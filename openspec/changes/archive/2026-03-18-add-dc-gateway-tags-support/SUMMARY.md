# Summary: Add Tags Support to DC Gateway Resource

**Change ID**: `add-dc-gateway-tags-support`  
**Status**: ✅ Validated - Ready for Review  
**Created**: 2026-02-02

---

## Quick Overview

This proposal adds comprehensive tags support to the `tencentcloud_dc_gateway` resource, enabling users to classify, organize, and manage Direct Connect Gateway resources using tags.

### What's Being Added
- ✅ `tags` parameter in resource schema (optional TypeMap)
- ✅ Tags set during gateway creation via `CreateDirectConnectGateway` API
- ✅ Tags read using universal tag service (`DescribeResourceTags`)
- ✅ Tags updatable using universal tag service (`ModifyTags`)
- ✅ Tags imported with resource
- ✅ Full backward compatibility (no breaking changes)

---

## The Problem

Currently, `tencentcloud_dc_gateway` doesn't support tags, preventing users from:
- Organizing and categorizing DC gateway resources
- Implementing cost allocation and chargeback
- Applying tag-based access control policies
- Automating operations based on resource tags

---

## The Solution

### API Strategy

Since the native DC gateway APIs have **limited tag support**:

| Operation | Native API Support | Solution |
|-----------|-------------------|----------|
| **Create** | ✅ Supports `Tags` parameter | Use `CreateDirectConnectGatewayRequest.Tags` |
| **Read** | ❌ No tag fields in response | Use universal `TagService.DescribeResourceTags()` |
| **Update** | ❌ No tag update API | Use universal `TagService.ModifyTags()` |
| **Delete** | N/A | Tags auto-deleted with gateway |

### Implementation Pattern

Following **proven patterns** from existing resources (CCN, CLB, VPC):

```go
// CREATE: Set tags during gateway creation
request.Tags = []*vpc.Tag{...}  // Native API support

// READ: Retrieve tags via universal service
tagService.DescribeResourceTags(ctx, "vpc", "dcg", region, id)

// UPDATE: Modify tags via universal service
tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
```

---

## Key Features

### 1. Schema Definition
```hcl
resource "tencentcloud_dc_gateway" "example" {
  name                = "my-gateway"
  network_type        = "VPC"
  network_instance_id = "vpc-123456"
  
  tags = {
    Environment = "production"
    Team        = "networking"
    Owner       = "ops-team"
  }
}
```

### 2. Full CRUD Support
- **Create**: Tags set immediately during gateway creation
- **Read**: Tags always fetched and displayed in state
- **Update**: Tags modified without gateway recreation
- **Delete**: Tags automatically removed with gateway

### 3. All Operations Supported
- ✅ Add new tags
- ✅ Modify existing tag values
- ✅ Remove tags
- ✅ Replace all tags
- ✅ Import gateway with tags

---

## Requirements Covered

9 comprehensive requirements defined in spec:

| ID | Requirement | Priority | Scenarios |
|----|-------------|----------|-----------|
| DCG-TAGS-001 | Schema Support | High | 1 |
| DCG-TAGS-002 | Create with Tags | High | 2 |
| DCG-TAGS-003 | Read Tags | High | 3 |
| DCG-TAGS-004 | Update Tags | High | 4 |
| DCG-TAGS-005 | Import with Tags | Medium | 1 |
| DCG-TAGS-006 | Validation & Constraints | Medium | 3 |
| DCG-TAGS-007 | Resource Name Format | High | 1 |
| DCG-TAGS-008 | Backward Compatibility | High | 2 |
| DCG-TAGS-009 | Error Handling | High | 3 |

**Total**: 9 requirements, 20 scenarios

---

## Implementation Tasks

18 tasks organized in 7 phases:

1. **Schema Definition** (2 tasks) - Add tags field and imports
2. **Create Enhancement** (3 tasks) - Set tags during creation
3. **Read Enhancement** (2 tasks) - Retrieve and display tags
4. **Update Enhancement** (4 tasks) - Enable tag modifications
5. **Testing** (4 tasks) - Comprehensive test coverage
6. **Documentation** (2 tasks) - User documentation and examples
7. **Code Quality** (3 tasks) - Formatting, linting, validation

**Estimated Effort**: ~3.5 days

---

## Benefits

### For Users
- **Resource Organization**: Tag-based categorization and search
- **Cost Management**: Track costs by project, team, or environment
- **Access Control**: Tag-based IAM policies
- **Automation**: Automated operations based on tags
- **Compliance**: Tag governance and policy enforcement

### For Provider
- **Consistency**: Aligns with other VPC resources (CCN, CLB, VPC)
- **Zero Breaking Changes**: Purely additive feature
- **Production Ready**: Uses established tag service infrastructure
- **Well Tested**: Proven pattern from existing resources

---

## Technical Details

### Resource Name Format
```
qcs::vpc:{region}:account:/dcg/{gateway-id}
```

### Tag Service Integration
```go
// Service type
ServiceType: "vpc"

// Resource type
ResourceType: "dcg"

// Example resource name
qcs::vpc:ap-guangzhou:account:/dcg/dcg-12345678
```

### Error Handling
- Transient errors: Retry with exponential backoff
- Fatal errors: Clear error messages with context
- Logging: Full operation context in logs

---

## Testing Strategy

### Acceptance Tests
1. **Create Test**: Create gateway with tags, verify tags set
2. **Read Test**: Read gateway, verify tags retrieved
3. **Update Test**: 
   - Add new tags
   - Modify existing tags
   - Remove tags
   - Replace all tags
4. **Import Test**: Import gateway, verify tags imported

### Manual Testing
- Verify tags visible in Tencent Cloud console
- Test with various tag combinations
- Verify tag service API compatibility

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Tag API rate limits | Medium | Existing retry logic handles this |
| Resource name format issues | High | Using proven `BuildTagResourceName()` |
| State inconsistency | Medium | Same pattern as CCN/CLB (battle-tested) |
| Backward compatibility | Critical | Tags are optional, fully backward compatible |

---

## Documentation Updates

### Resource Documentation
- Add `tags` to Argument Reference
- Add examples showing tag operations
- Document import behavior with tags

### Examples
1. Basic usage with tags
2. Adding tags to existing gateway
3. Updating tags
4. Importing gateway with tags

---

## Success Criteria

✅ **Functional**
- Users can set tags during creation
- Tags correctly read and displayed
- Tags updatable without recreation
- Tags work with import

✅ **Quality**
- All acceptance tests pass
- Code passes linting
- Documentation complete

✅ **Compatibility**
- No breaking changes
- Existing resources unaffected
- Optional feature

---

## References

- **API Documentation**: https://cloud.tencent.com/document/product/215/19192
- **Existing Pattern**: `tencentcloud/services/ccn/resource_tc_ccn.go`
- **Tag Service**: `tencentcloud/services/tag/service_tencentcloud_tag.go`
- **SDK Models**: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312/models.go`

---

## Next Steps

1. **Review**: Team review of proposal
2. **Approval**: Get stakeholder approval
3. **Implementation**: Execute tasks per `tasks.md`
4. **Testing**: Run full test suite
5. **Documentation**: Update user docs
6. **Release**: Include in next provider release

---

## Files Created

```
openspec/changes/add-dc-gateway-tags-support/
├── proposal.md          # This proposal document
├── tasks.md             # 18 implementation tasks
├── SUMMARY.md           # This summary
└── specs/
    └── dc-gateway-tags/
        └── spec.md      # 9 requirements, 20 scenarios
```

**Validation**: ✅ `openspec validate add-dc-gateway-tags-support --strict` PASSED

---

## Questions or Concerns?

Please review the detailed documents:
- `proposal.md` - Full technical proposal
- `tasks.md` - Detailed implementation plan
- `specs/dc-gateway-tags/spec.md` - Complete requirements

Ready for team review and approval! 🚀
