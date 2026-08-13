## Why

TencentCloud VOD (点播) has introduced AIGC (AI-Generated Content) quota management APIs (CreateAigcQuota, DescribeAigcQuotas, ModifyAigcQuota, DeleteAigcQuota) that allow customers to control usage limits for AI image/video/text generation tasks. Currently there is no Terraform resource to manage these quotas, requiring customers to use the console or API directly. A Terraform resource for this capability enables Infrastructure-as-Code management of AIGC quotas across VOD sub-applications.

## What Changes

- Add a new Terraform resource `tencentcloud_vod_aigc_quota` (resource type: RESOURCE_KIND_GENERAL) supporting full CRUD lifecycle management
- Implement Create/Read/Update/Delete operations backed by the VOD AIGC quota cloud APIs
- Support import of existing AIGC quotas using a composite ID format (`{sub_app_id}#{quota_type}#{api_token}`)
- Add service-layer methods in `service_tencentcloud_vod.go` for CreateAigcQuota, DescribeAigcQuotas, ModifyAigcQuota, and DeleteAigcQuota
- Register the new resource in `tencentcloud/provider.go`

## Capabilities

### New Capabilities
- `vod-aigc-quota-resource`: Terraform resource for managing VOD AIGC quotas, including create, read, update, delete, and import operations. Supports quota types Image, Video, and Text with optional API token filtering for Text quotas.

### Modified Capabilities
<!-- None -->

## Impact

- **New file**: `tencentcloud/services/vod/resource_tc_vod_aigc_quota.go` - Resource CRUD implementation
- **New file**: `tencentcloud/services/vod/resource_tc_vod_aigc_quota_test.go` - Unit tests
- **New file**: `tencentcloud/services/vod/resource_tc_vod_aigc_quota.md` - Documentation template
- **Modified file**: `tencentcloud/services/vod/service_tencentcloud_vod.go` - Add service-layer methods
- **Modified file**: `tencentcloud/provider.go` - Register new resource
- **Modified file**: `tencentcloud/provider.md` - Update provider documentation index
- **Cloud API dependency**: Uses existing SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717` (already vendored)
