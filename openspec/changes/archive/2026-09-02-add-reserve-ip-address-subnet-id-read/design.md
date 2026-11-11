## Context

The `tencentcloud_reserve_ip_address` resource manages VPC reserved IP addresses via the TencentCloud VPC API. The resource schema currently defines `subnet_id` as an `Optional` (non-Computed) string field. During Create, `subnet_id` is passed to `CreateReserveIpAddressesRequest.SubnetId`. However, the Read method (`resourceTencentCloudReserveIpAddressRead`) does not call `d.Set("subnet_id", ...)` from the `DescribeReserveIpAddresses` response, meaning the value is never refreshed from the cloud.

The cloud API `DescribeReserveIpAddresses` returns a `ReserveIpAddressSet` array where each `ReserveIpAddressInfo` item contains a `SubnetId` field (already present in the vendored SDK at `vendor/.../vpc/v20170312/models.go`). This field is available for read-back.

Additionally, the acceptance test explicitly ignores `subnet_id` during `ImportStateVerify` (`ImportStateVerifyIgnore: []string{"subnet_id"}`), which masks the fact that the field is not read back.

## Goals / Non-Goals

**Goals:**
- Ensure `subnet_id` is populated in Terraform state from the `DescribeReserveIpAddresses` API response during Read/refresh.
- Allow `subnet_id` to be refreshed even when the user did not explicitly set it in configuration (by marking it `Computed`).
- Enable import verification to cover `subnet_id`.

**Non-Goals:**
- Changing the Create or Update behavior for `subnet_id` (it remains `Optional`, not `ForceNew`; it is already in the `immutableArgs` list).
- Adding any new API calls — the existing `DescribeReserveIpAddresses` already returns `SubnetId`.
- Modifying the `tags` field or its handling (already fully implemented via tag service).

## Decisions

### 1. Mark `subnet_id` as `Computed` in addition to `Optional`

**Decision**: Change the schema from `Optional: true` to `Optional: true, Computed: true`.

**Rationale**: The cloud API may assign a `SubnetId` even when the user did not specify one (or the VPC default subnet applies). Marking `Computed` allows Terraform to store the value read from the API without requiring the user to set it. This is the standard pattern for fields that are user-settable on create but also returned by the read API.

**Alternatives considered**: Keep `Optional` only and always set from response — this would cause Terraform plan diffs when the user did not set `subnet_id` because the provider would write a value the user didn't declare. Adding `Computed` suppresses this diff.

### 2. Add `d.Set("subnet_id", ...)` in Read method

**Decision**: Add `_ = d.Set("subnet_id", reserveIpAddress.SubnetId)` in the Read method, after the existing `d.Set` calls (around line 200-208 of `resource_tc_reserve_ip_address.go`), guarded by a nil check on `reserveIpAddress.SubnetId` per project conventions.

**Rationale**: Follows the existing pattern used for `vpc_id`, `ip_address`, `name`, `description`, etc. The `ReserveIpAddressInfo.SubnetId` field is a `*string`, so a nil check is needed before setting.

### 3. Remove `subnet_id` from `ImportStateVerifyIgnore`

**Decision**: Remove `"subnet_id"` from the `ImportStateVerifyIgnore` list in both test cases in `resource_tc_reserve_ip_address_test.go`.

**Rationale**: Once `subnet_id` is `Computed` and read back, import verification should cover it. Keeping it ignored would hide regressions.

## Risks / Trade-offs

- **[Schema change to Computed]** → Adding `Computed` to an existing `Optional` field is backward compatible. Existing configs that set `subnet_id` continue to work; configs that don't set it will now have it populated from the API. No state migration is needed.
- **[Nil SubnetId in response]** → If the API returns `SubnetId` as nil (e.g., for a reserved IP not associated with a subnet), the nil guard prevents a panic. The field will simply not be set.
- **[Import test change]** → Removing `subnet_id` from `ImportStateVerifyIgnore` means the acceptance test now validates this field on import. If the API returns an empty SubnetId for some edge case, the test may fail. This is acceptable as it reflects real behavior.
