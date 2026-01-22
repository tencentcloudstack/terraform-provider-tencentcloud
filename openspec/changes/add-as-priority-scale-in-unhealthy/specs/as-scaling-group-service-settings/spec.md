# Spec Delta: AS Scaling Group Service Settings Configuration

**Capability**: `as-scaling-group-service-settings`  
**Change ID**: `add-as-priority-scale-in-unhealthy`  
**Type**: MODIFIED

---

## ADDED Requirements

### REQ-AS-SG-SS-010: Priority Scale-In for Unhealthy Instances

The `tencentcloud_as_scaling_group` resource SHALL support configuring priority for removing unhealthy instances during scale-in operations through the `priority_scale_in_unhealthy` parameter.

**Rationale**: Users need to control whether unhealthy instances are prioritized for removal during scale-in to optimize resource utilization and maintain service quality.

#### Scenario: Enable Priority Scale-In for Unhealthy Instances

**Given** a user wants to ensure unhealthy instances are removed first during scale-in  
**When** the user sets `priority_scale_in_unhealthy = true` in the scaling group configuration  
**Then** the Auto Scaling service SHALL prioritize removing unhealthy instances over healthy ones during scale-in operations

**Example Configuration**:
```hcl
resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name   = "example-asg"
  configuration_id     = "asc-xxxxx"
  max_size             = 10
  min_size             = 1
  vpc_id               = "vpc-xxxxx"
  subnet_ids           = ["subnet-xxxxx"]
  
  # Enable priority removal of unhealthy instances
  priority_scale_in_unhealthy = true
}
```

**Acceptance Criteria**:
- Parameter is accepted as boolean value
- Value is correctly sent to TencentCloud API in `ServiceSettings.PriorityScaleInUnhealthy`
- Terraform state reflects the configured value
- Parameter can be updated without recreating the resource

#### Scenario: Disable Priority Scale-In for Unhealthy Instances

**Given** a user wants standard scale-in behavior without prioritizing unhealthy instances  
**When** the user sets `priority_scale_in_unhealthy = false` in the scaling group configuration  
**Then** the Auto Scaling service SHALL use its default scale-in policy without special treatment for unhealthy instances

**Example Configuration**:
```hcl
resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name   = "example-asg"
  configuration_id     = "asc-xxxxx"
  max_size             = 10
  min_size             = 1
  vpc_id               = "vpc-xxxxx"
  subnet_ids           = ["subnet-xxxxx"]
  
  # Use standard scale-in behavior
  priority_scale_in_unhealthy = false
}
```

**Acceptance Criteria**:
- Parameter value `false` is correctly handled
- Value is correctly sent to TencentCloud API
- Behavior is consistent with other ServiceSettings boolean parameters

#### Scenario: Omit Priority Scale-In Configuration

**Given** a user does not specify the `priority_scale_in_unhealthy` parameter  
**When** the scaling group is created or updated  
**Then** the API default value SHALL apply (provider does not override)

**Example Configuration**:
```hcl
resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name   = "example-asg"
  configuration_id     = "asc-xxxxx"
  max_size             = 10
  min_size             = 1
  vpc_id               = "vpc-xxxxx"
  subnet_ids           = ["subnet-xxxxx"]
  
  # priority_scale_in_unhealthy not specified - API default applies
}
```

**Acceptance Criteria**:
- Parameter is optional
- Omitting the parameter does not cause errors
- API default behavior is preserved

#### Scenario: Update Priority Scale-In Setting

**Given** an existing scaling group with `priority_scale_in_unhealthy = false`  
**When** the user updates the configuration to `priority_scale_in_unhealthy = true`  
**Then** the scaling group SHALL be updated in-place without recreation  
**AND** the new setting SHALL apply to subsequent scale-in operations

**Acceptance Criteria**:
- In-place update is performed (no ForceNew)
- ModifyAutoScalingGroup API is called with updated ServiceSettings
- Terraform plan shows only the changed attribute
- No unnecessary resource recreation

#### Scenario: Read and Import Existing Configuration

**Given** a scaling group exists in TencentCloud with PriorityScaleInUnhealthy configured  
**When** the resource is imported or refreshed  
**Then** the `priority_scale_in_unhealthy` value SHALL be correctly read from the API response  
**AND** stored in Terraform state

**Acceptance Criteria**:
- Read operation correctly retrieves the value from DescribeAutoScalingGroups API
- Value is set in Terraform state via `d.Set()`
- Import functionality works correctly
- Nil values are handled gracefully

---

## MODIFIED Requirements

None. This change adds a new parameter without modifying existing requirements.

---

## REMOVED Requirements

None. This is a pure addition with no deprecations or removals.

---

## Cross-References

### Related Capabilities
- **as-scaling-group-configuration**: Base scaling group configuration capability
- **as-scaling-group-health-check**: Health check configuration (unhealthy instance detection)

### API References
- **TencentCloud API**: Auto Scaling Service v20180419
  - CreateAutoScalingGroup: https://cloud.tencent.com/document/product/377/20440
  - DescribeAutoScalingGroups: https://cloud.tencent.com/document/product/377/20438
  - ModifyAutoScalingGroup: https://cloud.tencent.com/document/product/377/20433

### Implementation Files
- `tencentcloud/services/as/resource_tc_as_scaling_group.go` - Main resource implementation
- `tencentcloud/services/as/resource_tc_as_scaling_group_test.go` - Acceptance tests
- `tencentcloud/services/as/resource_tc_as_scaling_group.md` - Documentation

---

## Notes

### Consistency with Existing ServiceSettings Parameters
This parameter follows the same pattern as existing ServiceSettings parameters:
- `replace_monitor_unhealthy` (bool)
- `replace_load_balancer_unhealthy` (bool)
- `scaling_mode` (string)
- `replace_mode` (string)
- `desired_capacity_sync_with_max_min_size` (bool)

### API Mapping
| Terraform Attribute | API Field | Type |
|---------------------|-----------|------|
| `priority_scale_in_unhealthy` | `ServiceSettings.PriorityScaleInUnhealthy` | *bool |

### Backward Compatibility
This change is fully backward compatible:
- Existing configurations without this parameter continue to work
- No state migration required
- No breaking changes to schema or behavior
