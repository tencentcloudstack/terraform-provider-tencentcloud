## Why

The `tencentcloud_teo_bind_security_template` resource currently exposes only `status` from the template binding query result. However, the cloud API `DescribeSecurityTemplateBindings` returns an additional `Message` field in `EntityStatus` that provides human-readable delivery status information (e.g., failure reasons). Users need access to this `Message` field to diagnose configuration delivery issues without having to check the cloud console.

## What Changes

- Add a new `message` computed attribute to the `tencentcloud_teo_bind_security_template` resource, sourced from `EntityStatus.Message` in the `DescribeSecurityTemplateBindings` API response.
- Modify the service layer `DescribeTeoBindSecurityTemplateById` function to use the `DescribeSecurityTemplateBindings` API instead of the current `DescribeWebSecurityTemplates` approach, which does not expose the `Message` field.

## Capabilities

### New Capabilities
- `teo-bind-security-template-message`: Expose the `Message` field from `DescribeSecurityTemplateBindings` API response as a computed attribute on `tencentcloud_teo_bind_security_template`.

### Modified Capabilities
- `teo-bind-security-template-read-api`: The `DescribeTeoBindSecurityTemplateById` service function SHALL use `DescribeSecurityTemplateBindings` API instead of `DescribeWebSecurityTemplates` to fetch the binding status, enabling access to the `Message` field.

## Impact

- **Affected files**:
  - `tencentcloud/services/teo/resource_tc_teo_bind_security_template.go` — add `message` schema field and set it in Read
  - `tencentcloud/services/teo/service_tencentcloud_teo.go` — modify `DescribeTeoBindSecurityTemplateById` to use `DescribeSecurityTemplateBindings`
  - `tencentcloud/services/teo/resource_tc_teo_bind_security_template.md` — add `message` to example usage
  - `tencentcloud/services/teo/resource_tc_teo_bind_security_template_test.go` — add test coverage for `message`
- **API change**: Switch from `DescribeWebSecurityTemplates` to `DescribeSecurityTemplateBindings` in the service layer
- **Backward compatibility**: Fully backward compatible — the new `message` field is computed-only, and the existing `status` field behavior is preserved