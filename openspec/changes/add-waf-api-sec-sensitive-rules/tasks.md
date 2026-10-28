# Implementation Tasks

> Reference resource for code style: `tencentcloud_igtm_monitor`. Doc/test naming: `resource_tc_config_compliance_pack.md` / `resource_tc_config_compliance_pack_test.go`. All API calls wrapped with `resource.Retry`; all response reads nil-safe. ID = `Domain#RuleName`. `status` exposes only `0`/`1`; Delete sends `Status=3` internally.

## 1. Service Layer (shared)

- [x] 1.1 In `tencentcloud/services/waf/service_tencentcloud_waf.go`, add `DescribeWafApiSecSensitiveRuleListByFilter(ctx, domain, ruleName, queryFlags...)` helper wrapping `DescribeApiSecSensitiveRuleList` (set `Domain`, the relevant `IsQuery*` flag, and `RuleName`), with `resource.Retry` + `ratelimit.Check` + nil-safe response, returning the full `DescribeApiSecSensitiveRuleListResponseParams`.
- [x] 1.2 Add small per-type getters (or let each resource match by `RuleName` locally) returning the matched sub-struct from the response list (`Data`/`ApiExtractRule`/`ApiSecPrivilegeRule`/`ApiSecSceneRule`/`ApiSecCustomEventRule`/`ApiExcludeRule`/`ApiSecSensitiveWhiteRule`).

## 2. Resource: tencentcloud_waf_api_sec_sensitive_custom_rule (CustomRule)

- [x] 2.1 Create `resource_tc_waf_api_sec_sensitive_custom_rule.go` with schema = `domain`, `rule_name`, `status` (0/1) + `ApiSecCustomSensitiveRule` fields (`position`, `match_key`, `match_value`, `level`, `match_cond`, `is_pan`) only; implement Create/Read/Update/Delete (Delete uses `Status=3`), ID `Domain#RuleName`, import passthrough.
- [x] 2.2 Read via Describe (default `Data` list), match by `rule_name`, set nil-safely; clear ID when not found.
- [x] 2.3 Create `resource_tc_waf_api_sec_sensitive_custom_rule.md` example doc.
- [x] 2.4 Create `resource_tc_waf_api_sec_sensitive_custom_rule_test.go` acceptance test (basic + update + import).

## 3. Resource: tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule (CustomApiExtractRule)

- [x] 3.1 Create `resource_tc_waf_api_sec_sensitive_custom_api_extract_rule.go` with schema = `domain`, `rule_name`, `status` + `ApiSecExtractRule` input fields (`api_name`, `methods`, `regex`) + computed `update_time`; full CRUD; unify `RuleName`/`Status` with top-level.
- [x] 3.2 Read via Describe with `IsQueryApiExtractRule=true`, match `ApiExtractRule` by `rule_name`, set nil-safely.
- [x] 3.3 Create `resource_tc_waf_api_sec_sensitive_custom_api_extract_rule.md` example doc.
- [x] 3.4 Create `resource_tc_waf_api_sec_sensitive_custom_api_extract_rule_test.go` acceptance test.

## 4. Resource: tencentcloud_waf_api_sec_sensitive_privilege_rule (ApiSecPrivilegeRule)

- [x] 4.1 Create `resource_tc_waf_api_sec_sensitive_privilege_rule.go` with schema = `domain`, `rule_name`, `status` + `ApiSecPrivilegeRule` input fields (`api_name`, `position`, `parameter_list`, `source`, `option`, nested `api_name_op` block with `value`/`op`/`api_name_method`) + computed `update_time`; full CRUD.
- [x] 4.2 Read via Describe with `IsQueryApiPrivilegeRule=true`, match by `rule_name`, flatten nested `api_name_op` nil-safely.
- [x] 4.3 Create `resource_tc_waf_api_sec_sensitive_privilege_rule.md` example doc.
- [x] 4.4 Create `resource_tc_waf_api_sec_sensitive_privilege_rule_test.go` acceptance test.

## 5. Resource: tencentcloud_waf_api_sec_sensitive_scene_rule (ApiSecSceneRule)

- [x] 5.1 Create `resource_tc_waf_api_sec_sensitive_scene_rule.go` with schema = `domain`, `rule_name`, `status` + `ApiSecSceneRule` input fields (`source`, nested `rule_list` of `ApiSecSceneRuleEntry`: `key`/`value`/`operate`/`name`) + computed `update_time`; full CRUD.
- [x] 5.2 Read via Describe with `IsQueryApiSceneRule=true`, match by `rule_name`, flatten `rule_list` nil-safely.
- [x] 5.3 Create `resource_tc_waf_api_sec_sensitive_scene_rule.md` example doc.
- [x] 5.4 Create `resource_tc_waf_api_sec_sensitive_scene_rule_test.go` acceptance test.

## 6. Resource: tencentcloud_waf_api_sec_sensitive_custom_event_rule (ApiSecCustomEventRuleRule)

- [x] 6.1 Create `resource_tc_waf_api_sec_sensitive_custom_event_rule.go` with schema = `domain`, `rule_name`, `status` + `ApiSecCustomEventRule` input fields (`description`, `req_frequency`, `risk_level`, `source`, nested `api_name_op`, `match_rule_list`/`stat_rule_list` of `ApiSecSceneRuleEntry`) + computed `update_time`; full CRUD.
- [x] 6.2 Read via Describe with `IsQueryApiCustomEventRule=true`, match by `rule_name`, flatten nested blocks nil-safely.
- [x] 6.3 Create `resource_tc_waf_api_sec_sensitive_custom_event_rule.md` example doc.
- [x] 6.4 Create `resource_tc_waf_api_sec_sensitive_custom_event_rule_test.go` acceptance test.

## 7. Resource: tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule (CustomApiExcludeRule)

- [x] 7.1 Create `resource_tc_waf_api_sec_sensitive_custom_api_exclude_rule.go` with schema = `domain`, `rule_name`, `status` + `ApiSecExcludeRule` input fields (`match_type`, `content`) + computed `update_time`; full CRUD.
- [x] 7.2 Read via Describe with `IsQueryApiExcludeRule=true`, match `ApiExcludeRule` by `rule_name`, set nil-safely.
- [x] 7.3 Create `resource_tc_waf_api_sec_sensitive_custom_api_exclude_rule.md` example doc.
- [x] 7.4 Create `resource_tc_waf_api_sec_sensitive_custom_api_exclude_rule_test.go` acceptance test.

## 8. Resource: tencentcloud_waf_api_sec_sensitive_white_rule (ApiSecSensitiveWhiteRuleRule)

- [x] 8.1 Create `resource_tc_waf_api_sec_sensitive_white_rule.go` with schema = `domain`, `rule_name`, `status` + `ApiSecSensitiveWhiteRule` input fields (`white_mode`, `description`, nested `api_name_op`, `white_fields` of `ApiSecSensitiveWhiteField`: `field_name`/`field_type`/`sensitive_types`) + computed `update_time`; full CRUD.
- [x] 8.2 Read via Describe with `IsQueryApiSensitiveWhiteRule=true`, match by `rule_name`, flatten `white_fields`/`api_name_op` nil-safely.
- [x] 8.3 Create `resource_tc_waf_api_sec_sensitive_white_rule.md` example doc.
- [x] 8.4 Create `resource_tc_waf_api_sec_sensitive_white_rule_test.go` acceptance test.

## 9. Provider Registration

- [x] 9.1 Register all 7 resources in `tencentcloud/provider.go` `ResourcesMap` (`waf.ResourceTencentCloudWafApiSecSensitive*()`).
- [x] 9.2 Add the 7 resource pages to `website/tencentcloud.erb` under the WAF section.

## 10. Documentation Generation

- [x] 10.1 Run `make doc` to generate the 7 `website/docs/r/waf_api_sec_sensitive_*.html.markdown` pages from the `.md` example files (do not hand-write the website docs).

## 11. Verification

- [x] 11.1 Run `gofmt -w` on all new files and `go build ./tencentcloud/...`.
- [x] 11.2 Run `go vet ./tencentcloud/services/waf/` and fix any newly introduced issues.
- [x] 11.3 Confirm each resource schema contains only its corresponding sub-struct fields (+ `domain`/`rule_name`/`status`) — strict field validation per requirement.
