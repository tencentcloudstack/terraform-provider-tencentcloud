## Why

The `tencentcloud_reserve_ip_address` resource currently defines `subnet_id` as an `Optional` schema field and sets it during Create, but the Read method does not populate it from the `DescribeReserveIpAddresses` API response. This means the `subnet_id` value is not refreshed from the cloud after creation, causing state drift on refresh/import. The cloud API `DescribeReserveIpAddresses` returns `SubnetId` in each `ReserveIpAddressSet` item, and this field should be read back into Terraform state.

## What Changes

- Modify the `subnet_id` schema field in `tencentcloud_reserve_ip_address` from `Optional` to `Optional + Computed` so that Terraform can refresh it from the cloud API when the user does not specify it.
- Add `d.Set("subnet_id", reserveIpAddress.SubnetId)` in the Read method of `resource_tc_reserve_ip_address.go` to populate `subnet_id` from the `DescribeReserveIpAddresses` response (`response.ReserveIpAddressSet[0].SubnetId`).
- Remove `subnet_id` from the `ImportStateVerifyIgnore` list in the acceptance test so import verification covers this field.
- Update the resource `.md` example file if needed (no structural change, `subnet_id` already documented).

## Capabilities

### New Capabilities

- `vpc-reserve-ip-address-subnet-id`: Read-back support for the `subnet_id` field of `tencentcloud_reserve_ip_address` from the `DescribeReserveIpAddresses` API response.

### Modified Capabilities

None.

## Impact

- **Code**: `tencentcloud/services/vpc/resource_tc_reserve_ip_address.go` — schema change (`subnet_id` becomes `Computed`) and Read method update.
- **Tests**: `tencentcloud/services/vpc/resource_tc_reserve_ip_address_test.go` — remove `subnet_id` from `ImportStateVerifyIgnore`.
- **Documentation**: `tencentcloud/services/vpc/resource_tc_reserve_ip_address.md` — no change needed (`subnet_id` already in example).
- **API**: Uses existing `DescribeReserveIpAddresses` response field `ReserveIpAddressSet[].SubnetId` (already present in vendor SDK `ReserveIpAddressInfo` struct).
- **Backward compatibility**: Fully compatible — `subnet_id` remains `Optional`, existing configs that set it are unaffected; adding `Computed` allows refresh when unset.
