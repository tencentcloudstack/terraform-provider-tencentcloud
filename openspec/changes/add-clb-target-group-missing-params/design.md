# Design Document: Add CLB Target Group Missing Parameters

## Architecture Overview

This change extends the existing `tencentcloud_clb_target_group` resource by adding 8 new optional parameters. The implementation follows the standard Terraform provider pattern: schema definition → CRUD operations → API integration.

### Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│ Terraform Resource: tencentcloud_clb_target_group           │
├─────────────────────────────────────────────────────────────┤
│ Schema (Existing)                                            │
│ - target_group_name, vpc_id, port                           │
│ - target_group_instances (deprecated)                        │
│ - type, protocol                                             │
├─────────────────────────────────────────────────────────────┤
│ Schema (NEW)                                                 │
│ - health_check (nested block)                               │
│ - schedule_algorithm (string)                               │
│ - tags (map)                                                 │
│ - weight (int)                                               │
│ - full_listen_switch (bool)                                  │
│ - keepalive_enable (bool)                                    │
│ - session_expire_time (int)                                  │
│ - ip_version (string)                                        │
└─────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│ Service Layer: ClbService                                    │
├─────────────────────────────────────────────────────────────┤
│ CreateTargetGroup() - Extended with new params              │
│ DescribeTargetGroups() - Read target group details          │
│ ModifyTargetGroup() - Update mutable attributes             │
│ DeleteTarget() - Delete target group (unchanged)            │
└─────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│ Tencent Cloud SDK                                            │
├─────────────────────────────────────────────────────────────┤
│ CreateTargetGroupRequest                                     │
│ - All new fields already supported in SDK                   │
│ DescribeTargetGroupsResponse                                 │
│ - Returns all target group attributes                       │
│ ModifyTargetGroupAttributeRequest                            │
│ - Supports updating some attributes                         │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

### Create Flow

```
User Config (HCL)
    │
    ▼
Schema Validation
    │
    ├─ Validate schedule_algorithm in [WRR, LEAST_CONN, IP_HASH]
    ├─ Validate weight in [0, 100]
    ├─ Validate session_expire_time in [30, 3600] or 0
    └─ Validate health check ranges
    │
    ▼
resourceTencentCloudClbTargetCreate()
    │
    ├─ Extract existing params (name, vpc, port, type, protocol)
    ├─ Extract health_check block → build TargetGroupHealthCheck
    ├─ Extract schedule_algorithm string
    ├─ Extract tags map → convert to []TagInfo
    ├─ Extract weight, full_listen_switch, keepalive_enable
    ├─ Extract session_expire_time, ip_version
    │
    ▼
ClbService.CreateTargetGroup()
    │
    ├─ Build CreateTargetGroupRequest
    ├─ Set all parameters
    ├─ Call API with retry logic
    │
    ▼
API Response
    │
    └─ Target Group ID returned
    │
    ▼
resourceTencentCloudClbTargetRead()
    │
    └─ Read and populate state (see Read Flow)
```

### Read Flow

```
resourceTencentCloudClbTargetRead()
    │
    ▼
ClbService.DescribeTargetGroups(id)
    │
    ▼
API Response: TargetGroupInfo
    │
    ├─ Set target_group_name
    ├─ Set vpc_id, port, type, protocol
    ├─ Flatten health_check → set nested block
    ├─ Set schedule_algorithm
    ├─ Flatten tags → set map
    ├─ Set weight, full_listen_switch
    ├─ Set keepalive_enable, session_expire_time
    └─ Set ip_version
    │
    ▼
State Updated
```

### Update Flow

```
Terraform Detects Changes
    │
    ▼
resourceTencentCloudClbTargetUpdate()
    │
    ├─ Check ForceNew parameters (schedule_algorithm, full_listen_switch, ip_version)
    │   └─ If changed → Terraform automatically recreates resource
    │
    ├─ Detect updatable parameter changes
    │   ├─ target_group_name (existing)
    │   ├─ port (existing)
    │   ├─ health_check (new - TBD if updatable)
    │   ├─ tags (new - likely updatable)
    │   ├─ weight (new - TBD)
    │   ├─ keepalive_enable (new - TBD)
    │   └─ session_expire_time (new - TBD)
    │
    ▼
ClbService.ModifyTargetGroup() or ModifyTargetGroupAttribute()
    │
    ├─ Call appropriate API for changed parameters
    ├─ Handle errors and retries
    │
    ▼
resourceTencentCloudClbTargetRead()
    │
    └─ Verify changes applied
```

## Schema Design

### Health Check Nested Block

```go
"health_check": {
    Type:     schema.TypeList,
    Optional: true,
    MaxItems: 1,
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "health_switch": {
                Type:        schema.TypeBool,
                Required:    true,
                Description: "Whether to enable health check.",
            },
            "protocol": {
                Type:         schema.TypeString,
                Optional:     true,
                Computed:     true,
                ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS", "PING", "CUSTOM", "GRPC"}, false),
                Description:  "Health check protocol. Valid for v2 target groups.",
            },
            "port": {
                Type:         schema.TypeInt,
                Optional:     true,
                Computed:     true,
                ValidateFunc: validation.IntBetween(1, 65535),
                Description:  "Health check port. Defaults to backend server port.",
            },
            "timeout": {
                Type:         schema.TypeInt,
                Optional:     true,
                Default:      2,
                ValidateFunc: validation.IntBetween(2, 30),
                Description:  "Health check timeout in seconds. Range: [2, 30]. Default: 2.",
            },
            // ... more fields ...
        },
    },
    Description: "Health check configuration.",
}
```

### Simple Parameters

```go
"schedule_algorithm": {
    Type:         schema.TypeString,
    Optional:     true,
    Computed:     true,
    ForceNew:     true,
    ValidateFunc: validation.StringInSlice([]string{"WRR", "LEAST_CONN", "IP_HASH"}, false),
    Description:  "Scheduling algorithm. Only valid for v2 target groups with HTTP/HTTPS/GRPC protocols. Valid values: WRR, LEAST_CONN, IP_HASH. Default: WRR.",
},

"tags": {
    Type:        schema.TypeMap,
    Optional:    true,
    Elem:        &schema.Schema{Type: schema.TypeString},
    Description: "Resource tags for the target group.",
},

"weight": {
    Type:         schema.TypeInt,
    Optional:     true,
    ValidateFunc: validation.IntBetween(0, 100),
    Description:  "Default backend server weight. Range: [0, 100]. Only valid for v2 target groups.",
},
```

## Type Conversion Helpers

### Health Check Expansion (Schema → SDK)

```go
func expandHealthCheck(d *schema.ResourceData) *clb.TargetGroupHealthCheck {
    if v, ok := d.GetOk("health_check"); !ok || len(v.([]interface{})) == 0 {
        return nil
    }
    
    hcList := v.([]interface{})
    hcMap := hcList[0].(map[string]interface{})
    
    hc := &clb.TargetGroupHealthCheck{}
    
    if v, ok := hcMap["health_switch"].(bool); ok {
        hc.HealthSwitch = helper.Bool(v)
    }
    
    if v, ok := hcMap["protocol"].(string); ok && v != "" {
        hc.Protocol = helper.String(v)
    }
    
    // ... map other fields ...
    
    return hc
}
```

### Health Check Flattening (SDK → Schema)

```go
func flattenHealthCheck(hc *clb.TargetGroupHealthCheck) []interface{} {
    if hc == nil {
        return nil
    }
    
    result := make(map[string]interface{})
    
    if hc.HealthSwitch != nil {
        result["health_switch"] = *hc.HealthSwitch
    }
    
    if hc.Protocol != nil {
        result["protocol"] = *hc.Protocol
    }
    
    // ... map other fields ...
    
    return []interface{}{result}
}
```

### Tags Conversion

```go
// Expansion: map[string]string → []*clb.TagInfo
func expandTags(tags map[string]interface{}) []*clb.TagInfo {
    tagInfos := make([]*clb.TagInfo, 0, len(tags))
    
    for k, v := range tags {
        tagInfo := &clb.TagInfo{
            TagKey:   helper.String(k),
            TagValue: helper.String(v.(string)),
        }
        tagInfos = append(tagInfos, tagInfo)
    }
    
    return tagInfos
}

// Flattening: []*clb.TagInfo → map[string]string
func flattenTags(tagInfos []*clb.TagInfo) map[string]string {
    tags := make(map[string]string, len(tagInfos))
    
    for _, tagInfo := range tagInfos {
        if tagInfo.TagKey != nil && tagInfo.TagValue != nil {
            tags[*tagInfo.TagKey] = *tagInfo.TagValue
        }
    }
    
    return tags
}
```

## Service Layer Signature

### Extended CreateTargetGroup Method

```go
func (me *ClbService) CreateTargetGroup(
    ctx context.Context,
    targetGroupName string,
    vpcId string,
    port uint64,
    targetGroupInstances []*clb.TargetGroupInstance,
    targetGroupType string,
    protocol string,
    healthCheck *clb.TargetGroupHealthCheck,      // NEW
    scheduleAlgorithm string,                     // NEW
    tags []*clb.TagInfo,                          // NEW
    weight *uint64,                               // NEW
    fullListenSwitch *bool,                       // NEW
    keepaliveEnable *bool,                        // NEW
    sessionExpireTime *uint64,                    // NEW
    ipVersion string,                             // NEW
) (targetGroupId string, err error) {
    
    request := clb.NewCreateTargetGroupRequest()
    
    // Set existing parameters
    request.TargetGroupName = &targetGroupName
    request.TargetGroupInstances = targetGroupInstances
    request.Port = &port
    
    if vpcId != "" {
        request.VpcId = &vpcId
    }
    if targetGroupType != "" {
        request.Type = &targetGroupType
    }
    if protocol != "" {
        request.Protocol = &protocol
    }
    
    // Set new parameters
    if healthCheck != nil {
        request.HealthCheck = healthCheck
    }
    if scheduleAlgorithm != "" {
        request.ScheduleAlgorithm = &scheduleAlgorithm
    }
    if len(tags) > 0 {
        request.Tags = tags
    }
    if weight != nil {
        request.Weight = weight
    }
    if fullListenSwitch != nil {
        request.FullListenSwitch = fullListenSwitch
    }
    if keepaliveEnable != nil {
        request.KeepaliveEnable = keepaliveEnable
    }
    if sessionExpireTime != nil {
        request.SessionExpireTime = sessionExpireTime
    }
    if ipVersion != "" {
        request.IpVersion = &ipVersion
    }
    
    // Execute request with retry logic
    err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        response, e := me.client.UseClbClient().CreateTargetGroup(request)
        if e != nil {
            return tccommon.RetryError(e)
        }
        if response.Response.TargetGroupId != nil {
            targetGroupId = *response.Response.TargetGroupId
        }
        return nil
    })
    
    return
}
```

## Update Strategy

### Research Needed: ModifyTargetGroupAttribute API

The `ModifyTargetGroupAttribute` API documentation needs review to determine which new parameters support in-place updates. Expected updatable parameters based on common patterns:

- `target_group_name` - ✓ Already supported
- `port` - ✓ Already supported
- `health_check` - Likely updatable
- `tags` - Likely updatable (standard for cloud resources)
- `weight` - Possibly updatable
- `keepalive_enable` - Possibly updatable
- `session_expire_time` - Possibly updatable

Parameters marked ForceNew (confirmed non-updatable):
- `schedule_algorithm` - Algorithm selection is a creation-time decision
- `full_listen_switch` - Full listener mode is structural
- `ip_version` - IP version is foundational

### Update Implementation Pattern

```go
func resourceTencentCloudClbTargetUpdate(d *schema.ResourceData, meta interface{}) error {
    // ... existing code ...
    
    // Handle updatable parameters
    needsUpdate := false
    request := clb.NewModifyTargetGroupAttributeRequest()
    request.TargetGroupId = &targetGroupId
    
    if d.HasChange("target_group_name") {
        request.TargetGroupName = helper.String(d.Get("target_group_name").(string))
        needsUpdate = true
    }
    
    if d.HasChange("port") {
        port := uint64(d.Get("port").(int))
        request.Port = &port
        needsUpdate = true
    }
    
    if d.HasChange("health_check") {
        // TBD: Check if API supports health check updates
        hc := expandHealthCheck(d)
        request.HealthCheck = hc
        needsUpdate = true
    }
    
    if d.HasChange("tags") {
        // Tags likely have separate update API
        // Use ModifyTargetGroupTags or similar
    }
    
    // ... handle other updatable parameters ...
    
    if needsUpdate {
        err := clbService.ModifyTargetGroup(ctx, request)
        if err != nil {
            return err
        }
    }
    
    return resourceTencentCloudClbTargetRead(d, meta)
}
```

## Testing Strategy

### Unit Tests

1. **Schema Validation Tests**
   - Test schedule_algorithm accepts only valid values
   - Test weight rejects values outside [0, 100]
   - Test session_expire_time validation
   - Test health check nested structure

### Integration Tests

1. **Basic Creation Tests**
   ```go
   func TestAccTencentCloudClbTargetGroup_HealthCheck(t *testing.T) {
       // Create v2 target group with TCP health check
       // Verify health check parameters are set
       // Verify plan is empty after apply
   }
   ```

2. **Advanced Features Test**
   ```go
   func TestAccTencentCloudClbTargetGroup_V2Advanced(t *testing.T) {
       // Create v2 HTTP target group with:
       // - schedule_algorithm = LEAST_CONN
       // - weight = 80
       // - session_expire_time = 1800
       // - keepalive_enable = true
       // - tags
       // Verify all parameters applied correctly
   }
   ```

3. **Full Listener Test**
   ```go
   func TestAccTencentCloudClbTargetGroup_FullListener(t *testing.T) {
       // Create full listener target group
       // Verify full_listen_switch = true
       // Verify port is not required
   }
   ```

4. **Update Tests**
   ```go
   func TestAccTencentCloudClbTargetGroup_Update(t *testing.T) {
       // Create target group
       // Update updatable parameters (name, port, tags, etc.)
       // Verify updates applied without recreation
       // Update ForceNew parameter
       // Verify resource is recreated
   }
   ```

5. **Import Test**
   ```go
   func TestAccTencentCloudClbTargetGroup_Import(t *testing.T) {
       // Create target group with all new parameters
       // Import by ID
       // Verify all parameters imported correctly
   }
   ```

## Error Handling

### Validation Errors

```go
// Example: Schedule algorithm validation
func validateScheduleAlgorithm(val interface{}, key string) (warns []string, errs []error) {
    v := val.(string)
    validAlgorithms := []string{"WRR", "LEAST_CONN", "IP_HASH"}
    
    for _, valid := range validAlgorithms {
        if v == valid {
            return
        }
    }
    
    errs = append(errs, fmt.Errorf(
        "%q must be one of %v, got: %s",
        key, validAlgorithms, v,
    ))
    return
}
```

### API Error Handling

Common error scenarios:
1. **Invalid parameter combination** (e.g., full_listen_switch=true with port set)
   - API returns error
   - Terraform displays API error message
   - User fixes configuration

2. **Version/protocol mismatch** (e.g., schedule_algorithm on v1 target group)
   - API ignores or rejects parameter
   - Document constraints clearly to prevent user confusion

3. **Update not supported** for a parameter
   - If parameter is updatable but API rejects: Surface error to user
   - If parameter is ForceNew: Terraform handles recreation automatically

## Backward Compatibility

### Compatibility Matrix

| Scenario | Impact | Mitigation |
|----------|--------|------------|
| Existing configs without new params | ✓ No impact | All params optional |
| State migration | ✓ No migration needed | New params simply absent in old state |
| Provider upgrade | ✓ Transparent | No user action required |
| Import existing resources | ⚠️ New params read | Import populates new parameters |

### Migration Path

Users can adopt new parameters incrementally:
1. Upgrade provider version
2. Optionally add new parameters to existing resources
3. Apply changes (may trigger recreation for ForceNew params)

No breaking changes are introduced. Existing resources continue working as before.

## Performance Considerations

- **API Calls**: One additional read may occur if new parameters are added (same as current behavior)
- **State Size**: Minimal increase due to new fields
- **Validation**: Negligible overhead from parameter validation

## Security Considerations

- **Tags**: May be used for access control and cost allocation - document tag security best practices
- **Health Check**: Credentials (if any) should not be stored in health check configuration
- **IP Version**: Ensure IPv6 security posture is documented

## Documentation Requirements

### Resource Documentation

1. **Parameter Reference**
   - Full description of each new parameter
   - Valid values and ranges
   - Default values
   - Version/protocol constraints
   - ForceNew indicators

2. **Examples**
   - Basic v1 target group with health check
   - v2 target group with HTTP and advanced features
   - Full listener target group example
   - IPv6 target group example

3. **Import**
   - Clarify that all parameters are imported

4. **Constraints Matrix**
   - Table showing parameter support by version/protocol

### User Migration Guide

Since there are no breaking changes, migration guide is not required. However, a blog post or changelog highlighting new features would be beneficial.

## Open Questions

1. **Update Support**: Which new parameters can be updated in-place via `ModifyTargetGroupAttribute`?
   - **Action**: Review API documentation
   - **Impact**: Determines ForceNew flag for each parameter

2. **Health Check Update**: Can health check configuration be modified after creation?
   - **Action**: Test with API or review documentation
   - **Impact**: Affects update implementation

3. **Tags Update**: Is there a separate tags update API (ModifyTargetGroupTags)?
   - **Action**: Check SDK for tag-specific methods
   - **Impact**: May require separate update logic for tags

4. **Full Listener + Port**: Should Terraform validation prevent setting both full_listen_switch=true and port?
   - **Action**: Test API behavior
   - **Impact**: May need custom validation logic

5. **Parameter Interaction**: Are there other parameter combinations that are invalid?
   - **Action**: Review API error responses during testing
   - **Impact**: May need additional validation

## References

- [CreateTargetGroup API](https://cloud.tencent.com/document/product/214/40559)
- [ModifyTargetGroupAttribute API](https://cloud.tencent.com/document/product/214/40558)
- [DescribeTargetGroups API](https://cloud.tencent.com/document/product/214/40562)
- [TargetGroupHealthCheck Structure](https://cloud.tencent.com/document/api/214/30694#TargetGroupHealthCheck)
- Current Implementation: `tencentcloud/services/clb/resource_tc_clb_target_group.go`
- Service Layer: `tencentcloud/services/clb/service_tencentcloud_clb.go`
