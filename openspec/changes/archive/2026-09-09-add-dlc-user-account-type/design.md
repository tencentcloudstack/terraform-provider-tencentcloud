## Context

The `tencentcloud_dlc_user` resource manages DLC (Data Lake Compute) users via the TencentCloud DLC API (`dlc/v20210125`). The resource currently supports `user_id`, `user_description`, `user_type`, `user_alias`, and `work_group_ids` parameters.

The DLC cloud API has been updated to support an `AccountType` parameter across all CRUD operations:
- `CreateUser`: accepts `AccountType` as an optional input (account type: `UserAccount` for user accounts, `RoleAccount` for role accounts; defaults to user account)
- `DescribeUsers`: returns `AccountType` in the `UserInfo` struct within `UserSet`
- `ModifyUser`: accepts `AccountType` as an optional input
- `DeleteUser`: accepts `AccountType` as an optional input

This change adds the `account_type` parameter to the Terraform resource to expose this capability.

### Current State
- Resource file: `tencentcloud/services/dlc/resource_tc_dlc_user.go`
- Service layer: `tencentcloud/services/dlc/service_tencentcloud_dlc.go`
  - `DescribeDlcUserById(ctx, userId)` returns `*dlc.UserInfo`
  - `DeleteDlcUserById(ctx, userId)` deletes a user
- The `DescribeDlcUserById` method does not pass `AccountType` in the `DescribeUsers` request
- The `DeleteDlcUserById` method does not pass `AccountType` in the `DeleteUser` request

## Goals / Non-Goals

**Goals:**
- Add `account_type` as an optional string parameter to the `tencentcloud_dlc_user` resource schema
- Pass `AccountType` in all four CRUD API calls (`CreateUser`, `DescribeUsers`, `ModifyUser`, `DeleteUser`)
- Read `AccountType` from the `DescribeUsers` response and set it in Terraform state
- Support in-place updates of `account_type` via `ModifyUser`
- Maintain full backward compatibility (parameter is optional)

**Non-Goals:**
- Do not modify any existing parameters or their behavior
- Do not change the resource ID format (still uses `user_id` as the sole ID)
- Do not add `account_type` to the `immutableArgs` list — it supports in-place updates
- Do not add `ForceNew` to `account_type`

## Decisions

### 1. Parameter Naming
**Decision**: Use `account_type` (snake_case)

**Rationale**:
- Follows existing Terraform naming conventions in this resource (`user_id`, `user_type`, etc.)
- Matches the cloud API field name `AccountType` in snake_case form
- Consistent with the project's naming pattern

### 2. Parameter Type and Optionality
**Decision**: `schema.TypeString`, `Optional: true`, no `ForceNew`

**Rationale**:
- The cloud API field `AccountType` is `*string` in the SDK
- The API accepts it as optional in `CreateUser`, `ModifyUser`, and `DeleteUser`
- `ModifyUser` supports updating `AccountType`, so it should not be `ForceNew`
- Optional ensures backward compatibility (existing configs without `account_type` still work)

### 3. Service Layer Changes
**Decision**: Update `DescribeDlcUserById` and `DeleteDlcUserById` to accept an `accountType` parameter

**Rationale**:
- `DescribeDlcUserById`: The `DescribeUsers` request supports `AccountType` as a filter/input. Passing it ensures the correct user is queried (relevant for role vs. user accounts).
- `DeleteDlcUserById`: The `DeleteUser` request supports `AccountType`. Passing it ensures the correct user type is deleted.

**Implementation**:
- `DescribeDlcUserById(ctx, userId, accountType string)` — set `request.AccountType` when non-empty
- `DeleteDlcUserById(ctx, userId, accountType string)` — set `request.AccountType` when non-empty
- Update all callers in `resource_tc_dlc_user.go` to pass the `account_type` value from the schema

**Alternatives Considered**:
- Keep service methods unchanged and only pass `AccountType` in resource CRUD: Not feasible because the service layer constructs the API requests, so the parameter must be threaded through.

### 4. Read Operation
**Decision**: Read `AccountType` from `UserInfo.AccountType` in the `DescribeUsers` response and set it via `d.Set("account_type", ...)` with nil check

**Rationale**:
- The `UserInfo` struct contains `AccountType *string`
- Follow the existing pattern in the read function: check nil before calling `d.Set()`

### 5. Update Operation
**Decision**: Include `account_type` in the `ModifyUser` request when `d.HasChange("account_type")` is true. Do NOT add `account_type` to the `immutableArgs` array.

**Rationale**:
- The `ModifyUser` API accepts `AccountType`, so in-place updates are supported
- The existing `immutableArgs` only contains `user_type` and `user_alias`, which cannot be changed via `ModifyUser`
- Since `account_type` is supported by `ModifyUser`, it should allow updates

### 6. Delete Operation
**Decision**: Pass `account_type` from the schema to `DeleteDlcUserById`

**Rationale**:
- The `DeleteUser` API accepts `AccountType`, which may be needed to identify the correct user to delete
- Thread the value through the service layer

### 7. Test Strategy
**Decision**: Use Terraform acceptance test suite (since this is a modification to an existing resource, per project conventions)

**Rationale**:
- Project guidelines specify that modifications to existing resources use the Terraform test suite
- Add test checks for `account_type` in create, update, and import scenarios

## Risks / Trade-offs

- **[Risk] `DescribeDlcUserById` signature change breaks other callers** → Mitigation: Search the codebase for all callers of `DescribeDlcUserById` and `DeleteDlcUserById` and update them. Based on code review, these methods are only called from `resource_tc_dlc_user.go`, so the blast radius is limited to that file.
- **[Risk] `AccountType` may be returned as empty string from API** → Mitigation: Use the nil-check pattern (`if user.AccountType != nil`) before `d.Set()`, consistent with existing code.
- **[Risk] `ModifyUser` may not actually support changing `AccountType` despite the API field existing** → Mitigation: The API SDK includes `AccountType` in `ModifyUserRequest`, indicating the API accepts it. If the backend rejects it, the retry/error handling will surface the error to the user.
