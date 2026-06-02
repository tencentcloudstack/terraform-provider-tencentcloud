## Why

The terraform-provider-tencentcloud currently lacks a resource for managing Lighthouse blueprint cross-account sharing. While the Lighthouse SDK provides `ShareBlueprintAcrossAccounts`, `DescribeBlueprintsShareAcrossAccountInfos`, and `CancelShareBlueprintAcrossAccounts` APIs for sharing blueprints to other TencentCloud accounts, there is no Terraform resource to manage these sharing relationships declaratively. Users must manually call the API or use the console to share/unshare blueprints across accounts.

## What Changes

- Add a new Terraform resource `tencentcloud_lighthouse_share_blueprint_across_account_attachment` (RESOURCE_KIND_ATTACHMENT) that manages the binding relationship of sharing a Lighthouse blueprint to target accounts
- Support Create (share blueprint across accounts), Read (query sharing info), and Delete (cancel sharing) operations
- The resource uses `blueprint_id` and `account_ids` to define the sharing relationship, with a composite ID combining both values

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
