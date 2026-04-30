## Why

The `tencentcloud_teo_content_identifier` resource currently lacks the `status` field from the cloud API's `ContentIdentifier` response struct. Users cannot observe the content identifier's status (active/deleted) through Terraform state, which limits their ability to monitor resource lifecycle and troubleshoot issues.

## What Changes

- Add a new computed field `status` (TypeString) to the `tencentcloud_teo_content_identifier` resource schema
- The `status` field maps to the `Status` field in the cloud API's `ContentIdentifier` struct, returned by the `DescribeContentIdentifiers` API
- Update the Read method to populate `status` from the API response
- No changes to Create/Update/Delete methods since `status` is a read-only computed field

## Capabilities

### New Capabilities
- `teo-content-identifier-status`: Adds the `status` computed field to the `tencentcloud_teo_content_identifier` resource, exposing the content identifier's lifecycle status (active/deleted) from the TEO cloud API

### Modified Capabilities

## Impact

- **Code**: `tencentcloud/services/teo/resource_tc_teo_content_identifier.go` — schema definition and Read method
- **Tests**: `tencentcloud/services/teo/resource_tc_teo_content_identifier_test.go` — add test coverage for the new `status` field
- **Docs**: `tencentcloud/services/teo/resource_tc_teo_content_identifier.md` — update example to reflect the new computed field
- **API**: Uses existing `DescribeContentIdentifiers` API — no new API calls required
- **Backward Compatibility**: Fully backward compatible — adding a computed field does not affect existing configurations
