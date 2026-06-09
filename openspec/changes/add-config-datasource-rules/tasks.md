## 1. Service Layer

- [x] 1.1 Append `DescribeConfigSystemRulesByFilter()` — wraps `ListSystemRules` with paged retry loop (Limit=100)
- [x] 1.2 Append `DescribeConfigRuleEvaluationResultsByFilter()` — wraps `ListConfigRuleEvaluationResults` with paged retry loop (Limit=100)
- [x] 1.3 Append `DescribeConfigRulesByFilter()` — wraps `ListConfigRules` with paged retry loop (Limit=200)

## 2. Data Source: tencentcloud_config_system_rules

- [x] 2.1 Create `data_source_tc_config_system_rules.go` with schema and read handler
- [x] 2.2 Schema: optional `keyword`, `risk_level` filters; computed `rule_list`
- [x] 2.3 Read handler: build paramMap, call service, flatten into `rule_list`

## 3. Data Source: tencentcloud_config_rule_evaluation_results

- [x] 3.1 Create `data_source_tc_config_rule_evaluation_results.go` with schema and read handler
- [x] 3.2 Schema: required `config_rule_id`, optional `resource_type`/`compliance_type`; computed `result_list`
- [x] 3.3 Read handler: build paramMap, call service, flatten including nested `annotation`

## 4. Data Source: tencentcloud_config_rules

- [x] 4.1 Create `data_source_tc_config_rules.go` with schema and read handler
- [x] 4.2 Schema: optional filters (`rule_name`, `risk_level`, `state`, `compliance_result`, `order_type`); computed `rule_list`
- [x] 4.3 Read handler: build paramMap, call service, flatten into `rule_list`

## 5. Provider Registration

- [x] 5.1 Register `tencentcloud_config_system_rules` in `provider.go` DataSourcesMap
- [x] 5.2 Register `tencentcloud_config_rule_evaluation_results` in `provider.go` DataSourcesMap
- [x] 5.3 Register `tencentcloud_config_rules` in `provider.go` DataSourcesMap

## 6. Documentation

- [x] 6.1 Create `data_source_tc_config_system_rules.md`
- [x] 6.2 Create `data_source_tc_config_rule_evaluation_results.md`
- [x] 6.3 Create `data_source_tc_config_rules.md`

## 7. Tests

- [x] 7.1 Create `data_source_tc_config_system_rules_test.go`
- [x] 7.2 Create `data_source_tc_config_rule_evaluation_results_test.go`
- [x] 7.3 Create `data_source_tc_config_rules_test.go`

