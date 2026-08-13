## Why

TencentCloud WAF (Web Application Firewall) provides rate limiting capabilities to protect web applications from abuse and DDoS attacks. Currently, there is no Terraform resource to manage WAF rate limit rules, requiring users to manually configure them through the console or API. Adding `tencentcloud_waf_rate_limit` enables Infrastructure-as-Code management of WAF rate limiting rules.

## What Changes

- Add a new Terraform resource `tencentcloud_waf_rate_limit` (RESOURCE_KIND_GENERAL) that manages the full lifecycle (CRUD) of WAF rate limit rules.
- The resource uses the following TencentCloud APIs:
  - `CreateRateLimitV2` - Create a rate limit rule
  - `DescribeRateLimitsV2` - Read/query rate limit rules
  - `UpdateRateLimitV2` - Update a rate limit rule
  - `DeleteRateLimitsV2` - Delete rate limit rules
- The resource supports configuring rate limiting by domain, API path, HTTP method, headers, GET/POST parameters, and IP location.
- The resource ID is a composite of `domain` and `limit_rule_id` (separated by `FILED_SP`).

## Capabilities

### New Capabilities
- `waf-rate-limit-crud`: Full CRUD lifecycle management for WAF rate limit rules, including creation, reading, updating, and deletion of rate limiting configurations.

### Modified Capabilities

## Impact

- New files:
  - `tencentcloud/services/waf/resource_tc_waf_rate_limit.go` - Resource implementation
  - `tencentcloud/services/waf/resource_tc_waf_rate_limit_test.go` - Unit tests
  - `tencentcloud/services/waf/resource_tc_waf_rate_limit.md` - Documentation
- Modified files:
  - `tencentcloud/provider.go` - Register the new resource
  - `tencentcloud/provider.md` - Add resource to provider documentation
- Dependencies: Uses existing `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125` package (already vendored)
