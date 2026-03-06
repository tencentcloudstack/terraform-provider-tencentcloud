# Design Document: TDMQ RabbitMQ VIP Instance Resource Tags

## Context

### Background
Tencent Cloud's TDMQ RabbitMQ service now supports resource tags through three API endpoints:
- `CreateRabbitMQVipInstance` - accepts `ResourceTags` parameter
- `DescribeRabbitMQVipInstances` - returns tags in `ClusterInfo.Tags`
- `ModifyRabbitMQVipInstance` - accepts `Tags` parameter for updates

### Current State
The Terraform provider resource `tencentcloud_tdmq_rabbitmq_vip_instance` does not expose any tag management functionality.

### Stakeholders
- Terraform users managing RabbitMQ VIP instances
- DevOps teams using tag-based cost allocation and resource organization
- Security teams implementing tag-based access control (CAM policies)

### Constraints
- Must maintain backward compatibility (no breaking changes)
- Must follow existing provider patterns for tag handling
- API has inconsistent naming: `ResourceTags` (create) vs `Tags` (modify/read)

## Goals / Non-Goals

### Goals
1. Enable users to set tags during instance creation
2. Allow users to view tags in Terraform state
3. Support tag updates through `terraform apply`
4. Support tag removal (setting to empty map)
5. Ensure state consistency with actual cloud resources

### Non-Goals
- Tag validation beyond API-level checks (key length, value constraints)
- Tag inheritance or propagation to related resources
- Special handling for system tags (if any exist)
- Migration of existing instances' tags from Tencent Cloud console

## Decisions

### Decision 1: Field Naming - `resource_tags` vs `tags`

**Choice**: Use `resource_tags` as the field name

**Rationale**:
- The API explicitly uses `ResourceTags` for the create operation, indicating these are resource-level tags
- This distinguishes from potential future support for other tag types
- Precedent exists in the provider (some resources use `resource_tags`)
- Aligns with API naming convention

**Alternatives Considered**:
1. **Use `tags`** - More concise and common in Terraform resources
   - **Rejected**: Could create confusion with future tag features; API uses "ResourceTags"
2. **Use `instance_tags`** - Domain-specific naming
   - **Rejected**: Less consistent with API; resource prefix already provides context

### Decision 2: Schema Type - TypeMap vs TypeList

**Choice**: Use `schema.TypeMap` with `Elem: &schema.Schema{Type: schema.TypeString}`

**Rationale**:
- Tags are inherently key-value pairs, map is the natural representation
- Consistent with tag handling in other provider resources (e.g., `resource_tc_instance`)
- Simpler user experience: `resource_tags = { "env" = "prod", "team" = "ops" }`
- Terraform automatically handles map diffs

**Alternatives Considered**:
1. **Use TypeList with nested blocks** - Provides more structure
   - **Rejected**: Overkill for simple key-value pairs; worse UX
2. **Use TypeSet** - Treats tags as unordered collection
   - **Rejected**: Maps already handle order insensitivity; less intuitive syntax

### Decision 3: Tag Conversion Pattern

**Choice**: Use inline conversion with `helper.String()` and iteration pattern

**Pattern**:
```go
// TF map → SDK array (Create/Update)
if v := helper.GetTags(d, "resource_tags"); len(v) > 0 {
    for tagKey, tagValue := range v {
        tag := tdmq.Tag{
            TagKey:   helper.String(tagKey),
            TagValue: helper.String(tagValue),
        }
        request.ResourceTags = append(request.ResourceTags, &tag)
    }
}

// SDK array → TF map (Read)
if len(rabbitmqVipInstance.ClusterInfo.Tags) > 0 {
    tags := make(map[string]string)
    for _, tag := range rabbitmqVipInstance.ClusterInfo.Tags {
        if tag.TagKey != nil && tag.TagValue != nil {
            tags[*tag.TagKey] = *tag.TagValue
        }
    }
    _ = d.Set("resource_tags", tags)
}
```

**Rationale**:
- Follows existing patterns in `resource_tc_instance.go` and `resource_tc_eip.go`
- `helper.GetTags()` already exists and handles nil checks
- Explicit nil checks prevent panics during Read
- Ignores nil tags gracefully (defensive programming)

**Alternatives Considered**:
1. **Extract to helper function** - Centralized conversion logic
   - **Rejected**: Tag structures vary across services (tdmq.Tag vs vpc.Tag); not worth abstraction for single use
2. **Use reflection-based conversion** - Generic approach
   - **Rejected**: Overkill; loses type safety; harder to debug

### Decision 4: Update Behavior

**Choice**: Full tag replacement on update

**Implementation**:
```go
if d.HasChange("resource_tags") {
    request.Tags = []*tdmq.Tag{} // Start fresh
    if v := helper.GetTags(d, "resource_tags"); len(v) > 0 {
        for tagKey, tagValue := range v {
            // ... append tags
        }
    }
    // Call ModifyRabbitMQVipInstance
}
```

**Rationale**:
- Terraform's state management expects full replacement semantics
- Simplifies logic (no need to calculate delta)
- API accepts full tag list, not incremental changes
- Aligns with how other resources handle similar updates

**Alternatives Considered**:
1. **Delta-based updates** - Only send changed tags
   - **Rejected**: API doesn't support partial updates; more complex; no benefit
2. **Delete-then-add** - Separate removal and addition
   - **Rejected**: API doesn't have separate tag delete endpoint; unnecessary complexity

### Decision 5: Empty Tags Handling

**Choice**: Sending empty array to API removes all tags

**Behavior**:
- User sets `resource_tags = {}` → API receives `Tags: []` → All tags removed
- User omits `resource_tags` field → Not included in update request → Tags unchanged

**Rationale**:
- Terraform semantics: explicit empty map means "no tags"
- API behavior: empty Tags array removes all tags
- Provides clear way to remove tags without resource recreation

## Data Model

### Terraform Schema
```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  # ... existing fields ...
  
  resource_tags = {
    "Environment" = "Production"
    "Team"        = "Platform"
    "CostCenter"  = "CC-1234"
  }
}
```

### API Request/Response Mapping

**Create API (CreateRabbitMQVipInstance)**:
```json
{
  "ResourceTags": [
    {"TagKey": "Environment", "TagValue": "Production"},
    {"TagKey": "Team", "TagValue": "Platform"}
  ]
}
```

**Read API (DescribeRabbitMQVipInstances)**:
```json
{
  "ClusterInfo": {
    "Tags": [
      {"TagKey": "Environment", "TagValue": "Production"},
      {"TagKey": "Team", "TagValue": "Platform"}
    ]
  }
}
```

**Update API (ModifyRabbitMQVipInstance)**:
```json
{
  "Tags": [
    {"TagKey": "Environment", "TagValue": "Staging"}
  ]
}
```

## Implementation Details

### Code Locations

**File**: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`

**Insertion Points**:
1. **Schema** (line ~125, after `cluster_version` field):
   ```go
   "resource_tags": {
       Optional:    true,
       Type:        schema.TypeMap,
       Elem:        &schema.Schema{Type: schema.TypeString},
       Description: "Instance resource tags. Key-value pairs for resource identification and management.",
   },
   ```

2. **Create Function** (line ~192, after `cluster_version` handling):
   ```go
   if v := helper.GetTags(d, "resource_tags"); len(v) > 0 {
       for tagKey, tagValue := range v {
           tag := tdmq.Tag{
               TagKey:   helper.String(tagKey),
               TagValue: helper.String(tagValue),
           }
           request.ResourceTags = append(request.ResourceTags, &tag)
       }
   }
   ```

3. **Read Function** (line ~306, after `cluster_version` set):
   ```go
   if rabbitmqVipInstance.ClusterInfo != nil && len(rabbitmqVipInstance.ClusterInfo.Tags) > 0 {
       tags := make(map[string]string)
       for _, tag := range rabbitmqVipInstance.ClusterInfo.Tags {
           if tag.TagKey != nil && tag.TagValue != nil {
               tags[*tag.TagKey] = *tag.TagValue
           }
       }
       _ = d.Set("resource_tags", tags)
   }
   ```

4. **Update Function** (line ~417, after `cluster_name` change handling):
   ```go
   if d.HasChange("resource_tags") {
       request.Tags = []*tdmq.Tag{}
       if v := helper.GetTags(d, "resource_tags"); len(v) > 0 {
           for tagKey, tagValue := range v {
               tag := tdmq.Tag{
                   TagKey:   helper.String(tagKey),
                   TagValue: helper.String(tagValue),
               }
               request.Tags = append(request.Tags, &tag)
           }
       }
       // ModifyRabbitMQVipInstance call already exists in cluster_name block
       // Will need to restructure to handle multiple changes
   }
   ```

### Update Function Restructuring

**Current Issue**: Update function only handles `cluster_name` changes in a single API call block.

**Solution**: Extract API call logic and call it when any supported field changes:

```go
func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
    // ... existing setup code ...
    
    needsUpdate := false
    
    if d.HasChange("cluster_name") {
        if v, ok := d.GetOk("cluster_name"); ok {
            request.ClusterName = helper.String(v.(string))
            needsUpdate = true
        }
    }
    
    if d.HasChange("resource_tags") {
        request.Tags = []*tdmq.Tag{}
        if v := helper.GetTags(d, "resource_tags"); len(v) > 0 {
            for tagKey, tagValue := range v {
                tag := tdmq.Tag{
                    TagKey:   helper.String(tagKey),
                    TagValue: helper.String(tagValue),
                }
                request.Tags = append(request.Tags, &tag)
            }
        }
        needsUpdate = true
    }
    
    if needsUpdate {
        request.InstanceId = &instanceId
        err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
            // ... existing API call ...
        })
        // ... existing error handling ...
    }
    
    return resourceTencentCloudTdmqRabbitmqVipInstanceRead(d, meta)
}
```

## Risks / Trade-offs

### Risk 1: API Naming Inconsistency
- **Issue**: Create uses `ResourceTags`, Modify uses `Tags`
- **Mitigation**: Verified with API documentation; no ambiguity in implementation
- **Impact**: Low - internal detail, users don't see this

### Risk 2: Large Tag Maps Performance
- **Issue**: Converting large maps to arrays could be slow
- **Likelihood**: Very low - cloud services typically limit tags to 50-100 per resource
- **Mitigation**: None needed; API enforces limits
- **Impact**: Negligible

### Risk 3: Tag Read Consistency
- **Issue**: Tags might not immediately appear after creation (eventual consistency)
- **Mitigation**: Existing retry logic in Read function handles this
- **Impact**: Low - covered by existing patterns

### Trade-off 1: Field Naming
- **Chosen**: `resource_tags` (more explicit)
- **Trade-off**: Slightly longer name vs clearer intent
- **Justification**: Clarity wins; precedent exists in provider

### Trade-off 2: Full vs Partial Updates
- **Chosen**: Full replacement
- **Trade-off**: Simpler code vs potential unnecessary API calls
- **Justification**: API doesn't support partial; Terraform expects full state

## Migration Plan

### For New Users
No migration needed - feature is opt-in (Optional field).

### For Existing Instances
**Scenario**: User has RabbitMQ instance created before this feature.

**Behavior**:
1. First `terraform plan` after upgrade shows no changes (no tags configured)
2. User adds `resource_tags = {...}` to configuration
3. `terraform plan` shows tag addition
4. `terraform apply` calls `ModifyRabbitMQVipInstance` to add tags

**State Migration**: Not required (new field is Optional, default behavior is safe).

### Rollback
If feature causes issues:
1. Users can remove `resource_tags` from configurations
2. Tags remain on cloud resources (API doesn't auto-delete)
3. Users can manually remove tags via Tencent Cloud console if needed

## Open Questions

### Q1: SDK Version Compatibility Issue (BLOCKING) ❌
**Status**: **SDK OUTDATED - BLOCKING IMPLEMENTATION**
**Discovered**: 2025-03-06 during implementation
**Issue**: The tencentcloud-sdk-go TDMQ module (v20200217) does NOT contain the Tags field in ModifyRabbitMQVipInstanceRequest, despite the official API documentation (updated 2026-01-14) clearly showing:
- `Tags.N` parameter for ModifyRabbitMQVipInstance
- `RemoveAllTags` parameter
- `EnableRiskWarning` parameter

**Current SDK Structure**:
```go
type ModifyRabbitMQVipInstanceRequest struct {
    InstanceId *string
    ClusterName *string
    Remark *string
    EnableDeletionProtection *bool
    // NO Tags field!
}
```

**Expected SDK Structure** (based on API docs):
```go
type ModifyRabbitMQVipInstanceRequest struct {
    InstanceId *string
    ClusterName *string
    Remark *string
    EnableDeletionProtection *bool
    RemoveAllTags *bool           // NEW - not in SDK
    Tags []*Tag                   // NEW - not in SDK
    EnableRiskWarning *bool       // NEW - not in SDK
}
```

**Impact**:
- ✅ CREATE operation: Works (ResourceTags field exists in CreateRabbitMQVipInstanceRequest)
- ✅ READ operation: Works (ClusterInfo.Tags exists in response)
- ❌ UPDATE operation: CANNOT be implemented with current SDK version

**Action Required**:
1. **Update vendor SDK**: Run `go get -u github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217`
2. **Verify new fields**: Check if Tags/RemoveAllTags/EnableRiskWarning are present
3. **Alternative**: If SDK update unavailable, implement Create+Read only, document Update limitation

**Workaround** (if SDK cannot be updated):
- Implement tags for Create and Read only
- Document that tag updates require resource recreation (ForceNew: true)
- Add this to schema: `ForceNew: true` for `resource_tags`

### Q2: Do system tags exist for RabbitMQ instances?
**Status**: To be verified with API behavior
**Action**: Test with real instance to see if API returns any system-managed tags
**Impact**: If yes, might need to filter system tags from state or document read-only tags

### Q2: Tag key/value validation rules?
**Status**: Documented in API reference
**Action**: Document common constraints in resource documentation
**Impact**: Users will get API-level errors; Terraform validation not needed (API is source of truth)

### Q3: Does tag update require instance restart?
**Status**: To be verified with API documentation
**Action**: Check if `ModifyRabbitMQVipInstance` with only tag changes triggers restart
**Impact**: If yes, document this behavior to warn users

---

**Decision Authority**: Development team
**Last Updated**: 2025-03-06
**Reviewers**: To be assigned
