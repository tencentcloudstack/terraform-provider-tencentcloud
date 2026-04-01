## 1. Schema Definition

- [x] 1.1 Add `rules` attribute to resource_teo_l7_acc_rule.go schema
- [x] 1.2 Define nested `rules` schema with basic fields (status, rule_id, rule_name, description, rule_priority)
- [x] 1.3 Define `branches` nested schema within each rule
- [x] 1.4 Define `actions` nested schema within branches
- [x] 1.5 Define `sub_rules` nested schema within branches
- [x] 1.6 Implement action parameter schemas for common types (Cache, CacheKey, CachePrefresh, AccessURLRedirect, UpstreamURLRewrite)
- [x] 1.7 Mark all rules-related fields as Computed and Optional

## 2. Read Function Implementation

- [x] 2.1 Add DescribeL7AccRules API call to service layer
- [x] 2.2 Create helper function to map API response to Rules schema structure
- [x] 2.3 Implement mapping for basic rule fields (status, rule_id, rule_name, description, rule_priority)
- [x] 2.4 Implement mapping for branches and their condition field
- [x] 2.5 Implement mapping for actions list with type-specific parameters
- [x] 2.6 Implement mapping for sub_rules and their nested structure
- [x] 2.7 Add error handling for nil/empty rules data from API
- [x] 2.8 Integrate rules mapping into Read function in resource_teo_l7_acc_rule.go

## 3. Action Parameter Mapping

- [x] 3.1 Implement CacheParameters mapping (FollowOrigin, NoCache, CustomTime)
- [x] 3.2 Implement CacheKeyParameters mapping (QueryString, IgnoreCase, Header, Scheme, Cookie)
- [x] 3.3 Implement CachePrefreshParameters mapping
- [x] 3.4 Implement AccessURLRedirectParameters mapping
- [x] 3.5 Implement UpstreamURLRewriteParameters mapping
- [x] 3.6 Add mapping for other common action types (QUIC, WebSocket, Authentication, MaxAge)

## 4. Create/Update Functions (Optional)

- [x] 4.1 Review Create function to ensure compatibility with new computed field
- [x] 4.2 Review Update function to ensure compatibility with new computed field
- [x] 4.3 Add placeholder logic for future rules creation/update support

## 5. Unit Tests

- [x] 5.1 Add unit test for rules schema definition
- [x] 5.2 Add unit test for rules mapping from API response
- [x] 5.3 Add unit test for basic rule fields (status, rule_id, rule_name, description)
- [x] 5.4 Add unit test for branches mapping
- [x] 5.5 Add unit test for actions mapping with different parameter types
- [x] 5.6 Add unit test for sub_rules nested structure
- [x] 5.7 Add unit test for empty/nil rules handling
- [x] 5.8 Run unit tests with `go test ./tencentcloud/services/teo -v`

## 6. Documentation

- [x] 6.1 Update resource_tc_teo_l7_acc_rule.md example file with rules documentation
- [x] 6.2 Add example showing rules structure in the example
- [x] 6.3 Run `make doc` to generate website/docs/ documentation
- [x] 6.4 Verify generated documentation includes rules attribute

## 7. Validation and Testing

- [ ] 7.1 Run `make build` to verify code compiles successfully
- [ ] 7.2 Run `make lint` to verify code passes linting checks
- [ ] 7.3 Run acceptance tests with `TF_ACC=1 go test ./tencentcloud/services/teo -v -run TestAccTencentCloudTeoL7AccRule`
- [ ] 7.4 Verify resource import works correctly with new rules attribute
- [ ] 7.5 Test backward compatibility with existing resources without rules in state

## 8. Code Review

- [x] 8.1 Self-review the schema definition for completeness
- [x] 8.2 Review Read function implementation for correctness
- [x] 8.3 Review error handling and edge cases
- [x] 8.4 Verify all helper functions are well-documented
- [x] 8.5 Ensure test coverage is adequate for new functionality
