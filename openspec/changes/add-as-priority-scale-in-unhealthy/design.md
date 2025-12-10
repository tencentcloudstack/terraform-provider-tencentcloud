# Design: Add PriorityScaleInUnhealthy Parameter to AS Scaling Group

## Overview
This document describes the design for adding the `PriorityScaleInUnhealthy` parameter to the `tencentcloud_as_scaling_group` resource.

## Architecture Context

### Current State
The `tencentcloud_as_scaling_group` resource already supports most `ServiceSettings` parameters:
- `replace_monitor_unhealthy` → `ServiceSettings.ReplaceMonitorUnhealthy`
- `scaling_mode` → `ServiceSettings.ScalingMode`
- `replace_load_balancer_unhealthy` → `ServiceSettings.ReplaceLoadBalancerUnhealthy`
- `replace_mode` → `ServiceSettings.ReplaceMode`
- `desired_capacity_sync_with_max_min_size` → `ServiceSettings.DesiredCapacitySyncWithMaxMinSize`

### Proposed State
Add one more parameter to complete the `ServiceSettings` support:
- `priority_scale_in_unhealthy` → `ServiceSettings.PriorityScaleInUnhealthy`

## Design Decisions

### 1. Parameter Naming
**Decision**: Use `priority_scale_in_unhealthy` (snake_case with underscores)

**Rationale**:
- Follows existing Terraform naming conventions in this resource
- Consistent with other similar parameters like `replace_load_balancer_unhealthy`
- Snake_case is the standard for Terraform resource attributes

**Alternatives Considered**:
- `prioritize_unhealthy_scale_in`: Less clear about the direction
- `unhealthy_priority`: Too generic

### 2. Parameter Type
**Decision**: `schema.TypeBool`

**Rationale**:
- API expects a boolean value
- Consistent with other ServiceSettings boolean fields
- Simple true/false semantics align with the feature behavior

### 3. Optional vs Required
**Decision**: Optional with no default value

**Rationale**:
- Not all users need this feature
- Allows API defaults to apply when not specified
- Consistent with other optional ServiceSettings parameters
- No breaking changes to existing configurations

### 4. Implementation Pattern
**Decision**: Follow the exact pattern used for other ServiceSettings parameters

**Implementation Points**:
1. Schema definition alongside other ServiceSettings fields (lines 163-188)
2. Read from `d.Get()` in create function (line 329-333 pattern)
3. Include in `ServiceSettings` struct (line 345-351 pattern)
4. Set value in read function (line 469-486 pattern)
5. Change detection in update function (line 610-614 pattern)
6. Include in update `ServiceSettings` struct (line 627 pattern)

**Rationale**:
- Proven pattern already working for 5 similar parameters
- Minimal risk of introducing bugs
- Easy to review and maintain
- Consistent codebase

## Data Flow

### Create Flow
```
User Config (HCL)
    ↓
priority_scale_in_unhealthy: true
    ↓
d.Get("priority_scale_in_unhealthy").(bool)
    ↓
ServiceSettings.PriorityScaleInUnhealthy = &priorityScaleInUnhealthy
    ↓
CreateAutoScalingGroup API Request
    ↓
TencentCloud Auto Scaling Service
```

### Read Flow
```
TencentCloud Auto Scaling Service
    ↓
DescribeAutoScalingGroups API Response
    ↓
response.AutoScalingGroupSet[0].ServiceSettings.PriorityScaleInUnhealthy
    ↓
d.Set("priority_scale_in_unhealthy", value)
    ↓
Terraform State
```

### Update Flow
```
User Config Change (HCL)
    ↓
d.HasChange("priority_scale_in_unhealthy")
    ↓
d.Get("priority_scale_in_unhealthy").(bool)
    ↓
ServiceSettings.PriorityScaleInUnhealthy = &priorityScaleInUnhealthy
    ↓
ModifyAutoScalingGroup API Request
    ↓
TencentCloud Auto Scaling Service
```

## Error Handling

### Nil Value Handling
- In read operation: Check if `ServiceSettings` or `PriorityScaleInUnhealthy` is nil before accessing
- Use safe access pattern: `if v, ok := ...; ok { d.Set(...) }`

### API Errors
- Rely on existing retry and error handling mechanisms in the resource
- No special error handling needed for this parameter

## Testing Strategy

### Unit Tests
Not applicable - Terraform provider tests are primarily acceptance tests.

### Acceptance Tests
Add test configuration with the parameter in `resource_tc_as_scaling_group_test.go`:

```hcl
resource "tencentcloud_as_scaling_group" "scaling_group" {
  # ... existing config ...
  priority_scale_in_unhealthy = true
}
```

Test assertions:
- Verify parameter is set correctly after creation
- Verify parameter can be updated
- Verify parameter is read correctly from state

### Manual Testing
1. Create scaling group with `priority_scale_in_unhealthy = true`
2. Verify in TencentCloud console that setting is applied
3. Update to `false` and verify change
4. Verify import functionality works correctly

## Documentation Strategy

### Resource Documentation
Update `resource_tc_as_scaling_group.md`:

1. **Argument Reference**: Add to ServiceSettings section
   ```markdown
   * `priority_scale_in_unhealthy` - (Optional, Bool) Whether to enable priority for unhealthy instances during scale-in operations. If set to `true`, unhealthy instances will be removed first when scaling in.
   ```

2. **Example Usage**: Add to complete example
   ```hcl
   resource "tencentcloud_as_scaling_group" "example" {
     # ... other settings ...
     priority_scale_in_unhealthy = true
   }
   ```

### Generated Documentation
- Run `make doc` to generate final provider documentation
- Verify generated docs include the new parameter

## Compatibility

### Backward Compatibility
✅ **Fully Compatible**
- New optional parameter - existing configurations work unchanged
- No default value changes existing behavior
- No state migration required

### Forward Compatibility
✅ **No Concerns**
- API supports the parameter
- Standard boolean type unlikely to change
- Follows established patterns

## Security Considerations
None. This is a configuration parameter that affects scaling behavior, not security settings.

## Performance Considerations
None. Adding one boolean field has negligible performance impact.

## Monitoring and Observability
- Standard Terraform logging will capture parameter usage
- TencentCloud API logs will show parameter in requests
- No additional monitoring required

## Rollback Strategy
If issues arise:
1. Users can remove the parameter from configuration (reverts to API default)
2. Provider can be rolled back to previous version
3. No data loss risk - only affects future scaling decisions

## Future Considerations
This completes the `ServiceSettings` parameter support. No additional related parameters are currently known from the API documentation.
