# Proposal: Add Missing Parameters to CLB Target Group Resource

## Overview

Add support for all currently missing parameters from the Tencent Cloud `CreateTargetGroup` API to the `tencentcloud_clb_target_group` Terraform resource. This change ensures feature parity with the cloud API and enables users to configure advanced target group features.

## Motivation

The current `tencentcloud_clb_target_group` resource only supports a subset of parameters available in the [CreateTargetGroup API](https://cloud.tencent.com/document/product/214/40559). Users cannot configure:

- Health check settings
- Scheduling algorithms (for v2 target groups with HTTP/HTTPS/GRPC protocols)
- Resource tags
- Default backend server weight
- Full listener target group mode
- Keep-alive connection settings (for HTTP/HTTPS target groups)
- Session persistence
- IP version
- SNAT (source IP replacement)

This limitation forces users to manually configure these features through the console or API, reducing Infrastructure as Code coverage.

## Current State

### Supported Parameters
- `target_group_name` - Target group name
- `vpc_id` - VPC ID
- `port` - Default port
- `target_group_instances` - Backend server bindings (deprecated)
- `type` - Target group type (v1/v2)
- `protocol` - Backend forwarding protocol (required for v2)

### Missing Parameters
1. **health_check** (object) - Health check configuration
2. **schedule_algorithm** (string) - Scheduling algorithm (v2 only, HTTP/HTTPS/GRPC)
3. **tags** (map) - Resource tags
4. **weight** (int) - Default backend server weight [0-100] (v2 only)
5. **full_listen_switch** (bool) - Full listener target group flag (v2 only)
6. **keepalive_enable** (bool) - Keep-alive connection enable (HTTP/HTTPS only)
7. **session_expire_time** (int) - Session persistence time in seconds [30-3600] (v2 only, HTTP/HTTPS/GRPC)
8. **ip_version** (string) - IP version type
9. **snat_enable** (bool) - Enable SNAT (source IP replacement)

## Proposed Solution

Add schema definitions for all missing parameters with appropriate:
- Type constraints
- Validation rules
- Computed flags where applicable
- Force new flags for immutable parameters
- Documentation describing version/protocol constraints

### Parameter Classification

#### V1 and V2 Target Groups
- `health_check` - Both versions support health checks
- `tags` - Both versions support tags
- `ip_version` - Both versions support IP version

#### V2-Only Parameters
- `schedule_algorithm` - v2 only, HTTP/HTTPS/GRPC protocols
- `weight` - v2 only
- `full_listen_switch` - v2 only
- `session_expire_time` - v2 only, HTTP/HTTPS/GRPC protocols

#### Protocol-Specific Parameters
- `keepalive_enable` - HTTP/HTTPS target groups only
- `schedule_algorithm` - HTTP/HTTPS/GRPC protocols only (v2)
- `session_expire_time` - HTTP/HTTPS/GRPC protocols only (v2)

### Implementation Approach

1. **Schema Extension**: Add new parameters to resource schema with proper validation
2. **Create Logic**: Pass new parameters to `CreateTargetGroup` API call
3. **Read Logic**: Read back new parameters from `DescribeTargetGroups` API response
4. **Update Logic**: Determine which parameters support updates via `ModifyTargetGroupAttribute` API
5. **Documentation**: Update resource documentation with examples and constraints
6. **Tests**: Add acceptance tests covering new parameter combinations

### API Compatibility

The SDK already supports all these fields in `CreateTargetGroupRequest`:
- `HealthCheck *TargetGroupHealthCheck`
- `ScheduleAlgorithm *string`
- `Tags []*TagInfo`
- `Weight *uint64`
- `FullListenSwitch *bool`
- `KeepaliveEnable *bool`
- `SessionExpireTime *uint64`
- `IpVersion *string`

No SDK upgrade is required.

## Impact Analysis

### User Impact
- **Positive**: Users gain access to advanced target group features through Terraform
- **Backward Compatible**: All new parameters are optional; existing configurations continue working
- **No Breaking Changes**: Existing resources are not affected

### Code Impact
- **Files Modified**:
  - `tencentcloud/services/clb/resource_tc_clb_target_group.go` - Schema and CRUD logic
  - `tencentcloud/services/clb/resource_tc_clb_target_group.md` - Source documentation
  - `tencentcloud/services/clb/service_tencentcloud_clb.go` - Service layer method signatures
  - `website/docs/r/clb_target_group.html.markdown` - Generated documentation

### Testing Requirements
- Unit tests for schema validation
- Acceptance tests for:
  - v1 target groups with health checks and tags
  - v2 target groups with all advanced features
  - Protocol-specific features (HTTP/HTTPS with keepalive, session persistence)
  - Full listener target group mode

## Risks and Mitigations

### Risk: Parameter Validation Complexity
**Impact**: Some parameters are only valid for specific type/protocol combinations

**Mitigation**: 
- Add clear documentation about constraints
- Implement custom validation where needed
- Rely on API validation as final check

### Risk: Update Support Uncertainty
**Impact**: Some parameters may not be modifiable after creation

**Mitigation**:
- Research `ModifyTargetGroupAttribute` API to determine updatable fields
- Mark non-updatable fields with `ForceNew: true`
- Document which parameters require resource recreation

## Alternatives Considered

### Alternative 1: Partial Implementation
Add only the most commonly used parameters (e.g., health_check, tags)

**Rejected**: Would leave gaps in API coverage and require future incremental changes

### Alternative 2: Separate Resources
Create separate resources for v1 and v2 target groups

**Rejected**: Would break existing user configurations and increase maintenance burden

## Success Criteria

1. All 8 missing parameters are supported in the resource schema
2. Parameters work correctly with v1 and v2 target groups according to API constraints
3. Resource documentation clearly explains version/protocol requirements
4. Acceptance tests pass for all parameter combinations
5. No breaking changes to existing configurations

## Timeline

- **Proposal**: 1 day
- **Implementation**: 2-3 days
- **Testing**: 1-2 days
- **Documentation**: 1 day
- **Total**: ~5-7 days

## References

- [CreateTargetGroup API Documentation](https://cloud.tencent.com/document/product/214/40559)
- [ModifyTargetGroupAttribute API](https://cloud.tencent.com/document/product/214/40558)
- [DescribeTargetGroups API](https://cloud.tencent.com/document/product/214/40562)
- Current resource: `tencentcloud/services/clb/resource_tc_clb_target_group.go`
