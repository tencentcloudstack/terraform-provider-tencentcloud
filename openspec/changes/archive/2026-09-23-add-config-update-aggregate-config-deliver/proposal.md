## Why

Users of Tencent Cloud Config who manage resources across member accounts via 账号组 (Aggregate / Account Group) currently have no Terraform support for configuring the 投递设置 (delivery settings) of an aggregate. They must use the console to manage how Config audit logs are delivered to COS or CLS for member accounts grouped under an account group. Adding this resource lets users declare and manage the aggregate delivery configuration as code, enabling compliance auditing across the account group.

## What Changes

- Add new Terraform resource `tencentcloud_config_update_aggregate_config_deliver` of kind RESOURCE_KIND_CONFIG for managing the delivery settings of a Tencent Cloud Config account group (账号组).
- The resource provides RU (Read + Update) operations — no dedicated Create or Delete API exists. First-time `terraform apply` calls the Update API to set the delivery configuration; the resource persists as long as the account group exists.
- Schema fields: `account_group_id` (required), `status` (required), and optional fields `deliver_name`, `target_arn`, `deliver_prefix`, `deliver_type`, `deliver_uin`, `deliver_content_type`; computed field `create_time`.
- Register the resource in `provider.go` / `provider.md`.
- Add a service-layer method wrapping `DescribeAggregateConfigDeliver` in `service_tencentcloud_config.go`.
- Add documentation (`resource_tc_config_update_aggregate_config_deliver.md`) and unit tests using gomonkey mocks.

## Capabilities

### New Capabilities
- `config-update-aggregate-config-deliver`: Manage the delivery settings (投递设置) of a Tencent Cloud Config account group (账号组), including delivery target (COS/CLS), delivery prefix, status switch, and content type. Supports RU lifecycle (Read via DescribeAggregateConfigDeliver, Update via UpdateAggregateConfigDeliver).

### Modified Capabilities
<!-- None -->

## Impact

- **New files**:
  - `tencentcloud/services/config/resource_tc_config_update_aggregate_config_deliver.go`
  - `tencentcloud/services/config/resource_tc_config_update_aggregate_config_deliver.md`
  - `tencentcloud/services/config/resource_tc_config_update_aggregate_config_deliver_test.go`
- **Modified files**:
  - `tencentcloud/services/config/service_tencentcloud_config.go` — append `DescribeAggregateConfigDeliver` service method
  - `tencentcloud/provider.go` — register `tencentcloud_config_update_aggregate_config_deliver` in ResourcesMap
  - `tencentcloud/provider.md` — add resource name to documentation list
- **Cloud APIs**: `DescribeAggregateConfigDeliver` (Read), `UpdateAggregateConfigDeliver` (Update) from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802`.
- **Resource ID**: Uses `account_group_id` as the resource ID (one delivery config per account group).
- **No breaking changes** — purely additive.
