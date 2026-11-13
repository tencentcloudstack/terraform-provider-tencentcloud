## Why

Tencent Cloud Config (配置审计) supports delivering audit logs to COS or CLS, but the current `tencentcloud_config_deliver_config` resource uses `helper.BuildToken()` as a singleton ID with a no-op Delete. As a dedicated RESOURCE_KIND_CONFIG resource, `tencentcloud_config_update_config_deliver` manages the delivery settings strictly through Read/Update — updating the configuration on apply and reading it back into state — without an artificial delete that would leave the remote config untouched. This gives users a clean config-type lifecycle that matches the "config exists as long as resource exists" model.

## What Changes

- Add a new Terraform resource `tencentcloud_config_update_config_deliver` (RESOURCE_KIND_CONFIG) under `tencentcloud/services/config/`.
  - File: `resource_tc_config_update_config_deliver.go` (Create=Update call to `UpdateConfigDeliver`, Read from `DescribeConfigDeliver`, Update via `UpdateConfigDeliver`, Delete=no-op).
  - Schema fields: `status` (required, int, 0/1), `deliver_name` (optional, string), `target_arn` (optional, string), `deliver_prefix` (optional, string), `deliver_type` (optional, string), `deliver_content_type` (optional, int), `create_time` (computed, string).
  - Resource ID: `helper.BuildToken()` (singleton global config — no natural unique key).
- Add service-layer method `DescribeConfigDeliver` already exists in `service_tencentcloud_config.go`; reuse it for Read.
- Add unit test file `resource_tc_config_update_config_deliver_test.go` using gomock (mock cloud API), per project rules for new resources.
- Add documentation file `resource_tc_config_update_config_deliver.md`.
- Register the new resource in `tencentcloud/provider.go` ResourcesMap and add it to the provider.go file comment index.

## Capabilities

### New Capabilities

- `config-update-config-deliver`: Manage the global Config delivery (投递设置) configuration via Read/Update only, as a RESOURCE_KIND_CONFIG resource.

### Modified Capabilities

<!-- None — this is a new resource, no existing spec requirements change. -->

## Impact

- **New files**: 
  - `tencentcloud/services/config/resource_tc_config_update_config_deliver.go`
  - `tencentcloud/services/config/resource_tc_config_update_config_deliver_test.go`
  - `tencentcloud/services/config/resource_tc_config_update_config_deliver.md`
- **Modified files**:
  - `tencentcloud/provider.go` (register resource + comment index)
- **Cloud APIs used**: `DescribeConfigDeliver` (Read), `UpdateConfigDeliver` (Create/Update) from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802`.
- **No breaking changes** — purely additive.