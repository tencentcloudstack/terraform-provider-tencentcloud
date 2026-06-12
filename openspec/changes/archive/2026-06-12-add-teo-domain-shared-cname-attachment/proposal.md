## Why

TencentCloud TEO (EdgeOne) supports binding acceleration domains to shared CNAMEs, but there is currently no Terraform resource to manage this binding relationship. Users need a declarative way to bind/unbind domains to shared CNAMEs through Terraform.

## What Changes

- Add a new Terraform resource `tencentcloud_teo_domain_shared_cname_attachment` of type RESOURCE_KIND_ATTACHMENT
- The resource manages the binding relationship between acceleration domains and shared CNAMEs using the `BindSharedCNAME` API (bind on create, unbind on delete)
- The resource reads binding status via the `DescribeSharedCNAME` API to verify the attachment exists
- Resource ID is a composite of `zone_id`, `shared_cname`, and `domain_names` joined by `tccommon.FILED_SP`

## Capabilities

### New Capabilities
- `teo-domain-shared-cname-attachment`: Terraform resource to manage the binding/unbinding of acceleration domains to a shared CNAME within a TEO zone

### Modified Capabilities

## Impact

- New resource file: `tencentcloud/services/teo/resource_tc_teo_domain_shared_cname_attachment.go`
- New test file: `tencentcloud/services/teo/resource_tc_teo_domain_shared_cname_attachment_test.go`
- New doc file: `tencentcloud/services/teo/resource_tc_teo_domain_shared_cname_attachment.md`
- Modified: `tencentcloud/provider.go` (register the new resource)
- Modified: `tencentcloud/provider.md` (add resource entry)
- Cloud API dependency: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` (already vendored)
