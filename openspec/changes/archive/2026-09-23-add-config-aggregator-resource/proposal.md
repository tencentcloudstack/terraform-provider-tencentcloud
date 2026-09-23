## Why

Tencent Cloud Config (配置审计) supports "账号组" (Aggregators) to group multiple member accounts for unified compliance management. Currently no Terraform resource exists to manage Aggregators, so users must create/edit/delete them manually. This change adds the `tencentcloud_config_aggregator` resource to manage the full Aggregator lifecycle via Terraform.

## What Changes

- Add new Terraform resource `tencentcloud_config_aggregator` (RESOURCE_KIND_GENERAL) under `tencentcloud/services/config/resource_tc_config_aggregator.go`.
- Implement full CRUD using config SDK v20220802:
  - Create → `CreateAggregator` (returns `AccountGroupId`).
  - Read → `DescribeAggregator` (query by `AccountGroupId` + `OwnerUin`).
  - Update → `UpdateAggregator` (edit name/description/aggregator_accounts).
  - Delete → `DeleteAggregators`.
- Resource ID is a composite `account_group_id#owner_uin` (joined with `tccommon.FILED_SP`).
- Register `tencentcloud_config_aggregator` in `provider.go` and `provider.md`.
- Add `resource_tc_config_aggregator.md` doc and `resource_tc_config_aggregator_test.go` (gomonkey-based unit tests).

## Capabilities

### New Capabilities
- `config-aggregator-resource`: Manage a Config Aggregator (账号组) lifecycle — create, read, update, delete — including nested member account list and creation status.

### Modified Capabilities
<!-- None -->

## Impact

- **New files**: `tencentcloud/services/config/resource_tc_config_aggregator.go`, `resource_tc_config_aggregator.md`, `resource_tc_config_aggregator_test.go`.
- **Modified files**: `tencentcloud/provider.go` (register resource), `tencentcloud/provider.md` (doc entry); service layer `service_tencentcloud_config.go` may add a `DescribeConfigAggregatorById` helper.
- **SDK**: uses already-vendored `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802`.
- **Breaking changes**: none (purely additive).