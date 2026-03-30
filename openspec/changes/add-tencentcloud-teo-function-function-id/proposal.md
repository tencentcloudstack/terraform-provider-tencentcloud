## Why

<!-- Explain the motivation for this change. What problem does it solve? Why now? -->
To support specifying the FunctionId parameter when creating a tencentcloud_teo_function resource. Currently, function_id is a computed field that is only returned by the CreateFunction API. Users need the ability to specify a custom FunctionId during resource creation for better integration with existing edge function workflows and to support import scenarios where the FunctionId is known in advance.

## What Changes

<!-- Describe what will change. Be specific about new capabilities, modifications, or removals. -->
- Update the `function_id` field in `tencentcloud_teo_function` resource from `Computed` to `Optional` to allow users to specify it during resource creation
- Modify the CreateFunction request to include the FunctionId parameter if provided by the user
- Ensure backward compatibility: existing resources without function_id specified will continue to work as before (API will generate the ID)

## Capabilities

### New Capabilities
<!-- Capabilities being introduced. Replace <name> with kebab-case identifier (e.g., user-auth, data-export, api-rate-limiting). Each creates specs/<name>/spec.md -->
- `teo-function-id-parameter`: Add support for optional FunctionId parameter in tencentcloud_teo_function resource, enabling users to specify a custom function ID during creation or import.

### Modified Capabilities
<!-- Existing capabilities whose REQUIREMENTS are changing (not just implementation).
     Only list here if spec-level behavior changes. Each needs a delta spec file.
     Use existing spec names from openspec/specs/. Leave empty if no requirement changes. -->

## Impact

<!-- Affected code, APIs, dependencies, systems -->
- **Resource**: `tencentcloud/services/teo/resource_tc_teo_function.go` - Schema and create function will be modified
- **Schema**: Change `function_id` from Computed to Optional while keeping Computed for API-generated values
- **API**: TeoV20220901Client CreateFunction API call will include optional FunctionId parameter
- **Tests**: `resource_tc_teo_function_test.go` - Test cases will be updated to cover both scenarios (with and without function_id)
- **Documentation**: `resource_tc_teo_function.md` - Documentation will be updated to reflect the new optional parameter
- **Backward Compatibility**: Fully backward compatible - existing configurations will continue to work without changes
