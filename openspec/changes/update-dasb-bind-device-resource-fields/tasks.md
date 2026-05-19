## 1. Bug Fixes in Read Function

- [x] 1.1 Fix `device_id_set` population: dereference `item.Id` (`*uint64`) and convert to `int` before appending to `tmpList`
- [x] 1.2 Fix `domain_id` assignment: move `d.Set("domain_id", ...)` outside the per-device loop; set it once from the first device with a non-nil `DomainId`

## 2. Schema — Add New Fields

- [x] 2.1 Add `domain_name` as `Computed: true, TypeString` to the resource schema
- [x] 2.2 Add `manage_dimension` as `Optional: true, TypeInt` to the resource schema (with description noting SDK upgrade prerequisite)
- [x] 2.3 Add `manage_account_id` as `Optional: true, TypeInt` to the resource schema
- [x] 2.4 Add `manage_account` as `Optional: true, TypeString` to the resource schema
- [x] 2.5 Add `manage_kubeconfig` as `Optional: true, Sensitive: true, TypeString` to the resource schema
- [x] 2.6 Add `namespace` as `Optional: true, TypeString` to the resource schema
- [x] 2.7 Add `workload` as `Optional: true, TypeString` to the resource schema

## 3. Read Function — Map New Response Fields

- [x] 3.1 In the Read loop, set `domain_name` from `item.DomainName` (first non-nil value) into state

## 4. Create / Update — Wire New Fields (requires SDK upgrade)

> **Blocked**: The vendored SDK `dasb/v20191018` `BindDeviceResourceRequest` does not yet contain `ManageDimension`, `ManageAccountId`, `ManageAccount`, `ManageKubeconfig`, `Namespace`, or `Workload` fields.
> **Prerequisite**: Run `go get github.com/tencentcloud/tencentcloud-sdk-go@<new-version>` and `go mod vendor` to update the SDK, then verify the new fields exist in `BindDeviceResourceRequest`.

- [x] 4.1 Upgrade vendored SDK to a version that includes the new `BindDeviceResourceRequest` fields
- [x] 4.2 In Create: add `GetOk` blocks for `manage_dimension`, `manage_account_id`, `manage_account`, `manage_kubeconfig`, `namespace`, `workload` and assign to `request`
- [x] 4.3 In Update (`d.HasChange("device_id_set")` add-branch): propagate the same six fields to the add-request
- [x] 4.4 Add the six new fields to the `mutableArgs` check list in Update so changes to them trigger a re-bind

## 5. Verification

- [x] 5.1 Run `go build ./tencentcloud/services/bh/` to confirm no compilation errors
- [ ] 5.2 Manually verify `terraform plan` shows no spurious drift for `device_id_set` and `domain_id` on an existing imported resource
