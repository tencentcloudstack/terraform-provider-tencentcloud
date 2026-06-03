## Why

The terraform-provider-tencentcloud currently lacks a resource for managing Lighthouse blueprint cross-account sharing. While the Lighthouse SDK provides `ShareBlueprintAcrossAccounts`, `DescribeBlueprintsShareAcrossAccountInfos`, and `CancelShareBlueprintAcrossAccounts` APIs for sharing blueprints to other TencentCloud accounts, there is no Terraform resource to manage these sharing relationships declaratively. Users must manually call the API or use the console to share/unshare blueprints across accounts.

## What Changes

- Add a new Terraform resource `tencentcloud_lighthouse_share_blueprint_across_account_attachment` (RESOURCE_KIND_ATTACHMENT) that manages the binding relationship of sharing a Lighthouse blueprint to target accounts
- Support Create (share blueprint across accounts), Read (query sharing info with pagination), Update (incrementally add/remove accounts), and Delete (cancel sharing) operations
- The resource uses `blueprint_id` as the resource ID for simple import, with `account_ids` as an updatable field
- **Design Decision**: `account_ids` is defined as `TypeSet` (not `TypeList`) to avoid ordering issues — since the API returns account IDs in non-deterministic order, using Set ensures Terraform state remains stable without manual sorting
- **Pagination**: Read function implements pagination when calling `DescribeBlueprintsShareAcrossAccountInfos`, looping through pages using `Offset`/`Limit` until all shared account IDs are collected
- **Batch Processing**: Create, Update, and Delete functions process account IDs in batches of 10 per API call, since both `ShareBlueprintAcrossAccounts` and `CancelShareBlueprintAcrossAccounts` APIs limit each request to a maximum of 10 account IDs

## Capabilities

### New Capabilities
- `lighthouse-share-blueprint-across-account-attachment`: New attachment resource for managing Lighthouse blueprint cross-account sharing relationships through Terraform

### Modified Capabilities
<!-- No existing capabilities are modified by this change -->

## Impact

- **New files**:
  - `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_attachment.go` - Resource implementation
  - `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_attachment_test.go` - Unit tests
  - `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_attachment.md` - Documentation example
- **Modified files**:
  - `tencentcloud/provider.go` - Resource registration
  - `tencentcloud/provider.md` - Resource documentation entry
- **SDK Dependencies**: Uses existing `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324` package (already vendored)
- **Breaking Changes**: None - this is a new resource
