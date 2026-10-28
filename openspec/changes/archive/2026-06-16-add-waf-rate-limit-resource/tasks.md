## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/waf/resource_tc_waf_rate_limit.go` with schema definition including all parameters (domain, name, priority, status, limit_window, limit_object, limit_strategy, limit_method, limit_paths, limit_headers, limit_header_name, get_params_name, get_params_value, post_params_name, post_params_value, ip_location, redirect_info, block_page, object_src, quota_share, paths_option, order, limit_rule_id) and CRUD functions
- [x] 1.2 Implement Create function calling `CreateRateLimitV2` API with retry logic, response nil-check, and composite ID (`domain#limit_rule_id`) assignment
- [x] 1.3 Implement Read function calling `DescribeRateLimitsV2` API with `Domain` and `Id` filter, retry logic, nil-check with logging before `d.SetId("")`, and setting all fields into state
- [x] 1.4 Implement Update function calling `UpdateRateLimitV2` API with retry logic, passing all changed parameters
- [x] 1.5 Implement Delete function calling `DeleteRateLimitsV2` API with retry logic, passing `Domain` and single-element `LimitRuleIds` array

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_waf_rate_limit` resource in `tencentcloud/provider.go`
- [x] 2.2 Add `tencentcloud_waf_rate_limit` entry in `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/waf/resource_tc_waf_rate_limit.md` with one-line description, Example Usage section, and Import section showing composite ID format `domain#limit_rule_id`

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/waf/resource_tc_waf_rate_limit_test.go` with gomonkey-based unit tests covering Create, Read, Update, and Delete operations
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass
