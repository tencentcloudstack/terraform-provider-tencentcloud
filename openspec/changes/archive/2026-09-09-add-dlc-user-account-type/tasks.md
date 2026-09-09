## 1. Schema Definition

- [x] 1.1 Add `account_type` field to resource schema in `tencentcloud/services/dlc/resource_tc_dlc_user.go`
  - Type: `schema.TypeString`
  - Optional: true
  - Description: "Account type. Valid values: `UserAccount` (user account), `RoleAccount` (role account). Default is `UserAccount`."

## 2. Service Layer

- [x] 2.1 Update `DescribeDlcUserById` method in `tencentcloud/services/dlc/service_tencentcloud_dlc.go` to accept an `accountType string` parameter and set `request.AccountType` when non-empty
- [x] 2.2 Update `DeleteDlcUserById` method in `tencentcloud/services/dlc/service_tencentcloud_dlc.go` to accept an `accountType string` parameter and set `request.AccountType` when non-empty

## 3. CRUD Functions

- [x] 3.1 Update `resourceTencentCloudDlcUserCreate` to read `account_type` from schema and set `request.AccountType` in the `CreateUser` API request
- [x] 3.2 Update `resourceTencentCloudDlcUserRead` to read `AccountType` from the `UserInfo` response (nil-checked) and set `account_type` in Terraform state via `d.Set()`
- [x] 3.3 Update `resourceTencentCloudDlcUserUpdate` to include `account_type` in the `ModifyUser` API request when `d.HasChange("account_type")` is true; do NOT add `account_type` to the `immutableArgs` array
- [x] 3.4 Update `resourceTencentCloudDlcUserDelete` to read `account_type` from schema and pass it to the updated `DeleteDlcUserById` service method
- [x] 3.5 Update `DescribeDlcUserById` callers in read function to pass `account_type` from schema

## 4. Testing

- [x] 4.1 Add test case in `tencentcloud/services/dlc/resource_tc_dlc_user_test.go` for creating a DLC user with `account_type` set
- [x] 4.2 Add test case for updating `account_type` in-place and verifying the value is correctly updated
- [x] 4.3 Add test assertion for reading `account_type` from state after create/update
- [x] 4.4 Add test case for importing a DLC user and verifying `account_type` is populated

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/dlc/resource_tc_dlc_user.md` to include `account_type` parameter in example usage
- [ ] 5.2 Run `make doc` (in finalize phase) to generate website documentation

## 6. Code Quality (Finalize Phase)

- [ ] 6.1 Run `gofmt` to format all modified Go files
- [x] 6.2 Verify all error returns are properly handled (use `_ =` for functions that cannot fail)
- [ ] 6.3 Create changelog file via `tfpacer-finalize` skill
