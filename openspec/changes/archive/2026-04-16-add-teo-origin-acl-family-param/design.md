## Context

The `tencentcloud_teo_origin_acl` Terraform resource currently manages TEO (TencentCloud EdgeOne) origin ACL configuration with three parameters: `zone_id`, `l7_hosts`, and `l4_proxy_ids`. The resource uses the TE SDK v20220901 package with three main cloud API operations:

- **Create**: `EnableOriginACL` - enables origin ACL for a zone
- **Read**: `DescribeOriginACL` - retrieves current origin ACL configuration
- **Update**: `ModifyOriginACL` - modifies origin ACL entities
- **Delete**: `DisableOriginACL` - disables origin ACL for a zone

The cloud API supports an additional `OriginACLFamily` parameter (string type) in both `EnableOriginACL` and `ModifyOriginACL` operations, which specifies the control domain (gaz, mlc, emc, plat-gaz, plat-mlc, plat-emc) for the origin ACL configuration. This parameter is also returned in the `DescribeOriginACL` response via `OriginACLInfo.OriginACLFamily`.

Current implementation lacks support for this parameter, preventing users from configuring the ACL control domain through Terraform.

## Goals / Non-Goals

**Goals:**
- Add `origin_acl_family` parameter to the `tencentcloud_teo_origin_acl` resource schema
- Implement parameter mapping to cloud API in Create, Read, and Update operations
- Ensure backward compatibility with existing configurations
- Add appropriate documentation and tests

**Non-Goals:**
- Modify existing parameters or their behavior
- Change the resource's fundamental CRUD logic
- Add new cloud API operations or modify existing API calls
- Implement data source changes (only resource changes required)

## Decisions

### 1. Schema Definition
- **Decision**: Define `origin_acl_family` as `Type: schema.TypeString`, `Optional: true`, `Computed: true`
- **Rationale**: The parameter is optional in the cloud API with default behavior. Making it computed allows the API's default value to be read back and stored in state, ensuring users can see the actual value configured by the service.
- **Alternatives Considered**:
  - `Optional: true, Computed: false` - Would not allow reading back API defaults, potentially causing configuration drift
  - `Required: true` - Would break backward compatibility and prevent existing configurations from working

### 2. Create Operation Mapping
- **Decision**: Map `origin_acl_family` to `EnableOriginACLRequest.OriginACLFamily` in the Create function
- **Rationale**: Direct mapping to the cloud API parameter. If user provides a value, set it; if not, let the API use its default.
- **Alternatives Considered**: None - this is the standard pattern for adding new parameters

### 3. Read Operation Mapping
- **Decision**: Map `OriginACLInfo.OriginACLFamily` from `DescribeOriginACL` response to state in the Read function
- **Rationale**: Ensures state reflects the actual API configuration, including any default values applied by the service.
- **Alternatives Considered**: None - this is the standard pattern for computed fields

### 4. Update Operation Handling
- **Decision**: Support `origin_acl_family` updates via `ModifyOriginACL` API when the parameter changes
- **Rationale**: The ModifyOriginACL API supports this parameter, allowing users to change the control domain without recreating the resource.
- **Implementation**: Check `d.HasChange("origin_acl_family")` and call ModifyOriginACL with the new value if changed
- **Alternatives Considered**: Require resource recreation (ForceNew) - rejected because the API supports updates, requiring recreation would be a poor user experience

### 5. Error Handling
- **Decision**: Follow existing error handling patterns in the resource (defer tccommon.LogElapsed(), defer tccommon.InconsistentCheck())
- **Rationale**: Consistency with existing codebase and proper error tracking
- **Alternatives Considered**: None - maintain existing patterns

### 6. Documentation
- **Decision**: Update `resource_tc_teo_origin_acl.md` to document the new parameter with examples
- **Rationale**: All Terraform resources require documentation in the website/docs/ directory (per project constraints)
- **Alternatives Considered**: None - required by project standards

## Risks / Trade-offs

### Risk 1: State Migration for Existing Resources
- **Risk**: Existing resources may have state mismatches if they were created before this parameter was added
- **Mitigation**: The parameter is computed and optional, so Terraform will automatically read the API value and update state without user intervention
- **Trade-off**: No state migration needed - computed fields handle this gracefully

### Risk 2: API Compatibility Changes
- **Risk**: The cloud API may change the behavior of `OriginACLFamily` parameter (e.g., new values, deprecated values)
- **Mitigation**: The parameter is a pass-through to the API; Terraform does not validate values, allowing the API to handle validation
- **Trade-off**: Accepting this risk is reasonable as it mirrors the API contract

### Risk 3: Update Conflicts
- **Risk**: Concurrent updates to `origin_acl_family` and other parameters (l7_hosts, l4_proxy_ids) might cause API conflicts
- **Mitigation**: The existing code already handles batch updates separately; extend this pattern to handle `origin_acl_family` updates independently or together with entity updates
- **Trade-off**: Minimal risk - the parameter update logic can be integrated into existing update handling

### Risk 4: Default Value Visibility
- **Risk**: Users may not be aware of the default value if they don't explicitly set `origin_acl_family`
- **Mitigation**: The parameter is computed, so the API's default will be visible in state after creation; documentation should explain this behavior
- **Trade-off**: Users will see the actual value in state, improving transparency

### Risk 5: Performance Impact
- **Risk**: Adding a new parameter to Read operation adds minimal overhead (one additional string field)
- **Mitigation**: Negligible impact - string field reading is lightweight
- **Trade-off**: No meaningful performance degradation

## Migration Plan

### Deployment Steps
1. Update resource schema in `resource_tc_teo_origin_acl.go`
2. Modify Create function to set `origin_acl_family` if provided
3. Modify Read function to read `origin_acl_family` from API response
4. Modify Update function to handle `origin_acl_family` changes via ModifyOriginACL
5. Add unit tests in `resource_tc_teo_origin_acl_test.go`
6. Update documentation in `resource_tc_teo_origin_acl.md`
7. Run acceptance tests (TF_ACC=1) to verify functionality

### Rollback Strategy
- If issues arise, the change can be rolled back by removing the schema field and related code
- Since the parameter is optional and computed, existing resources will continue to work (the API will use defaults)
- No data loss or state corruption risk

### Backward Compatibility
- Fully backward compatible: optional parameter means existing configurations continue to work
- State migration not required: computed field will be populated automatically on next read
- No breaking changes to existing functionality

## Open Questions

None - the implementation is straightforward and follows established patterns in the codebase. All cloud API mappings are verified, and the design decisions are clear based on the API contract and project conventions.
