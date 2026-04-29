## Why

Terraform Provider for TencentCloud currently lacks a data source to query TEO (TencentCloud EdgeOne) web security policy templates. Users need to discover available security policy templates by zone IDs in order to reference them in other resources (e.g., binding templates to domains via `tencentcloud_teo_bind_security_template`). Without this data source, users must manually look up template IDs from the console.

## What Changes

- Add a new data source `tencentcloud_teo_web_security_templates` that wraps the `DescribeWebSecurityTemplates` API
- The data source accepts `zone_ids` as input and returns a list of `security_policy_templates` with their details (zone_id, template_id, template_name, bind_domains)
- Register the new data source in `tencentcloud/provider.go` and `tencentcloud/provider.md`
- Add corresponding `.md` documentation file for the data source
- Add unit tests using gomonkey mock approach

## Capabilities

### New Capabilities
- `teo-web-security-templates-datasource`: Data source to query TEO web security policy templates by zone IDs using the DescribeWebSecurityTemplates API

### Modified Capabilities

## Impact

- New files in `tencentcloud/services/teo/`: data source Go file, test file, md file
- Modified files: `tencentcloud/provider.go` (registration), `tencentcloud/provider.md` (documentation entry)
- Cloud API dependency: `teo/v20220901.DescribeWebSecurityTemplates`
