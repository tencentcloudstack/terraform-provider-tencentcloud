# Implementation Tasks: Add CLB Target Group Missing Parameters

## Phase 1: Schema Definition (Core Task)

- [ ] **Task 1.1**: Add `health_check` schema block with nested fields
  - Add `health_switch` (bool) - Enable/disable health check
  - Add `protocol` (string) - Health check protocol (TCP/HTTP/HTTPS/PING/CUSTOM)
  - Add `port` (int) - Health check port
  - Add `timeout` (int) - Timeout in seconds [2-30]
  - Add `gap_time` (int) - Interval in seconds [2-300]
  - Add `good_limit` (int) - Healthy threshold [2-10]
  - Add `bad_limit` (int) - Unhealthy threshold [2-10]
  - Add `http_code` (int) - HTTP status codes (for HTTP/HTTPS)
  - Add `http_check_domain` (string) - Health check domain (for HTTP/HTTPS)
  - Add `http_check_path` (string) - Health check path (for HTTP/HTTPS)
  - Add `http_check_method` (string) - HTTP method (for HTTP/HTTPS)
  - Add `http_version` (string) - HTTP version (for HTTP/HTTPS)
  - Add `source_ip_type` (int) - Source IP type
  - Add `extended_code` (string) - Extended status code
  - **Validation**: Ensure nested structure is properly defined
  - **Dependencies**: None

- [ ] **Task 1.2**: Add `schedule_algorithm` schema field
  - Type: String, Optional, Computed
  - Valid values: WRR, LEAST_CONN, IP_HASH
  - Default: WRR
  - ForceNew: true (cannot be modified after creation)
  - Description: "Scheduling algorithm. Only valid for v2 target groups with HTTP/HTTPS/GRPC protocols. Options: WRR (weighted round robin), LEAST_CONN (least connections), IP_HASH (IP hash). Default: WRR."
  - **Validation**: Add ValidateFunc to check allowed values
  - **Dependencies**: None

- [ ] **Task 1.3**: Add `tags` schema field
  - Type: Map, Optional
  - Description: "Resource tags for the target group."
  - **Validation**: Standard tag validation
  - **Dependencies**: Import helper package for tags

- [ ] **Task 1.4**: Add `weight` schema field
  - Type: Int, Optional
  - Valid range: [0, 100]
  - Description: "Default backend server weight. Valid for v2 target groups only. When set, backend servers added to the target group will use this default weight if not specified."
  - **Validation**: Add ValidateFunc to check range [0, 100]
  - **Dependencies**: None

- [ ] **Task 1.5**: Add `full_listen_switch` schema field
  - Type: Bool, Optional
  - ForceNew: true
  - Description: "Whether this is a full listener target group. Only valid for v2 target groups. true = full listener target group, false = normal target group."
  - **Validation**: None needed
  - **Dependencies**: None

- [ ] **Task 1.6**: Add `keepalive_enable` schema field
  - Type: Bool, Optional
  - Default: false
  - Description: "Enable keep-alive connections. Only valid for HTTP/HTTPS target groups. true = enable, false = disable. Default: false."
  - **Validation**: None needed
  - **Dependencies**: None

- [ ] **Task 1.7**: Add `session_expire_time` schema field
  - Type: Int, Optional
  - Valid range: [30, 3600]
  - Default: 0 (disabled)
  - Description: "Session persistence time in seconds. Only valid for v2 target groups with HTTP/HTTPS/GRPC protocols. Range: 30-3600. Default: 0 (disabled)."
  - **Validation**: Add ValidateFunc to check range [30, 3600] or 0
  - **Dependencies**: None

- [ ] **Task 1.8**: Add `ip_version` schema field
  - Type: String, Optional, Computed
  - ForceNew: true
  - Description: "IP version type. Common values: IPv4, IPv6, IPv6FullChain."
  - **Validation**: None (rely on API validation)
  - **Dependencies**: None

## Phase 2: Create Logic Implementation

- [ ] **Task 2.1**: Update `resourceTencentCloudClbTargetCreate` function
  - Extract new parameters from resource data
  - Handle `health_check` nested object conversion
  - Handle `tags` array conversion
  - Pass parameters to service layer
  - **Validation**: Verify all new fields are extracted correctly
  - **Dependencies**: Task 1.1-1.8, Task 2.2

- [ ] **Task 2.2**: Update `ClbService.CreateTargetGroup` method signature
  - Add parameters: healthCheck, scheduleAlgorithm, tags, weight, fullListenSwitch, keepaliveEnable, sessionExpireTime, ipVersion
  - Build `CreateTargetGroupRequest` with new fields
  - Map Terraform types to SDK types (e.g., tags map to TagInfo array)
  - **Validation**: Ensure proper type conversions
  - **Dependencies**: Task 1.1-1.8

- [ ] **Task 2.3**: Implement health check object construction
  - Create helper function to build `TargetGroupHealthCheck` from schema
  - Handle all nested fields properly
  - Ensure nil-safety for optional fields
  - **Validation**: Test with various health check configurations
  - **Dependencies**: Task 1.1, Task 2.2

- [ ] **Task 2.4**: Implement tags conversion
  - Convert Terraform tags map to SDK TagInfo array
  - Use existing helper functions if available
  - **Validation**: Verify tag format matches SDK expectations
  - **Dependencies**: Task 1.3, Task 2.2

## Phase 3: Read Logic Implementation

- [ ] **Task 3.1**: Update `resourceTencentCloudClbTargetRead` function
  - Read `health_check` from API response
  - Read `schedule_algorithm` from response
  - Read `tags` from response
  - Read `weight` from response
  - Read `full_listen_switch` from response
  - Read `keepalive_enable` from response
  - Read `session_expire_time` from response
  - Read `ip_version` from response
  - **Validation**: Verify all fields are set correctly in state
  - **Dependencies**: Task 1.1-1.8, Task 3.2

- [ ] **Task 3.2**: Implement health check flattening
  - Create helper function to flatten `TargetGroupHealthCheck` to schema
  - Handle all nested fields
  - Ensure nil-safety for optional response fields
  - **Validation**: Test read with various API responses
  - **Dependencies**: Task 1.1

- [ ] **Task 3.3**: Implement tags flattening
  - Convert SDK TagInfo array to Terraform tags map
  - Use existing helper functions if available
  - **Validation**: Verify tags are correctly read back
  - **Dependencies**: Task 1.3

## Phase 4: Update Logic Implementation

- [ ] **Task 4.1**: Research `ModifyTargetGroupAttribute` API capabilities
  - Determine which new parameters can be updated
  - Identify parameters requiring ForceNew
  - Document findings in code comments
  - **Validation**: Review API documentation
  - **Dependencies**: None

- [ ] **Task 4.2**: Update `resourceTencentCloudClbTargetUpdate` function
  - Add change detection for updatable new parameters
  - Implement update logic for each supported parameter
  - For non-updatable fields, rely on ForceNew flag
  - **Validation**: Test update scenarios
  - **Dependencies**: Task 4.1, Task 4.3

- [ ] **Task 4.3**: Update service layer update methods if needed
  - Extend `ClbService.ModifyTargetGroup` or create new methods
  - Handle API calls for parameter updates
  - **Validation**: Verify API calls are correct
  - **Dependencies**: Task 4.1

## Phase 5: Documentation

- [ ] **Task 5.1**: Update source documentation file
  - Update `resource_tc_clb_target_group.md`
  - Add example using health_check
  - Add example using v2 advanced features (schedule_algorithm, weight, session_expire_time)
  - Add example with tags
  - Add example showing full listener target group
  - Document version/protocol constraints clearly
  - **Validation**: Review examples for accuracy
  - **Dependencies**: None

- [ ] **Task 5.2**: Generate website documentation
  - Run `make doc` to generate `website/docs/r/clb_target_group.html.markdown`
  - Verify all new parameters appear in generated docs
  - Check formatting and clarity
  - **Validation**: Review generated documentation
  - **Dependencies**: Task 5.1

- [ ] **Task 5.3**: Add parameter constraint documentation
  - Document which parameters work with v1 vs v2
  - Document protocol-specific parameters
  - Add notes about ForceNew parameters
  - Include validation rules (ranges, allowed values)
  - **Validation**: Ensure users understand limitations
  - **Dependencies**: Task 5.1

## Phase 6: Testing

- [ ] **Task 6.1**: Add unit tests for schema validation
  - Test health_check nested structure
  - Test schedule_algorithm validation
  - Test weight range validation
  - Test session_expire_time range validation
  - **Validation**: All tests pass
  - **Dependencies**: Task 1.1-1.8

- [ ] **Task 6.2**: Add acceptance test for v1 target group with new features
  - Create v1 target group with health_check
  - Create v1 target group with tags
  - Verify parameters are set correctly
  - Verify read operations work
  - **Validation**: Test passes
  - **Dependencies**: Tasks 2.x, 3.x

- [ ] **Task 6.3**: Add acceptance test for v2 target group with advanced features
  - Create v2 target group with all new parameters
  - Test schedule_algorithm (WRR, LEAST_CONN, IP_HASH)
  - Test weight parameter
  - Test session_expire_time
  - Test keepalive_enable
  - Verify all features work together
  - **Validation**: Test passes
  - **Dependencies**: Tasks 2.x, 3.x

- [ ] **Task 6.4**: Add acceptance test for full listener target group
  - Create v2 target group with full_listen_switch=true
  - Verify constraint: port parameter should not be set
  - **Validation**: Test passes
  - **Dependencies**: Tasks 2.x, 3.x

- [ ] **Task 6.5**: Add acceptance test for update scenarios
  - Test updating target_group_name (existing functionality)
  - Test updating port (existing functionality)
  - Test updating updatable new parameters (if any)
  - Verify ForceNew parameters trigger recreation
  - **Validation**: Test passes
  - **Dependencies**: Tasks 4.x

- [ ] **Task 6.6**: Add negative tests for validation
  - Test invalid schedule_algorithm value
  - Test weight out of range
  - Test session_expire_time out of range
  - Test incompatible parameter combinations (e.g., full_listen_switch=true with port set)
  - **Validation**: Tests correctly reject invalid inputs
  - **Dependencies**: Task 1.1-1.8

## Phase 7: Code Quality and Finalization

- [ ] **Task 7.1**: Run linters and fix issues
  - Run `golangci-lint` on modified files
  - Fix any linting warnings or errors
  - Ensure code style consistency
  - **Validation**: Linter passes with no errors
  - **Dependencies**: All code tasks complete

- [ ] **Task 7.2**: Run `go fmt` and `goimports`
  - Format all modified Go files
  - Organize imports correctly
  - **Validation**: No formatting changes needed
  - **Dependencies**: All code tasks complete

- [ ] **Task 7.3**: Compile and verify build
  - Run `go build` on the provider
  - Ensure no compilation errors
  - **Validation**: Build succeeds
  - **Dependencies**: All code tasks complete

- [ ] **Task 7.4**: Create changelog entry
  - Create `.changelog/<PR_NUMBER>.txt`
  - Format: `resource/tencentcloud_clb_target_group: support health_check, schedule_algorithm, tags, weight, full_listen_switch, keepalive_enable, session_expire_time, and ip_version parameters`
  - **Validation**: Changelog format is correct
  - **Dependencies**: All tasks complete

- [ ] **Task 7.5**: Final validation with `openspec validate`
  - Run `openspec validate add-clb-target-group-missing-params --strict`
  - Fix any validation issues
  - **Validation**: Validation passes
  - **Dependencies**: All tasks complete

## Summary

- **Total Tasks**: 36
- **Estimated Effort**: 5-7 days
- **Critical Path**: Schema → Create → Read → Update → Tests → Documentation
- **Parallel Work Possible**: Documentation can start after schema is defined; unit tests can be written alongside implementation

## Dependencies Graph

```
Phase 1 (Schema) → Phase 2 (Create) → Phase 3 (Read) → Phase 4 (Update)
                ↓
              Phase 5 (Docs)
                ↓
              Phase 6 (Tests) → Phase 7 (Finalization)
```

## Notes

- All new parameters are optional to maintain backward compatibility
- Some parameters have type/protocol constraints - document clearly
- Tags integration may require importing helper package
- Health check is the most complex parameter (nested object with many fields)
- Schedule algorithm, full_listen_switch, and ip_version are likely ForceNew (cannot be updated)
