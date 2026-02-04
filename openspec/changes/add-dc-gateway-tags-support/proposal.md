# Proposal: Add Tags Support to DC Gateway Resource

**Change ID**: `add-dc-gateway-tags-support`  
**Status**: Proposal  
**Author**: AI Agent  
**Date**: 2026-02-02

## Executive Summary

Add tags support to the `tencentcloud_dc_gateway` resource to enable resource classification, management, and cost allocation. The implementation will use the `CreateDirectConnectGateway` API's `Tags` parameter for creation, and leverage the universal tag service for read, update, and delete operations since the native API does not support tag operations.

## Problem Statement

Currently, the `tencentcloud_dc_gateway` resource does not support tags, limiting users' ability to:
- Classify and organize DC gateway resources
- Implement cost allocation and billing analysis
- Apply tag-based access control policies
- Enable automated operations based on tags

## Background

### API Support Analysis

1. **CreateDirectConnectGateway** (https://cloud.tencent.com/document/product/215/19192)
   - **Supports**: `Tags` parameter (Array of Tag)
   - Can set tags during resource creation

2. **DescribeDirectConnectGateways**
   - **Does NOT support**: Reading tags from response
   - Cannot retrieve tags via native API

3. **ModifyDirectConnectGatewayAttribute**
   - **Does NOT support**: Tag modification
   - Cannot update tags via native API

### Solution Approach

Following the established pattern in the codebase (e.g., `resource_tc_ccn.go`, `resource_tc_clb_instance.go`):
- **Create**: Use `CreateDirectConnectGatewayRequest.Tags` field
- **Read**: Use universal `TagService.DescribeResourceTags()` with resource type `vpc:dcg`
- **Update**: Use universal `TagService.ModifyTags()` 
- **Delete**: Tags automatically deleted when resource is destroyed

## Proposed Solution

### Schema Changes

Add a new optional field to `tencentcloud_dc_gateway` schema:

```go
"tags": {
    Type:        schema.TypeMap,
    Optional:    true,
    Description: "Tag key-value pairs for the DC gateway. Multiple tags can be set.",
},
```

### Implementation Components

1. **Create Logic Enhancement**
   - Extract tags from schema
   - Convert to `[]*vpc.Tag` format
   - Set `request.Tags` before API call

2. **Read Logic Enhancement**
   - After reading gateway info, call `TagService.DescribeResourceTags()`
   - Resource type: `"vpc"`, Resource prefix: `"dcg"`
   - Set tags in Terraform state

3. **Update Logic Enhancement**
   - Detect tag changes using `d.HasChange("tags")`
   - Calculate diff using `svctag.DiffTags()`
   - Call `TagService.ModifyTags()` with resource name format: `qcs::vpc:{region}:account:/{resource-type}/{resource-id}`

4. **Delete Logic**
   - No changes needed (tags automatically removed when gateway is deleted)

### Code Patterns

Following established patterns from `resource_tc_ccn.go`:

**Create**:
```go
// After CreateDirectConnectGateway succeeds
if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
    tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
    tagService := svctag.NewTagService(tcClient)
    resourceName := tccommon.BuildTagResourceName("vpc", "dcg", tcClient.Region, d.Id())
    if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
        return err
    }
}
```

**Read**:
```go
tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
tagService := svctag.NewTagService(tcClient)
tags, err := tagService.DescribeResourceTags(ctx, "vpc", "dcg", tcClient.Region, d.Id())
if err != nil {
    return err
}
_ = d.Set("tags", tags)
```

**Update**:
```go
if d.HasChange("tags") {
    oldValue, newValue := d.GetChange("tags")
    replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
    tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
    tagService := svctag.NewTagService(tcClient)
    resourceName := tccommon.BuildTagResourceName("vpc", "dcg", tcClient.Region, d.Id())
    err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
    if err != nil {
        return err
    }
}
```

## Alternatives Considered

### Alternative 1: Wait for Native API Support
**Rejected**: No timeline for when native APIs will support tag operations. Users need this functionality now.

### Alternative 2: Tags in Create Only
**Rejected**: Not user-friendly. Users expect full CRUD support for tags as with other resources.

### Alternative 3: Custom Tag Implementation
**Rejected**: Reinventing the wheel. Universal tag service is production-ready and widely used.

## Benefits

1. **Consistency**: Aligns with other VPC resources (CCN, CLB, VPC, etc.)
2. **User Experience**: Enables standard Terraform tag management workflows
3. **Resource Management**: Improves resource organization and cost tracking
4. **Zero Breaking Changes**: Purely additive feature

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Tag API rate limits | Medium | Use existing retry logic and rate limiting |
| Resource name format mismatch | High | Follow established `BuildTagResourceName()` pattern |
| Inconsistent tag state | Medium | Use same pattern as CCN/CLB with proven reliability |

## Success Criteria

1. Users can set tags during DC gateway creation
2. Tags are correctly read and displayed in Terraform state
3. Tags can be updated without recreating the gateway
4. Tags are properly removed when gateway is destroyed
5. All acceptance tests pass
6. Documentation is complete and accurate

## Testing Strategy

1. **Unit Tests**: Schema validation
2. **Acceptance Tests**:
   - Create gateway with tags
   - Read and verify tags
   - Update tags (add, modify, remove)
   - Import gateway and verify tags
3. **Manual Tests**: Verify in Tencent Cloud console

## Documentation

Update `resource_tc_dc_gateway.md` with:
- New `tags` argument description
- Usage examples showing tag operations
- Import behavior with tags

## Timeline

- **Proposal**: 1 day
- **Implementation**: 1 day
- **Testing**: 1 day
- **Documentation**: 0.5 day
- **Total**: ~3.5 days

## Dependencies

- No external dependencies
- Leverages existing `TagService` infrastructure
- Compatible with current SDK version

## References

- CreateDirectConnectGateway API: https://cloud.tencent.com/document/product/215/19192
- Existing implementation: `tencentcloud/services/ccn/resource_tc_ccn.go`
- Tag Service: `tencentcloud/services/tag/service_tencentcloud_tag.go`
