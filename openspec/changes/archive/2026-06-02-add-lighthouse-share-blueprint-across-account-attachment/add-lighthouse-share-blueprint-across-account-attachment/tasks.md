## 1. Resource Implementation (Core)

- [x] 1.1 Create resource file `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_attachment.go` with package, imports, and `ResourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachment()` function
- [x] 1.2 Define Schema with `blueprint_id` (Required, ForceNew, TypeString) and `account_ids` (Required, ForceNew: false, TypeSet of TypeString)
- [x] 1.3 Implement Create function: call `ShareBlueprintAcrossAccounts` API with retry in batches of 10 accounts per request, set ID as `blueprint_id`
- [x] 1.4 Implement Read function: use `blueprint_id` from ID, call `DescribeBlueprintsShareAcrossAccountInfos` API with pagination support (Offset/Limit loop), update or clear state
- [x] 1.5 Implement Update function: calculate diff of account_ids, batch-call Cancel/Share APIs for removed/added accounts
- [x] 1.6 Implement Delete function: get account_ids from state, batch-call `CancelShareBlueprintAcrossAccounts` API with retry
- [x] 1.7 Add `Importer` with `schema.ImportStatePassthrough`
- [x] 1.8 Add constant `shareBlueprintBatchSize = 10` for write API batch size limit

## 2. Provider Registration

- [x] 2.1 Register resource in `tencentcloud/provider.go` under Lighthouse resources section as `tencentcloud_lighthouse_share_blueprint_across_account_attachment`
- [x] 2.2 Add resource entry in `tencentcloud/provider.md`

## 3. Unit Testing

- [x] 3.1 Create test file `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_attachment_test.go`
- [x] 3.2 Write unit test for Create scenario using gomonkey to mock `ShareBlueprintAcrossAccounts` API
- [x] 3.3 Write unit test for Read scenario using gomonkey to mock `DescribeBlueprintsShareAcrossAccountInfos` API
- [x] 3.4 Write unit test for Update scenario using gomonkey to mock both Add and Remove APIs
- [x] 3.5 Write unit test for Delete scenario using gomonkey to mock `CancelShareBlueprintAcrossAccounts` API
- [x] 3.6 Write unit test for Import scenario
- [x] 3.7 Run unit tests with `go test -gcflags=all=-l` and ensure all tests pass

## 4. Documentation

- [x] 4.1 Create `.md` documentation file at `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_attachment.md` with description, Example Usage, and Import sections
- [ ] 4.2 Run `make doc` to generate website documentation at `website/docs/r/lighthouse_share_blueprint_across_account_attachment.html.markdown`

## 5. Code Quality Verification

- [x] 5.1 Verify the resource file compiles correctly
- [x] 5.2 Verify all imports are used and no unused variables
- [x] 5.3 Verify all error returns are handled properly
- [x] 5.4 Verify nil checks on all API response pointer fields
