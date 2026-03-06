# Change: Add Public Network Access Fields to TDMQ RabbitMQ VIP Instance

## Why
The Tencent Cloud API `CreateRabbitMQVipInstance` supports public network access configuration through `Bandwidth` and `EnablePublicAccess` fields. Currently, the Terraform resource `tencentcloud_tdmq_rabbitmq_vip_instance` does not expose these fields, preventing users from configuring public network access during instance creation. This gap limits the resource's functionality and requires manual configuration after creation.

## What Changes
- Add `band_width` field to control public network bandwidth (in Mbps)
- Add `enable_public_access` field to enable/disable public network access
- Both fields are immutable after creation (cannot be modified through Update operation)
- Fields are optional with default behavior:
  - `enable_public_access` defaults to `false` (no public access)
  - `band_width` only takes effect when `enable_public_access` is `true`
- Read operation maps values from API response:
  - `band_width` reads from `rabbitmqVipInstance.ClusterSpecInfo.PublicNetworkTps`
  - `enable_public_access` reads from `rabbitmqVipInstance.ClusterNetInfo.PublicDataStreamStatus` (ON → true, OFF → false)

## Impact
- **Affected specs**: `tdmq-rabbitmq-vip-instance`
- **Affected code**:
  - `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - Add schema fields, Create/Read logic, Update validation
  - `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown` - Update documentation with new fields
- **Breaking change**: No - New optional fields with backward-compatible defaults
- **API compatibility**: Requires SDK version supporting these fields (already upgraded in current codebase)
