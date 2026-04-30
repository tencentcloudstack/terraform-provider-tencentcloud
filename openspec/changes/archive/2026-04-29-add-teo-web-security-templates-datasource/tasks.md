## 1. Service Layer

- [x] 1.1 Add `DescribeTeoWebSecurityTemplatesByFilter` method in `tencentcloud/services/teo/service_tencentcloud_teo.go` that calls `DescribeWebSecurityTemplates` API with ZoneIds parameter, using `resource.Retry(tccommon.ReadRetryTimeout, ...)` for retry, and returns `[]*teov20220901.SecurityPolicyTemplateInfo`

## 2. Data Source Implementation

- [x] 2.1 Create `tencentcloud/services/teo/data_source_tc_teo_web_security_templates.go` with schema definition: `zone_ids` (required TypeList of TypeString), `security_policy_templates` (computed TypeList of TypeResource with nested fields: zone_id, template_id, template_name, bind_domains), `result_output_file` (optional TypeString)
- [x] 2.2 Implement `dataSourceTencentCloudTeoWebSecurityTemplatesRead` function following the pattern from `data_source_tc_igtm_instance_list.go`: call service method with retry, flatten response into schema format, set resource ID using `helper.DataResourceIdsHash()`, handle nil field checks before setting values
- [x] 2.3 Ensure bind_domains nested structure properly maps BindDomainInfo fields (domain, zone_id, status) with nil checks

## 3. Provider Registration

- [x] 3.1 Add `"tencentcloud_teo_web_security_templates": teo.DataSourceTencentCloudTeoWebSecurityTemplates()` entry in `tencentcloud/provider.go` data sources map
- [x] 3.2 Add `tencentcloud_teo_web_security_templates` entry in `tencentcloud/provider.md` data sources section

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/data_source_tc_teo_web_security_templates.md` with one-line description mentioning TEO, example usage with zone_ids, and no Argument/Attribute Reference sections

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/teo/data_source_tc_teo_web_security_templates_test.go` with gomonkey mock tests for the read operation, covering normal case with templates returned and empty result case
- [x] 5.2 Run unit tests with `go test -gcflags=all=-l` to verify they pass
