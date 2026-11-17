## Why

The `tencentcloud_teo_bind_security_template` resource currently relies on the `DescribeSecurityTemplateBindings` API to read binding state. This API is being phased out / does not reliably return the binding status for all scenarios, causing read failures and state drift. The `DescribeWebSecurityTemplates` API provides the same binding information (template-to-domain associations with delivery status) and is the recommended replacement. Switching the read path to `DescribeZones` + `DescribeWebSecurityTemplates` keeps the resource's user-facing behavior unchanged while using a supported, more reliable API.

## What Changes

- Replace the `DescribeTeoBindSecurityTemplateById` implementation in `tencentcloud/services/teo/service_tencentcloud_teo.go`:
  - Remove the `DescribeSecurityTemplateBindings` API call.
  - Add a new helper `describeTeoAllZoneIds` that pages through `DescribeZones` (Limit=100) to collect all zone IDs under the account.
  - Call `DescribeWebSecurityTemplates` in batches of at most 100 zone IDs per request (the API's documented upper limit), filter the returned `SecurityPolicyTemplates` by `TemplateId` and the bound `BindDomains.Domain` (entity), and return the matching `EntityStatus` (entity + status).
- Update `resource_tc_teo_bind_security_template_extension.go` state refresh function to drop the now-unused `DescribeSecurityTemplateBindingsRequest` reference and to return an empty state string when `resp.Status` is nil (avoids nil dereference during create polling).
- Update `resource_tc_teo_bind_security_template.go` read method to print the resource id (`[CRUD] teo_bind_security_template id=%s`) before clearing the id when the binding is not found, preserving the id for log diagnostics.
- Add gomonkey-based unit tests in `resource_tc_teo_bind_security_template_test.go` covering: read success, read not-found (id cleared), read with no zones, read with >100 zones (batching verified), and schema validation.
- Update `.changelog/4261.txt` to reflect the behavior-preserving API replacement (bug-fix / enhancement to read reliability).

## Capabilities

### New Capabilities
- `teo-bind-security-template-read-api`: Read path for the `tencentcloud_teo_bind_security_template` resource, using `DescribeZones` + `DescribeWebSecurityTemplates` to locate a template-to-domain binding and return its delivery status.

### Modified Capabilities
<!-- None. The user-facing schema and lifecycle of the resource are unchanged. -->

## Impact

- Affected code:
  - `tencentcloud/services/teo/service_tencentcloud_teo.go` (`DescribeTeoBindSecurityTemplateById`, new `describeTeoAllZoneIds`)
  - `tencentcloud/services/teo/resource_tc_teo_bind_security_template.go` (read log ordering)
  - `tencentcloud/services/teo/resource_tc_teo_bind_security_template_extension.go` (state refresh cleanup)
  - `tencentcloud/services/teo/resource_tc_teo_bind_security_template_test.go` (new unit tests)
- Cloud APIs: drops usage of `DescribeSecurityTemplateBindings`; uses `DescribeZones` and `DescribeWebSecurityTemplates` (both already used elsewhere in the provider).
- Backward compatible: no schema changes, no state migration, no breaking changes to Terraform configurations.
