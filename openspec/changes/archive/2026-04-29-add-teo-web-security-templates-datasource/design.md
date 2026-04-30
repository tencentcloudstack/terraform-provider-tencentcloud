## Context

The TencentCloud Terraform Provider currently has a resource `tencentcloud_teo_web_security_template` for managing individual TEO web security policy templates, and `tencentcloud_teo_bind_security_template` for binding templates to domains. However, there is no data source to query/list available security policy templates by zone IDs. Users need to reference template IDs from the data source when configuring bindings.

The `DescribeWebSecurityTemplates` API (in `teo/v20220901`) already exists and is used internally in `service_tencentcloud_teo.go` (in `DescribeTeoWebSecurityTemplateNameById`). This design exposes it as a proper Terraform data source.

## Goals / Non-Goals

**Goals:**
- Provide a `tencentcloud_teo_web_security_templates` data source that queries security policy templates by zone IDs
- Follow the existing datasource pattern (reference: `data_source_tc_igtm_instance_list.go`)
- Return nested structure matching the API response (SecurityPolicyTemplateInfo with BindDomains)
- Register in provider.go and provider.md
- Add corresponding .md documentation and unit tests

**Non-Goals:**
- Modifying the existing `tencentcloud_teo_web_security_template` resource
- Adding create/update/delete operations (this is a read-only data source)
- Exposing pagination parameters (the API does not have pagination for this endpoint)

## Decisions

1. **Data source name**: `tencentcloud_teo_web_security_templates` (plural form, following the naming convention for list-type data sources)
   - Alternative: `tencentcloud_teo_web_security_template_list` - but the API returns templates, and the plural form is more natural

2. **Schema design for input**: `zone_ids` as `TypeList` of `TypeString` (required) - matches the API's `ZoneIds` parameter which accepts a list of zone IDs (max 100)

3. **Schema design for output**: `security_policy_templates` as `TypeList` of `TypeResource` with nested fields:
   - `zone_id` (string, computed) - from SecurityPolicyTemplateInfo.ZoneId
   - `template_id` (string, computed) - from SecurityPolicyTemplateInfo.TemplateId
   - `template_name` (string, computed) - from SecurityPolicyTemplateInfo.TemplateName
   - `bind_domains` (TypeList of TypeResource, computed) - from SecurityPolicyTemplateInfo.BindDomains
     - `domain` (string, computed) - from BindDomainInfo.Domain
     - `zone_id` (string, computed) - from BindDomainInfo.ZoneId
     - `status` (string, computed) - from BindDomainInfo.Status

4. **No pagination needed**: The `DescribeWebSecurityTemplates` API does not have Offset/Limit parameters, so no pagination logic is required

5. **Resource ID**: Use `helper.DataResourceIdsHash()` with collected template IDs, consistent with other list data sources

6. **Retry logic**: Use `resource.Retry(tccommon.ReadRetryTimeout, ...)` for API calls, consistent with existing patterns

7. **Test approach**: Use gomonkey mock approach for unit tests (not Terraform acceptance tests), as per the code generation requirements for new resources

## Risks / Trade-offs

- [API returns all templates for given zone IDs in one call] → No pagination risk, but large number of templates could mean large response. Mitigation: The API limits zone IDs to 100 per call which constrains the response size.
- [Existing service method `DescribeTeoWebSecurityTemplateNameById` already calls this API] → The new data source will add a separate service method for the full list query, avoiding coupling with the existing method's specific use case.
