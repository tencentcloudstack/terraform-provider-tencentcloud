## 1. Service Layer

- [x] 1.1 Append `DescribeConfigCompliancePacksByFilter()` to `service_tencentcloud_config.go` — wraps `ListCompliancePacks` with paged retry loop
- [x] 1.2 Append `DescribeSystemConfigCompliancePacks()` to `service_tencentcloud_config.go` — wraps `ListSystemCompliancePacks` with paged retry loop

## 2. Data Source: tencentcloud_config_compliance_packs

- [x] 2.1 Create `data_source_tc_config_compliance_packs.go` with `DataSourceTencentCloudConfigCompliancePacks()` schema and read handler
- [x] 2.2 Schema: optional filters (`compliance_pack_name`, `risk_level`, `status`, `compliance_result`, `order_type`), computed `compliance_pack_list`, optional `result_output_file`
- [x] 2.3 Read handler: build paramMap from schema, call service, flatten result into `compliance_pack_list`

## 3. Data Source: tencentcloud_system_config_compliance_packs

- [x] 3.1 Create `data_source_tc_system_config_compliance_packs.go` with `DataSourceTencentCloudSystemConfigCompliancePacks()` schema and read handler
- [x] 3.2 Schema: computed `compliance_pack_list` (with nested `config_rules`), optional `result_output_file`
- [x] 3.3 Read handler: call service, flatten result including nested config_rules

## 4. Provider Registration

- [x] 4.1 Register `tencentcloud_config_compliance_packs` in `provider.go` DataSourcesMap
- [x] 4.2 Register `tencentcloud_system_config_compliance_packs` in `provider.go` DataSourcesMap

## 5. Documentation

- [x] 5.1 Create `data_source_tc_config_compliance_packs.md` with usage example
- [x] 5.2 Create `data_source_tc_system_config_compliance_packs.md` with usage example

## 6. Tests

- [x] 6.1 Create `data_source_tc_config_compliance_packs_test.go` with basic acceptance test
- [x] 6.2 Create `data_source_tc_system_config_compliance_packs_test.go` with basic acceptance test

